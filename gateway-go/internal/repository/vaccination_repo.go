package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/balitravelhealth/platform/gateway-go/internal/models"
)

type VaccinationRepo struct {
	db *pgxpool.Pool
}

func NewVaccinationRepo(db *pgxpool.Pool) *VaccinationRepo {
	return &VaccinationRepo{db: db}
}

func (r *VaccinationRepo) ListByUserID(ctx context.Context, userID int64) ([]models.VaccinationRecord, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, user_id, jenis_vaksin, tanggal, dosis, catatan, created_at
		 FROM vaccination_records
		 WHERE user_id = $1
		 ORDER BY tanggal DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.VaccinationRecord
	for rows.Next() {
		var v models.VaccinationRecord
		if err := rows.Scan(&v.ID, &v.UserID, &v.JenisVaksin, &v.Tanggal, &v.Dosis, &v.Catatan, &v.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, v)
	}
	return list, rows.Err()
}

type VaccinationInput struct {
	UserID      int64
	JenisVaksin string
	Tanggal     time.Time
	Dosis       *string
	Catatan     *string
}

func (r *VaccinationRepo) Create(ctx context.Context, in VaccinationInput) (*models.VaccinationRecord, error) {
	v := &models.VaccinationRecord{}
	err := r.db.QueryRow(ctx,
		`INSERT INTO vaccination_records (user_id, jenis_vaksin, tanggal, dosis, catatan)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, user_id, jenis_vaksin, tanggal, dosis, catatan, created_at`,
		in.UserID, in.JenisVaksin, in.Tanggal, in.Dosis, in.Catatan,
	).Scan(&v.ID, &v.UserID, &v.JenisVaksin, &v.Tanggal, &v.Dosis, &v.Catatan, &v.CreatedAt)
	return v, err
}

func (r *VaccinationRepo) Delete(ctx context.Context, id, userID int64) error {
	tag, err := r.db.Exec(ctx,
		`DELETE FROM vaccination_records WHERE id = $1 AND user_id = $2`,
		id, userID,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// FindByID returns ErrNotFound if the record doesn't exist or belongs to a different user.
func (r *VaccinationRepo) FindByID(ctx context.Context, id, userID int64) (*models.VaccinationRecord, error) {
	v := &models.VaccinationRecord{}
	err := r.db.QueryRow(ctx,
		`SELECT id, user_id, jenis_vaksin, tanggal, dosis, catatan, created_at
		 FROM vaccination_records WHERE id = $1 AND user_id = $2`,
		id, userID,
	).Scan(&v.ID, &v.UserID, &v.JenisVaksin, &v.Tanggal, &v.Dosis, &v.Catatan, &v.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return v, err
}
