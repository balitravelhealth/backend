package repository

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/balitravelhealth/platform/gateway-go/internal/models"
)

type AssessmentRepo struct {
	db *pgxpool.Pool
}

func NewAssessmentRepo(db *pgxpool.Pool) *AssessmentRepo {
	return &AssessmentRepo{db: db}
}

func (r *AssessmentRepo) ListByUserID(ctx context.Context, userID int64, limit, offset int) ([]models.HealthAssessment, int64, error) {
	var total int64
	if err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM health_assessments WHERE user_id = $1`, userID,
	).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(ctx,
		`SELECT id, user_id, symptoms, ai_analysis_raw, diagnosis, confidence_score, risk_level, created_at
		 FROM health_assessments
		 WHERE user_id = $1
		 ORDER BY created_at DESC
		 LIMIT $2 OFFSET $3`,
		userID, limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []models.HealthAssessment
	for rows.Next() {
		var a models.HealthAssessment
		if err := rows.Scan(
			&a.ID, &a.UserID, &a.Symptoms, &a.AIAnalysisRaw,
			&a.Diagnosis, &a.ConfidenceScore, &a.RiskLevel, &a.CreatedAt,
		); err != nil {
			return nil, 0, err
		}
		list = append(list, a)
	}
	return list, total, rows.Err()
}

func (r *AssessmentRepo) ListAll(ctx context.Context, limit, offset int) ([]models.HealthAssessment, int64, error) {
	var total int64
	if err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM health_assessments`).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := r.db.Query(ctx,
		`SELECT id, user_id, symptoms, ai_analysis_raw, diagnosis, confidence_score, risk_level, created_at
		 FROM health_assessments ORDER BY created_at DESC LIMIT $1 OFFSET $2`,
		limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var list []models.HealthAssessment
	for rows.Next() {
		var a models.HealthAssessment
		if err := rows.Scan(&a.ID, &a.UserID, &a.Symptoms, &a.AIAnalysisRaw,
			&a.Diagnosis, &a.ConfidenceScore, &a.RiskLevel, &a.CreatedAt); err != nil {
			return nil, 0, err
		}
		list = append(list, a)
	}
	return list, total, rows.Err()
}

type AssessmentInput struct {
	UserID          int64
	Symptoms        json.RawMessage
	AIAnalysisRaw   json.RawMessage
	Diagnosis       *string
	ConfidenceScore *float64
	RiskLevel       *string
}

func (r *AssessmentRepo) Create(ctx context.Context, in AssessmentInput) (*models.HealthAssessment, error) {
	a := &models.HealthAssessment{}
	err := r.db.QueryRow(ctx,
		`INSERT INTO health_assessments
		    (user_id, symptoms, ai_analysis_raw, diagnosis, confidence_score, risk_level)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING id, user_id, symptoms, ai_analysis_raw, diagnosis, confidence_score, risk_level, created_at`,
		in.UserID, in.Symptoms, in.AIAnalysisRaw,
		in.Diagnosis, in.ConfidenceScore, in.RiskLevel,
	).Scan(&a.ID, &a.UserID, &a.Symptoms, &a.AIAnalysisRaw,
		&a.Diagnosis, &a.ConfidenceScore, &a.RiskLevel, &a.CreatedAt)
	return a, err
}
