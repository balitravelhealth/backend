package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/balitravelhealth/platform/gateway-go/internal/models"
	"github.com/balitravelhealth/platform/gateway-go/internal/repository"
)

type VaccinationService struct {
	repo *repository.VaccinationRepo
}

func NewVaccinationService(repo *repository.VaccinationRepo) *VaccinationService {
	return &VaccinationService{repo: repo}
}

func (s *VaccinationService) List(ctx context.Context, userID int64) ([]models.VaccinationRecord, error) {
	list, err := s.repo.ListByUserID(ctx, userID)
	if list == nil {
		list = []models.VaccinationRecord{}
	}
	return list, err
}

type VaccinationInput struct {
	JenisVaksin string
	Tanggal     string
	Dosis       *string
	Catatan     *string
}

func (s *VaccinationService) Create(ctx context.Context, userID int64, in VaccinationInput) (*models.VaccinationRecord, error) {
	if in.JenisVaksin == "" {
		return nil, fmt.Errorf("jenis_vaksin is required")
	}
	tgl, err := parseDate(&in.Tanggal)
	if err != nil || tgl == nil {
		return nil, fmt.Errorf("tanggal: format harus YYYY-MM-DD")
	}
	return s.repo.Create(ctx, repository.VaccinationInput{
		UserID:      userID,
		JenisVaksin: in.JenisVaksin,
		Tanggal:     *tgl,
		Dosis:       in.Dosis,
		Catatan:     in.Catatan,
	})
}

func (s *VaccinationService) Delete(ctx context.Context, id, userID int64) error {
	if err := s.repo.Delete(ctx, id, userID); errors.Is(err, repository.ErrNotFound) {
		return ErrNotFound
	} else {
		return err
	}
}
