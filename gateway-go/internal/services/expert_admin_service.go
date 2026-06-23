package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/balitravelhealth/platform/gateway-go/internal/models"
	"github.com/balitravelhealth/platform/gateway-go/internal/repository"
)

type ExpertAdminService struct {
	repo *repository.ExpertRepo
}

func NewExpertAdminService(repo *repository.ExpertRepo) *ExpertAdminService {
	return &ExpertAdminService{repo: repo}
}

// ── Symptoms ─────────────────────────────────────────────────────────────────

func (s *ExpertAdminService) ListSymptoms(ctx context.Context) ([]models.ExpertSymptom, error) {
	list, err := s.repo.ListSymptoms(ctx)
	if list == nil {
		list = []models.ExpertSymptom{}
	}
	return list, err
}

func (s *ExpertAdminService) ListPublicSymptoms(ctx context.Context, kategori string) ([]models.ExpertSymptom, error) {
	if kategori == "" {
		return s.ListSymptoms(ctx)
	}
	if kategori != "pre_travel" && kategori != "post_travel" {
		return nil, fmt.Errorf("kategori must be pre_travel or post_travel")
	}
	list, err := s.repo.ListPublishedSymptomsByKategori(ctx, kategori)
	if list == nil {
		list = []models.ExpertSymptom{}
	}
	return list, err
}

func (s *ExpertAdminService) CreateSymptom(ctx context.Context, kode, labelID, labelEN string, deskripsiID, deskripsiEN *string) (*models.ExpertSymptom, error) {
	if kode == "" || labelID == "" || labelEN == "" {
		return nil, fmt.Errorf("kode, label_id, and label_en are required")
	}
	sym, err := s.repo.CreateSymptom(ctx, kode, labelID, labelEN, deskripsiID, deskripsiEN)
	if errors.Is(err, repository.ErrAlreadyExists) {
		return nil, ErrAlreadyExists
	}
	return sym, err
}

func (s *ExpertAdminService) UpdateSymptom(ctx context.Context, id int64, kode, labelID, labelEN string, deskripsiID, deskripsiEN *string) (*models.ExpertSymptom, error) {
	if kode == "" || labelID == "" || labelEN == "" {
		return nil, fmt.Errorf("kode, label_id, and label_en are required")
	}
	sym, err := s.repo.UpdateSymptom(ctx, id, kode, labelID, labelEN, deskripsiID, deskripsiEN)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrNotFound
	}
	if errors.Is(err, repository.ErrAlreadyExists) {
		return nil, ErrAlreadyExists
	}
	return sym, err
}

func (s *ExpertAdminService) DeleteSymptom(ctx context.Context, id int64) error {
	if err := s.repo.DeleteSymptom(ctx, id); errors.Is(err, repository.ErrNotFound) {
		return ErrNotFound
	} else {
		return err
	}
}

// ── Diseases ─────────────────────────────────────────────────────────────────

func (s *ExpertAdminService) ListDiseases(ctx context.Context) ([]models.ExpertDisease, error) {
	list, err := s.repo.ListDiseases(ctx)
	if list == nil {
		list = []models.ExpertDisease{}
	}
	return list, err
}

func (s *ExpertAdminService) CreateDisease(ctx context.Context, namaID, namaEN string, deskripsiID, deskripsiEN *string, rekDefaultID, rekDefaultEN json.RawMessage) (*models.ExpertDisease, error) {
	if namaID == "" || namaEN == "" {
		return nil, fmt.Errorf("nama_id and nama_en are required")
	}
	d, err := s.repo.CreateDisease(ctx, namaID, namaEN, deskripsiID, deskripsiEN, rekDefaultID, rekDefaultEN)
	if errors.Is(err, repository.ErrAlreadyExists) {
		return nil, ErrAlreadyExists
	}
	return d, err
}

func (s *ExpertAdminService) UpdateDisease(ctx context.Context, id int64, namaID, namaEN string, deskripsiID, deskripsiEN *string, rekDefaultID, rekDefaultEN json.RawMessage) (*models.ExpertDisease, error) {
	if namaID == "" || namaEN == "" {
		return nil, fmt.Errorf("nama_id and nama_en are required")
	}
	d, err := s.repo.UpdateDisease(ctx, id, namaID, namaEN, deskripsiID, deskripsiEN, rekDefaultID, rekDefaultEN)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrNotFound
	}
	return d, err
}

func (s *ExpertAdminService) DeleteDisease(ctx context.Context, id int64) error {
	if err := s.repo.DeleteDisease(ctx, id); errors.Is(err, repository.ErrNotFound) {
		return ErrNotFound
	} else {
		return err
	}
}

// ── Rules (GO-28 — validasi ketat) ───────────────────────────────────────────

type RuleInput struct {
	Nama      string
	Premis    []int64 // array symptom_id
	DiseaseID int64
	BobotCF   float64
	MB        float64
	MD        float64
	Kategori  string
	CreatedBy int64
}

var ErrRuleValidation = errors.New("rule validation failed")

func (s *ExpertAdminService) validateRule(ctx context.Context, in RuleInput) error {
	if in.Nama == "" {
		return fmt.Errorf("%w: nama is required", ErrRuleValidation)
	}
	if len(in.Premis) == 0 {
		return fmt.Errorf("%w: premis must contain at least one symptom", ErrRuleValidation)
	}
	if in.BobotCF < 0 || in.BobotCF > 1 {
		return fmt.Errorf("%w: bobot_cf must be between 0 and 1", ErrRuleValidation)
	}
	if in.MB < 0 || in.MB > 1 {
		return fmt.Errorf("%w: mb must be between 0 and 1", ErrRuleValidation)
	}
	if in.MD < 0 || in.MD > 1 {
		return fmt.Errorf("%w: md must be between 0 and 1", ErrRuleValidation)
	}
	if in.Kategori != "pre_travel" && in.Kategori != "post_travel" {
		return fmt.Errorf("%w: kategori must be pre_travel or post_travel", ErrRuleValidation)
	}
	// Verify all symptom IDs exist
	ok, err := s.repo.SymptomIDsExist(ctx, in.Premis)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("%w: one or more symptom_ids in premis do not exist", ErrRuleValidation)
	}
	// Verify disease exists
	ok, err = s.repo.DiseaseExists(ctx, in.DiseaseID)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("%w: disease_id %d does not exist", ErrRuleValidation, in.DiseaseID)
	}
	return nil
}

func (s *ExpertAdminService) ListRules(ctx context.Context) ([]models.ExpertRule, error) {
	list, err := s.repo.ListRules(ctx)
	if list == nil {
		list = []models.ExpertRule{}
	}
	return list, err
}

func (s *ExpertAdminService) CreateRule(ctx context.Context, in RuleInput) (*models.ExpertRule, error) {
	if err := s.validateRule(ctx, in); err != nil {
		return nil, err
	}
	premisJSON, _ := json.Marshal(in.Premis)
	return s.repo.CreateRule(ctx, repository.RuleInput{
		Nama: in.Nama, Premis: premisJSON, DiseaseID: in.DiseaseID,
		BobotCF: in.BobotCF, MB: in.MB, MD: in.MD,
		Kategori: in.Kategori, CreatedBy: in.CreatedBy,
	})
}

func (s *ExpertAdminService) UpdateRule(ctx context.Context, ruleID int64, in RuleInput) (*models.ExpertRule, error) {
	if err := s.validateRule(ctx, in); err != nil {
		return nil, err
	}
	premisJSON, _ := json.Marshal(in.Premis)
	ru, err := s.repo.UpdateRule(ctx, ruleID, repository.RuleInput{
		Nama: in.Nama, Premis: premisJSON, DiseaseID: in.DiseaseID,
		BobotCF: in.BobotCF, MB: in.MB, MD: in.MD,
		Kategori: in.Kategori, CreatedBy: in.CreatedBy,
	})
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrNotFound
	}
	return ru, err
}

func (s *ExpertAdminService) PublishRule(ctx context.Context, ruleID int64) (*models.ExpertRule, error) {
	ru, err := s.repo.SetRuleStatus(ctx, ruleID, "published")
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrNotFound
	}
	return ru, err
}

func (s *ExpertAdminService) UnpublishRule(ctx context.Context, ruleID int64) (*models.ExpertRule, error) {
	ru, err := s.repo.SetRuleStatus(ctx, ruleID, "draft")
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrNotFound
	}
	return ru, err
}

func (s *ExpertAdminService) DeleteRule(ctx context.Context, ruleID int64) error {
	if err := s.repo.DeleteRule(ctx, ruleID); errors.Is(err, repository.ErrNotFound) {
		return ErrNotFound
	} else {
		return err
	}
}
