package v1_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/saintbyte/home-ctrl/internal/auth"
	"github.com/saintbyte/home-ctrl/internal/config"
	"github.com/saintbyte/home-ctrl/internal/database"
	"github.com/saintbyte/home-ctrl/internal/server/v1"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() *v1.Router {
	// Load test configuration
	cfg, err := config.LoadConfig("../../../../../test/config.yaml")
	if err != nil {
		panic(err)
	}

	// Initialize authentication
	authService := auth.NewAuth(cfg.Auth)

	// Initialize database
	db, err := database.NewDatabase("../../../../../test_data")
	if err != nil {
		panic(err)
	}

	// Initialize router
	router := v1.NewRouter(cfg, authService, db)
	router.SetupRoutes()

	return router
}

func TestRouterPublicRoutes(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	r := setupTestRouter()

	// Create a request to test health endpoint
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/health", nil)
	r.GetRouter().ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "ok", response["status"])
}

func TestRouterVersionRoute(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	r := setupTestRouter()

	// Create a request to test version endpoint
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/version", nil)
	r.GetRouter().ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "0.1.0", response["version"])
	assert.Equal(t, "home-ctrl", response["name"])
}

func TestRouterProtectedRoutesWithoutAuth(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	r := setupTestRouter()

	// Test accessing protected endpoint without authentication
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/example", nil)
	r.GetRouter().ServeHTTP(w, req)

	// Assert unauthorized access
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Unauthorized", response["error"])
}

func TestRouterProtectedRoutesWithAuth(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	r := setupTestRouter()

	// First, login to get a token
	w := httptest.NewRecorder()
	loginData := map[string]string{
		"username": "testuser",
		"password": "testpass",
	}
	loginBody, _ := json.Marshal(loginData)
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(loginBody))
	req.Header.Set("Content-Type", "application/json")
	r.GetRouter().ServeHTTP(w, req)

	// Extract token from response
	var loginResponse map[string]string
	json.Unmarshal(w.Body.Bytes(), &loginResponse)
	token := loginResponse["token"]
	assert.NotEmpty(t, token)

	// Now test accessing protected endpoint with token
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/example", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.GetRouter().ServeHTTP(w, req)

	// Assert successful access
	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Hello from protected endpoint!", response["message"])
	assert.Equal(t, "testuser", response["user"])
}

func TestRouter404Handler(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	r := setupTestRouter()

	// Create a request to a non-existent endpoint
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/nonexistent", nil)
	r.GetRouter().ServeHTTP(w, req)

	// Assert 404 response
	assert.Equal(t, http.StatusNotFound, w.Code)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Not found", response["error"])
	assert.Equal(t, "/nonexistent", response["path"])
}
