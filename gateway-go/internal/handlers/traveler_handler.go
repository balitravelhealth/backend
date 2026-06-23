package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/balitravelhealth/platform/gateway-go/internal/middleware"
	"github.com/balitravelhealth/platform/gateway-go/internal/services"
)

type travelerRequest struct {
	NamaLengkap   string  `json:"nama_lengkap"`
	TanggalLahir  *string `json:"tanggal_lahir"`
	KontakDarurat *string `json:"kontak_darurat"`
}

func (h *Handler) GetTravelerProfile(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(int64)
	t, err := h.TravelerService.Get(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "traveler profile not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, t)
}

func (h *Handler) CreateTravelerProfile(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(int64)
	var req travelerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	t, err := h.TravelerService.Create(c.Request.Context(), userID, services.TravelerInput{
		NamaLengkap:   req.NamaLengkap,
		TanggalLahir:  req.TanggalLahir,
		KontakDarurat: req.KontakDarurat,
	})
	if err != nil {
		if errors.Is(err, services.ErrAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "traveler profile already exists"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, t)
}

func (h *Handler) UpdateTravelerProfile(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(int64)
	var req travelerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	t, err := h.TravelerService.Update(c.Request.Context(), userID, services.TravelerInput{
		NamaLengkap:   req.NamaLengkap,
		TanggalLahir:  req.TanggalLahir,
		KontakDarurat: req.KontakDarurat,
	})
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "traveler profile not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, t)
}
