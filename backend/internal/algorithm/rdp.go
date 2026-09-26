package algorithm

import (
	"fmt"
	"math"
	"slices"

	"github.com/bramlew/NEA/backend/internal/models"
	"github.com/twpayne/go-polyline"
)

func DecimateLine(polylineStr string, tolerance float64) (string, error) {
	// Decimate a polyline using the Ramer-Douglas-Peucker algorithm
	coords, err := DecodePolyline(polylineStr)
	if err != nil {
		return "", err
	}
	decimated, err := RDP(coords, tolerance)
	if err != nil {
		return "", err
	}
	decimatedPolyline := EncodePolyline(decimated)
	return decimatedPolyline, nil
}

func RDP(locations []models.Coords, epsilon float64) ([]models.Coords, error) {
	// Run the Ramer-Douglas-Peucker algorithm on a set of locations
	dMax, index, err := calcDMax(locations)
	if err != nil {
		return nil, err
	}
	if dMax <= epsilon {
		// If the max distance is within the specified tolerance, remove all inner points on the line
		return []models.Coords{locations[0], locations[len(locations)-1]}, nil
	}
	// If the max distance is not within the specified tolerance, recursively call the RDP algorithm for points either
	// side of the max distance point's index
	before, err := RDP(locations[:index], epsilon)
	if err != nil {
		return nil, err
	}
	after, err := RDP(locations[index:], epsilon)
	if err != nil {
		return nil, err
	}
	full := slices.Concat(before, after)
	return full, nil
}

func calcDMax(locations []models.Coords) (float64, int, error) {
	// Using spherical trigonometry, calculate the point at which there is a maximum distance between itself and the
	// line from the start to the end of the provided locations array
	dMax, index := 0.0, 0
	start, end := 0, len(locations)-1
	for i := start + 1; i < end; i++ {
		// For every location excluding the start and end one, calculate the distance from itself to the line
		// In these calculations, A is the start location, B is the end location, and C is the point from which we are
		// trying to calculate the distance from
		// ALoc, BLoc, and CLoc are the coordinates of the three distinctive vertices of the spherical triangle
		// a, b, and c are the three sides opposite their respective angle on the spherical triangle
		// A and B are the angles at their respective location on the spherical triangle (C not needed for this case)
		// d is the side opposite A in a right spherical triangle between A, C, and the point at which the perpendicular
		// between C and the line AB intersects AB
		ALoc := locations[start]
		BLoc := locations[end]
		CLoc := locations[i]
		a, err := GreatCircleDistance(BLoc, CLoc)
		if err != nil {
			return 0.0, 0, fmt.Errorf("error calculating great circle distance between %v and %v: %v", BLoc, CLoc, err)
		}
		b, err := GreatCircleDistance(ALoc, CLoc)
		if err != nil {
			return 0.0, 0, fmt.Errorf("error calculating great circle distance between %v and %v: %v", ALoc, CLoc, err)
		}
		c, err := GreatCircleDistance(ALoc, BLoc)
		if err != nil {
			return 0.0, 0, fmt.Errorf("error calculating great circle distance between %v and %v: %v", ALoc, BLoc, err)
		}
		d := 0.0
		if A := SLCAngle(b, c, a); A >= math.Pi/2 {
			d = b
		} else if B := SLCAngle(a, c, b); B >= math.Pi/2 {
			d = a
		} else {
			d = SLSSide(A, math.Pi/2, b)
		}
		if d > dMax {
			dMax, index = d, i
		}
	}
	return dMax, index, nil
}

func EncodePolyline(coords []models.Coords) string {
	// Encode a coordinates array into a polyline string
	coords2D := make([][]float64, len(coords))
	for i, coord := range coords {
		// Convert the coordinates array from 1D coords struct to 2D float
		coords2D[i] = []float64{coord.Lat, coord.Lon}
	}
	return string(polyline.EncodeCoords(coords2D))
}

func DecodePolyline(polylineStr string) ([]models.Coords, error) {
	// Decode a polyline string into a coordinates array
	coords2D, _, err := polyline.DecodeCoords([]byte(polylineStr))
	if err != nil {
		return nil, fmt.Errorf("error decoding polyline string: %v", err)
	}
	coords := make([]models.Coords, len(coords2D))
	for i := range coords2D {
		// Convert the coordinates array from 2D float to 1D coords struct
		coords[i].Lat, coords[i].Lon = coords2D[i][0], coords2D[i][1]
	}
	return coords, nil
}
