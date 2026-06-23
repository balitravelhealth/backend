package services

import (
	"context"

	"github.com/balitravelhealth/platform/gateway-go/internal/models"
	"github.com/balitravelhealth/platform/gateway-go/internal/repository"
)

type LocationService struct {
	repo *repository.LocationRepo
}

func NewLocationService(repo *repository.LocationRepo) *LocationService {
	return &LocationService{repo: repo}
}

func (s *LocationService) ListDestinations(ctx context.Context) ([]models.Destination, error) {
	list, err := s.repo.ListDestinations(ctx)
	if list == nil {
		list = []models.Destination{}
	}
	return list, err
}

func (s *LocationService) ListHealthRisks(ctx context.Context, destinationID int64) ([]models.HealthRisk, error) {
	list, err := s.repo.ListHealthRisks(ctx, destinationID)
	if list == nil {
		list = []models.HealthRisk{}
	}
	return list, err
}

func (s *LocationService) ListEmergencyGuides(ctx context.Context, kategori string) ([]models.EmergencyGuide, error) {
	list, err := s.repo.ListEmergencyGuides(ctx, kategori)
	if list == nil {
		list = []models.EmergencyGuide{}
	}
	return list, err
}
