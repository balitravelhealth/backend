package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var ErrExpertServiceUnavailable = errors.New("expert service unavailable")

type ExpertUserProfile struct {
	Age      *int     `json:"age,omitempty"`
	WeightKg *float64 `json:"weight_kg,omitempty"`
	Gender   *string  `json:"gender,omitempty"`
}

type ExpertRequest struct {
	Symptoms    []int64            `json:"symptoms"`
	Kategori    string             `json:"kategori"`
	Language    string             `json:"language,omitempty"` // "id" | "en"
	UserProfile *ExpertUserProfile `json:"user_profile,omitempty"`
}

type ExpertResponse struct {
	Diagnosis       string  `json:"diagnosis"`
	ConfidenceScore float64 `json:"confidence_score"`
	RiskLevel       string  `json:"risk_level"`
	Recommendation  string  `json:"recommendation"`
}

type ExpertService struct {
	client  *http.Client
	baseURL string
}

func NewExpertService() *ExpertService {
	return &ExpertService{
		client:  &http.Client{Timeout: 10 * time.Second},
		baseURL: os.Getenv("EXPERT_SERVICE_URL"),
	}
}

// Diagnose forwards the request to the Python expert service.
// Returns ErrExpertServiceUnavailable if the service is down or times out.
func (s *ExpertService) Diagnose(ctx context.Context, req ExpertRequest) (*ExpertResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost,
		s.baseURL+"/diagnose", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(httpReq)
	if err != nil {
		// GO-21: fail-safe — Python is down or timed out
		log.Printf("expert service error: %v", err)
		return nil, ErrExpertServiceUnavailable
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(resp.Body)
		log.Printf("expert service returned %d: %s", resp.StatusCode, raw)
		return nil, ErrExpertServiceUnavailable
	}

	var result ExpertResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode expert response: %w", err)
	}
	return &result, nil
}

// Ping checks if the expert service health endpoint is reachable.
func (s *ExpertService) Ping(ctx context.Context) bool {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.baseURL+"/health", nil)
	if err != nil {
		return false
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return false
	}
	resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}
