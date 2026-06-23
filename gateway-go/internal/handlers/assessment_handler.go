package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/balitravelhealth/platform/gateway-go/internal/middleware"
	"github.com/balitravelhealth/platform/gateway-go/internal/services"
)

func (h *Handler) ListAssessments(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(int64)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	result, err := h.AssessmentService.List(c.Request.Context(), userID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, result)
}

type assessmentRequest struct {
	Symptoms []int64 `json:"symptoms"`
	Kategori string  `json:"kategori"`
	Language string  `json:"language"` // "id" | "en" — defaults to "id"
}

// PostAssessment — GO-20/21/22
func (h *Handler) PostAssessment(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(int64)

	var req assessmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(req.Symptoms) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "symptoms must not be empty"})
		return
	}
	if req.Kategori == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "kategori is required (pre_travel / post_travel)"})
		return
	}

	lang := req.Language
	if lang != "en" {
		lang = "id"
	}

	// Build optional user_profile from health_profile
	expertReq := services.ExpertRequest{
		Symptoms: req.Symptoms,
		Kategori: req.Kategori,
		Language: lang,
	}
	if profile, err := h.ProfileService.Get(c.Request.Context(), userID); err == nil {
		ep := &services.ExpertUserProfile{
			WeightKg: profile.BeratKg,
			Gender:   profile.JenisKelamin,
		}
		if profile.TanggalLahir != nil {
			age := ageFromDOB(*profile.TanggalLahir)
			ep.Age = &age
		}
		expertReq.UserProfile = ep
	}

	assessment, err := h.AssessmentService.Submit(c.Request.Context(), userID, expertReq)
	if err != nil {
		// GO-21: fail-safe — expert service down, other endpoints still work
		if errors.Is(err, services.ErrExpertServiceUnavailable) {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error": "layanan diagnosis sementara tidak tersedia, coba lagi nanti",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, assessment)
}

func ageFromDOB(dob time.Time) int {
	now := time.Now()
	age := now.Year() - dob.Year()
	if now.YearDay() < dob.YearDay() {
		age--
	}
	return age
}
