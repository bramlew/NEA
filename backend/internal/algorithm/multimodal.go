package algorithm

import (
	"log"
	"sync"

	"github.com/bramlew/NEA/backend/internal/client"
	"github.com/bramlew/NEA/backend/internal/models"
)

func ConstructMultimodalMatrix(locations []models.Coords) (*models.Matrix, error) {
	// Calculate the multimodal distance between two sets of coordinates
	const MaxRoadDist = 3000.0 // Maximum distance before defaulting to using air travel

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
			currentLegReversed := mat.Matrix[LookupIndex(j, i, mat.Cols)]
			if currentLeg.Distance < MaxRoadDist {
				// Road route
				origins = append(origins, i)
				dests = append(dests, j)
				currentLeg.IsRoad = true
				currentLegReversed.IsRoad = true
			}
			// Do nothing for air route, as it is already in the correct form
			if len(origins) >= 59 {
				// Execute the request (reached max capacity)
				wg.Go(func() {
					updateDists(mat, locations, origins, dests)
				})
				go updateDists(mat, locations, origins, dests)
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

func updateDists(mat *models.Matrix, locations []models.Coords, origins []int, dests []int) {
	// Update the distances to be road distances in a matrix by using an ORS request
	originsLength := len(origins)
	destsLength := len(dests)
	if originsLength != destsLength {
		log.Printf("origins array is not same length as dests array, lengths %d and %d respectively", originsLength, destsLength)
		return
	}
	requestOrigins := make([]models.Coords, len(origins))
	requestDests := make([]models.Coords, len(origins))
	for i := range origins {
		requestOrigins[i] = locations[origins[i]]
		requestDests[i] = locations[dests[i]]
	}
	newDists, err := client.MatrixRequest(requestOrigins, requestDests)
	if err != nil {
		log.Printf("error calculating road distances: %v", err)
		return
	}
	for i := range origins {
		newDist := newDists.Matrix[LookupIndex(i, i, newDists.Cols)].Distance
		mat.Matrix[LookupIndex(origins[i], dests[i], mat.Cols)].Distance = newDist
		mat.Matrix[LookupIndex(dests[i], origins[i], mat.Cols)].Distance = newDist
	}
}
