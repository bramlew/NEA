package algorithm

import (
	"testing"

	"github.com/bramlew/NEA/backend/internal/client"
	"github.com/bramlew/NEA/backend/internal/models"
)

func TestDecimation(t *testing.T) {
	// Request polylines between random European coordinates and test the decimation algorithm on them
	const N = 10 // No. of times to execute the test
	const Tolerance = RDPTolerance
	for range N {
		oldPolyline, _, err := client.PolylineRequest(&models.Leg{
			Origin: RandEuropeanLocation(),
			Dest:   RandEuropeanLocation(),
		})
		if err != nil {
			t.Errorf("error obtaining polyline: %v", err)
		}
		oldLength := len(oldPolyline)
		newPolyline, err := DecimateLine(oldPolyline, Tolerance)
		if err != nil {
			t.Errorf("error decimating polyline: %v", err)
		}
		newLength := len(newPolyline)
		if newLength >= oldLength {
			t.Errorf("failed to produce shorter polyline with tolerance %.2f: old length %d vs new length %d", Tolerance, oldLength, newLength)
		} else {
			t.Logf("successfully shortened polyline length with tolerance %.2f: old length %d vs new length %d", Tolerance, oldLength, newLength)
		}
	}
}
