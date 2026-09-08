package handlers

import (
	"net/http"

	"github.com/bramlew/NEA/backend/internal/algorithm"
	"github.com/bramlew/NEA/backend/internal/models"
	"github.com/gin-gonic/gin"
)

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
	mat, err := algorithm.ConstructMatrix(locations)
	if err != nil {
		// Return an error if the matrix construction fails
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	// algorithm.ParseJsonResponse(mat)
	tour, err := algorithm.MultiTSP(mat)
	if err != nil {
		// Return an error if the TSP algorithm fails
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	// Make a locations list in optimised order
	length := len(locations)
	order := make([]models.Coords, length)
	for i, locationIndex := range tour.Order {
		order[i] = locations[locationIndex]
	}
	// Make a weight list of individual legs
	weights := make([]float64, length)
	for i := 0; i < length-1; i++ {
		weights[i] = mat.Matrix[algorithm.LookupIndex(tour.Order[i], tour.Order[i+1], mat.Cols)].Distance
	}
	response := models.Response{
		Locations:   order,
		Weights:     weights,
		TotalWeight: tour.TotalWeight,
	}
	c.JSON(http.StatusOK, gin.H{"route": response})
}
