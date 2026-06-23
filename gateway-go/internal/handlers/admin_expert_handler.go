package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/balitravelhealth/platform/gateway-go/internal/services"
)

// ── GO-28: Symptoms ───────────────────────────────────────────────────────────

func (h *Handler) ListExpertSymptoms(c *gin.Context) {
	lang := resolveLang(c)
	list, err := h.ExpertAdminService.ListPublicSymptoms(c.Request.Context(), c.Query("kategori"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out := make([]map[string]any, 0, len(list))
	for _, s := range list {
		label := s.LabelID
		if lang == "en" && s.LabelEN != "" {
			label = s.LabelEN
		}
		out = append(out, map[string]any{
			"id":         s.SymptomID,
			"symptom_id": s.SymptomID,
			"kode":       s.Kode,
			"label":      label,
			"label_id":   s.LabelID,
			"label_en":   s.LabelEN,
		})
	}
	c.JSON(http.StatusOK, gin.H{"data": out})
}


func (h *Handler) AdminListSymptoms(c *gin.Context) {
	list, err := h.ExpertAdminService.ListSymptoms(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list})
}

func (h *Handler) AdminCreateSymptom(c *gin.Context) {
	var req struct {
		Kode        string  `json:"kode" binding:"required"`
		LabelID     string  `json:"label_id" binding:"required"`
		LabelEN     string  `json:"label_en" binding:"required"` // Required
		DeskripsiID *string `json:"deskripsi_id"`
		DeskripsiEN *string `json:"deskripsi_en"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sym, err := h.ExpertAdminService.CreateSymptom(c.Request.Context(), req.Kode, req.LabelID, req.LabelEN, req.DeskripsiID, req.DeskripsiEN)
	if errors.Is(err, services.ErrAlreadyExists) {
		c.JSON(http.StatusConflict, gin.H{"error": "symptom with this kode already exists"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusCreated, sym)
}

func (h *Handler) AdminUpdateSymptom(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req struct {
		Kode        string  `json:"kode" binding:"required"`
		LabelID     string  `json:"label_id" binding:"required"`
		LabelEN     string  `json:"label_en" binding:"required"` // Required
		DeskripsiID *string `json:"deskripsi_id"`
		DeskripsiEN *string `json:"deskripsi_en"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sym, err := h.ExpertAdminService.UpdateSymptom(c.Request.Context(), id, req.Kode, req.LabelID, req.LabelEN, req.DeskripsiID, req.DeskripsiEN)
	if errors.Is(err, services.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "symptom not found"})
		return
	}
	if errors.Is(err, services.ErrAlreadyExists) {
		c.JSON(http.StatusConflict, gin.H{"error": "symptom with this kode already exists"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, sym)
}

func (h *Handler) AdminDeleteSymptom(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.ExpertAdminService.DeleteSymptom(c.Request.Context(), id); errors.Is(err, services.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "symptom not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.Status(http.StatusNoContent)
}

// ── GO-28: Diseases ───────────────────────────────────────────────────────────

func (h *Handler) AdminListDiseases(c *gin.Context) {
	list, err := h.ExpertAdminService.ListDiseases(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list})
}

type diseaseAdminRequest struct {
	NamaID            string          `json:"nama_id" binding:"required"`
	NamaEN            string          `json:"nama_en" binding:"required"` // Required
	DeskripsiID       *string         `json:"deskripsi_id"`
	DeskripsiEN       *string         `json:"deskripsi_en"`
	RekDefaultID      json.RawMessage `json:"rekomendasi_default_id"`
	RekDefaultEN      json.RawMessage `json:"rekomendasi_default_en"`
}

func (h *Handler) AdminCreateDisease(c *gin.Context) {
	var req diseaseAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	d, err := h.ExpertAdminService.CreateDisease(c.Request.Context(), req.NamaID, req.NamaEN, req.DeskripsiID, req.DeskripsiEN, req.RekDefaultID, req.RekDefaultEN)
	if errors.Is(err, services.ErrAlreadyExists) {
		c.JSON(http.StatusConflict, gin.H{"error": "disease with this name already exists"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusCreated, d)
}

func (h *Handler) AdminUpdateDisease(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req diseaseAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	d, err := h.ExpertAdminService.UpdateDisease(c.Request.Context(), id, req.NamaID, req.NamaEN, req.DeskripsiID, req.DeskripsiEN, req.RekDefaultID, req.RekDefaultEN)
	if errors.Is(err, services.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "disease not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, d)
}

func (h *Handler) AdminDeleteDisease(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.ExpertAdminService.DeleteDisease(c.Request.Context(), id); errors.Is(err, services.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "disease not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.Status(http.StatusNoContent)
}

// ── GO-28: Rules ──────────────────────────────────────────────────────────────

type ruleAdminRequest struct {
	Nama      string  `json:"nama" binding:"required"`
	Premis    []int64 `json:"premis" binding:"required"`
	DiseaseID int64   `json:"disease_id" binding:"required"`
	BobotCF   float64 `json:"bobot_cf"`
	MB        float64 `json:"mb"`
	MD        float64 `json:"md"`
	Kategori  string  `json:"kategori" binding:"required"`
}

func (h *Handler) AdminListRules(c *gin.Context) {
	list, err := h.ExpertAdminService.ListRules(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list})
}

func (h *Handler) AdminCreateRule(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var req ruleAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	rule, err := h.ExpertAdminService.CreateRule(c.Request.Context(), services.RuleInput{
		Nama:      req.Nama,
		Premis:    req.Premis,
		DiseaseID: req.DiseaseID,
		BobotCF:   req.BobotCF,
		MB:        req.MB,
		MD:        req.MD,
		Kategori:  req.Kategori,
		CreatedBy: userID.(int64),
	})
	if errors.Is(err, services.ErrRuleValidation) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusCreated, rule)
}

func (h *Handler) AdminUpdateRule(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	userID, _ := c.Get("user_id")
	var req ruleAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	rule, err := h.ExpertAdminService.UpdateRule(c.Request.Context(), id, services.RuleInput{
		Nama:      req.Nama,
		Premis:    req.Premis,
		DiseaseID: req.DiseaseID,
		BobotCF:   req.BobotCF,
		MB:        req.MB,
		MD:        req.MD,
		Kategori:  req.Kategori,
		CreatedBy: userID.(int64),
	})
	if errors.Is(err, services.ErrRuleValidation) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	if errors.Is(err, services.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "rule not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, rule)
}

func (h *Handler) AdminPublishRule(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	rule, err := h.ExpertAdminService.PublishRule(c.Request.Context(), id)
	if errors.Is(err, services.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "rule not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, rule)
}

func (h *Handler) AdminUnpublishRule(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	rule, err := h.ExpertAdminService.UnpublishRule(c.Request.Context(), id)
	if errors.Is(err, services.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "rule not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, rule)
}

func (h *Handler) AdminDeleteRule(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.ExpertAdminService.DeleteRule(c.Request.Context(), id); errors.Is(err, services.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "rule not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.Status(http.StatusNoContent)
}
