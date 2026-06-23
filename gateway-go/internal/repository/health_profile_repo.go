package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/balitravelhealth/platform/gateway-go/internal/models"
)

type HealthProfileRepo struct {
	db *pgxpool.Pool
}

func NewHealthProfileRepo(db *pgxpool.Pool) *HealthProfileRepo {
	return &HealthProfileRepo{db: db}
}

func (r *HealthProfileRepo) FindByUserID(ctx context.Context, userID int64) (*models.HealthProfile, error) {
	p := &models.HealthProfile{}
	var tgl pgtype.Date

	err := r.db.QueryRow(ctx,
		`SELECT id, user_id, tanggal_lahir, jenis_kelamin,
		        tinggi_cm, berat_kg, golongan_darah, riwayat_alergi,
		        created_at, updated_at
		 FROM health_profiles WHERE user_id = $1`,
		userID,
	).Scan(
		&p.ID, &p.UserID, &tgl, &p.JenisKelamin,
		&p.TinggiCm, &p.BeratKg, &p.GolonganDarah, &p.RiwayatAlergi,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if tgl.Valid {
		t := tgl.Time
		p.TanggalLahir = &t
	}
	return p, nil
}

type HealthProfileInput struct {
	UserID        int64
	TanggalLahir  *time.Time
	JenisKelamin  *string
	TinggiCm      *float64
	BeratKg       *float64
	GolonganDarah *string
	RiwayatAlergi *string
}

func (r *HealthProfileRepo) Create(ctx context.Context, in HealthProfileInput) (*models.HealthProfile, error) {
	p := &models.HealthProfile{}
	var tgl pgtype.Date

	err := r.db.QueryRow(ctx,
		`INSERT INTO health_profiles
		    (user_id, tanggal_lahir, jenis_kelamin, tinggi_cm, berat_kg, golongan_darah, riwayat_alergi)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING id, user_id, tanggal_lahir, jenis_kelamin,
		           tinggi_cm, berat_kg, golongan_darah, riwayat_alergi,
		           created_at, updated_at`,
		in.UserID, in.TanggalLahir, in.JenisKelamin, in.TinggiCm,
		in.BeratKg, in.GolonganDarah, in.RiwayatAlergi,
	).Scan(
		&p.ID, &p.UserID, &tgl, &p.JenisKelamin,
		&p.TinggiCm, &p.BeratKg, &p.GolonganDarah, &p.RiwayatAlergi,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if isUniqueViolation(err) {
			return nil, ErrAlreadyExists
		}
		return nil, err
	}
	if tgl.Valid {
		t := tgl.Time
		p.TanggalLahir = &t
	}
	return p, nil
}

func (r *HealthProfileRepo) Update(ctx context.Context, in HealthProfileInput) (*models.HealthProfile, error) {
	p := &models.HealthProfile{}
	var tgl pgtype.Date

	err := r.db.QueryRow(ctx,
		`UPDATE health_profiles SET
		    tanggal_lahir  = $2,
		    jenis_kelamin  = $3,
		    tinggi_cm      = $4,
		    berat_kg       = $5,
		    golongan_darah = $6,
		    riwayat_alergi = $7,
		    updated_at     = NOW()
		 WHERE user_id = $1
		 RETURNING id, user_id, tanggal_lahir, jenis_kelamin,
		           tinggi_cm, berat_kg, golongan_darah, riwayat_alergi,
		           created_at, updated_at`,
		in.UserID, in.TanggalLahir, in.JenisKelamin, in.TinggiCm,
		in.BeratKg, in.GolonganDarah, in.RiwayatAlergi,
	).Scan(
		&p.ID, &p.UserID, &tgl, &p.JenisKelamin,
		&p.TinggiCm, &p.BeratKg, &p.GolonganDarah, &p.RiwayatAlergi,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if tgl.Valid {
		t := tgl.Time
		p.TanggalLahir = &t
	}
	return p, nil
}
