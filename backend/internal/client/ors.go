package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"

	"github.com/bramlew/NEA/backend/internal/algorithm"
	"github.com/bramlew/NEA/backend/internal/models"
	"github.com/joho/godotenv"
)

const EndpointMatrix = "https://api.heigit.org/openrouteservice/v2/matrix/driving-hgv"
const EndpointDirections = "https://api.heigit.org/openrouteservice/v2/directions/driving-hgv"

func ORSRequest(payload []byte, endpoint string) (*http.Response, error) {
	// Make an API request to the given ORS endpoint
	if err := godotenv.Load(); err != nil {
		return nil, errors.New("error loading .env file")
	}
	apiKey := os.Getenv("ORS_API_KEY")
	client := &http.Client{Timeout: 5000}

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

	return res, nil
}

func MatrixRequest(origins []models.Coords, dests []models.Coords) (*models.Matrix, error) {
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
			Units:     "km",
		}
	} else {
		combined := slices.Concat(origins, dests)
		payload = models.MatrixPayload{
			Locations:    parseCoordsList(combined),
			Destinations: genRangeSlice(originsLength, len(combined)),
			Metrics:      []string{"distance"},
			Sources:      genRangeSlice(0, originsLength),
			Units:        "km",
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

	// Construct a matrix from the returned 2D slice
	cols := len(resStruct.Distances)
	mat := &models.Matrix{
		Matrix: make([]*models.Leg, cols*cols),
		Cols:   cols,
	}

	// Make the 2D slice into a 1D flat slice for faster indexing later
	for i, row := range resStruct.Distances {
		for j, dist := range row {
			mat.Matrix[algorithm.LookupIndex(i, j, cols)] = &models.Leg{Distance: dist}
		}
	}
	return mat, nil
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
	for i := 0; i < length; i++ {
		ints[i] = start + i
	}
	return ints
}
