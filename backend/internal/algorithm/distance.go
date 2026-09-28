package algorithm

import (
	"fmt"
	"math"

	"github.com/bramlew/NEA/backend/internal/models"
)

const EarthRadKm = 6371.0084

func GreatCircleDistance(origin models.Coords, dest models.Coords) (float64, error) {
	// Calculate the great circle distance between two sets of coordinates in km
	c, err := AngularGCD(origin, dest)
	if err != nil {
		return 0.0, err
	}
	// If the great circle distance is requested as a distance, return the central angle
	d := EarthRadKm * c
	return d, nil
}

func AngularGCD(origin models.Coords, dest models.Coords) (float64, error) {
	// Using the spherical law of cosines, compute the central angle between two sets of provided coordinates

	// Firstly, ensure that the longitude and latitude values in degrees are within valid boundaries
	if !(origin.Lon >= -180 && origin.Lon <= 180 && dest.Lon >= -180 && dest.Lon <= 180 && origin.Lat >= -90 && origin.Lat <= 90 && dest.Lat >= -90 && dest.Lat <= 90) {
		return 0.0, fmt.Errorf("coordinates not within valid range, given origin %+v and destination %+v", origin, dest)
	}

	// Convert longitude and latitude coordinates (given in degrees) into radians to be used in trig functions
	lambda1, lambda2, phi1, phi2 := toRad(origin.Lon), toRad(dest.Lon), toRad(origin.Lat), toRad(dest.Lat)

	// Compute the central angle between the two points
	c := SLCSide(math.Pi/2-phi1, math.Pi/2-phi2, math.Abs(lambda2-lambda1))
	return c, nil
}
