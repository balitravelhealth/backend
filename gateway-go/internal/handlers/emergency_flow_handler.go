package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/balitravelhealth/platform/gateway-go/internal/models"
)

// resolveLang returns "en" or "id" from the ?lang= query param.
func resolveLang(c *gin.Context) string {
	if c.Query("lang") == "en" {
		return "en"
	}
	return "id"
}

// ── Public: GET /emergency-guide-flows ───────────────────────────────────────

func (h *Handler) ListEmergencyFlows(c *gin.Context) {
	lang := resolveLang(c)

	rows, err := h.DB.Query(c.Request.Context(),
		`SELECT id, title_id, title_en, kategori, deskripsi, created_at, updated_at
		 FROM emergency_guide_flows ORDER BY id`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	defer rows.Close()

	list := []map[string]any{}
	for rows.Next() {
		var (
			id        int64
			titleID   string
			titleEN   string
			kategori  string
			deskripsi *string
			createdAt any
			updatedAt any
		)
		if err := rows.Scan(&id, &titleID, &titleEN, &kategori, &deskripsi, &createdAt, &updatedAt); err != nil {
			continue
		}
		title := titleID
		if lang == "en" && titleEN != "" {
			title = titleEN
		}
		list = append(list, map[string]any{
			"id":         id,
			"title":      title,
			"kategori":   kategori,
			"deskripsi":  deskripsi,
			"created_at": createdAt,
			"updated_at": updatedAt,
		})
	}
	c.JSON(http.StatusOK, gin.H{"data": list})
}

// ── Public: GET /emergency-guide-flows/:id ────────────────────────────────────

func (h *Handler) GetEmergencyFlow(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	lang := resolveLang(c)
	flow, err := h.getFlow(c, id, lang)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "flow not found"})
		return
	}
	c.JSON(http.StatusOK, flow)
}

// getFlow fetches a single flow and returns it with language-appropriate title/nodes.
func (h *Handler) getFlow(c *gin.Context, id int64, lang string) (map[string]any, error) {
	var (
		f          models.EmergencyGuideFlow
		nodesIDRaw []byte
		nodesENRaw []byte
	)
	err := h.DB.QueryRow(c.Request.Context(),
		`SELECT id, title_id, title_en, kategori, deskripsi, nodes_id, nodes_en, created_at, updated_at
		 FROM emergency_guide_flows WHERE id = $1`, id,
	).Scan(&f.ID, &f.TitleID, &f.TitleEN, &f.Kategori, &f.Deskripsi,
		&nodesIDRaw, &nodesENRaw, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		return nil, err
	}

	_ = json.Unmarshal(nodesIDRaw, &f.NodesID)
	_ = json.Unmarshal(nodesENRaw, &f.NodesEN)

	title := f.TitleID
	nodes := f.NodesID
	if lang == "en" {
		if f.TitleEN != "" {
			title = f.TitleEN
		}
		if len(f.NodesEN) > 0 {
			nodes = f.NodesEN
		}
	}

	return map[string]any{
		"id":         f.ID,
		"title":      title,
		"kategori":   f.Kategori,
		"deskripsi":  f.Deskripsi,
		"nodes":      nodes,
		"created_at": f.CreatedAt,
		"updated_at": f.UpdatedAt,
	}, nil
}

// ── Admin: GET /admin/emergency-guide-flows ───────────────────────────────────

func (h *Handler) AdminListEmergencyFlows(c *gin.Context) {
	rows, err := h.DB.Query(c.Request.Context(),
		`SELECT id, title_id, title_en, kategori, deskripsi, nodes_id, nodes_en, created_at, updated_at
		 FROM emergency_guide_flows ORDER BY id`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	defer rows.Close()

	var list []models.EmergencyGuideFlow
	for rows.Next() {
		var f models.EmergencyGuideFlow
		var nodesIDRaw, nodesENRaw []byte
		if err := rows.Scan(&f.ID, &f.TitleID, &f.TitleEN, &f.Kategori, &f.Deskripsi,
			&nodesIDRaw, &nodesENRaw, &f.CreatedAt, &f.UpdatedAt); err != nil {
			continue
		}
		_ = json.Unmarshal(nodesIDRaw, &f.NodesID)
		_ = json.Unmarshal(nodesENRaw, &f.NodesEN)
		list = append(list, f)
	}
	if list == nil {
		list = []models.EmergencyGuideFlow{}
	}
	c.JSON(http.StatusOK, gin.H{"data": list})
}

// ── Admin: POST /admin/emergency-guide-flows ──────────────────────────────────

func (h *Handler) AdminCreateEmergencyFlow(c *gin.Context) {
	var req struct {
		TitleID   string             `json:"title_id" binding:"required"`
		TitleEN   string             `json:"title_en" binding:"required"`
		Kategori  string             `json:"kategori" binding:"required"`
		Deskripsi *string            `json:"deskripsi"`
		NodesID   []models.GuideNode `json:"nodes_id"`
		NodesEN   []models.GuideNode `json:"nodes_en"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.NodesID == nil {
		req.NodesID = []models.GuideNode{}
	}
	if req.NodesEN == nil {
		req.NodesEN = []models.GuideNode{}
	}
	nodesIDJSON, _ := json.Marshal(req.NodesID)
	nodesENJSON, _ := json.Marshal(req.NodesEN)

	var f models.EmergencyGuideFlow
	var nodesIDRaw, nodesENRaw []byte
	err := h.DB.QueryRow(c.Request.Context(),
		`INSERT INTO emergency_guide_flows (title_id, title_en, kategori, deskripsi, nodes_id, nodes_en)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING id, title_id, title_en, kategori, deskripsi, nodes_id, nodes_en, created_at, updated_at`,
		req.TitleID, req.TitleEN, req.Kategori, req.Deskripsi, nodesIDJSON, nodesENJSON,
	).Scan(&f.ID, &f.TitleID, &f.TitleEN, &f.Kategori, &f.Deskripsi,
		&nodesIDRaw, &nodesENRaw, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	_ = json.Unmarshal(nodesIDRaw, &f.NodesID)
	_ = json.Unmarshal(nodesENRaw, &f.NodesEN)
	c.JSON(http.StatusCreated, f)
}

// ── Admin: PUT /admin/emergency-guide-flows/:id ───────────────────────────────

func (h *Handler) AdminUpdateEmergencyFlow(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req struct {
		TitleID   string             `json:"title_id" binding:"required"`
		TitleEN   string             `json:"title_en" binding:"required"`
		Kategori  string             `json:"kategori" binding:"required"`
		Deskripsi *string            `json:"deskripsi"`
		NodesID   []models.GuideNode `json:"nodes_id"`
		NodesEN   []models.GuideNode `json:"nodes_en"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.NodesID == nil {
		req.NodesID = []models.GuideNode{}
	}
	if req.NodesEN == nil {
		req.NodesEN = []models.GuideNode{}
	}
	nodesIDJSON, _ := json.Marshal(req.NodesID)
	nodesENJSON, _ := json.Marshal(req.NodesEN)

	var f models.EmergencyGuideFlow
	var nodesIDRaw, nodesENRaw []byte
	err = h.DB.QueryRow(c.Request.Context(),
		`UPDATE emergency_guide_flows
		 SET title_id=$1, title_en=$2, kategori=$3, deskripsi=$4, nodes_id=$5, nodes_en=$6, updated_at=now()
		 WHERE id=$7
		 RETURNING id, title_id, title_en, kategori, deskripsi, nodes_id, nodes_en, created_at, updated_at`,
		req.TitleID, req.TitleEN, req.Kategori, req.Deskripsi, nodesIDJSON, nodesENJSON, id,
	).Scan(&f.ID, &f.TitleID, &f.TitleEN, &f.Kategori, &f.Deskripsi,
		&nodesIDRaw, &nodesENRaw, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "flow not found"})
		return
	}
	_ = json.Unmarshal(nodesIDRaw, &f.NodesID)
	_ = json.Unmarshal(nodesENRaw, &f.NodesEN)
	c.JSON(http.StatusOK, f)
}

// ── Admin: DELETE /admin/emergency-guide-flows/:id ────────────────────────────

func (h *Handler) AdminDeleteEmergencyFlow(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	tag, err := h.DB.Exec(c.Request.Context(),
		`DELETE FROM emergency_guide_flows WHERE id=$1`, id)
	if err != nil || tag.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "flow not found"})
		return
	}
	c.Status(http.StatusNoContent)
}
