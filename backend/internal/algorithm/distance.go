package algorithm

import (
	"fmt"
	"math"

	"github.com/bramlew/NEA/backend/internal/models"
)

const earthRadKm = 6371.0084

func GreatCircleDistance(origin models.Coords, dest models.Coords) (float64, error) {
	// Using the spherical law of cosines, compute the great circle distance between two sets of provided coordinates

	// Firstly, ensure that the longitude and latitude values in degrees are within valid boundaries
	if !(origin.Lon >= -180 && origin.Lon <= 180 && dest.Lon >= -180 && dest.Lon <= 180 && origin.Lat >= -90 && origin.Lat <= 90 && dest.Lat >= -90 && dest.Lat <= 90) {
		return 0.0, fmt.Errorf("coordinates not within valid range, given origin %+v and destination %+v", origin, dest)
	}

	// Convert longitude and latitude coordinates (given in degrees) into radians to be used in trig functions
	lambda1, lambda2, phi1, phi2 := toRad(origin.Lon), toRad(dest.Lon), toRad(origin.Lat), toRad(dest.Lat)

	// Compute the central angle between the two points, and using the arc length formula, compute the great circle distance
	c := math.Acos(math.Sin(phi1)*math.Sin(phi2) + math.Cos(phi1)*math.Cos(phi2)*math.Cos(math.Abs(lambda2-lambda1)))
	d := earthRadKm * c

	return d, nil

}

func toRad(deg float64) float64 {
	// Convert an angle in degrees to radians

	return deg * math.Pi / 180
}
