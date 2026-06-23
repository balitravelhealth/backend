package models

import "time"

type Nurse struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"user_id"`
	NamaLengkap  string    `json:"nama_lengkap"`
	NomorLisensi string    `json:"nomor_lisensi"`
	Sertifikasi  *string   `json:"sertifikasi,omitempty"`
	Aktif        bool      `json:"aktif"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type NursingCareRecord struct {
	ID                    int64     `json:"id"`
	TravelerID            int64     `json:"traveler_id"`
	NurseID               int64     `json:"nurse_id"`
	TanggalKunjungan      time.Time `json:"tanggal_kunjungan"`
	NursingAssessment     *string   `json:"nursing_assessment,omitempty"`
	NursingDiagnosis      *string   `json:"nursing_diagnosis,omitempty"`
	NursingPlanning       *string   `json:"nursing_planning,omitempty"`
	NursingImplementation *string   `json:"nursing_implementation,omitempty"`
	NursingEvaluation     *string   `json:"nursing_evaluation,omitempty"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}
