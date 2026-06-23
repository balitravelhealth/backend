package handlers_test

// SRV-T-6: RBAC + knowledge base security tests (TEST-5, TEST-6)
//
// Test matrix:
//   - Unauthenticated → 401 on all protected endpoints
//   - Traveler JWT (role=traveler) → 403 on all /admin/* endpoints
//   - Admin JWT (role=admin) → 200/201/204 on /admin/* endpoints
//   - GO-28 rule validation → 422 on invalid CF / missing fields

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/balitravelhealth/platform/gateway-go/internal/handlers"
	"github.com/balitravelhealth/platform/gateway-go/internal/middleware"
)

const testJWTSecret = "b951354077e26851a3fa510cdbcdc2903a4376b809a73267514a447784840106"

func TestMain(m *testing.M) {
	// Ensure middleware and helpers share the same JWT secret during tests.
	os.Setenv("JWT_SECRET", testJWTSecret)
	os.Exit(m.Run())
}

// ── helpers ──────────────────────────────────────────────────────────────────

func mustPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		getenv("DB_HOST", "localhost"),
		getenv("DB_PORT", "5432"),
		getenv("POSTGRES_USER", "balitravelhealthdb"),
		getenv("POSTGRES_PASSWORD", "anjaysemhas"),
		getenv("POSTGRES_DB", "balitravelhealth"),
	)
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		t.Skipf("DB not available, skipping integration test: %v", err)
	}
	if err := pool.Ping(context.Background()); err != nil {
		t.Skipf("DB ping failed, skipping: %v", err)
	}
	return pool
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func makeJWT(t *testing.T, userID int64) string {
	t.Helper()
	secret := testJWTSecret
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   fmt.Sprintf("%d", userID),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("sign JWT: %v", err)
	}
	return signed
}

func newRouter(t *testing.T) (*gin.Engine, *pgxpool.Pool) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	pool := mustPool(t)
	h := handlers.New(pool)

	r := gin.New()
	r.GET("/health", h.Health)
	r.GET("/location/classify", h.ClassifyLocation)
	r.GET("/destinations", h.ListDestinations)
	r.GET("/emergency-guides", h.ListEmergencyGuides)
	r.GET("/facilities/nearby", h.NearbyFacilities)
	r.POST("/auth/google", h.GoogleLogin)

	r.POST("/admin/auth/login", h.AdminLogin)
	r.POST("/admin/bootstrap", h.AdminBootstrap)

	protected := r.Group("/")
	protected.Use(middleware.Auth())
	protected.GET("/assessments", h.ListAssessments)
	protected.POST("/assessment", h.PostAssessment)
	protected.GET("/health-profile", h.GetHealthProfile)
	protected.GET("/vaccinations", h.ListVaccinations)

	admin := r.Group("/admin")
	admin.Use(middleware.Auth())
	admin.Use(middleware.RequireRole(pool, "admin", "nurse"))
	admin.GET("/facilities", h.AdminListFacilities)
	admin.POST("/facilities", h.AdminCreateFacility)
	admin.GET("/nurses", h.AdminListNurses)
	admin.POST("/nurses", h.AdminCreateNurse)
	admin.GET("/assessments", h.AdminListAssessments)
	admin.GET("/expert/symptoms", h.AdminListSymptoms)
	admin.GET("/expert/diseases", h.AdminListDiseases)
	admin.GET("/expert/rules", h.AdminListRules)
	admin.POST("/expert/rules", h.AdminCreateRule)
	admin.POST("/expert/rules/:id/publish", h.AdminPublishRule)

	return r, pool
}

func doRequest(r *gin.Engine, method, path string, body any, token string) *httptest.ResponseRecorder {
	var b *bytes.Buffer
	if body != nil {
		raw, _ := json.Marshal(body)
		b = bytes.NewBuffer(raw)
	} else {
		b = bytes.NewBuffer(nil)
	}
	req, _ := http.NewRequest(method, path, b)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// ── TEST-6A: Unauthenticated cannot access protected endpoints ───────────────

func TestRBAC_Unauthenticated_CannotAccessProtected(t *testing.T) {
	r, _ := newRouter(t)

	endpoints := []struct{ method, path string }{
		{"GET", "/assessments"},
		{"GET", "/health-profile"},
		{"GET", "/vaccinations"},
	}
	for _, e := range endpoints {
		w := doRequest(r, e.method, e.path, nil, "")
		if w.Code != http.StatusUnauthorized {
			t.Errorf("%s %s: want 401, got %d", e.method, e.path, w.Code)
		}
	}
}

// ── TEST-6B: Unauthenticated cannot access admin endpoints ───────────────────

func TestRBAC_Unauthenticated_CannotAccessAdmin(t *testing.T) {
	r, _ := newRouter(t)

	adminEndpoints := []struct{ method, path string }{
		{"GET", "/admin/facilities"},
		{"GET", "/admin/nurses"},
		{"GET", "/admin/assessments"},
		{"GET", "/admin/expert/symptoms"},
		{"GET", "/admin/expert/rules"},
	}
	for _, e := range adminEndpoints {
		w := doRequest(r, e.method, e.path, nil, "")
		if w.Code != http.StatusUnauthorized {
			t.Errorf("%s %s: want 401, got %d", e.method, e.path, w.Code)
		}
	}
}

// ── TEST-6C: Traveler JWT cannot access admin endpoints (403) ────────────────

func TestRBAC_TravelerCannotAccessAdmin(t *testing.T) {
	r, pool := newRouter(t)

	// Find a user with only the traveler role (or use a non-existent user ID
	// that would pass JWT validation but fail role check).
	// We use userID=99999 (does not exist in roles table → no admin/nurse role).
	travelerToken := makeJWT(t, 99999)

	adminEndpoints := []struct{ method, path string }{
		{"GET", "/admin/facilities"},
		{"GET", "/admin/nurses"},
		{"GET", "/admin/assessments"},
		{"GET", "/admin/expert/symptoms"},
		{"GET", "/admin/expert/diseases"},
		{"GET", "/admin/expert/rules"},
	}
	_ = pool
	for _, e := range adminEndpoints {
		w := doRequest(r, e.method, e.path, nil, travelerToken)
		if w.Code != http.StatusForbidden {
			t.Errorf("%s %s: want 403, got %d (body: %s)",
				e.method, e.path, w.Code, w.Body.String())
		}
	}
}

// ── TEST-6D: Admin JWT can access admin endpoints (200) ──────────────────────

func TestRBAC_AdminCanAccessAdminEndpoints(t *testing.T) {
	r, pool := newRouter(t)

	// Find actual admin user from DB
	var adminID int64
	err := pool.QueryRow(context.Background(),
		`SELECT u.id FROM users u
		 JOIN user_roles ur ON ur.user_id = u.id
		 JOIN roles ro ON ro.role_id = ur.role_id
		 WHERE ro.nama_role = 'admin' LIMIT 1`,
	).Scan(&adminID)
	if err != nil {
		t.Skipf("no admin user in DB, skipping: %v", err)
	}

	adminToken := makeJWT(t, adminID)

	readEndpoints := []struct{ method, path string }{
		{"GET", "/admin/facilities"},
		{"GET", "/admin/nurses"},
		{"GET", "/admin/assessments"},
		{"GET", "/admin/expert/symptoms"},
		{"GET", "/admin/expert/diseases"},
		{"GET", "/admin/expert/rules"},
	}
	for _, e := range readEndpoints {
		w := doRequest(r, e.method, e.path, nil, adminToken)
		if w.Code != http.StatusOK {
			t.Errorf("%s %s: want 200, got %d (body: %s)",
				e.method, e.path, w.Code, w.Body.String())
		}
	}
}

// ── TEST-6E: GO-28 rule validation — server rejects invalid rules ─────────────

func TestRBAC_ExpertRuleValidation_InvalidCF(t *testing.T) {
	r, pool := newRouter(t)

	var adminID int64
	err := pool.QueryRow(context.Background(),
		`SELECT u.id FROM users u
		 JOIN user_roles ur ON ur.user_id = u.id
		 JOIN roles ro ON ro.role_id = ur.role_id
		 WHERE ro.nama_role = 'admin' LIMIT 1`,
	).Scan(&adminID)
	if err != nil {
		t.Skipf("no admin user in DB, skipping: %v", err)
	}
	adminToken := makeJWT(t, adminID)

	cases := []struct {
		name string
		body map[string]any
		want int
	}{
		{
			name: "bobot_cf > 1",
			body: map[string]any{
				"nama": "Test Rule", "premis": []int{1}, "disease_id": 1,
				"bobot_cf": 1.5, "mb": 0.8, "md": 0.1, "kategori": "pre_travel",
			},
			want: http.StatusUnprocessableEntity,
		},
		{
			name: "mb negatif",
			body: map[string]any{
				"nama": "Test Rule", "premis": []int{1}, "disease_id": 1,
				"bobot_cf": 0.7, "mb": -0.1, "md": 0.0, "kategori": "pre_travel",
			},
			want: http.StatusUnprocessableEntity,
		},
		{
			name: "premis kosong",
			body: map[string]any{
				"nama": "Test Rule", "premis": []int{}, "disease_id": 1,
				"bobot_cf": 0.7, "mb": 0.8, "md": 0.1, "kategori": "pre_travel",
			},
			// validateRule rejects empty premis with ErrRuleValidation → 422
			want: http.StatusUnprocessableEntity,
		},
		{
			name: "kategori tidak valid",
			body: map[string]any{
				"nama": "Test Rule", "premis": []int{1}, "disease_id": 1,
				"bobot_cf": 0.7, "mb": 0.8, "md": 0.1, "kategori": "invalid_kategori",
			},
			want: http.StatusUnprocessableEntity,
		},
		{
			name: "symptom_id tidak ada di master",
			body: map[string]any{
				"nama": "Test Rule", "premis": []int{999999}, "disease_id": 1,
				"bobot_cf": 0.7, "mb": 0.8, "md": 0.1, "kategori": "pre_travel",
			},
			want: http.StatusUnprocessableEntity,
		},
		{
			name: "disease_id tidak ada di master",
			body: map[string]any{
				"nama": "Test Rule", "premis": []int{1}, "disease_id": 999999,
				"bobot_cf": 0.7, "mb": 0.8, "md": 0.1, "kategori": "pre_travel",
			},
			want: http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := doRequest(r, "POST", "/admin/expert/rules", tc.body, adminToken)
			if w.Code != tc.want {
				t.Errorf("want %d, got %d (body: %s)", tc.want, w.Code, w.Body.String())
			}
		})
	}
}

// ── TEST-6F: Public endpoints accessible without auth ────────────────────────

func TestPublicEndpoints_NoAuthRequired(t *testing.T) {
	r, _ := newRouter(t)

	public := []struct{ method, path string }{
		{"GET", "/health"},
		{"GET", "/destinations"},
		{"GET", "/emergency-guides"},
		{"GET", "/location/classify?lat=-8.5&lng=115.2"},
	}
	for _, e := range public {
		w := doRequest(r, e.method, e.path, nil, "")
		if w.Code == http.StatusUnauthorized || w.Code == http.StatusForbidden {
			t.Errorf("%s %s: should be public, got %d", e.method, e.path, w.Code)
		}
	}
}
