package algorithm

import (
	"io"
	"os"
	"testing"

	"github.com/bramlew/NEA/backend/internal/models"
	"github.com/goccy/go-json"
)

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
}

func TestAPI(t *testing.T) {
	// Test the API for errors whilst running multimodal algorithm
	jsonFile, err := os.Open("..\\..\\..\\json\\api\\first\\request.json")
	if err != nil {
		t.Fatalf("unable to open json file: %v", err)
	}
	defer jsonFile.Close()
	jsonBody, err := io.ReadAll(jsonFile)
	if err != nil {
		t.Fatalf("unable to read json file: %v", err)
	}
	var req models.Request
	err = json.Unmarshal(jsonBody, &req)
	if err != nil {
		t.Fatalf("unable to unmarshal json: %v", err)
	}
	locations := req.Locations
	mat, err := ConstructMultimodalMatrix(locations)
	if err != nil {
		t.Errorf("error constructing multimodal matrix: %v", err)
	}
	_, err = MultiTSP(mat)
	if err != nil {
		t.Errorf("error running tsp algorithms: %v", err)
	}
}
