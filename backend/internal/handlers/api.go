package handlers

import (
	"fmt"
	"net/http"
	"slices"

	"github.com/bramlew/NEA/backend/internal/algorithm"
	"github.com/bramlew/NEA/backend/internal/models"
	"github.com/gin-gonic/gin"
)

const MaxNoLocations = 30 // Maximum no. of locations which can be provided to the API
const MinNoLocations = 2  // Minimum no. of locations which can be provided to the API

func Start() {
	// Start the Gin server and optimisation route
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
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{
			"error": fmt.Sprintf("too many locations provided: %d > %d (max)", length, MaxNoLocations),
		})
		return
	} else if length < MinNoLocations {
		// Return an error if too few locations are given
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": fmt.Sprintf("too few locations provided: %d < %d (min)", length, MinNoLocations),
		})
		return
	}

	if checkIfAnyEqual(locations) {
		// Return an error if any of the coordinates are equal
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": fmt.Sprintf("two or more locations had equal coordinates")})
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
	// Add the polylines to the legs array
	algorithm.AddPolylines(legs)

	response := models.Response{
		Route:       legs,
		TotalWeight: tour.TotalWeight,
	}
	c.JSON(http.StatusOK, gin.H{"details": response})
}

func checkIfAnyEqual(coords []models.Coords) bool {
	// Check if any coords in a slice are equal
	checked := make([]models.Coords, 0, len(coords))
	for _, coord := range coords {
		if !slices.Contains(checked, coord) {
			checked = append(checked, coord)
		} else {
			return true
		}
	}
	return false
}
