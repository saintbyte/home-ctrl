package handlers_test

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
	"github.com/saintbyte/home-ctrl/internal/server/v1/handlers"
	"github.com/stretchr/testify/assert"
)

func setupKeyValueTestRouter() (*gin.Engine, *database.Database) {
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

	// Setup test data
	setupTestData(db)

	return router.GetRouter(), db
}

func setupTestData(db *database.Database) {
	// Read and execute test data SQL
	testData, err := os.ReadFile("../../../../../test/testdata.sql")
	if err != nil {
		panic(err)
	}

	statements := strings.Split(string(testData), ";")
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		if _, err := db.GetDB().Exec(stmt); err != nil {
			panic(err)
		}
	}
}

func TestKeyValueHandlerCreateKeyValue(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router, db := setupKeyValueTestRouter()

	// Login to get token
	w := httptest.NewRecorder()
	loginData := map[string]string{
		"username": "testuser",
		"password": "testpass",
	}
	loginBody, _ := json.Marshal(loginData)
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(loginBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Extract token
	var loginResponse map[string]string
	json.Unmarshal(w.Body.Bytes(), &loginResponse)
	token := loginResponse["token"]
	assert.NotEmpty(t, token)

	// Test creating a new key-value pair
	w = httptest.NewRecorder()
	createData := map[string]string{
		"key":   "new-test-key",
		"value": "new-test-value",
	}
	createBody, _ := json.Marshal(createData)
	req, _ = http.NewRequest("POST", "/api/v1/keyvalue", bytes.NewBuffer(createBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	// Assert creation success
	assert.Equal(t, http.StatusCreated, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "new-test-key", response["key"])
	assert.Equal(t, "new-test-value", response["value"])
	assert.Equal(t, "unread", response["status"])
	assert.Equal(t, false, response["hidden"])

	// Verify it was actually created in database
	kv, err := db.GetKeyValue("new-test-key")
	assert.NoError(t, err)
	assert.NotNil(t, kv)
	assert.Equal(t, "new-test-key", kv.Key)
	assert.Equal(t, "new-test-value", kv.Value)
}

func TestKeyValueHandlerGetKeyValue(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router, _ := setupKeyValueTestRouter()

	// Login to get token
	w := httptest.NewRecorder()
	loginData := map[string]string{
		"username": "testuser",
		"password": "testpass",
	}
	loginBody, _ := json.Marshal(loginData)
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(loginBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Extract token
	var loginResponse map[string]string
	json.Unmarshal(w.Body.Bytes(), &loginResponse)
	token := loginResponse["token"]
	assert.NotEmpty(t, token)

	// Test getting existing key-value pair
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/keyvalue/test-key-1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	// Assert successful retrieval
	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "test-key-1", response["key"])
	assert.Equal(t, "test-value-1", response["value"])
	assert.Equal(t, "unread", response["status"])
	assert.Equal(t, false, response["hidden"])

	// Test getting non-existent key
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/keyvalue/non-existent-key", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	// Assert not found
	assert.Equal(t, http.StatusNotFound, w.Code)
	var errorResponse map[string]string
	json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.Equal(t, "Not Found", errorResponse["error"])
	assert.Equal(t, "Key not found", errorResponse["message"])
}

func TestKeyValueHandlerUpdateKeyValue(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router, _ := setupKeyValueTestRouter()

	// Login to get token
	w := httptest.NewRecorder()
	loginData := map[string]string{
		"username": "testuser",
		"password": "testpass",
	}
	loginBody, _ := json.Marshal(loginData)
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(loginBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Extract token
	var loginResponse map[string]string
	json.Unmarshal(w.Body.Bytes(), &loginResponse)
	token := loginResponse["token"]
	assert.NotEmpty(t, token)

	// Test updating existing key-value pair
	w = httptest.NewRecorder()
	updateData := map[string]string{
		"value": "updated-test-value",
	}
	updateBody, _ := json.Marshal(updateData)
	req, _ = http.NewRequest("PUT", "/api/v1/keyvalue/test-key-1", bytes.NewBuffer(updateBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	// Assert update success
	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "test-key-1", response["key"])
	assert.Equal(t, "updated-test-value", response["value"])
	assert.Equal(t, "unread", response["status"])
	assert.Equal(t, false, response["hidden"])
}

func TestKeyValueHandlerDeleteKeyValue(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router, db := setupKeyValueTestRouter()

	// Login to get token
	w := httptest.NewRecorder()
	loginData := map[string]string{
		"username": "testuser",
		"password": "testpass",
	}
	loginBody, _ := json.Marshal(loginData)
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(loginBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Extract token
	var loginResponse map[string]string
	json.Unmarshal(w.Body.Bytes(), &loginResponse)
	token := loginResponse["token"]
	assert.NotEmpty(t, token)

	// Test deleting existing key-value pair
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/api/v1/keyvalue/test-key-1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	// Assert deletion success
	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Key-value pair deleted successfully", response["message"])

	// Verify it was actually deleted from database
	kv, err := db.GetKeyValue("test-key-1")
	assert.NoError(t, err)
	assert.Nil(t, kv)
}

func TestKeyValueHandlerListKeyValues(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router, _ := setupKeyValueTestRouter()

	// Login to get token
	w := httptest.NewRecorder()
	loginData := map[string]string{
		"username": "testuser",
		"password": "testpass",
	}
	loginBody, _ := json.Marshal(loginData)
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(loginBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Extract token
	var loginResponse map[string]string
	json.Unmarshal(w.Body.Bytes(), &loginResponse)
	token := loginResponse["token"]
	assert.NotEmpty(t, token)

	// Test listing key-value pairs (without hidden)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/keyvalue", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	// Assert successful listing
	assert.Equal(t, http.StatusOK, w.Code)
	var response []map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Len(t, response, 2) // Should have 2 non-hidden items

	// Test listing key-value pairs (with hidden)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/keyvalue?include_hidden=true", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	// Assert successful listing with hidden
	assert.Equal(t, http.StatusOK, w.Code)
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Len(t, response, 3) // Should have 3 items including hidden
}
