package services

// F-10: Bali bounding box (Bab 3.10)
const (
	baliLatMin = -8.90
	baliLatMax = -8.00
	baliLngMin = 114.40
	baliLngMax = 115.75
)

// F-11: Approximate bounding boxes per kabupaten/kota Bali
type region struct {
	Name   string
	LatMin float64
	LatMax float64
	LngMin float64
	LngMax float64
}

var baliRegions = []region{
	{"Kota Denpasar", -8.730, -8.580, 115.170, 115.270},
	{"Kabupaten Badung", -8.850, -8.620, 115.060, 115.270},
	{"Kabupaten Gianyar", -8.590, -8.430, 115.200, 115.430},
	{"Kabupaten Tabanan", -8.680, -8.320, 114.890, 115.200},
	{"Kabupaten Buleleng", -8.290, -8.030, 114.500, 115.400},
	{"Kabupaten Karangasem", -8.590, -8.240, 115.380, 115.750},
	{"Kabupaten Klungkung", -8.640, -8.430, 115.290, 115.570},
	{"Kabupaten Bangli", -8.530, -8.200, 115.200, 115.500},
	{"Kabupaten Jembrana", -8.630, -8.300, 114.400, 114.900},
}

type LocationClassification struct {
	InBali    bool    `json:"in_bali"`
	Region    *string `json:"region,omitempty"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func ClassifyLocation(lat, lng float64) LocationClassification {
	inBali := lat >= baliLatMin && lat <= baliLatMax &&
		lng >= baliLngMin && lng <= baliLngMax

	result := LocationClassification{
		InBali:    inBali,
		Latitude:  lat,
		Longitude: lng,
	}

	if !inBali {
		return result
	}

	for _, r := range baliRegions {
		if lat >= r.LatMin && lat <= r.LatMax && lng >= r.LngMin && lng <= r.LngMax {
			name := r.Name
			result.Region = &name
			return result
		}
	}
	// In Bali bounding box but no sub-region matched (overlap gaps)
	generic := "Bali"
	result.Region = &generic
	return result
}
