package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/balitravelhealth/platform/gateway-go/internal/middleware"
	"github.com/balitravelhealth/platform/gateway-go/internal/repository"
	"github.com/balitravelhealth/platform/gateway-go/internal/services"
)

// GO-17
func (h *Handler) ListNurses(c *gin.Context) {
	list, err := h.NursingService.ListNurses(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list})
}

// GO-18 — traveler creates appointment
type appointmentRequest struct {
	NurseID          int64  `json:"nurse_id"`
	TanggalKunjungan string `json:"tanggal_kunjungan"`
}

func (h *Handler) CreateAppointment(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(int64)

	var req appointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.NurseID == 0 || req.TanggalKunjungan == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "nurse_id and tanggal_kunjungan are required"})
		return
	}

	rec, err := h.NursingService.CreateAppointment(c.Request.Context(), userID, req.NurseID, req.TanggalKunjungan)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, rec)
}

// Traveler sees their own nursing records
func (h *Handler) ListMyNursingRecords(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(int64)
	list, err := h.NursingService.ListTravelerRecords(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list})
}

// GO-19 — nurse fills in care record
type careRecordRequest struct {
	NursingAssessment     *string `json:"nursing_assessment"`
	NursingDiagnosis      *string `json:"nursing_diagnosis"`
	NursingPlanning       *string `json:"nursing_planning"`
	NursingImplementation *string `json:"nursing_implementation"`
	NursingEvaluation     *string `json:"nursing_evaluation"`
}

func (h *Handler) UpdateCareRecord(c *gin.Context) {
	nurseUserID := c.MustGet(middleware.UserIDKey).(int64)

	recordID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid record id"})
		return
	}

	var req careRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rec, err := h.NursingService.UpdateCareRecord(c.Request.Context(), nurseUserID, recordID,
		repository.NursingRecordUpdate{
			Assessment:     req.NursingAssessment,
			Diagnosis:      req.NursingDiagnosis,
			Planning:       req.NursingPlanning,
			Implementation: req.NursingImplementation,
			Evaluation:     req.NursingEvaluation,
		},
	)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "record not found or not assigned to you"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rec)
}

// Nurse sees records assigned to them
func (h *Handler) ListNurseRecords(c *gin.Context) {
	nurseUserID := c.MustGet(middleware.UserIDKey).(int64)
	list, err := h.NursingService.ListNurseRecords(c.Request.Context(), nurseUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list})
}
