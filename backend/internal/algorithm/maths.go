package algorithm

import "math"

func toRad(deg float64) float64 {
	// Convert an angle in degrees to radians
	return deg * math.Pi / 180
}

func SLCAngle(a float64, b float64, c float64) float64 {
	// General formula for the spherical law of cosines rearranged for an angle, C
	return math.Acos((math.Cos(c) - math.Cos(a)*math.Cos(b)) / (math.Sin(a) * math.Sin(b)))
}

func SLCSide(a float64, b float64, C float64) float64 {
	// General formula for the spherical law of cosines rearranged for a side, c
	return math.Acos(math.Cos(a)*math.Cos(b) + math.Sin(a)*math.Sin(b)*math.Cos(C))
}

func SLSSide(A float64, B float64, b float64) float64 {
	// General formula for the spherical law of sines rearranged for a side, a
	return math.Asin(math.Sin(A) * math.Sin(b) / math.Sin(B))
}
