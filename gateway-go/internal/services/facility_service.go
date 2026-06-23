package services

import (
	"context"
	"errors"
	"math"
	"sort"

	"github.com/balitravelhealth/platform/gateway-go/internal/models"
	"github.com/balitravelhealth/platform/gateway-go/internal/repository"
)

const earthRadiusKm = 6371.0

type FacilityService struct {
	repo *repository.FacilityRepo
}

func NewFacilityService(repo *repository.FacilityRepo) *FacilityService {
	return &FacilityService{repo: repo}
}

var ErrInvalidCoordinates = errors.New("lat and lng are required")

func (s *FacilityService) Nearby(ctx context.Context, lat, lng float64, radiusKm float64, limit int) ([]models.FacilityNearby, error) {
	if limit <= 0 || limit > 50 {
		limit = 10
	}
	if radiusKm <= 0 {
		radiusKm = 20
	}

	// Bounding box pre-filter (1° lat ≈ 111km)
	deltaLat := radiusKm / 111.0
	deltaLng := radiusKm / (111.0 * math.Cos(lat*math.Pi/180))

	candidates, err := s.repo.FindInBoundingBox(ctx,
		lat-deltaLat, lat+deltaLat,
		lng-deltaLng, lng+deltaLng,
	)
	if err != nil {
		return nil, err
	}

	// Compute exact Haversine distances and filter by radius
	results := make([]models.FacilityNearby, 0, len(candidates))
	for _, f := range candidates {
		if f.Latitude == nil || f.Longitude == nil {
			continue
		}
		dist := haversineKm(lat, lng, *f.Latitude, *f.Longitude)
		if dist <= radiusKm {
			results = append(results, models.FacilityNearby{
				MedicalFacility: f,
				DistanceKm:      math.Round(dist*100) / 100,
			})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].DistanceKm < results[j].DistanceKm
	})

	if len(results) > limit {
		results = results[:limit]
	}
	return results, nil
}

func haversineKm(lat1, lng1, lat2, lng2 float64) float64 {
	dLat := (lat2 - lat1) * math.Pi / 180
	dLng := (lng2 - lng1) * math.Pi / 180
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLng/2)*math.Sin(dLng/2)
	return earthRadiusKm * 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
}
