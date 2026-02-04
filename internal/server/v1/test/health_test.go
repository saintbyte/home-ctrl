package v1_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/saintbyte/home-ctrl/internal/server/v1"
	"github.com/stretchr/testify/assert"
)

func TestHealthHandler(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Create handler and router group
	handler := v1.NewHealthHandler()
	group := gin.New()

	// Setup routes
	handler.SetupRoutes(group)

	// Create request
	req, _ := http.NewRequest("GET", "/health", nil)
	group.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "ok", response["status"])
	assert.Equal(t, "Service is running", response["message"])
}
