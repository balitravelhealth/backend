package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/balitravelhealth/platform/gateway-go/internal/models"
	"github.com/balitravelhealth/platform/gateway-go/internal/repository"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
)

type ProfileService struct {
	repo *repository.HealthProfileRepo
}

func NewProfileService(repo *repository.HealthProfileRepo) *ProfileService {
	return &ProfileService{repo: repo}
}

func (s *ProfileService) Get(ctx context.Context, userID int64) (*models.HealthProfile, error) {
	p, err := s.repo.FindByUserID(ctx, userID)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrNotFound
	}
	return p, err
}

type ProfileInput struct {
	TanggalLahir  *string
	JenisKelamin  *string
	TinggiCm      *float64
	BeratKg       *float64
	GolonganDarah *string
	RiwayatAlergi *string
}

func (s *ProfileService) Create(ctx context.Context, userID int64, in ProfileInput) (*models.HealthProfile, error) {
	tgl, err := parseDate(in.TanggalLahir)
	if err != nil {
		return nil, fmt.Errorf("tanggal_lahir: %w", err)
	}

	p, err := s.repo.Create(ctx, repository.HealthProfileInput{
		UserID:        userID,
		TanggalLahir:  tgl,
		JenisKelamin:  in.JenisKelamin,
		TinggiCm:      in.TinggiCm,
		BeratKg:       in.BeratKg,
		GolonganDarah: in.GolonganDarah,
		RiwayatAlergi: in.RiwayatAlergi,
	})
	if errors.Is(err, repository.ErrAlreadyExists) {
		return nil, ErrAlreadyExists
	}
	return p, err
}

func (s *ProfileService) Update(ctx context.Context, userID int64, in ProfileInput) (*models.HealthProfile, error) {
	tgl, err := parseDate(in.TanggalLahir)
	if err != nil {
		return nil, fmt.Errorf("tanggal_lahir: %w", err)
	}

	p, err := s.repo.Update(ctx, repository.HealthProfileInput{
		UserID:        userID,
		TanggalLahir:  tgl,
		JenisKelamin:  in.JenisKelamin,
		TinggiCm:      in.TinggiCm,
		BeratKg:       in.BeratKg,
		GolonganDarah: in.GolonganDarah,
		RiwayatAlergi: in.RiwayatAlergi,
	})
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrNotFound
	}
	return p, err
}

func parseDate(s *string) (*time.Time, error) {
	if s == nil {
		return nil, nil
	}
	t, err := time.Parse("2006-01-02", *s)
	if err != nil {
		return nil, errors.New("format harus YYYY-MM-DD")
	}
	return &t, nil
}
