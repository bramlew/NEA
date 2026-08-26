package algorithm

import (
	"math"
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

	mat, err := ConstructMatrix(locations)
	if err != nil {
		t.Errorf("error constructing matrix: %v", err)
	}
	for i, origin := range locations {
		for j, dest := range locations {
			matDist := mat.Matrix[LookupIndex(i, j, mat.Cols)].Distance
			if i != j {
				dist, err := GreatCircleDistance(origin, dest)
				if err != nil {
					t.Errorf("error calculating great circle distance: %v", err)
				}
				if dist != matDist {
					t.Errorf("expected distance %v at index [%v, %v], got %v", dist, i, j, matDist)
				}
			} else if !math.IsInf(matDist, 1) {
				t.Errorf("expected infinite distance at index [%v, %v], got %v", i, j, matDist)
			}
		}
	}
}
