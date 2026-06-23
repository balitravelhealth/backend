package models

import (
	"encoding/json"
	"time"
)

type ExpertSymptom struct {
	SymptomID    int64      `json:"symptom_id"`
	Kode         string     `json:"kode"`
	LabelID      string     `json:"label_id"`
	LabelEN      string     `json:"label_en"` // Required
	DeskripsiID  *string    `json:"deskripsi_id,omitempty"`
	DeskripsiEN  *string    `json:"deskripsi_en,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type ExpertDisease struct {
	ID                    int64           `json:"id"`
	NamaID                string          `json:"nama_id"`
	NamaEN                string          `json:"nama_en"` // Required
	DeskripsiID           *string         `json:"deskripsi_id,omitempty"`
	DeskripsiEN           *string         `json:"deskripsi_en,omitempty"`
	RekomendasiDefaultID  json.RawMessage `json:"rekomendasi_default_id,omitempty"`
	RekomendasiDefaultEN  json.RawMessage `json:"rekomendasi_default_en,omitempty"`
	CreatedAt             time.Time       `json:"created_at"`
	UpdatedAt             time.Time       `json:"updated_at"`
}

type ExpertRule struct {
	RuleID    int64           `json:"rule_id"`
	Nama      string          `json:"nama"`
	Premis    json.RawMessage `json:"premis"`
	DiseaseID int64           `json:"disease_id"`
	BobotCF   float64         `json:"bobot_cf"`
	MB        float64         `json:"mb"`
	MD        float64         `json:"md"`
	Kategori  string          `json:"kategori"`
	Status    string          `json:"status"`
	CreatedBy int64           `json:"created_by"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}
