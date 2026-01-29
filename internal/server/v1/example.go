package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/saintbyte/home-ctrl/internal/config"
)

// ExampleHandler handles example endpoints
type ExampleHandler struct {
	config *config.Config
}

// NewExampleHandler creates a new example handler
func NewExampleHandler(cfg *config.Config) *ExampleHandler {
	return &ExampleHandler{
		config: cfg,
	}
}

// SetupRoutes sets up example-related routes
func (h *ExampleHandler) SetupRoutes(router *gin.RouterGroup) {
	router.GET("/example", h.exampleEndpoint)
	router.GET("/me", h.userInfoEndpoint)
}

// exampleEndpoint handles GET /example
func (h *ExampleHandler) exampleEndpoint(c *gin.Context) {
	username, _ := c.Get("username")
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello from protected endpoint!",
		"user":    username,
		"config": gin.H{
			"host": h.config.Server.Host,
			"port": h.config.Server.Port,
		},
	})
}

// userInfoEndpoint handles GET /me
func (h *ExampleHandler) userInfoEndpoint(c *gin.Context) {
	username, _ := c.Get("username")
	c.JSON(http.StatusOK, gin.H{
		"username": username,
		"message": "You are authenticated!",
	})
}