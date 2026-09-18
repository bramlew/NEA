package algorithm

import "testing"

func TestMultimodal(t *testing.T) {
	// Test the multimodal matrix construction algorithm
	const N = 10 // Number of random locations to choose
	locations := RandEuropeanLocations(N)
	multimodal, err := ConstructMultimodalMatrix(locations)
	if err != nil {
		t.Errorf("error constructing multimodal matrix: %v", err)
	}
	standard, err := ConstructMatrix(locations)
	if err != nil {
		t.Errorf("error constructing standard matrix: %v", err)
	}
	for i, multimodalLeg := range multimodal.Matrix {
		if multimodalLeg.IsRoad && multimodalLeg.Distance == standard.Matrix[i].Distance {
			t.Errorf("multimodal road distance same as standard distance: %.2f", multimodalLeg.Distance)
		}
	}
	t.Logf("multimodal matrix: %v", *multimodal)
}
