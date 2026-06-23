package models

import "time"

type Traveler struct {
	ID            int64      `json:"id"`
	UserID        int64      `json:"user_id"`
	NamaLengkap   string     `json:"nama_lengkap"`
	TanggalLahir  *time.Time `json:"tanggal_lahir,omitempty"`
	KontakDarurat *string    `json:"kontak_darurat,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}
