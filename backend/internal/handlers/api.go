package handlers

import (
	"net/http"

	"github.com/bramlew/NEA/backend/internal/algorithm"
	"github.com/bramlew/NEA/backend/internal/models"
	"github.com/gin-gonic/gin"
)

func Start() {
	// Start the Gin server and route(s)
	r := gin.Default()

	r.POST("/optimise", optimise)

	err := r.Run(":8080")
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
	m := algorithm.ConstructMatrix(coords)
	algorithm.ParseJsonResponse(m)
	c.JSON(http.StatusOK, gin.H{"matrix": m})
}
