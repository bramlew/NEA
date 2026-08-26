package algorithm

import (
	"fmt"
	"math"

	"github.com/bramlew/NEA/backend/internal/models"
)

func ConstructMatrix(coords []models.Coords) (*models.Matrix, error) {
	// Construct and return a distance adjacency matrix for a given slice of coordinates

	// Define initial constants/values
	inf := math.Inf(1)
	length := len(coords)
	// Store the matrix as a flat array instead of a 2D array, as this speeds up indexing time, you just need to store the
	// no. of columns along with the actual matrix.
	m := make([]*models.Route, length*length)
	for i, origin := range coords {
		for j, dest := range coords {
			// For every possible origin and destination coordinate pair, calculate the great-circle distance between them.
			// If the origin and destination are the same, set the distance to infinity.
			index := LookupIndex(i, j, length)
			if i != j {
				distance, err := GreatCircleDistance(origin, dest)
				if err != nil {
					return nil, err
				}
				m[index] = &models.Route{Distance: distance}
			} else {
				m[index] = &models.Route{Distance: inf}
			}
		}
	}
	return &models.Matrix{
		Matrix: m,
		Cols:   length,
	}, nil
}

func LookupIndex(i int, j int, cols int) int {
	// Returns the flat index from a given 2D index
	return i*cols + j
}

func PrintMatrix(m *models.Matrix) {
	// Output every item in the matrix to the console
	var printStr string
	for i, route := range m.Matrix {
		if i%m.Cols == 0 {
			printStr += "\n"
		}
		printStr += fmt.Sprintf("%g", route.Distance) + "      "
	}
	fmt.Println(printStr)
}

func ParseJsonResponse(m *models.Matrix) {
	// Replace all occurrences of infinity length as zero, as infinity cannot be passed back into JSON
	for _, r := range m.Matrix {
		if math.IsInf(r.Distance, 1) {
			r.Distance = 0
		}
	}
}
