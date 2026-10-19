package client

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"slices"
	"time"

	"github.com/bramlew/NEA/backend/internal/models"
	"github.com/goccy/go-json"
)

const EndpointMatrix = "https://api.heigit.org/openrouteservice/v2/matrix/driving-car"       // ORS API URL for Matrix endpoint
const EndpointPolyline = "https://api.heigit.org/openrouteservice/v2/directions/driving-car" // ORS API URL for Polyline (directions) endpoint
const EndpointGeocode = "https://api.heigit.org/pelias/v1/search"                            // ORS API URL for Geocoding
const StandardDistanceUnit = "km"                                                            // Standard distance unit used throughout payloads
const GeocodeRequestSize = 1                                                                 // No. of pairs of coordinates to be returned by the geocoding request

var client = &http.Client{
	// Create an HTTP client to be used for all ORS requests
	// Transport values are all default values from the http.DefaultTransport variable, with MaxIdleConnsPerHost being
	// modified to be 100, rather than 2, which crucially allows for lots of active TCP connections to the ORS API.
	Timeout: 5 * time.Second,
	Transport: &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	},
}

func ORSRequest(payload []byte, endpoint string) (*http.Response, error) {
	// Make an API request to the given ORS endpoint
	apiKey := os.Getenv("ORS_API_KEY")

	// Create the request and handle errors
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Add the required headers to the request
	req.Header.Add("Accept", "application/json, application/geo+json, application/gpx+xml, img/png; charset=utf-8")
	req.Header.Add("Authorization", apiKey)
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	// Execute the request and handle errors
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	statusCode := res.StatusCode
	if statusCode != http.StatusOK {
		res.Body.Close()
		return nil, fmt.Errorf("non-200 status code: %d", statusCode)
	}

	return res, nil
}

func MatrixRequest(origins []models.Coords, dests []models.Coords, originPairs []int, destPairs []int) ([]float64, error) {
	// Make an ORS Matrix request with given origins and destinations.
	originsLength := len(origins)
	size := originsLength * len(dests)

	// Firstly ensure that the matrix size is within the API limits, i.e. max 3500 elements in the matrix
	if size > 3500 {
		return nil, fmt.Errorf("matrix size must be less than 3500, got %v", size)
	}

	// Set the payload appropriately - if the destinations are not provided, just make a square matrix from the origins
	var payload models.MatrixPayload
	if dests == nil {
		payload = models.MatrixPayload{
			Locations: parseCoordsList(origins),
			Metrics:   []string{"distance"},
			Units:     StandardDistanceUnit,
		}
	} else {
		combined := slices.Concat(origins, dests)
		payload = models.MatrixPayload{
			Locations:    parseCoordsList(combined),
			Destinations: genRangeSlice(originsLength, len(combined)),
			Metrics:      []string{"distance"},
			Sources:      genRangeSlice(0, originsLength),
			Units:        StandardDistanceUnit,
		}
	}

	// Send the request
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error parsing payload: %v", err)
	}
	res, err := ORSRequest(jsonPayload, EndpointMatrix)
	if err != nil {
		return nil, err
	}

	// Read the returned HTTP response and parse it
	defer res.Body.Close()
	var resStruct models.MatrixResponse
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}
	err = json.Unmarshal(resBody, &resStruct)
	if err != nil {
		return nil, fmt.Errorf("error parsing json: %v", err)
	}

	// Get only the distances which were requested rather than the entire matrix
	resDists := resStruct.Distances
	dists := make([]float64, len(originPairs))
	for i, originPair := range originPairs {
		destPair := destPairs[i]
		dists[i] = resDists[originPair][destPair]
	}
	return dists, nil
}

func PolylineRequest(leg *models.Leg) (string, float64, error) {
	// Make an ORS polyline (i.e. directions) request for a given leg

	// Initialise the request payload
	payload := models.PolylinePayload{
		Coordinates:      parseCoordsList([]models.Coords{leg.Origin, leg.Dest}),
		GeometrySimplify: true,
		Instructions:     false,
		Units:            StandardDistanceUnit,
	}

	// Send the request
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", 0.0, fmt.Errorf("error parsing payload: %v", err)
	}
	res, err := ORSRequest(jsonPayload, EndpointPolyline)
	if err != nil {
		return "", 0.0, err
	}

	// Read the returned HTTP response and parse it
	defer res.Body.Close()
	var resStruct models.PolylineResponse
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return "", 0.0, fmt.Errorf("error reading response: %v", err)
	}
	err = json.Unmarshal(resBody, &resStruct)
	if err != nil {
		return "", 0.0, fmt.Errorf("error parsing json: %v", err)
	}

	// Return the details of the requested polyline
	routeDetails := resStruct.Routes[0]
	polyline, distance := routeDetails.Geometry, routeDetails.Summary.Distance
	return polyline, distance, nil
}

func GeocodeRequest(input string) (models.Coords, error) {
	// Make an ORS request to geocode some string
	apiKey := os.Getenv("ORS_API_KEY")

	// Create the request and handle errors
	fullUrl := fmt.Sprintf("%s?api_key=%s&text=%s&size=%d", EndpointGeocode, apiKey, url.QueryEscape(input), GeocodeRequestSize)
	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		return models.Coords{}, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Add("Accept", "application/json, application/geo+json, application/gpx+xml, img/png; charset=utf-8")

	// Execute the request and handle errors
	res, err := client.Do(req)
	if err != nil {
		return models.Coords{}, fmt.Errorf("error sending request: %v", err)
	}
	defer res.Body.Close()
	statusCode := res.StatusCode
	if statusCode != http.StatusOK {
		return models.Coords{}, fmt.Errorf("non-200 status code: %d", statusCode)
	}

	// Read the returned HTTP response and parse it
	var resStruct models.GeocodeResponse
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return models.Coords{}, fmt.Errorf("error reading response: %v", err)
	}
	err = json.Unmarshal(resBody, &resStruct)
	if err != nil {
		return models.Coords{}, fmt.Errorf("error parsing json: %v", err)
	}

	// Return the coordinate pair corresponding to the string requested
	coord := resStruct.Features[0].Geometry.Coordinates
	return models.Coords{
		Lon: coord[0],
		Lat: coord[1],
	}, nil
}

func parseCoordsList(coords []models.Coords) [][]float64 {
	// Parse a coords list to be passed into JSON
	parsed := make([][]float64, len(coords))
	for i, coord := range coords {
		parsed[i] = []float64{coord.Lon, coord.Lat}
	}
	return parsed
}

func genRangeSlice(start int, end int) []int {
	// Generate an int slice of a given range, with end being exclusive
	length := end - start
	ints := make([]int, length)
	for i := range length {
		ints[i] = start + i
	}
	return ints
}
