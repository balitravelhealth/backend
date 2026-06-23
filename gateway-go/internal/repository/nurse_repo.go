package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/balitravelhealth/platform/gateway-go/internal/models"
)

type NurseRepo struct {
	db *pgxpool.Pool
}

func NewNurseRepo(db *pgxpool.Pool) *NurseRepo {
	return &NurseRepo{db: db}
}

const nurseCols = `id, user_id, nama_lengkap, nomor_lisensi, sertifikasi, aktif, created_at, updated_at`

func scanNurse(row pgx.Row) (*models.Nurse, error) {
	n := &models.Nurse{}
	err := row.Scan(&n.ID, &n.UserID, &n.NamaLengkap, &n.NomorLisensi, &n.Sertifikasi, &n.Aktif, &n.CreatedAt, &n.UpdatedAt)
	return n, err
}

func (r *NurseRepo) ListActive(ctx context.Context) ([]models.Nurse, error) {
	rows, err := r.db.Query(ctx,
		`SELECT `+nurseCols+` FROM nurses WHERE aktif = true ORDER BY nama_lengkap`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Nurse
	for rows.Next() {
		n := models.Nurse{}
		if err := rows.Scan(&n.ID, &n.UserID, &n.NamaLengkap, &n.NomorLisensi, &n.Sertifikasi, &n.Aktif, &n.CreatedAt, &n.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, n)
	}
	return list, rows.Err()
}

func (r *NurseRepo) FindByUserID(ctx context.Context, userID int64) (*models.Nurse, error) {
	n, err := scanNurse(r.db.QueryRow(ctx,
		`SELECT `+nurseCols+` FROM nurses WHERE user_id = $1`, userID))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return n, err
}

func (r *NurseRepo) FindByID(ctx context.Context, id int64) (*models.Nurse, error) {
	n, err := scanNurse(r.db.QueryRow(ctx,
		`SELECT `+nurseCols+` FROM nurses WHERE id = $1`, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return n, err
}

func (r *NurseRepo) ListAll(ctx context.Context) ([]models.Nurse, error) {
	rows, err := r.db.Query(ctx, `SELECT `+nurseCols+` FROM nurses ORDER BY nama_lengkap`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.Nurse
	for rows.Next() {
		n := models.Nurse{}
		if err := rows.Scan(&n.ID, &n.UserID, &n.NamaLengkap, &n.NomorLisensi, &n.Sertifikasi, &n.Aktif, &n.CreatedAt, &n.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, n)
	}
	return list, rows.Err()
}

func (r *NurseRepo) Create(ctx context.Context, userID int64, namaLengkap, nomorLisensi string, sertifikasi *string) (*models.Nurse, error) {
	n, err := scanNurse(r.db.QueryRow(ctx,
		`INSERT INTO nurses (user_id, nama_lengkap, nomor_lisensi, sertifikasi)
		 VALUES ($1,$2,$3,$4)
		 RETURNING `+nurseCols,
		userID, namaLengkap, nomorLisensi, sertifikasi,
	))
	if err != nil && isUniqueViolation(err) {
		return nil, ErrAlreadyExists
	}
	return n, err
}

func (r *NurseRepo) ToggleAktif(ctx context.Context, id int64) (*models.Nurse, error) {
	n, err := scanNurse(r.db.QueryRow(ctx,
		`UPDATE nurses SET aktif = NOT aktif, updated_at = NOW()
		 WHERE id = $1 RETURNING `+nurseCols, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return n, err
}
