package models

import "time"

type VaccinationRecord struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	JenisVaksin string    `json:"jenis_vaksin"`
	Tanggal     time.Time `json:"tanggal"`
	Dosis       *string   `json:"dosis,omitempty"`
	Catatan     *string   `json:"catatan,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}
