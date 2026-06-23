package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/balitravelhealth/platform/gateway-go/internal/models"
	"github.com/balitravelhealth/platform/gateway-go/internal/repository"
)

type AdminService struct {
	userRepo     *repository.UserRepo
	nurseRepo    *repository.NurseRepo
	facilityRepo *repository.FacilityRepo
	locationRepo *repository.LocationRepo
	assessRepo   *repository.AssessmentRepo
}

func NewAdminService(
	userRepo *repository.UserRepo,
	nurseRepo *repository.NurseRepo,
	facilityRepo *repository.FacilityRepo,
	locationRepo *repository.LocationRepo,
	assessRepo *repository.AssessmentRepo,
) *AdminService {
	return &AdminService{
		userRepo:     userRepo,
		nurseRepo:    nurseRepo,
		facilityRepo: facilityRepo,
		locationRepo: locationRepo,
		assessRepo:   assessRepo,
	}
}

// GO-23: admin/nurse login
func (s *AdminService) AdminLogin(ctx context.Context, email, password, deviceInfo string, authSvc *AuthService) (*TokenPair, *models.User, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, nil, ErrInvalidCredentials
	}
	if err != nil {
		return nil, nil, err
	}
	if user.Provider != "email" || user.PasswordHash == nil {
		return nil, nil, ErrInvalidCredentials
	}
	if err := bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(password)); err != nil {
		return nil, nil, ErrInvalidCredentials
	}
	ok, err := s.userRepo.HasRole(ctx, user.ID, "admin", "nurse")
	if err != nil || !ok {
		return nil, nil, ErrInvalidCredentials
	}
	pair, err := authSvc.generateTokenPair(ctx, user.ID, deviceInfo)
	return pair, user, err
}

// Bootstrap: create first admin (only works if no admin exists yet)
func (s *AdminService) Bootstrap(ctx context.Context, email, password string, authSvc *AuthService) (*TokenPair, *models.User, error) {
	count, err := s.userRepo.AdminCount(ctx)
	if err != nil {
		return nil, nil, err
	}
	if count > 0 {
		return nil, nil, fmt.Errorf("admin sudah ada")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil, nil, err
	}
	user, err := s.userRepo.CreateEmailUser(ctx, email, string(hash), "admin")
	if err != nil {
		return nil, nil, err
	}
	pair, err := authSvc.generateTokenPair(ctx, user.ID, "bootstrap")
	return pair, user, err
}

// GO-26: create nurse account
type CreateNurseInput struct {
	Email        string
	Password     string
	NamaLengkap  string
	NomorLisensi string
	Sertifikasi  *string
}

func (s *AdminService) CreateNurse(ctx context.Context, in CreateNurseInput) (*models.Nurse, error) {
	if in.NamaLengkap == "" || in.NomorLisensi == "" {
		return nil, fmt.Errorf("nama_lengkap and nomor_lisensi are required")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), 12)
	if err != nil {
		return nil, err
	}
	user, err := s.userRepo.CreateEmailUser(ctx, in.Email, string(hash), "nurse")
	if errors.Is(err, repository.ErrAlreadyExists) {
		return nil, ErrAlreadyExists
	}
	if err != nil {
		return nil, err
	}
	return s.nurseRepo.Create(ctx, user.ID, in.NamaLengkap, in.NomorLisensi, in.Sertifikasi)
}

func (s *AdminService) ListAllNurses(ctx context.Context) ([]models.Nurse, error) {
	list, err := s.nurseRepo.ListAll(ctx)
	if list == nil {
		list = []models.Nurse{}
	}
	return list, err
}

func (s *AdminService) ToggleNurse(ctx context.Context, nurseID int64) (*models.Nurse, error) {
	n, err := s.nurseRepo.ToggleAktif(ctx, nurseID)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrNotFound
	}
	return n, err
}

// GO-27: list all assessments
func (s *AdminService) ListAllAssessments(ctx context.Context, page, limit int) (*models.AssessmentPage, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	data, total, err := s.assessRepo.ListAll(ctx, limit, (page-1)*limit)
	if err != nil {
		return nil, err
	}
	if data == nil {
		data = []models.HealthAssessment{}
	}
	return &models.AssessmentPage{Data: data, Total: total, Page: page, Limit: limit}, nil
}

// GO-24: admin facility CRUD
type FacilityInput struct {
	DestinationID  int64
	Nama           string
	Kategori       *string
	Alamat         *string
	Latitude       *float64
	Longitude      *float64
	Kontak         *string
	JamOperasional *string
}

func (s *AdminService) CreateFacility(ctx context.Context, in FacilityInput) (*models.MedicalFacility, error) {
	return s.facilityRepo.Create(ctx, repository.FacilityInput(in))
}

func (s *AdminService) UpdateFacility(ctx context.Context, id int64, in FacilityInput) (*models.MedicalFacility, error) {
	f, err := s.facilityRepo.Update(ctx, id, repository.FacilityInput(in))
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrNotFound
	}
	return f, err
}

func (s *AdminService) DeleteFacility(ctx context.Context, id int64) error {
	if err := s.facilityRepo.Delete(ctx, id); errors.Is(err, repository.ErrNotFound) {
		return ErrNotFound
	} else {
		return err
	}
}

func (s *AdminService) ListFacilities(ctx context.Context) ([]models.MedicalFacility, error) {
	return s.facilityRepo.ListAll(ctx)
}

// GO-25: destination CRUD
func (s *AdminService) CreateDestination(ctx context.Context, namaDaerah string) (*models.Destination, error) {
	return s.locationRepo.CreateDestination(ctx, namaDaerah)
}

func (s *AdminService) DeleteDestination(ctx context.Context, id int64) error {
	if err := s.locationRepo.DeleteDestination(ctx, id); errors.Is(err, repository.ErrNotFound) {
		return ErrNotFound
	} else {
		return err
	}
}

func (s *AdminService) UpdateDestination(ctx context.Context, id int64, namaDaerah string) (*models.Destination, error) {
	d, err := s.locationRepo.UpdateDestination(ctx, id, namaDaerah)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrNotFound
	}
	return d, err
}

// GO-25: health risk CRUD
type HealthRiskInput struct {
	DestinationID          int64
	NamaRisikoID           string
	NamaRisikoEN           string
	SaranPencegahanID      *string
	SaranPencegahanEN      *string
	RekomendasiVaksinasiID *string
	RekomendasiVaksinasiEN *string
}

func (s *AdminService) CreateHealthRisk(ctx context.Context, in HealthRiskInput) (*models.HealthRisk, error) {
	return s.locationRepo.CreateHealthRisk(ctx, repository.HealthRiskInput(in))
}

func (s *AdminService) UpdateHealthRisk(ctx context.Context, id int64, in HealthRiskInput) (*models.HealthRisk, error) {
	h, err := s.locationRepo.UpdateHealthRisk(ctx, id, repository.HealthRiskInput(in))
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrNotFound
	}
	return h, err
}

func (s *AdminService) DeleteHealthRisk(ctx context.Context, id int64) error {
	if err := s.locationRepo.DeleteHealthRisk(ctx, id); errors.Is(err, repository.ErrNotFound) {
		return ErrNotFound
	} else {
		return err
	}
}

// GO-25: emergency guide CRUD
type EmergencyGuideInput struct {
	Kategori   string
	Langkah    int
	IsiMediaID json.RawMessage
	IsiMediaEN json.RawMessage
}

func (s *AdminService) CreateEmergencyGuide(ctx context.Context, in EmergencyGuideInput) (*models.EmergencyGuide, error) {
	g, err := s.locationRepo.CreateEmergencyGuide(ctx, repository.EmergencyGuideInput(in))
	if errors.Is(err, repository.ErrAlreadyExists) {
		return nil, ErrAlreadyExists
	}
	return g, err
}

func (s *AdminService) UpdateEmergencyGuide(ctx context.Context, id int64, in EmergencyGuideInput) (*models.EmergencyGuide, error) {
	g, err := s.locationRepo.UpdateEmergencyGuide(ctx, id, repository.EmergencyGuideInput(in))
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrNotFound
	}
	return g, err
}

func (s *AdminService) DeleteEmergencyGuide(ctx context.Context, id int64) error {
	if err := s.locationRepo.DeleteEmergencyGuide(ctx, id); errors.Is(err, repository.ErrNotFound) {
		return ErrNotFound
	} else {
		return err
	}
}
