package algorithm

import (
	"log"
	"slices"
	"sync"

	"github.com/bramlew/NEA/backend/internal/client"
	"github.com/bramlew/NEA/backend/internal/models"
)

const LengthTolerance = 1.0 // Maximum error tolerance between the polyline distance and matrix distance
const RDPTolerance = 0.1    // Tolerance to be passed into RDP function
const MaxRoadDist = 2000.0  // Maximum distance before defaulting to using air travel
const MaxMatrixLength = 59  // Maximum length of a square matrix for an ORS request

func ConstructMultimodalMatrix(locations []models.Coords) (*models.Matrix, error) {
	// Calculate the multimodal distance between two sets of coordinates
	mat, err := ConstructMatrix(locations)
	if err != nil {
		return nil, err
	}
	origins := make([]int, 0, MaxMatrixLength)
	dests := make([]int, 0, MaxMatrixLength)
	originPairs := make([]int, 0, MaxMatrixLength*MaxMatrixLength)
	destPairs := make([]int, 0, MaxMatrixLength*MaxMatrixLength)

	// For every item on one side of the diagonal of the matrix, calculate the multimodal distance
	var wg sync.WaitGroup
	for i := 0; i < mat.Cols; i++ {
		for j := 0; j < i; j++ {
			currentLeg := mat.Matrix[LookupIndex(i, j, mat.Cols)]
			if currentLeg.Distance < MaxRoadDist {
				// Road route
				if !slices.Contains(origins, i) {
					origins = append(origins, i)
				}
				originPairs = append(originPairs, i)
				if !slices.Contains(dests, j) {
					dests = append(dests, j)
				}
				destPairs = append(destPairs, j)
			}
			// Do nothing for air route, as it is already in the correct form
			if len(origins) >= MaxMatrixLength || len(dests) >= MaxMatrixLength {
				// Execute the request (reached max capacity)
				reqOrigins, reqDests := slices.Clone(origins), slices.Clone(dests)
				wg.Go(func() {
					updateDists(mat, locations, reqOrigins, reqDests, originPairs, destPairs)
				})
				origins = make([]int, 0, MaxMatrixLength)
				dests = make([]int, 0, MaxMatrixLength)
			}
		}
	}
	if len(origins) > 0 {
		// If there are still locations left, request them now
		log.Printf("requesting origins (length > 0)")
		wg.Go(func() {
			updateDists(mat, locations, origins, dests, originPairs, destPairs)
		})
	}
	wg.Wait()
	return mat, nil
}

func AddPolylines(order []*models.Leg) {
	// Add the polylines to a given tour
	var wg sync.WaitGroup
	for _, leg := range order {
		if leg.IsRoad {
			// If the current leg is a road leg, asynchronously request a polyline for it
			wg.Go(func() {
				changePolyline(leg)
			})
		}
	}
	wg.Wait()
}

func changePolyline(leg *models.Leg) {
	// Update a polyline to be one requested by an ORS request

	polylineStr, dist, err := client.PolylineRequest(leg)
	if err != nil {
		log.Printf("error requesting polyline: %v", err)
		return
	}

	// Calculate the percentage error and log if it is greater than the tolerance
	pErr := PErr(dist, leg.Distance)
	if pErr > LengthTolerance {
		log.Printf("error between polyline length and matrix length greater than tolerance: %.2f error", pErr)
	}
	decimated, err := DecimateLine(polylineStr, RDPTolerance)
	if err != nil {
		log.Printf("error decimating polyline: %v", err)
	}
	leg.Polyline = decimated
}

func updateDists(mat *models.Matrix, locations []models.Coords, origins []int, dests []int, originPairs []int, destPairs []int) {
	// Update the distances to be road distances in a matrix by using an ORS request

	originsLength := len(originPairs)
	destsLength := len(destPairs)

	if originsLength != destsLength {
		log.Printf("origins array is not same length as dests array, lengths %d and %d respectively", originsLength, destsLength)
		return
	}

	// Make arrays for the origins and destinations as coordinates
	requestOrigins := make([]models.Coords, len(origins))
	for i, origin := range origins {
		requestOrigins[i] = locations[origin]
	}
	requestDests := make([]models.Coords, len(dests))
	for i, dest := range dests {
		requestOrigins[i] = locations[dest]
	}

	// Calculate the road distances
	newDists, err := client.MatrixRequest(requestOrigins, requestDests)
	if err != nil {
		log.Printf("error calculating road distances: %v", err)
		return
	}
	for i := range originPairs {
		// For every pair of locations, correct the current great-circle (i.e. air) distance to the road distance if a
		// road distance was calculated
		leg := newDists.Matrix[LookupIndex(slices.Index(origins, originPairs[i]), slices.Index(dests, destPairs[i]), mat.Cols)]
		if leg.Distance != 0 {
			matLeg := mat.Matrix[LookupIndex(originPairs[i], destPairs[i], mat.Cols)]
			matLegReversed := mat.Matrix[LookupIndex(destPairs[i], originPairs[i], mat.Cols)]
			matLeg.Distance = leg.Distance
			matLegReversed.Distance = leg.Distance
			matLeg.IsRoad = true
			matLegReversed.IsRoad = true
		}
	}
}
