package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/balitravelhealth/platform/gateway-go/internal/models"
)

type NursingRecordRepo struct {
	db *pgxpool.Pool
}

func NewNursingRecordRepo(db *pgxpool.Pool) *NursingRecordRepo {
	return &NursingRecordRepo{db: db}
}

const ncrCols = `id, traveler_id, nurse_id, tanggal_kunjungan,
	nursing_assessment, nursing_diagnosis, nursing_planning,
	nursing_implementation, nursing_evaluation, created_at, updated_at`

func scanNCR(row pgx.Row) (*models.NursingCareRecord, error) {
	r := &models.NursingCareRecord{}
	err := row.Scan(
		&r.ID, &r.TravelerID, &r.NurseID, &r.TanggalKunjungan,
		&r.NursingAssessment, &r.NursingDiagnosis, &r.NursingPlanning,
		&r.NursingImplementation, &r.NursingEvaluation,
		&r.CreatedAt, &r.UpdatedAt,
	)
	return r, err
}

func (r *NursingRecordRepo) Create(ctx context.Context, travelerID, nurseID int64, tanggal time.Time) (*models.NursingCareRecord, error) {
	rec, err := scanNCR(r.db.QueryRow(ctx,
		`INSERT INTO nursing_care_records (traveler_id, nurse_id, tanggal_kunjungan)
		 VALUES ($1, $2, $3)
		 RETURNING `+ncrCols,
		travelerID, nurseID, tanggal,
	))
	return rec, err
}

func (r *NursingRecordRepo) FindByID(ctx context.Context, id int64) (*models.NursingCareRecord, error) {
	rec, err := scanNCR(r.db.QueryRow(ctx,
		`SELECT `+ncrCols+` FROM nursing_care_records WHERE id = $1`, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return rec, err
}

func (r *NursingRecordRepo) ListByTraveler(ctx context.Context, travelerID int64) ([]models.NursingCareRecord, error) {
	return r.listBy(ctx, `traveler_id = $1`, travelerID)
}

func (r *NursingRecordRepo) ListByNurse(ctx context.Context, nurseID int64) ([]models.NursingCareRecord, error) {
	return r.listBy(ctx, `nurse_id = $1`, nurseID)
}

func (r *NursingRecordRepo) listBy(ctx context.Context, where string, arg int64) ([]models.NursingCareRecord, error) {
	rows, err := r.db.Query(ctx,
		`SELECT `+ncrCols+` FROM nursing_care_records WHERE `+where+` ORDER BY tanggal_kunjungan DESC`,
		arg,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.NursingCareRecord
	for rows.Next() {
		rec := models.NursingCareRecord{}
		if err := rows.Scan(
			&rec.ID, &rec.TravelerID, &rec.NurseID, &rec.TanggalKunjungan,
			&rec.NursingAssessment, &rec.NursingDiagnosis, &rec.NursingPlanning,
			&rec.NursingImplementation, &rec.NursingEvaluation,
			&rec.CreatedAt, &rec.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, rec)
	}
	return list, rows.Err()
}

type NursingRecordUpdate struct {
	Assessment     *string
	Diagnosis      *string
	Planning       *string
	Implementation *string
	Evaluation     *string
}

// UpdateCareRecord updates the care fields. Only the assigned nurse can call this
// (nurse_id check is enforced here to prevent unauthorized updates).
func (r *NursingRecordRepo) UpdateCareRecord(ctx context.Context, id, nurseID int64, u NursingRecordUpdate) (*models.NursingCareRecord, error) {
	rec, err := scanNCR(r.db.QueryRow(ctx,
		`UPDATE nursing_care_records SET
		    nursing_assessment     = $3,
		    nursing_diagnosis      = $4,
		    nursing_planning       = $5,
		    nursing_implementation = $6,
		    nursing_evaluation     = $7,
		    updated_at             = NOW()
		 WHERE id = $1 AND nurse_id = $2
		 RETURNING `+ncrCols,
		id, nurseID,
		u.Assessment, u.Diagnosis, u.Planning, u.Implementation, u.Evaluation,
	))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return rec, err
}
