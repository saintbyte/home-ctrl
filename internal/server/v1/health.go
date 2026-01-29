package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthHandler handles the health check endpoint
type HealthHandler struct{}

// NewHealthHandler creates a new health handler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// SetupRoutes sets up health-related routes
func (h *HealthHandler) SetupRoutes(router *gin.RouterGroup) {
	router.GET("/health", h.healthCheck)
}

// healthCheck handles GET /health
func (h *HealthHandler) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Service is running",
	})
}