package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/balitravelhealth/platform/gateway-go/internal/services"
)

// GET /location/classify?lat=-8.5&lng=115.2  — F-10, F-11
func (h *Handler) ClassifyLocation(c *gin.Context) {
	lat, err := strconv.ParseFloat(c.Query("lat"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "lat is required and must be a number"})
		return
	}
	lng, err := strconv.ParseFloat(c.Query("lng"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "lng is required and must be a number"})
		return
	}

	result := services.ClassifyLocation(lat, lng)
	c.JSON(http.StatusOK, result)
}
