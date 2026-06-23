package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/balitravelhealth/platform/gateway-go/internal/database"
	"github.com/balitravelhealth/platform/gateway-go/internal/handlers"
	"github.com/balitravelhealth/platform/gateway-go/internal/middleware"
)

func main() {
	ctx := context.Background()

	pool, err := database.NewPool(ctx)
	if err != nil {
		log.Fatalf("cannot connect to database: %v", err)
	}
	defer pool.Close()

	gin.SetMode(getenv("GIN_MODE", "debug"))

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())

	h := handlers.New(pool)

	r.GET("/health", h.Health)
	r.Static("/uploads", "./uploads")

	auth := r.Group("/auth")
	{
		auth.POST("/google", h.GoogleLogin)
		auth.POST("/refresh", h.Refresh)
		auth.POST("/logout", h.Logout)
	}

	// Public endpoints (no auth)
	r.GET("/location/classify", h.ClassifyLocation) // F-10, F-11
	r.GET("/facilities/nearby", h.NearbyFacilities)
	r.GET("/destinations", h.ListDestinations)
	r.GET("/destinations/:id/health-risks", h.ListHealthRisks)
	r.GET("/emergency-guides", h.ListEmergencyGuides)
	r.GET("/emergency-guide-flows", h.ListEmergencyFlows)
	r.GET("/emergency-guide-flows/:id", h.GetEmergencyFlow)
	r.GET("/expert/symptoms", h.ListExpertSymptoms)

	protected := r.Group("/")
	protected.Use(middleware.Auth())
	{
		// Health profile (GO-12)
		protected.GET("/health-profile", h.GetHealthProfile)
		protected.POST("/health-profile", h.CreateHealthProfile)
		protected.PUT("/health-profile", h.UpdateHealthProfile)

		// Traveler profile (GO-13)
		protected.GET("/traveler-profile", h.GetTravelerProfile)
		protected.POST("/traveler-profile", h.CreateTravelerProfile)
		protected.PUT("/traveler-profile", h.UpdateTravelerProfile)

		// Assessment — submit (GO-20/21/22) + history (GO-14)
		protected.POST("/assessment", h.PostAssessment)
		protected.GET("/assessments", h.ListAssessments)

		// Vaccination records (GO-14b)
		protected.GET("/vaccinations", h.ListVaccinations)
		protected.POST("/vaccinations", h.CreateVaccination)
		protected.DELETE("/vaccinations/:id", h.DeleteVaccination)

		// Nursing (GO-17, GO-18)
		protected.GET("/nurses", h.ListNurses)
		protected.POST("/nursing/appointments", h.CreateAppointment)
		protected.GET("/nursing/my-records", h.ListMyNursingRecords)

		// Nurse-side endpoints (GO-19) — RBAC enforced inside service via nurse profile lookup
		protected.GET("/nursing/nurse-records", h.ListNurseRecords)
		protected.PUT("/nursing/records/:id", h.UpdateCareRecord)
	}

	// Admin auth (public — before RBAC group)
	r.POST("/admin/auth/login", h.AdminLogin)
	r.POST("/admin/bootstrap", h.AdminBootstrap)

	admin := r.Group("/admin")
	admin.Use(middleware.Auth())
	admin.Use(middleware.RequireRole(pool, "admin", "nurse"))
	admin.POST("/upload", h.UploadImage)
	{
		// GO-24: facilities
		admin.GET("/facilities", h.AdminListFacilities)
		admin.POST("/facilities", h.AdminCreateFacility)
		admin.PUT("/facilities/:id", h.AdminUpdateFacility)
		admin.DELETE("/facilities/:id", h.AdminDeleteFacility)

		// GO-25: destinations
		admin.POST("/destinations", h.AdminCreateDestination)
		admin.PUT("/destinations/:id", h.AdminUpdateDestination)
		admin.DELETE("/destinations/:id", h.AdminDeleteDestination)

		// GO-25: health risks
		admin.POST("/health-risks", h.AdminCreateHealthRisk)
		admin.PUT("/health-risks/:id", h.AdminUpdateHealthRisk)
		admin.DELETE("/health-risks/:id", h.AdminDeleteHealthRisk)

		// GO-25: emergency guides (sequential, legacy)
		admin.POST("/emergency-guides", h.AdminCreateEmergencyGuide)
		admin.PUT("/emergency-guides/:id", h.AdminUpdateEmergencyGuide)
		admin.DELETE("/emergency-guides/:id", h.AdminDeleteEmergencyGuide)

		// Emergency guide flows (decision-tree, new)
		admin.GET("/emergency-guide-flows", h.AdminListEmergencyFlows)
		admin.POST("/emergency-guide-flows", h.AdminCreateEmergencyFlow)
		admin.PUT("/emergency-guide-flows/:id", h.AdminUpdateEmergencyFlow)
		admin.DELETE("/emergency-guide-flows/:id", h.AdminDeleteEmergencyFlow)

		// GO-26: nurse management (admin only — nurses can't create other nurses)
		admin.GET("/nurses", h.AdminListNurses)
		admin.POST("/nurses", h.AdminCreateNurse)
		admin.PUT("/nurses/:id/toggle", h.AdminToggleNurse)

		// GO-27: view all assessments
		admin.GET("/assessments", h.AdminListAssessments)

		// GO-28: expert knowledge base
		admin.GET("/expert/symptoms", h.AdminListSymptoms)
		admin.POST("/expert/symptoms", h.AdminCreateSymptom)
		admin.PUT("/expert/symptoms/:id", h.AdminUpdateSymptom)
		admin.DELETE("/expert/symptoms/:id", h.AdminDeleteSymptom)

		admin.GET("/expert/diseases", h.AdminListDiseases)
		admin.POST("/expert/diseases", h.AdminCreateDisease)
		admin.PUT("/expert/diseases/:id", h.AdminUpdateDisease)
		admin.DELETE("/expert/diseases/:id", h.AdminDeleteDisease)

		admin.GET("/expert/rules", h.AdminListRules)
		admin.POST("/expert/rules", h.AdminCreateRule)
		admin.PUT("/expert/rules/:id", h.AdminUpdateRule)
		admin.DELETE("/expert/rules/:id", h.AdminDeleteRule)
		admin.POST("/expert/rules/:id/publish", h.AdminPublishRule)
		admin.POST("/expert/rules/:id/unpublish", h.AdminUnpublishRule)
	}

	if err := r.Run(":" + getenv("PORT", "8080")); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
