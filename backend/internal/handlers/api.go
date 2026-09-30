package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bramlew/NEA/backend/internal/algorithm"
	"github.com/bramlew/NEA/backend/internal/models"
	"github.com/gin-gonic/gin"
)

const MaxNoLocations = 30 // Maximum no. of locations which can be provided to the API

func Start() {
	// Start the Gin server and route(s)
	router := gin.Default()

	router.POST("/optimise", optimise)

	err := router.Run(":8080")
	if err != nil {
		return
	}
}

func optimise(c *gin.Context) {
	// Optimisation API route
	var req models.Request
	if err := c.ShouldBindJSON(&req); err != nil {
		// Return an error if the JSON syntax is invalid
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// If the request does have valid syntax, construct a matrix of distances between the coordinates sent
	locations := req.Locations
	length := len(locations)
	if length > MaxNoLocations {
		// Return an error if too many locations are given
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": fmt.Sprintf("too many locations provided: %d > %d (max)", length, MaxNoLocations),
		})
		return
	}
	mat, err := algorithm.ConstructMultimodalMatrix(locations)
	if err != nil {
		// Return an error if the matrix construction fails
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	log.Print("finished constructing mm matrix, starting tsp")

	// Verify matrix is asymmetrical
	for i := 0; i < mat.Cols; i++ {
		for j := 0; j < mat.Cols; j++ {
			leg := mat.Matrix[algorithm.LookupIndex(i, j, mat.Cols)]
			legReversed := mat.Matrix[algorithm.LookupIndex(j, i, mat.Cols)]
			if leg.Distance != legReversed.Distance {
				log.Printf("leg distance not equal to reversed distance: %2.f vs %.2f", leg.Distance, legReversed.Distance)
			}
			if leg.Origin.Lon != legReversed.Dest.Lon || leg.Origin.Lat != legReversed.Dest.Lat {
				log.Printf("leg origins are incorrect: %+v vs %+v reversed", leg.Origin, legReversed.Dest)
			}
			if leg.Dest.Lon != legReversed.Origin.Lon || leg.Dest.Lat != legReversed.Origin.Lat {
				log.Printf("leg dests are incorrect: %+v vs %+v reversed", leg.Dest, legReversed.Origin)
			}
			if leg.Distance <= 0.01 && i != j {
				log.Printf("zero leg found outside diagonal: index [%d, %d], distance %.2f", i, j, leg.Distance)
			}
		}
	}

	// algorithm.ParseJsonResponse(mat)
	tour, err := algorithm.MultiTSP(mat)
	if err != nil {
		// Return an error if the TSP algorithm fails
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Print("finished tsp, starting polyline requests")
	// Make a list of individual leg details
	legs := make([]*models.Leg, length-1)
	for i := 0; i < length-1; i++ {
		legs[i] = mat.Matrix[algorithm.LookupIndex(tour.Order[i], tour.Order[i+1], mat.Cols)]
	}
	algorithm.AddPolylines(legs) // Add the polylines to the legs array
	response := models.Response{
		Route:       legs,
		TotalWeight: tour.TotalWeight,
	}
	c.JSON(http.StatusOK, gin.H{"route": response})
}
