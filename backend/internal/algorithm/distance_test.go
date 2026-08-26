package algorithm

import (
	"math"
	"testing"

	"github.com/pymaxion/geographiclib-go/v2/geodesic"
)

func TestInit(t *testing.T) {
	// Test that the great circle distance calculation algorithm does not produce an error using random valid and invalid sets of coordinates
	const N = 10 // No. of times to randomly execute the tests
	for i := 0; i < N; i++ {
		// Test some valid coordinates
		_, err := GreatCircleDistance(RandValidCoords(), RandValidCoords())
		if err != nil {
			t.Errorf("error calculating great circle distance: %v", err)
		}
	}
	for i := 0; i < N; i++ {
		// Test some invalid coordinates
		_, err := GreatCircleDistance(RandInvalidCoords(), RandInvalidCoords())
		if err == nil {
			t.Error("expected error for invalid coordinates, got none")
		}
	}
}

func TestAccuracy(t *testing.T) {
	// Test the accuracy of the great circle distance calculation algorithm with random sets of valid coordinates
	const N = 1000000      // No. of times to randomly execute the test
	const MaxDiff = 0.0057 // Maximum percentage difference between calculated and actual values of the distance as a decimal
	for i := 0; i < N; i++ {
		origin, dest := RandValidCoords(), RandValidCoords()
		dist, err := GreatCircleDistance(origin, dest)
		actualDist := geodesic.WGS84.Inverse(origin.Lat, origin.Lon, dest.Lat, dest.Lon).S12 / 1000
		if err != nil {
			t.Errorf("error calculating great circle distance: %v", err)
		}
		if diff := math.Abs(dist-actualDist) / actualDist; diff >= MaxDiff {
			t.Errorf("calculated distance is not within expected error range, expected maximum 0.56%% error, got %.2f%% error (calculated %.2fkm, expected %.2fkm)", diff*100, dist, actualDist)
		}
	}
}
