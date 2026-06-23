package models

import "time"

type MedicalFacility struct {
	ID             int64     `json:"id"`
	DestinationID  int64     `json:"destination_id"`
	Nama           string    `json:"nama"`
	Kategori       *string   `json:"kategori,omitempty"`
	Alamat         *string   `json:"alamat,omitempty"`
	Latitude       *float64  `json:"latitude,omitempty"`
	Longitude      *float64  `json:"longitude,omitempty"`
	Kontak         *string   `json:"kontak,omitempty"`
	JamOperasional *string   `json:"jam_operasional,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type FacilityNearby struct {
	MedicalFacility
	DistanceKm float64 `json:"distance_km"`
}
