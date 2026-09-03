package algorithm

import (
	"testing"

	"github.com/bramlew/NEA/backend/internal/models"
)

func TestNearestNeighbourInit(t *testing.T) {
	// Test that the NN algorithm produces a reasonable result without errors
	const N = 10 // No. of locations to randomly choose
	locations := make([]models.Coords, N)
	for i := range locations {
		locations[i] = RandValidCoords()
	}
	mat, err := ConstructMatrix(locations)
	if err != nil {
		t.Errorf("error constructing matrix: %v", err)
	}
	for i := range locations {
		result := NearestNeighbour(mat, i)
		if result == nil {
			t.Error("error running nearest neighbour algorithm")
		} else {
			t.Logf("nearest neighbour algorithm produced order: %v", result)
		}
	}

}
