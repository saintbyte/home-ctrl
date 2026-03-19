package v1

import (
	"github.com/saintbyte/home-ctrl/internal/version"
	"net/http"

	"github.com/gin-gonic/gin"
)

// VersionHandler handles the version endpoint
type VersionHandler struct{}

// NewVersionHandler creates a new version handler
func NewVersionHandler() *VersionHandler {
	return &VersionHandler{}
}

// SetupRoutes sets up version-related routes
func (h *VersionHandler) SetupRoutes(router *gin.RouterGroup) {
	router.GET("/version", h.getVersion)
}

// getVersion handles GET /version
func (h *VersionHandler) getVersion(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"version": version.Version,
		"name":    version.AppName,
	})
}
