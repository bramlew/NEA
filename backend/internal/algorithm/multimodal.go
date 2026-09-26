package algorithm

import (
	"log"
	"sync"

	"github.com/bramlew/NEA/backend/internal/client"
	"github.com/bramlew/NEA/backend/internal/models"
)

func ConstructMultimodalMatrix(locations []models.Coords) (*models.Matrix, error) {
	// Calculate the multimodal distance between two sets of coordinates
	const MaxRoadDist = 2000.0 // Maximum distance before defaulting to using air travel

	mat, err := ConstructMatrix(locations)
	if err != nil {
		return nil, err
	}
	origins := make([]int, 0, 59)
	dests := make([]int, 0, 59)

	var wg sync.WaitGroup
	for i := 0; i < mat.Cols; i++ {
		for j := 0; j < i+1; j++ {
			currentLeg := mat.Matrix[LookupIndex(i, j, mat.Cols)]
			if currentLeg.Distance < MaxRoadDist {
				// Road route
				origins = append(origins, i)
				dests = append(dests, j)
			}
			// Do nothing for air route, as it is already in the correct form
			if len(origins) >= 59 {
				// Execute the request (reached max capacity)
				wg.Go(func() {
					updateDists(mat, locations, origins, dests)
				})
				origins = make([]int, 0, 59)
				dests = make([]int, 0, 59)
			}
		}
	}
	if len(origins) > 0 {
		// If there are still locations left, request them now
		wg.Go(func() {
			updateDists(mat, locations, origins, dests)
		})
	}
	wg.Wait()
	return mat, nil
}

func GetPolylines(order []int, locations []models.Coords) {

}

func updateDists(mat *models.Matrix, locations []models.Coords, origins []int, dests []int) {
	// Update the distances to be road distances in a matrix by using an ORS request
	originsLength := len(origins)
	destsLength := len(dests)
	if originsLength != destsLength {
		log.Printf("origins array is not same length as dests array, lengths %d and %d respectively", originsLength, destsLength)
		return
	}

	// Make arrays for the origins and destinations as coordinates
	requestOrigins := make([]models.Coords, originsLength)
	requestDests := make([]models.Coords, destsLength)
	for i := range origins {
		requestOrigins[i] = locations[origins[i]]
		requestDests[i] = locations[dests[i]]
	}

	// Construct a great-circle distance adjacency matrix
	newDists, err := client.MatrixRequest(requestOrigins, requestDests)
	if err != nil {
		log.Printf("error calculating road distances: %v", err)
		return
	}
	for i := range origins {
		// For every pair of locations, correct the current great-circle (i.e. air) distance to the road distance if a
		// road distance was calculated
		newDist := newDists.Matrix[LookupIndex(i, i, newDists.Cols)].Distance
		if newDist != 0 {
			leg := mat.Matrix[LookupIndex(origins[i], dests[i], mat.Cols)]
			legReversed := mat.Matrix[LookupIndex(dests[i], origins[i], mat.Cols)]
			leg.Distance = newDist
			leg.IsRoad = true
			legReversed.Distance = newDist
			legReversed.IsRoad = true
		}
	}
}
