package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/balitravelhealth/platform/gateway-go/internal/models"
)

type FacilityRepo struct {
	db *pgxpool.Pool
}

func NewFacilityRepo(db *pgxpool.Pool) *FacilityRepo {
	return &FacilityRepo{db: db}
}

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

const facilityCols = `id, destination_id, nama, kategori, alamat, latitude, longitude, kontak, jam_operasional, created_at, updated_at`

func scanFacility(row pgx.Row) (*models.MedicalFacility, error) {
	f := &models.MedicalFacility{}
	err := row.Scan(&f.ID, &f.DestinationID, &f.Nama, &f.Kategori, &f.Alamat,
		&f.Latitude, &f.Longitude, &f.Kontak, &f.JamOperasional, &f.CreatedAt, &f.UpdatedAt)
	return f, err
}

func (r *FacilityRepo) ListAll(ctx context.Context) ([]models.MedicalFacility, error) {
	rows, err := r.db.Query(ctx, `SELECT `+facilityCols+` FROM medical_facilities ORDER BY nama`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.MedicalFacility
	for rows.Next() {
		f := models.MedicalFacility{}
		if err := rows.Scan(&f.ID, &f.DestinationID, &f.Nama, &f.Kategori, &f.Alamat,
			&f.Latitude, &f.Longitude, &f.Kontak, &f.JamOperasional, &f.CreatedAt, &f.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, f)
	}
	return list, rows.Err()
}

func (r *FacilityRepo) Create(ctx context.Context, in FacilityInput) (*models.MedicalFacility, error) {
	return scanFacility(r.db.QueryRow(ctx,
		`INSERT INTO medical_facilities (destination_id,nama,kategori,alamat,latitude,longitude,kontak,jam_operasional)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING `+facilityCols,
		in.DestinationID, in.Nama, in.Kategori, in.Alamat,
		in.Latitude, in.Longitude, in.Kontak, in.JamOperasional,
	))
}

func (r *FacilityRepo) Update(ctx context.Context, id int64, in FacilityInput) (*models.MedicalFacility, error) {
	f, err := scanFacility(r.db.QueryRow(ctx,
		`UPDATE medical_facilities SET destination_id=$2,nama=$3,kategori=$4,alamat=$5,
		 latitude=$6,longitude=$7,kontak=$8,jam_operasional=$9,updated_at=NOW()
		 WHERE id=$1 RETURNING `+facilityCols,
		id, in.DestinationID, in.Nama, in.Kategori, in.Alamat,
		in.Latitude, in.Longitude, in.Kontak, in.JamOperasional,
	))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return f, err
}

func (r *FacilityRepo) Delete(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM medical_facilities WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// FindInBoundingBox returns all facilities whose coordinates fall within the given bounding box.
// Used as a pre-filter before Haversine sort in the service layer.
func (r *FacilityRepo) FindInBoundingBox(ctx context.Context, minLat, maxLat, minLng, maxLng float64) ([]models.MedicalFacility, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, destination_id, nama, kategori, alamat,
		        latitude, longitude, kontak, jam_operasional,
		        created_at, updated_at
		 FROM medical_facilities
		 WHERE latitude  IS NOT NULL AND longitude IS NOT NULL
		   AND latitude  BETWEEN $1 AND $2
		   AND longitude BETWEEN $3 AND $4`,
		minLat, maxLat, minLng, maxLng,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.MedicalFacility
	for rows.Next() {
		var f models.MedicalFacility
		if err := rows.Scan(
			&f.ID, &f.DestinationID, &f.Nama, &f.Kategori, &f.Alamat,
			&f.Latitude, &f.Longitude, &f.Kontak, &f.JamOperasional,
			&f.CreatedAt, &f.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, f)
	}
	return list, rows.Err()
}
