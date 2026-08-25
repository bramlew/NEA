package algorithm

import (
	"math"

	"github.com/bramlew/NEA/backend/internal/models"
)

const earthRadKm = 6371.0084

func GreatCircleDistance(origin models.Coords, dest models.Coords) float64 {
	// Using the spherical law of cosines, compute the great circle distance between two sets of provided coordinates

	// Convert longitude and latitude coordinates (given in degrees) into radians to be used in trig functions
	lambda1, lambda2, phi1, phi2 := toRad(origin.Lon), toRad(dest.Lon), toRad(origin.Lat), toRad(dest.Lat)

	// Compute the central angle between the two points, and using the arc length formula, compute the great circle distance
	c := math.Acos(math.Sin(phi1)*math.Sin(phi2) + math.Cos(phi1)*math.Cos(phi2)*math.Cos(math.Abs(lambda2-lambda1)))
	d := earthRadKm * c

	return d

}

func toRad(deg float64) float64 {
	// Convert an angle in degrees to radians

	return deg * math.Pi / 180
}
