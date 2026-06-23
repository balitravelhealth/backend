package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/balitravelhealth/platform/gateway-go/internal/middleware"
	"github.com/balitravelhealth/platform/gateway-go/internal/services"
)

type healthProfileRequest struct {
	TanggalLahir  *string  `json:"tanggal_lahir"`
	JenisKelamin  *string  `json:"jenis_kelamin"`
	TinggiCm      *float64 `json:"tinggi_cm"`
	BeratKg       *float64 `json:"berat_kg"`
	GolonganDarah *string  `json:"golongan_darah"`
	RiwayatAlergi *string  `json:"riwayat_alergi"`
}

func (h *Handler) GetHealthProfile(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(int64)

	profile, err := h.ProfileService.Get(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "health profile not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

func (h *Handler) CreateHealthProfile(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(int64)

	var req healthProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	profile, err := h.ProfileService.Create(c.Request.Context(), userID, services.ProfileInput{
		TanggalLahir:  req.TanggalLahir,
		JenisKelamin:  req.JenisKelamin,
		TinggiCm:      req.TinggiCm,
		BeratKg:       req.BeratKg,
		GolonganDarah: req.GolonganDarah,
		RiwayatAlergi: req.RiwayatAlergi,
	})
	if err != nil {
		if errors.Is(err, services.ErrAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "health profile already exists"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, profile)
}

func (h *Handler) UpdateHealthProfile(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(int64)

	var req healthProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	profile, err := h.ProfileService.Update(c.Request.Context(), userID, services.ProfileInput{
		TanggalLahir:  req.TanggalLahir,
		JenisKelamin:  req.JenisKelamin,
		TinggiCm:      req.TinggiCm,
		BeratKg:       req.BeratKg,
		GolonganDarah: req.GolonganDarah,
		RiwayatAlergi: req.RiwayatAlergi,
	})
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "health profile not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, profile)
}
