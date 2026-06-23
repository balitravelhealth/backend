package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/balitravelhealth/platform/gateway-go/internal/services"
)

func (h *Handler) NearbyFacilities(c *gin.Context) {
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

	radius, _ := strconv.ParseFloat(c.DefaultQuery("radius_km", "20"), 64)
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	results, err := h.FacilityService.Nearby(c.Request.Context(), lat, lng, radius, limit)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCoordinates) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": results})
}
