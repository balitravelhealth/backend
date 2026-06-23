package repository

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/balitravelhealth/platform/gateway-go/internal/models"
)

type LocationRepo struct {
	db *pgxpool.Pool
}

func NewLocationRepo(db *pgxpool.Pool) *LocationRepo {
	return &LocationRepo{db: db}
}

func (r *LocationRepo) ListDestinations(ctx context.Context) ([]models.Destination, error) {
	rows, err := r.db.Query(ctx,
		`SELECT destination_id, nama_daerah, created_at FROM destinations ORDER BY nama_daerah`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Destination
	for rows.Next() {
		var d models.Destination
		if err := rows.Scan(&d.ID, &d.NamaDaerah, &d.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, d)
	}
	return list, rows.Err()
}

func (r *LocationRepo) ListHealthRisks(ctx context.Context, destinationID int64) ([]models.HealthRisk, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, destination_id, nama_risiko_id, nama_risiko_en, saran_pencegahan_id, saran_pencegahan_en,
		        rekomendasi_vaksinasi_id, rekomendasi_vaksinasi_en, created_at, updated_at
		 FROM health_risks WHERE destination_id = $1 ORDER BY id`,
		destinationID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.HealthRisk
	for rows.Next() {
		var h models.HealthRisk
		if err := rows.Scan(&h.ID, &h.DestinationID, &h.NamaRisikoID, &h.NamaRisikoEN,
			&h.SaranPencegahanID, &h.SaranPencegahanEN,
			&h.RekomendasiVaksinasiID, &h.RekomendasiVaksinasiEN, &h.CreatedAt, &h.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, h)
	}
	return list, rows.Err()
}

// ListEmergencyGuides returns all guides, or only for the given kategori if non-empty.
func (r *LocationRepo) ListEmergencyGuides(ctx context.Context, kategori string) ([]models.EmergencyGuide, error) {
	if kategori != "" {
		return r.listEmergencyByKategori(ctx, kategori)
	}
	rows, err := r.db.Query(ctx,
		`SELECT id, kategori, langkah, isi_media_id, isi_media_en, created_at, updated_at
		 FROM emergency_guides ORDER BY kategori, langkah`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanGuides(rows)
}

func (r *LocationRepo) listEmergencyByKategori(ctx context.Context, kategori string) ([]models.EmergencyGuide, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, kategori, langkah, isi_media_id, isi_media_en, created_at, updated_at
		 FROM emergency_guides WHERE kategori = $1 ORDER BY langkah`,
		kategori,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanGuides(rows)
}

type guideRows interface {
	Next() bool
	Scan(dest ...any) error
	Close()
	Err() error
}

// ── Admin CRUD ────────────────────────────────────────────────────────────────

func (r *LocationRepo) CreateDestination(ctx context.Context, namaDaerah string) (*models.Destination, error) {
	d := &models.Destination{}
	err := r.db.QueryRow(ctx,
		`INSERT INTO destinations (nama_daerah) VALUES ($1) RETURNING destination_id, nama_daerah, created_at`,
		namaDaerah,
	).Scan(&d.ID, &d.NamaDaerah, &d.CreatedAt)
	return d, err
}

func (r *LocationRepo) DeleteDestination(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM destinations WHERE destination_id=$1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *LocationRepo) UpdateDestination(ctx context.Context, id int64, namaDaerah string) (*models.Destination, error) {
	d := &models.Destination{}
	err := r.db.QueryRow(ctx,
		`UPDATE destinations SET nama_daerah=$1 WHERE destination_id=$2
		 RETURNING destination_id, nama_daerah, created_at`,
		namaDaerah, id,
	).Scan(&d.ID, &d.NamaDaerah, &d.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return d, err
}

type HealthRiskInput struct {
	DestinationID          int64
	NamaRisikoID           string
	NamaRisikoEN           string
	SaranPencegahanID      *string
	SaranPencegahanEN      *string
	RekomendasiVaksinasiID *string
	RekomendasiVaksinasiEN *string
}

func (r *LocationRepo) CreateHealthRisk(ctx context.Context, in HealthRiskInput) (*models.HealthRisk, error) {
	h := &models.HealthRisk{}
	err := r.db.QueryRow(ctx,
		`INSERT INTO health_risks (destination_id, nama_risiko_id, nama_risiko_en, saran_pencegahan_id, saran_pencegahan_en, rekomendasi_vaksinasi_id, rekomendasi_vaksinasi_en)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING id, destination_id, nama_risiko_id, nama_risiko_en, saran_pencegahan_id, saran_pencegahan_en, rekomendasi_vaksinasi_id, rekomendasi_vaksinasi_en, created_at, updated_at`,
		in.DestinationID, in.NamaRisikoID, in.NamaRisikoEN, in.SaranPencegahanID, in.SaranPencegahanEN, in.RekomendasiVaksinasiID, in.RekomendasiVaksinasiEN,
	).Scan(&h.ID, &h.DestinationID, &h.NamaRisikoID, &h.NamaRisikoEN, &h.SaranPencegahanID, &h.SaranPencegahanEN, &h.RekomendasiVaksinasiID, &h.RekomendasiVaksinasiEN, &h.CreatedAt, &h.UpdatedAt)
	return h, err
}

func (r *LocationRepo) UpdateHealthRisk(ctx context.Context, id int64, in HealthRiskInput) (*models.HealthRisk, error) {
	h := &models.HealthRisk{}
	err := r.db.QueryRow(ctx,
		`UPDATE health_risks
		 SET nama_risiko_id=$2, nama_risiko_en=$3, saran_pencegahan_id=$4, saran_pencegahan_en=$5, rekomendasi_vaksinasi_id=$6, rekomendasi_vaksinasi_en=$7, updated_at=NOW()
		 WHERE id=$1
		 RETURNING id, destination_id, nama_risiko_id, nama_risiko_en, saran_pencegahan_id, saran_pencegahan_en, rekomendasi_vaksinasi_id, rekomendasi_vaksinasi_en, created_at, updated_at`,
		id, in.NamaRisikoID, in.NamaRisikoEN, in.SaranPencegahanID, in.SaranPencegahanEN, in.RekomendasiVaksinasiID, in.RekomendasiVaksinasiEN,
	).Scan(&h.ID, &h.DestinationID, &h.NamaRisikoID, &h.NamaRisikoEN, &h.SaranPencegahanID, &h.SaranPencegahanEN, &h.RekomendasiVaksinasiID, &h.RekomendasiVaksinasiEN, &h.CreatedAt, &h.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
	}
	return h, err
}

func (r *LocationRepo) DeleteHealthRisk(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM health_risks WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

type EmergencyGuideInput struct {
	Kategori    string
	Langkah     int
	IsiMediaID  json.RawMessage
	IsiMediaEN  json.RawMessage
}

func (r *LocationRepo) CreateEmergencyGuide(ctx context.Context, in EmergencyGuideInput) (*models.EmergencyGuide, error) {
	g := &models.EmergencyGuide{}
	err := r.db.QueryRow(ctx,
		`INSERT INTO emergency_guides (kategori, langkah, isi_media_id, isi_media_en)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, kategori, langkah, isi_media_id, isi_media_en, created_at, updated_at`,
		in.Kategori, in.Langkah, in.IsiMediaID, in.IsiMediaEN,
	).Scan(&g.ID, &g.Kategori, &g.Langkah, &g.IsiMediaID, &g.IsiMediaEN, &g.CreatedAt, &g.UpdatedAt)
	if err != nil && isUniqueViolation(err) {
		return nil, ErrAlreadyExists
	}
	return g, err
}

func (r *LocationRepo) UpdateEmergencyGuide(ctx context.Context, id int64, in EmergencyGuideInput) (*models.EmergencyGuide, error) {
	g := &models.EmergencyGuide{}
	err := r.db.QueryRow(ctx,
		`UPDATE emergency_guides
		 SET kategori=$2, langkah=$3, isi_media_id=$4, isi_media_en=$5, updated_at=NOW()
		 WHERE id=$1
		 RETURNING id, kategori, langkah, isi_media_id, isi_media_en, created_at, updated_at`,
		id, in.Kategori, in.Langkah, in.IsiMediaID, in.IsiMediaEN,
	).Scan(&g.ID, &g.Kategori, &g.Langkah, &g.IsiMediaID, &g.IsiMediaEN, &g.CreatedAt, &g.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
	}
	return g, err
}

func (r *LocationRepo) DeleteEmergencyGuide(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM emergency_guides WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func scanGuides(rows guideRows) ([]models.EmergencyGuide, error) {
	var list []models.EmergencyGuide
	for rows.Next() {
		var g models.EmergencyGuide
		if err := rows.Scan(&g.ID, &g.Kategori, &g.Langkah, &g.IsiMediaID, &g.IsiMediaEN, &g.CreatedAt, &g.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, g)
	}
	return list, rows.Err()
}
