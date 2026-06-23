package models

import (
	"encoding/json"
	"time"
)

type Destination struct {
	ID         int64     `json:"id"`
	NamaDaerah string    `json:"nama_daerah"`
	CreatedAt  time.Time `json:"created_at"`
}

type HealthRisk struct {
	ID                      int64     `json:"id"`
	DestinationID           int64     `json:"destination_id"`
	NamaRisikoID            string    `json:"nama_risiko_id"`
	NamaRisikoEN            string    `json:"nama_risiko_en"` // Required
	SaranPencegahanID       *string   `json:"saran_pencegahan_id,omitempty"`
	SaranPencegahanEN       *string   `json:"saran_pencegahan_en,omitempty"`
	RekomendasiVaksinasiID  *string   `json:"rekomendasi_vaksinasi_id,omitempty"`
	RekomendasiVaksinasiEN  *string   `json:"rekomendasi_vaksinasi_en,omitempty"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
}

type EmergencyGuide struct {
	ID          int64           `json:"id"`
	Kategori    string          `json:"kategori"`
	Langkah     int             `json:"langkah"`
	IsiMediaID  json.RawMessage `json:"isi_media_id"`
	IsiMediaEN  json.RawMessage `json:"isi_media_en"` // Required
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

// ── EmergencyGuideFlow ─────────────────────────────────────────────────────

type GuideChoice struct {
	Label   string  `json:"label"`
	NextID  *string `json:"next_id"`
	Variant string  `json:"variant"` // "yes" | "no" | "neutral"
}

type GuideNode struct {
	ID          string        `json:"id"`
	Title       string        `json:"title"`
	Instruction string        `json:"instruction"`
	ImageURL    string        `json:"image_url"`
	IsEntry     bool          `json:"is_entry"`
	Choices     []GuideChoice `json:"choices"`
}

type EmergencyGuideFlow struct {
	ID        int64       `json:"id"`
	TitleID   string      `json:"title_id"`
	TitleEN   string      `json:"title_en"` // Required
	Kategori  string      `json:"kategori"`
	Deskripsi *string     `json:"deskripsi,omitempty"`
	NodesID   []GuideNode `json:"nodes_id"`
	NodesEN   []GuideNode `json:"nodes_en"` // Required
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}
