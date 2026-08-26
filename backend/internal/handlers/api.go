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
	coords := req.Coordinates
	mat, err := algorithm.ConstructMatrix(coords)
	if err != nil {
		// Return an error if the matrix construction fails
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	algorithm.ParseJsonResponse(mat)
	c.JSON(http.StatusOK, gin.H{"matrix": mat})
}
