package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/saintbyte/home-ctrl/internal/auth"
	"github.com/saintbyte/home-ctrl/internal/config"
	"github.com/saintbyte/home-ctrl/internal/database"
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
		_ = db.GetDB().Exec("DROP TABLE IF EXISTS api_keys")
		_ = db.GetDB().Exec("DROP TABLE IF EXISTS sessions")
	}()

	authService := auth.NewAuth(config.DefaultConfig(), db)
	authService.AddUser("test", "test123")

	// Create server with default config, auth, and database
	cfg := config.DefaultConfig()
	srv := NewServer(cfg, authService, db)
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
		req, _ := http.NewRequest("GET", "/api/version", nil)
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
}