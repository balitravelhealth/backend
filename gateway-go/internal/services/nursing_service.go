package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/balitravelhealth/platform/gateway-go/internal/models"
	"github.com/balitravelhealth/platform/gateway-go/internal/repository"
)

type NursingService struct {
	nurseRepo  *repository.NurseRepo
	recordRepo *repository.NursingRecordRepo
	travelerRepo *repository.TravelerRepo
}

func NewNursingService(
	nurseRepo *repository.NurseRepo,
	recordRepo *repository.NursingRecordRepo,
	travelerRepo *repository.TravelerRepo,
) *NursingService {
	return &NursingService{
		nurseRepo:    nurseRepo,
		recordRepo:   recordRepo,
		travelerRepo: travelerRepo,
	}
}

// GO-17: list active nurses
func (s *NursingService) ListNurses(ctx context.Context) ([]models.Nurse, error) {
	list, err := s.nurseRepo.ListActive(ctx)
	if list == nil {
		list = []models.Nurse{}
	}
	return list, err
}

// GO-18: traveler creates an appointment
func (s *NursingService) CreateAppointment(ctx context.Context, userID, nurseID int64, tanggalStr string) (*models.NursingCareRecord, error) {
	// Resolve traveler_id from user_id
	traveler, err := s.travelerRepo.FindByUserID(ctx, userID)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, fmt.Errorf("buat traveler profile terlebih dahulu")
	}
	if err != nil {
		return nil, err
	}

	// Validate nurse exists and is active
	nurse, err := s.nurseRepo.FindByID(ctx, nurseID)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, fmt.Errorf("nurse tidak ditemukan")
	}
	if err != nil {
		return nil, err
	}
	if !nurse.Aktif {
		return nil, fmt.Errorf("nurse sedang tidak tersedia")
	}

	tanggal, err := time.Parse(time.RFC3339, tanggalStr)
	if err != nil {
		// try date-only format
		tanggal, err = time.Parse("2006-01-02T15:04", tanggalStr)
		if err != nil {
			return nil, fmt.Errorf("tanggal_kunjungan: format harus RFC3339 atau YYYY-MM-DDTHH:MM")
		}
	}

	return s.recordRepo.Create(ctx, traveler.ID, nurseID, tanggal)
}

// GO-19: nurse fills in the care record
func (s *NursingService) UpdateCareRecord(ctx context.Context, nurseUserID, recordID int64, u repository.NursingRecordUpdate) (*models.NursingCareRecord, error) {
	nurse, err := s.nurseRepo.FindByUserID(ctx, nurseUserID)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, fmt.Errorf("nurse profile tidak ditemukan")
	}
	if err != nil {
		return nil, err
	}

	rec, err := s.recordRepo.UpdateCareRecord(ctx, recordID, nurse.ID, u)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrNotFound
	}
	return rec, err
}

// ListMyRecords — traveler sees their own records
func (s *NursingService) ListTravelerRecords(ctx context.Context, userID int64) ([]models.NursingCareRecord, error) {
	traveler, err := s.travelerRepo.FindByUserID(ctx, userID)
	if errors.Is(err, repository.ErrNotFound) {
		return []models.NursingCareRecord{}, nil
	}
	if err != nil {
		return nil, err
	}
	list, err := s.recordRepo.ListByTraveler(ctx, traveler.ID)
	if list == nil {
		list = []models.NursingCareRecord{}
	}
	return list, err
}

// ListNurseRecords — nurse sees records assigned to them
func (s *NursingService) ListNurseRecords(ctx context.Context, nurseUserID int64) ([]models.NursingCareRecord, error) {
	nurse, err := s.nurseRepo.FindByUserID(ctx, nurseUserID)
	if errors.Is(err, repository.ErrNotFound) {
		return []models.NursingCareRecord{}, nil
	}
	if err != nil {
		return nil, err
	}
	list, err := s.recordRepo.ListByNurse(ctx, nurse.ID)
	if list == nil {
		list = []models.NursingCareRecord{}
	}
	return list, err
}
