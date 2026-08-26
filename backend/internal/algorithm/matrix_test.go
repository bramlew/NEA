package algorithm

import (
	"testing"

	"github.com/bramlew/NEA/backend/internal/models"
)

func TestMatrixValidity(t *testing.T) {
	// Test that the constructed matrix has correct values
	const N int = 10 // No. of random locations to generate
	locations := make([]models.Coords, N)
	for i := 0; i < N; i++ {
		locations[i] = RandValidCoords()
	}

	m, err := ConstructMatrix(locations)
	if err != nil {
		t.Errorf("error constructing matrix: %v", err)
	}
	for i, origin := range locations {
		for j, dest := range locations {
			index := LookupIndex(i, j, m.Cols)
			d, err := GreatCircleDistance(origin, dest)
			if err != nil {
				t.Errorf("error calculating great circle distance: %v", err)
			}
			if d != m.Matrix[index].Distance {
				t.Errorf("expected value %v at index [%v, %v], got %v", d, i, j, m.Matrix[index].Distance)
			}
		}
	}
}
