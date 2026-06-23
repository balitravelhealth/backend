package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/balitravelhealth/platform/gateway-go/internal/models"
)

type TravelerRepo struct {
	db *pgxpool.Pool
}

func NewTravelerRepo(db *pgxpool.Pool) *TravelerRepo {
	return &TravelerRepo{db: db}
}

const travelerCols = `id, user_id, nama_lengkap, tanggal_lahir, kontak_darurat, created_at, updated_at`

func scanTraveler(row pgx.Row) (*models.Traveler, error) {
	t := &models.Traveler{}
	var tgl pgtype.Date
	if err := row.Scan(&t.ID, &t.UserID, &t.NamaLengkap, &tgl, &t.KontakDarurat, &t.CreatedAt, &t.UpdatedAt); err != nil {
		return nil, err
	}
	if tgl.Valid {
		v := tgl.Time
		t.TanggalLahir = &v
	}
	return t, nil
}

func (r *TravelerRepo) FindByUserID(ctx context.Context, userID int64) (*models.Traveler, error) {
	row := r.db.QueryRow(ctx,
		`SELECT `+travelerCols+` FROM travelers WHERE user_id = $1`, userID)
	t, err := scanTraveler(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return t, err
}

type TravelerInput struct {
	UserID        int64
	NamaLengkap   string
	TanggalLahir  interface{} // *time.Time or nil
	KontakDarurat *string
}

func (r *TravelerRepo) Create(ctx context.Context, in TravelerInput) (*models.Traveler, error) {
	row := r.db.QueryRow(ctx,
		`INSERT INTO travelers (user_id, nama_lengkap, tanggal_lahir, kontak_darurat)
		 VALUES ($1, $2, $3, $4)
		 RETURNING `+travelerCols,
		in.UserID, in.NamaLengkap, in.TanggalLahir, in.KontakDarurat,
	)
	t, err := scanTraveler(row)
	if err != nil {
		if isUniqueViolation(err) {
			return nil, ErrAlreadyExists
		}
		return nil, err
	}
	return t, nil
}

func (r *TravelerRepo) Update(ctx context.Context, in TravelerInput) (*models.Traveler, error) {
	row := r.db.QueryRow(ctx,
		`UPDATE travelers SET
		    nama_lengkap   = $2,
		    tanggal_lahir  = $3,
		    kontak_darurat = $4,
		    updated_at     = NOW()
		 WHERE user_id = $1
		 RETURNING `+travelerCols,
		in.UserID, in.NamaLengkap, in.TanggalLahir, in.KontakDarurat,
	)
	t, err := scanTraveler(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return t, err
}
