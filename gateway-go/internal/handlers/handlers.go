package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/balitravelhealth/platform/gateway-go/internal/repository"
	"github.com/balitravelhealth/platform/gateway-go/internal/services"
)

type Handler struct {
	DB                 *pgxpool.Pool
	AuthService        *services.AuthService
	ProfileService     *services.ProfileService
	TravelerService    *services.TravelerService
	AssessmentService  *services.AssessmentService
	VaccinationService *services.VaccinationService
	FacilityService    *services.FacilityService
	LocationService    *services.LocationService
	NursingService     *services.NursingService
	AdminService       *services.AdminService
	ExpertAdminService *services.ExpertAdminService
}

func New(db *pgxpool.Pool) *Handler {
	travelerRepo := repository.NewTravelerRepo(db)
	userRepo := repository.NewUserRepo(db)
	nurseRepo := repository.NewNurseRepo(db)
	facilityRepo := repository.NewFacilityRepo(db)
	locationRepo := repository.NewLocationRepo(db)
	assessRepo := repository.NewAssessmentRepo(db)
	expertRepo := repository.NewExpertRepo(db)
	return &Handler{
		DB:                 db,
		AuthService:        services.NewAuthService(userRepo, repository.NewTokenRepo(db)),
		ProfileService:     services.NewProfileService(repository.NewHealthProfileRepo(db)),
		TravelerService:    services.NewTravelerService(travelerRepo),
		AssessmentService:  services.NewAssessmentService(assessRepo, services.NewExpertService()),
		VaccinationService: services.NewVaccinationService(repository.NewVaccinationRepo(db)),
		FacilityService:    services.NewFacilityService(facilityRepo),
		LocationService:    services.NewLocationService(locationRepo),
		NursingService: services.NewNursingService(
			nurseRepo,
			repository.NewNursingRecordRepo(db),
			travelerRepo,
		),
		AdminService: services.NewAdminService(userRepo, nurseRepo, facilityRepo, locationRepo, assessRepo),
		ExpertAdminService: services.NewExpertAdminService(expertRepo),
	}
}

func (h *Handler) Health(c *gin.Context) {
	if err := h.DB.Ping(c.Request.Context()); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "error", "detail": "database unreachable"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
