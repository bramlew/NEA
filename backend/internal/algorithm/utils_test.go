package algorithm

import (
	"math/rand/v2"

	"github.com/bramlew/NEA/backend/internal/models"
)

func RandValidCoords() models.Coords {
	// Generate a set of random valid coordinates
	return models.Coords{
		Lon: (rand.Float64() - 0.5) * 360,
		Lat: (rand.Float64() - 0.5) * 180,
	}
}

func RandInvalidCoords() models.Coords {
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

func RandLocations(n int) []models.Coords {
	locations := make([]models.Coords, n)
	for i := range locations {
		locations[i] = RandValidCoords()
	}
	return locations
}

func RandEuropeanLocations(n int) []models.Coords {
	// Generate random locatiosn found in continental Europe, i.e. from longitude 5-23, latitude 46-52
	locations := make([]models.Coords, n)
	for i := range locations {
		locations[i] = models.Coords{
			Lon: rand.Float64()*18 + 5,
			Lat: rand.Float64()*6 + 46,
		}
	}
	return locations
}
