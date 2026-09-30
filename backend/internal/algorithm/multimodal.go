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
	// Origins and dests are the arrays containing the indexes of the locations in the locations array
	origins := make([]int, 0, MaxMatrixLength)
	dests := make([]int, 0, MaxMatrixLength)
	// The pairs arrays contain the indexes of the location in the respective origins/dests array
	// For example, if we are trying to work out A -> B and A -> C, assuming locations = [A, B, C],
	// origins will be [0], dests will be [1, 2], originPairs will be [0, 0], and destPairs will be [0, 1].
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
					originPairs = append(originPairs, len(origins))
					origins = append(origins, i)
				} else {
					originPairs = append(originPairs, slices.Index(origins, i))
				}
				if !slices.Contains(dests, j) {
					destPairs = append(destPairs, len(dests))
					dests = append(dests, j)
				} else {
					destPairs = append(destPairs, slices.Index(dests, j))
				}
			}
			// Do nothing for air route, as it is already in the correct form
			if len(origins) >= MaxMatrixLength || len(dests) >= MaxMatrixLength {
				// Execute the request (reached max capacity)
				reqOrigins, reqDests, reqOriginPairs, reqDestPairs := slices.Clone(origins), slices.Clone(dests), slices.Clone(originPairs), slices.Clone(destPairs)
				wg.Go(func() {
					updateDists(mat, locations, reqOrigins, reqDests, reqOriginPairs, reqDestPairs)
				})
				// Reset the arrays
				origins = make([]int, 0, MaxMatrixLength)
				dests = make([]int, 0, MaxMatrixLength)
				originPairs = make([]int, 0, MaxMatrixLength*MaxMatrixLength)
				destPairs = make([]int, 0, MaxMatrixLength*MaxMatrixLength)
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

	originPairsLength := len(originPairs)
	destPairsLength := len(destPairs)
	if originPairsLength != destPairsLength {
		log.Printf("origin pairs array length not equal dests pairs array length: %d vs %d", originPairsLength, destPairsLength)
		return
	}

	// Make arrays for the origins and destinations as coordinates
	requestOrigins := make([]models.Coords, len(origins))
	for i, origin := range origins {
		requestOrigins[i] = locations[origin]
		log.Printf("origin %d: %+v", i, locations[origin])
	}
	requestDests := make([]models.Coords, len(dests))
	for i, dest := range dests {
		log.Printf("dest %d: %+v", i, locations[dest])
		requestDests[i] = locations[dest]
	}

	// Calculate the road distances
	log.Printf("originPairs:%+v\ndestPairs:%+v", originPairs, destPairs)
	dists, err := client.MatrixRequest(requestOrigins, requestDests, originPairs, destPairs)
	if err != nil {
		log.Printf("error calculating road distances: %v", err)
		return
	}
	log.Printf("%+v", dists)
	for i, dist := range dists {
		// For every pair of locations, correct the current great-circle (i.e. air) distance to the road distance if a
		// road distance was calculated
		if dist != 0 {
			log.Printf("number leg: %.2f", dist)
			matLeg := mat.Matrix[LookupIndex(origins[originPairs[i]], dests[destPairs[i]], mat.Cols)]
			matLegReversed := mat.Matrix[LookupIndex(dests[destPairs[i]], origins[originPairs[i]], mat.Cols)]
			matLeg.Distance = dist
			matLegReversed.Distance = dist
			matLeg.IsRoad = true
			matLegReversed.IsRoad = true
		}
	}
}
