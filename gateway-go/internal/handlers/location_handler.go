package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) ListDestinations(c *gin.Context) {
	list, err := h.LocationService.ListDestinations(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list})
}

func (h *Handler) ListHealthRisks(c *gin.Context) {
	destID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid destination id"})
		return
	}
	lang := resolveLang(c)

	list, err := h.LocationService.ListHealthRisks(c.Request.Context(), destID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// Return language-appropriate field names (mobile clients get "nama_risiko", not "_id"/"_en")
	result := make([]map[string]any, 0, len(list))
	for _, r := range list {
		nama := r.NamaRisikoID
		if lang == "en" && r.NamaRisikoEN != "" {
			nama = r.NamaRisikoEN
		}

		var saran *string = r.SaranPencegahanID
		if lang == "en" && r.SaranPencegahanEN != nil {
			saran = r.SaranPencegahanEN
		}

		var vaksinasi *string = r.RekomendasiVaksinasiID
		if lang == "en" && r.RekomendasiVaksinasiEN != nil {
			vaksinasi = r.RekomendasiVaksinasiEN
		}

		result = append(result, map[string]any{
			"id":                    r.ID,
			"destination_id":        r.DestinationID,
			"nama_risiko":           nama,
			"saran_pencegahan":      saran,
			"rekomendasi_vaksinasi": vaksinasi,
			"created_at":            r.CreatedAt,
			"updated_at":            r.UpdatedAt,
		})
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}

func (h *Handler) ListEmergencyGuides(c *gin.Context) {
	kategori := c.Query("kategori")
	lang := resolveLang(c)

	list, err := h.LocationService.ListEmergencyGuides(c.Request.Context(), kategori)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// Return language-appropriate isi_media
	result := make([]map[string]any, 0, len(list))
	for _, g := range list {
		isiMedia := g.IsiMediaID
		if lang == "en" && len(g.IsiMediaEN) > 0 {
			isiMedia = g.IsiMediaEN
		}
		result = append(result, map[string]any{
			"id":         g.ID,
			"kategori":   g.Kategori,
			"langkah":    g.Langkah,
			"isi_media":  isiMedia,
			"created_at": g.CreatedAt,
			"updated_at": g.UpdatedAt,
		})
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}
