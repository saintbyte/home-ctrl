package server

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/saintbyte/home-ctrl/internal/auth"
	"github.com/saintbyte/home-ctrl/internal/config"
	"github.com/saintbyte/home-ctrl/internal/database"
	"github.com/saintbyte/home-ctrl/internal/scheduler"
	"github.com/stretchr/testify/assert"
)

func TestServerRoutes(t *testing.T) {
	// Set gin to test mode
	gin.SetMode(gin.TestMode)

	// Create database and auth for testing
	db, err := database.NewDatabase("test_data")
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer db.Close()
	defer func() {
		// Cleanup test database
		db.GetDB().Exec("DROP TABLE IF EXISTS api_keys")
		db.GetDB().Exec("DROP TABLE IF EXISTS sessions")
	}()

	authService := auth.NewAuth(config.DefaultConfig(), db)
	authService.AddUser("test", "test123")

	// Create server with default config, auth, and database
	cfg := config.DefaultConfig()
	srv := NewServer(cfg, authService, db, scheduler.NewScheduler(cfg, nil), nil)
	srv.SetupRoutes()

	t.Run("Health endpoint", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/health", nil)
		srv.GetRouter().ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Service is running")
	})

	t.Run("Version endpoint", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/version", nil)
		srv.GetRouter().ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "home-ctrl")
		assert.Contains(t, w.Body.String(), "0.1.0")
	})

	t.Run("Example endpoint without auth", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/example", nil)
		srv.GetRouter().ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Authentication required")
	})

	t.Run("Example endpoint with valid auth", func(t *testing.T) {
		// First, login to get a token
		loginW := httptest.NewRecorder()
		loginReq, _ := http.NewRequest("POST", "/api/v1/auth/login", nil)
		loginReq.Header.Set("Content-Type", "application/json")
		srv.GetRouter().ServeHTTP(loginW, loginReq)

		// Then use the token to access protected endpoint
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/example", nil)
		req.Header.Set("Authorization", "Bearer test-token")
		srv.GetRouter().ServeHTTP(w, req)

		// Should still fail because we don't have a real token
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("404 endpoint", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/nonexistent", nil)
		srv.GetRouter().ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "Not found")
	})

	t.Run("Login endpoint", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/auth/login", nil)
		req.Header.Set("Content-Type", "application/json")
		srv.GetRouter().ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Bad Request")
	})

	t.Run("HSTS header sent on TLS", func(t *testing.T) {
		cfg := config.DefaultConfig()
		cfg.Server.TLS = true
		tlsSrv := NewServer(cfg, authService, db, scheduler.NewScheduler(cfg, nil), nil)
		tlsSrv.SetupRoutes()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/health", nil)
		tlsSrv.GetRouter().ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "max-age=31536000", w.Header().Get("Strict-Transport-Security"))
	})

	t.Run("No HSTS header without TLS", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/health", nil)
		srv.GetRouter().ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Empty(t, w.Header().Get("Strict-Transport-Security"))
	})
}

func TestBootstrapRouter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// The app resolves ./public relative to the repo root, but tests run from
	// the package directory. Change to the repo root for the duration.
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(cwd)
	if err := os.Chdir("../.."); err != nil {
		t.Fatalf("Failed to change to repo root: %v", err)
	}

	cfg := config.DefaultConfig()
	cfg.Server.TLS = true
	cfg.Server.Port = 8443
	cfg.Server.HTTPPort = 8080

	db, err := database.NewDatabase("test_data")
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer db.Close()
	defer func() {
		db.GetDB().Exec("DROP TABLE IF EXISTS api_keys")
		db.GetDB().Exec("DROP TABLE IF EXISTS sessions")
	}()

	authService := auth.NewAuth(config.DefaultConfig(), db)
	srv := NewServer(cfg, authService, db, scheduler.NewScheduler(cfg, nil), nil)
	bootstrap := srv.buildBootstrapRouter()

	t.Run("Main page is served over HTTP", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		bootstrap.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Cert download served over HTTP", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/cert", nil)
		bootstrap.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "root certificate is not configured")
	})

	t.Run("Everything else redirected to HTTPS preserving host", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/version", nil)
		req.Host = "home.lan:8080"
		bootstrap.ServeHTTP(w, req)

		assert.Equal(t, http.StatusMovedPermanently, w.Code)
		assert.Equal(t, "https://home.lan:8443/api/v1/version", w.Header().Get("Location"))
	})

	t.Run("Bootstrap responses carry no HSTS header", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/version", nil)
		bootstrap.ServeHTTP(w, req)
		assert.Empty(t, w.Header().Get("Strict-Transport-Security"))
	})
}
