package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/balitravelhealth/platform/gateway-go/internal/models"
	"github.com/balitravelhealth/platform/gateway-go/internal/repository"
)

type AssessmentService struct {
	repo          *repository.AssessmentRepo
	expertService *ExpertService
}

func NewAssessmentService(repo *repository.AssessmentRepo, expertService *ExpertService) *AssessmentService {
	return &AssessmentService{repo: repo, expertService: expertService}
}

// Submit calls the expert service then persists the result. GO-20/21/22.
func (s *AssessmentService) Submit(ctx context.Context, userID int64, req ExpertRequest) (*models.HealthAssessment, error) {
	result, err := s.expertService.Diagnose(ctx, req)
	if err != nil {
		return nil, err // ErrExpertServiceUnavailable propagates to handler
	}

	// Marshal symptoms and full expert response for storage
	symptomsJSON, err := json.Marshal(req.Symptoms)
	if err != nil {
		return nil, fmt.Errorf("marshal symptoms: %w", err)
	}
	rawJSON, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("marshal expert response: %w", err)
	}

	score := result.ConfidenceScore
	return s.repo.Create(ctx, repository.AssessmentInput{
		UserID:          userID,
		Symptoms:        symptomsJSON,
		AIAnalysisRaw:   rawJSON,
		Diagnosis:       &result.Diagnosis,
		ConfidenceScore: &score,
		RiskLevel:       &result.RiskLevel,
	})
}

func (s *AssessmentService) List(ctx context.Context, userID int64, page, limit int) (*models.AssessmentPage, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 10
	}
	offset := (page - 1) * limit

	data, total, err := s.repo.ListByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	if data == nil {
		data = []models.HealthAssessment{}
	}
	return &models.AssessmentPage{
		Data:  data,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}
