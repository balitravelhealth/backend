package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/balitravelhealth/platform/gateway-go/internal/services"
)

// ── GO-24: Facilities ─────────────────────────────────────────────────────────

type facilityAdminRequest struct {
	DestinationID  int64    `json:"destination_id" binding:"required"`
	Nama           string   `json:"nama" binding:"required"`
	Kategori       *string  `json:"kategori"`
	Alamat         *string  `json:"alamat"`
	Latitude       *float64 `json:"latitude"`
	Longitude      *float64 `json:"longitude"`
	Kontak         *string  `json:"kontak"`
	JamOperasional *string  `json:"jam_operasional"`
}

func (h *Handler) AdminListFacilities(c *gin.Context) {
	list, err := h.AdminService.ListFacilities(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list})
}

func (h *Handler) AdminCreateFacility(c *gin.Context) {
	var req facilityAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	f, err := h.AdminService.CreateFacility(c.Request.Context(), services.FacilityInput{
		DestinationID:  req.DestinationID,
		Nama:           req.Nama,
		Kategori:       req.Kategori,
		Alamat:         req.Alamat,
		Latitude:       req.Latitude,
		Longitude:      req.Longitude,
		Kontak:         req.Kontak,
		JamOperasional: req.JamOperasional,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusCreated, f)
}

func (h *Handler) AdminUpdateFacility(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req facilityAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	f, err := h.AdminService.UpdateFacility(c.Request.Context(), id, services.FacilityInput{
		DestinationID:  req.DestinationID,
		Nama:           req.Nama,
		Kategori:       req.Kategori,
		Alamat:         req.Alamat,
		Latitude:       req.Latitude,
		Longitude:      req.Longitude,
		Kontak:         req.Kontak,
		JamOperasional: req.JamOperasional,
	})
	if errors.Is(err, services.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "facility not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, f)
}

func (h *Handler) AdminDeleteFacility(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.AdminService.DeleteFacility(c.Request.Context(), id); errors.Is(err, services.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "facility not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.Status(http.StatusNoContent)
}

// ── GO-25: Destinations ───────────────────────────────────────────────────────

func (h *Handler) AdminCreateDestination(c *gin.Context) {
	var req struct {
		NamaDaerah string `json:"nama_daerah" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "nama_daerah is required"})
		return
	}
	d, err := h.AdminService.CreateDestination(c.Request.Context(), req.NamaDaerah)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusCreated, d)
}

func (h *Handler) AdminUpdateDestination(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req struct {
		NamaDaerah string `json:"nama_daerah" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "nama_daerah is required"})
		return
	}
	d, err := h.AdminService.UpdateDestination(c.Request.Context(), id, req.NamaDaerah)
	if errors.Is(err, services.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "destination not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, d)
}

func (h *Handler) AdminDeleteDestination(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.AdminService.DeleteDestination(c.Request.Context(), id); errors.Is(err, services.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "destination not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.Status(http.StatusNoContent)
}

// ── GO-25: Health Risks ───────────────────────────────────────────────────────

type healthRiskAdminRequest struct {
	DestinationID          int64   `json:"destination_id" binding:"required"`
	NamaRisikoID           string  `json:"nama_risiko_id" binding:"required"`
	NamaRisikoEN           string  `json:"nama_risiko_en" binding:"required"` // Required
	SaranPencegahanID      *string `json:"saran_pencegahan_id"`
	SaranPencegahanEN      *string `json:"saran_pencegahan_en"`
	RekomendasiVaksinasiID *string `json:"rekomendasi_vaksinasi_id"`
	RekomendasiVaksinasiEN *string `json:"rekomendasi_vaksinasi_en"`
}

func (h *Handler) AdminCreateHealthRisk(c *gin.Context) {
	var req healthRiskAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hr, err := h.AdminService.CreateHealthRisk(c.Request.Context(), services.HealthRiskInput{
		DestinationID:          req.DestinationID,
		NamaRisikoID:           req.NamaRisikoID,
		NamaRisikoEN:           req.NamaRisikoEN,
		SaranPencegahanID:      req.SaranPencegahanID,
		SaranPencegahanEN:      req.SaranPencegahanEN,
		RekomendasiVaksinasiID: req.RekomendasiVaksinasiID,
		RekomendasiVaksinasiEN: req.RekomendasiVaksinasiEN,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusCreated, hr)
}

func (h *Handler) AdminUpdateHealthRisk(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req healthRiskAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hr, err := h.AdminService.UpdateHealthRisk(c.Request.Context(), id, services.HealthRiskInput{
		DestinationID:          req.DestinationID,
		NamaRisikoID:           req.NamaRisikoID,
		NamaRisikoEN:           req.NamaRisikoEN,
		SaranPencegahanID:      req.SaranPencegahanID,
		SaranPencegahanEN:      req.SaranPencegahanEN,
		RekomendasiVaksinasiID: req.RekomendasiVaksinasiID,
		RekomendasiVaksinasiEN: req.RekomendasiVaksinasiEN,
	})
	if errors.Is(err, services.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "health risk not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, hr)
}

func (h *Handler) AdminDeleteHealthRisk(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.AdminService.DeleteHealthRisk(c.Request.Context(), id); errors.Is(err, services.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "health risk not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.Status(http.StatusNoContent)
}

// ── GO-25: Emergency Guides ───────────────────────────────────────────────────

type emergencyGuideAdminRequest struct {
	Kategori   string          `json:"kategori" binding:"required"`
	Langkah    int             `json:"langkah" binding:"required"`
	IsiMediaID json.RawMessage `json:"isi_media_id" binding:"required"`
	IsiMediaEN json.RawMessage `json:"isi_media_en" binding:"required"` // Required
}

func (h *Handler) AdminCreateEmergencyGuide(c *gin.Context) {
	var req emergencyGuideAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	g, err := h.AdminService.CreateEmergencyGuide(c.Request.Context(), services.EmergencyGuideInput{
		Kategori:   req.Kategori,
		Langkah:    req.Langkah,
		IsiMediaID: req.IsiMediaID,
		IsiMediaEN: req.IsiMediaEN,
	})
	if errors.Is(err, services.ErrAlreadyExists) {
		c.JSON(http.StatusConflict, gin.H{"error": "guide with this kategori and langkah already exists"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusCreated, g)
}

func (h *Handler) AdminUpdateEmergencyGuide(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req emergencyGuideAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	g, err := h.AdminService.UpdateEmergencyGuide(c.Request.Context(), id, services.EmergencyGuideInput{
		Kategori:   req.Kategori,
		Langkah:    req.Langkah,
		IsiMediaID: req.IsiMediaID,
		IsiMediaEN: req.IsiMediaEN,
	})
	if errors.Is(err, services.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "emergency guide not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, g)
}

func (h *Handler) AdminDeleteEmergencyGuide(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.AdminService.DeleteEmergencyGuide(c.Request.Context(), id); errors.Is(err, services.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "emergency guide not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.Status(http.StatusNoContent)
}

// ── GO-26: Nurse management ───────────────────────────────────────────────────

type createNurseRequest struct {
	Email        string  `json:"email" binding:"required"`
	Password     string  `json:"password" binding:"required,min=8"`
	NamaLengkap  string  `json:"nama_lengkap" binding:"required"`
	NomorLisensi string  `json:"nomor_lisensi" binding:"required"`
	Sertifikasi  *string `json:"sertifikasi"`
}

func (h *Handler) AdminListNurses(c *gin.Context) {
	list, err := h.AdminService.ListAllNurses(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list})
}

func (h *Handler) AdminCreateNurse(c *gin.Context) {
	var req createNurseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	nurse, err := h.AdminService.CreateNurse(c.Request.Context(), services.CreateNurseInput{
		Email:        req.Email,
		Password:     req.Password,
		NamaLengkap:  req.NamaLengkap,
		NomorLisensi: req.NomorLisensi,
		Sertifikasi:  req.Sertifikasi,
	})
	if errors.Is(err, services.ErrAlreadyExists) {
		c.JSON(http.StatusConflict, gin.H{"error": "email already registered"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusCreated, nurse)
}

func (h *Handler) AdminToggleNurse(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	nurse, err := h.AdminService.ToggleNurse(c.Request.Context(), id)
	if errors.Is(err, services.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "nurse not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, nurse)
}

// ── GO-27: View all assessments ───────────────────────────────────────────────

func (h *Handler) AdminListAssessments(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	result, err := h.AdminService.ListAllAssessments(c.Request.Context(), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, result)
}
