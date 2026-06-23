package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/balitravelhealth/platform/gateway-go/internal/middleware"
	"github.com/balitravelhealth/platform/gateway-go/internal/services"
)

type vaccinationRequest struct {
	JenisVaksin string  `json:"jenis_vaksin"`
	Tanggal     string  `json:"tanggal"`
	Dosis       *string `json:"dosis"`
	Catatan     *string `json:"catatan"`
}

func (h *Handler) ListVaccinations(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(int64)
	list, err := h.VaccinationService.List(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list})
}

func (h *Handler) CreateVaccination(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(int64)
	var req vaccinationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	v, err := h.VaccinationService.Create(c.Request.Context(), userID, services.VaccinationInput{
		JenisVaksin: req.JenisVaksin,
		Tanggal:     req.Tanggal,
		Dosis:       req.Dosis,
		Catatan:     req.Catatan,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": v})
}

func (h *Handler) DeleteVaccination(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(int64)
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.VaccinationService.Delete(c.Request.Context(), id, userID); err != nil {
		if errors.Is(err, services.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "record not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
