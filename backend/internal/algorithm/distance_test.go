package algorithm

import (
	"math"
	"math/rand/v2"
	"testing"

	"github.com/bramlew/NEA/backend/internal/models"
	"github.com/pymaxion/geographiclib-go/v2/geodesic"
)

func TestInit(t *testing.T) {
	// Test that the great circle distance calculation algorithm does not produce an error using random valid and invalid sets of coordinates
	const N int = 10 // No. of times to randomly execute the tests
	for i := 0; i < N; i++ {
		// Test some valid coordinates
		_, err := GreatCircleDistance(randValidCoords(), randValidCoords())
		if err != nil {
			t.Errorf("error calculating great circle distance: %v", err)
		}
	}
	for i := 0; i < N; i++ {
		// Test some invalid coordinates
		_, err := GreatCircleDistance(randInvalidCoords(), randInvalidCoords())
		if err == nil {
			t.Errorf("expected error for invalid coordinates, got none")
		}
	}
}

func TestAccuracy(t *testing.T) {
	// Test the accuracy of the great circle distance calculation algorithm with random valid and invalid sets of coordinates
	const N int = 10000 // No. of times to randomly execute the test
	for i := 0; i < N; i++ {
		origin, dest := randValidCoords(), randValidCoords()
		d, err := GreatCircleDistance(origin, dest)
		actD := geodesic.WGS84.Inverse(origin.Lat, origin.Lon, dest.Lat, dest.Lon).S12 / 1000
		if err != nil {
			t.Errorf("error calculating great circle distance: %v", err)
		}
		if math.Abs(d-actD)/actD >= 0.005 {
			t.Errorf("calculated distance is not within expected error range, expected maximum 0.5%% error, got %.2f%%", math.Abs(d-actD)/actD*100)
		}
	}
}

func randValidCoords() models.Coords {
	// Generate a set of random valid coordinates
	return models.Coords{
		Lon: (rand.Float64() - 0.5) * 360,
		Lat: (rand.Float64() - 0.5) * 180,
	}
}

func randInvalidCoords() models.Coords {
	// Generate a set of random invalid coordinates
	return models.Coords{
		Lon: randSign() * (rand.Float64()*180 + 180),
		Lat: randSign() * (rand.Float64()*90 + 90),
	}
}

func randSign() float64 {
	// Generates either positive 1 or negative 1, by using the remainder when dividing by two of a random unsigned 32-bit integer, hence 2 equally probable outputs
	if rand.Uint32()%2 == 0 {
		return 1
	}
	return -1
}
