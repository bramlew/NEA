package algorithm

import (
	"math"
	"slices"
	"testing"
)

func TestNearestNeighbourInit(t *testing.T) {
	// Test that the NN algorithm produces a reasonable result without errors
	const N = 10 // No. of locations to randomly choose
	locations := RandLocations(N)
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

func TestNearestNeighbourAccuracy(t *testing.T) {
	// Test that the NN algorithm produces an accurate result - i.e. it actually goes to the nearest neighbour
	const N = 10 // No. of locations to randomly choose
	locations := RandLocations(N)
	mat, err := ConstructMatrix(locations)
	if err != nil {
		t.Errorf("error constructing matrix: %v", err)
	}
	for i := range locations {
		result := NearestNeighbour(mat, i)
		if result != nil {
			visitedIndexes := make([]int, 0, N)
			visitedIndexes = append(visitedIndexes, i)
			currentNode := i
			for j := 0; j < N-1; j++ {
				nextDist, nextIndex := mat.Matrix[LookupIndex(result[j], result[j+1], N)].Distance, result[j+1]
				bestDist, bestIndex := math.Inf(1), -1
				for k := range locations {
					dist := mat.Matrix[LookupIndex(currentNode, k, N)].Distance
					if dist < bestDist && !slices.Contains(visitedIndexes, k) {
						bestDist, bestIndex = dist, k
					}
				}
				if bestIndex != nextIndex && bestDist != nextDist {
					t.Errorf("expected %v to be next node index, got %v", bestIndex, nextIndex)
				}
				visitedIndexes = append(visitedIndexes, nextIndex)
				currentNode = nextIndex
			}
		} else {
			t.Error("error running nearest neighbour algorithm")
		}
	}
}
