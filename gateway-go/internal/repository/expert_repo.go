package repository

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/balitravelhealth/platform/gateway-go/internal/models"
)

type ExpertRepo struct {
	db *pgxpool.Pool
}

func NewExpertRepo(db *pgxpool.Pool) *ExpertRepo {
	return &ExpertRepo{db: db}
}

// ── Symptoms ─────────────────────────────────────────────────────────────────

func (r *ExpertRepo) ListSymptoms(ctx context.Context) ([]models.ExpertSymptom, error) {
	rows, err := r.db.Query(ctx,
		`SELECT symptom_id, kode, label_id, label_en, deskripsi_id, deskripsi_en, created_at, updated_at
		 FROM expert_symptoms ORDER BY kode`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.ExpertSymptom
	for rows.Next() {
		var s models.ExpertSymptom
		if err := rows.Scan(&s.SymptomID, &s.Kode, &s.LabelID, &s.LabelEN, &s.DeskripsiID, &s.DeskripsiEN, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, s)
	}
	return list, rows.Err()
}

func (r *ExpertRepo) ListPublishedSymptomsByKategori(ctx context.Context, kategori string) ([]models.ExpertSymptom, error) {
	rows, err := r.db.Query(ctx,
		`SELECT DISTINCT s.symptom_id, s.kode, s.label_id, s.label_en, s.deskripsi_id, s.deskripsi_en, s.created_at, s.updated_at
		 FROM expert_symptoms s
		 JOIN expert_rules r ON r.status = 'published'
		  AND EXISTS (
		      SELECT 1
		      FROM jsonb_array_elements_text(r.premis) AS premis(symptom_id)
		      WHERE premis.symptom_id::BIGINT = s.symptom_id
		  )
		 WHERE r.kategori = $1
		 ORDER BY s.kode`,
		kategori,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.ExpertSymptom
	for rows.Next() {
		var s models.ExpertSymptom
		if err := rows.Scan(&s.SymptomID, &s.Kode, &s.LabelID, &s.LabelEN, &s.DeskripsiID, &s.DeskripsiEN, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, s)
	}
	return list, rows.Err()
}

func (r *ExpertRepo) SymptomIDsExist(ctx context.Context, ids []int64) (bool, error) {
	var count int
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM expert_symptoms WHERE symptom_id = ANY($1)`, ids,
	).Scan(&count)
	return count == len(ids), err
}

func (r *ExpertRepo) CreateSymptom(ctx context.Context, kode, labelID, labelEN string, deskripsiID, deskripsiEN *string) (*models.ExpertSymptom, error) {
	s := &models.ExpertSymptom{}
	err := r.db.QueryRow(ctx,
		`INSERT INTO expert_symptoms (kode, label_id, label_en, deskripsi_id, deskripsi_en)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING symptom_id, kode, label_id, label_en, deskripsi_id, deskripsi_en, created_at, updated_at`,
		kode, labelID, labelEN, deskripsiID, deskripsiEN,
	).Scan(&s.SymptomID, &s.Kode, &s.LabelID, &s.LabelEN, &s.DeskripsiID, &s.DeskripsiEN, &s.CreatedAt, &s.UpdatedAt)
	if err != nil && isUniqueViolation(err) {
		return nil, ErrAlreadyExists
	}
	return s, err
}

func (r *ExpertRepo) UpdateSymptom(ctx context.Context, id int64, kode, labelID, labelEN string, deskripsiID, deskripsiEN *string) (*models.ExpertSymptom, error) {
	s := &models.ExpertSymptom{}
	err := r.db.QueryRow(ctx,
		`UPDATE expert_symptoms
		 SET kode=$1, label_id=$2, label_en=$3, deskripsi_id=$4, deskripsi_en=$5, updated_at=NOW()
		 WHERE symptom_id=$6
		 RETURNING symptom_id, kode, label_id, label_en, deskripsi_id, deskripsi_en, created_at, updated_at`,
		kode, labelID, labelEN, deskripsiID, deskripsiEN, id,
	).Scan(&s.SymptomID, &s.Kode, &s.LabelID, &s.LabelEN, &s.DeskripsiID, &s.DeskripsiEN, &s.CreatedAt, &s.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil && isUniqueViolation(err) {
		return nil, ErrAlreadyExists
	}
	return s, err
}

func (r *ExpertRepo) DeleteSymptom(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM expert_symptoms WHERE symptom_id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// ── Diseases ─────────────────────────────────────────────────────────────────

func (r *ExpertRepo) ListDiseases(ctx context.Context) ([]models.ExpertDisease, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, nama_id, nama_en, deskripsi_id, deskripsi_en, rekomendasi_default_id, rekomendasi_default_en, created_at, updated_at
		 FROM expert_diseases ORDER BY nama_id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.ExpertDisease
	for rows.Next() {
		var d models.ExpertDisease
		if err := rows.Scan(&d.ID, &d.NamaID, &d.NamaEN, &d.DeskripsiID, &d.DeskripsiEN,
			&d.RekomendasiDefaultID, &d.RekomendasiDefaultEN, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, d)
	}
	return list, rows.Err()
}

func (r *ExpertRepo) DiseaseExists(ctx context.Context, id int64) (bool, error) {
	var count int
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM expert_diseases WHERE id = $1`, id).Scan(&count)
	return count > 0, err
}

func (r *ExpertRepo) CreateDisease(ctx context.Context, namaID, namaEN string, deskripsiID, deskripsiEN *string, rekDefaultID, rekDefaultEN json.RawMessage) (*models.ExpertDisease, error) {
	d := &models.ExpertDisease{}
	err := r.db.QueryRow(ctx,
		`INSERT INTO expert_diseases (nama_id, nama_en, deskripsi_id, deskripsi_en, rekomendasi_default_id, rekomendasi_default_en)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING id, nama_id, nama_en, deskripsi_id, deskripsi_en, rekomendasi_default_id, rekomendasi_default_en, created_at, updated_at`,
		namaID, namaEN, deskripsiID, deskripsiEN, rekDefaultID, rekDefaultEN,
	).Scan(&d.ID, &d.NamaID, &d.NamaEN, &d.DeskripsiID, &d.DeskripsiEN,
		&d.RekomendasiDefaultID, &d.RekomendasiDefaultEN, &d.CreatedAt, &d.UpdatedAt)
	if err != nil && isUniqueViolation(err) {
		return nil, ErrAlreadyExists
	}
	return d, err
}

func (r *ExpertRepo) UpdateDisease(ctx context.Context, id int64, namaID, namaEN string, deskripsiID, deskripsiEN *string, rekDefaultID, rekDefaultEN json.RawMessage) (*models.ExpertDisease, error) {
	d := &models.ExpertDisease{}
	err := r.db.QueryRow(ctx,
		`UPDATE expert_diseases
		 SET nama_id=$2, nama_en=$3, deskripsi_id=$4, deskripsi_en=$5, rekomendasi_default_id=$6, rekomendasi_default_en=$7, updated_at=NOW()
		 WHERE id=$1
		 RETURNING id, nama_id, nama_en, deskripsi_id, deskripsi_en, rekomendasi_default_id, rekomendasi_default_en, created_at, updated_at`,
		id, namaID, namaEN, deskripsiID, deskripsiEN, rekDefaultID, rekDefaultEN,
	).Scan(&d.ID, &d.NamaID, &d.NamaEN, &d.DeskripsiID, &d.DeskripsiEN,
		&d.RekomendasiDefaultID, &d.RekomendasiDefaultEN, &d.CreatedAt, &d.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return d, err
}

func (r *ExpertRepo) DeleteDisease(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM expert_diseases WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// ── Rules ─────────────────────────────────────────────────────────────────────

func (r *ExpertRepo) ListRules(ctx context.Context) ([]models.ExpertRule, error) {
	rows, err := r.db.Query(ctx,
		`SELECT rule_id, nama, premis, disease_id, bobot_cf, mb, md, kategori, status, created_by, created_at, updated_at
		 FROM expert_rules ORDER BY rule_id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.ExpertRule
	for rows.Next() {
		var ru models.ExpertRule
		if err := rows.Scan(&ru.RuleID, &ru.Nama, &ru.Premis, &ru.DiseaseID, &ru.BobotCF, &ru.MB, &ru.MD,
			&ru.Kategori, &ru.Status, &ru.CreatedBy, &ru.CreatedAt, &ru.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, ru)
	}
	return list, rows.Err()
}

type RuleInput struct {
	Nama      string
	Premis    json.RawMessage
	DiseaseID int64
	BobotCF   float64
	MB        float64
	MD        float64
	Kategori  string
	CreatedBy int64
}

func (r *ExpertRepo) CreateRule(ctx context.Context, in RuleInput) (*models.ExpertRule, error) {
	ru := &models.ExpertRule{}
	err := r.db.QueryRow(ctx,
		`INSERT INTO expert_rules (nama, premis, disease_id, bobot_cf, mb, md, kategori, created_by)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		 RETURNING rule_id, nama, premis, disease_id, bobot_cf, mb, md, kategori, status, created_by, created_at, updated_at`,
		in.Nama, in.Premis, in.DiseaseID, in.BobotCF, in.MB, in.MD, in.Kategori, in.CreatedBy,
	).Scan(&ru.RuleID, &ru.Nama, &ru.Premis, &ru.DiseaseID, &ru.BobotCF, &ru.MB, &ru.MD,
		&ru.Kategori, &ru.Status, &ru.CreatedBy, &ru.CreatedAt, &ru.UpdatedAt)
	return ru, err
}

func (r *ExpertRepo) UpdateRule(ctx context.Context, ruleID int64, in RuleInput) (*models.ExpertRule, error) {
	ru := &models.ExpertRule{}
	err := r.db.QueryRow(ctx,
		`UPDATE expert_rules SET nama=$2, premis=$3, disease_id=$4, bobot_cf=$5, mb=$6, md=$7, kategori=$8, updated_at=NOW()
		 WHERE rule_id=$1
		 RETURNING rule_id, nama, premis, disease_id, bobot_cf, mb, md, kategori, status, created_by, created_at, updated_at`,
		ruleID, in.Nama, in.Premis, in.DiseaseID, in.BobotCF, in.MB, in.MD, in.Kategori,
	).Scan(&ru.RuleID, &ru.Nama, &ru.Premis, &ru.DiseaseID, &ru.BobotCF, &ru.MB, &ru.MD,
		&ru.Kategori, &ru.Status, &ru.CreatedBy, &ru.CreatedAt, &ru.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return ru, err
}

func (r *ExpertRepo) SetRuleStatus(ctx context.Context, ruleID int64, status string) (*models.ExpertRule, error) {
	ru := &models.ExpertRule{}
	err := r.db.QueryRow(ctx,
		`UPDATE expert_rules SET status=$2, updated_at=NOW() WHERE rule_id=$1
		 RETURNING rule_id, nama, premis, disease_id, bobot_cf, mb, md, kategori, status, created_by, created_at, updated_at`,
		ruleID, status,
	).Scan(&ru.RuleID, &ru.Nama, &ru.Premis, &ru.DiseaseID, &ru.BobotCF, &ru.MB, &ru.MD,
		&ru.Kategori, &ru.Status, &ru.CreatedBy, &ru.CreatedAt, &ru.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return ru, err
}

func (r *ExpertRepo) DeleteRule(ctx context.Context, ruleID int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM expert_rules WHERE rule_id = $1`, ruleID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
