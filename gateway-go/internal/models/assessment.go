package models

import (
	"encoding/json"
	"time"
)

type HealthAssessment struct {
	ID              int64           `json:"id"`
	UserID          int64           `json:"user_id"`
	Symptoms        json.RawMessage `json:"symptoms"`
	AIAnalysisRaw   json.RawMessage `json:"ai_analysis_raw,omitempty"`
	Diagnosis       *string         `json:"diagnosis,omitempty"`
	ConfidenceScore *float64        `json:"confidence_score,omitempty"`
	RiskLevel       *string         `json:"risk_level,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
}

type AssessmentPage struct {
	Data  []HealthAssessment `json:"data"`
	Total int64              `json:"total"`
	Page  int                `json:"page"`
	Limit int                `json:"limit"`
}
