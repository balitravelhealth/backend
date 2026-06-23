package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/balitravelhealth/platform/gateway-go/internal/models"
	"github.com/balitravelhealth/platform/gateway-go/internal/repository"
)

type TravelerService struct {
	repo *repository.TravelerRepo
}

func NewTravelerService(repo *repository.TravelerRepo) *TravelerService {
	return &TravelerService{repo: repo}
}

func (s *TravelerService) Get(ctx context.Context, userID int64) (*models.Traveler, error) {
	t, err := s.repo.FindByUserID(ctx, userID)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrNotFound
	}
	return t, err
}

type TravelerInput struct {
	NamaLengkap   string
	TanggalLahir  *string
	KontakDarurat *string
}

func (s *TravelerService) Create(ctx context.Context, userID int64, in TravelerInput) (*models.Traveler, error) {
	if in.NamaLengkap == "" {
		return nil, fmt.Errorf("nama_lengkap is required")
	}
	tgl, err := parseDate(in.TanggalLahir)
	if err != nil {
		return nil, fmt.Errorf("tanggal_lahir: %w", err)
	}
	t, err := s.repo.Create(ctx, repository.TravelerInput{
		UserID:        userID,
		NamaLengkap:   in.NamaLengkap,
		TanggalLahir:  tgl,
		KontakDarurat: in.KontakDarurat,
	})
	if errors.Is(err, repository.ErrAlreadyExists) {
		return nil, ErrAlreadyExists
	}
	return t, err
}

func (s *TravelerService) Update(ctx context.Context, userID int64, in TravelerInput) (*models.Traveler, error) {
	if in.NamaLengkap == "" {
		return nil, fmt.Errorf("nama_lengkap is required")
	}
	tgl, err := parseDate(in.TanggalLahir)
	if err != nil {
		return nil, fmt.Errorf("tanggal_lahir: %w", err)
	}
	t, err := s.repo.Update(ctx, repository.TravelerInput{
		UserID:        userID,
		NamaLengkap:   in.NamaLengkap,
		TanggalLahir:  tgl,
		KontakDarurat: in.KontakDarurat,
	})
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrNotFound
	}
	return t, err
}
