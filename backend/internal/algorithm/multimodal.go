package algorithm

import (
	"github.com/bramlew/NEA/backend/internal/models"
)

func CalcDistances(locations []models.Coords) (*models.Matrix, error) {
	// Calculate the multimodal distance between two sets of coordinates
	const MaxRoadDist = 3000.0 // Maximum distance before defaulting to using air travel

	directDistMat, err := ConstructMatrix(locations)
	if err != nil {
		return nil, err
	}

	if directDist < MaxRoadDist {
		// Road route

	} else {
		// Air route
	}
}
