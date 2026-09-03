package algorithm

import (
	"math"
	"slices"

	"github.com/bramlew/NEA/backend/internal/models"
)

func NearestNeighbour(mat *models.Matrix, startNode int) []int {
	// Nearest neighbour greedy algorithm for route optimisation
	route := make([]int, 0, mat.Cols)
	route = append(route, startNode)
	currentNode := startNode
	for i := 0; i < mat.Cols-1; i++ {
		// For every index in the list of locations apart from the start node, attempt to find the closest node
		bestRoute, bestIndex := math.Inf(1), -1
		for j := 0; j < mat.Cols; j++ {
			// Check the distance between the current node and every single other node, hence finding the nearest node
			index := LookupIndex(currentNode, j, mat.Cols)
			dist := mat.Matrix[index].Distance
			if dist < bestRoute && !slices.Contains(route, j) {
				bestRoute = dist
				bestIndex = j
			}
		}
		// Add the nearest neighbour to the route and set the current node to this node
		route = append(route, bestIndex)
		currentNode = bestIndex
	}
	return route
}

func threeOpt(mat *models.Matrix, baseline []int) []int {
	// Perform 3-opt optimisation logic
	route := baseline
	for {
		// Until an explicit break of the loop (i.e. it cannot be optimised any more), keep trying to optimise it
		optimal, optimalDelta := route, 0.0
		for i := 0; i < len(baseline)-3; i++ {
			for j := i + 1; j < len(baseline)-2; j++ {
				for k := j + 1; k < len(baseline)-1; k++ {
					// The for loops above initialise iterators that can be used to generate every possible triplet of edges

					// The below constants are arrays of the indexes of where the end of split segments of the route are
					segments := [4][2]int{
						{0, i},
						{i + 1, j},
						{j + 1, k},
						{k + 1, len(baseline) - 1},
					}
					// Note that the below 2 constants only contain the indexes for the inner two segments, as the outer two are constant and can be found in the previously declared constant
					reversedSegments := [2][2]int{
						{j, i + 1},
						{k, j + 1},
					}
					orders := [7][2][2]int{
						{reversedSegments[0], segments[2]},
						{segments[1], reversedSegments[1]},
						reversedSegments,
						{segments[2], segments[1]},
						{reversedSegments[1], segments[1]},
						{segments[2], reversedSegments[0]},
						{reversedSegments[1], reversedSegments[0]},
					}

					// Determine the old/current weight of the 3 edges that are going to be reordered
					oldWeight := mat.Matrix[LookupIndex(segments[0][1], segments[1][0], mat.Cols)].Distance + mat.Matrix[LookupIndex(segments[1][1], segments[2][0], mat.Cols)].Distance + mat.Matrix[LookupIndex(segments[2][1], segments[3][0], mat.Cols)].Distance
					for _, order := range orders {
						// For every possible way the edges can be reordered, determine the weight of those three edges and find the delta
						newWeight := mat.Matrix[LookupIndex(segments[0][1], order[0][0], mat.Cols)].Distance + mat.Matrix[LookupIndex(order[0][1], order[1][0], mat.Cols)].Distance + mat.Matrix[LookupIndex(order[1][1], segments[3][0], mat.Cols)].Distance
						delta := newWeight - oldWeight
						if delta < optimalDelta {
							// If the delta is better than the current best delta in the whole loop, update it
							optimal = slices.Concat(route[:i+1], route[order[0][0]:order[0][1]], route[order[1][0]:order[1][1]], route[k+1:])
							optimalDelta = delta
						}
					}
				}
			}
		}
		if optimalDelta == 0 {
			// If there was no improvement in the delta, break the loop
			break
		}
		route = optimal
	}
	return route
}

func calcWeight(mat *models.Matrix, route []int) float64 {
	weight := 0.0
	for i := range route {
		weight += mat.Matrix[LookupIndex(route[i], route[i+1], mat.Cols)].Distance
	}
	return weight
}

/*
	Segments are:
	A: 0   ->  i
	B: i+1 ->  j
	C: j+1 ->  k
	D: k+1 -> -1

	By default, the route is in order A -> B -> C -> D
	It can be reordered in 7 other ways (i.e. 8 total):
	i:   A -> B' -> C  -> D
	ii:  A -> B  -> C' -> D
	iii: A -> B' -> C' -> D
	iv:  A -> C  -> B  -> D
	v:   A -> C' -> B  -> D
	vi:  A -> C  -> B' -> D
	vii: A -> C' -> B' -> D

	Note that ' means that a segment has been reversed

	Take, for example, order v. Decoded, this looks like:
	0 -> i --> k -> j+1 --> i+1 -> j --> k+1 -> -1

	In order to implement this in code in the context of 3-opt, we must figure out which of the 7
	rearranged orders has the lowest weight. Whilst we could construct a new route for each
	rearrangement of edges and calculate its overall weight linearly, this would be inefficient as
	the only weight that is actually changing is the weight of the 3 edges. By doing a simple
	comparison between the old weight of the 3 edges and the new weight, we can determine which
	route's overall weight will be the smallest. We can also, for the outer loops, instead of
	comparing the best of the 7 route's overall weight, we can just use the difference between the
	original 3 edge weight and new 3 edge weight to find the greatest overall reduction in weight,
	since the weight of the edges that aren't chosen is constant, so we just want the largest delta
	between the 3 modified edges.

*/
