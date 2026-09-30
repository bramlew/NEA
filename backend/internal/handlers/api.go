package handlers

import (
	"fmt"
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
	tour, err := algorithm.MultiTSP(mat)
	if err != nil {
		// Return an error if the TSP algorithm fails
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
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
