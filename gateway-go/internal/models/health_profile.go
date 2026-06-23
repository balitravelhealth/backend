package models

import "time"

type HealthProfile struct {
	ID            int64      `json:"id"`
	UserID        int64      `json:"user_id"`
	TanggalLahir  *time.Time `json:"tanggal_lahir,omitempty"`
	JenisKelamin  *string    `json:"jenis_kelamin,omitempty"`
	TinggiCm      *float64   `json:"tinggi_cm,omitempty"`
	BeratKg       *float64   `json:"berat_kg,omitempty"`
	GolonganDarah *string    `json:"golongan_darah,omitempty"`
	RiwayatAlergi *string    `json:"riwayat_alergi,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}
