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
	"github.com/saintbyte/home-ctrl/internal/server/v1"
	"github.com/stretchr/testify/assert"
)

func setupExampleTestRouter() (*v1.Router, *config.Config) {
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

	return router, cfg
}

func TestExampleHandlerExampleEndpoint(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	r, cfg := setupExampleTestRouter()

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

	// Now test accessing example endpoint with token
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

	// Check config is properly returned
	configMap := response["config"].(map[string]interface{})
	assert.Equal(t, cfg.Server.Host, configMap["host"])
	assert.Equal(t, float64(cfg.Server.Port), configMap["port"]) // JSON numbers are float64
}

func TestExampleHandlerUserInfoEndpoint(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	r, _ := setupExampleTestRouter()

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

	// Now test accessing user info endpoint with token
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.GetRouter().ServeHTTP(w, req)

	// Assert successful access
	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "testuser", response["username"])
	assert.Equal(t, "You are authenticated!", response["message"])
}
