package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/saintbyte/home-ctrl/internal/auth"
	"github.com/saintbyte/home-ctrl/internal/config"
)

// Router represents the v1 API router
type Router struct {
	config    *config.Config
	auth      *auth.Auth
	database  *database.Database
	router    *gin.Engine
}

// NewRouter creates a new v1 router
func NewRouter(cfg *config.Config, authService *auth.Auth, db *database.Database) *Router {
	return &Router{
		config:    cfg,
		auth:      authService,
		database:  db,
		router:    gin.Default(),
	}
}

// SetupRoutes sets up all v1 routes
func (r *Router) SetupRoutes() {
	// Public routes (no authentication required)
	r.setupPublicRoutes()
	
	// Auth routes
	r.setupAuthRoutes()
	
	// Protected routes (require authentication)
	r.setupProtectedRoutes()
	
	// 404 handler
	r.router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"error": "Not found",
			"path":  c.Request.URL.Path,
		})
	})
}

// setupPublicRoutes sets up routes that don't require authentication
func (r *Router) setupPublicRoutes() {
	publicGroup := r.router.Group("/api/v1")
	
	healthHandler := NewHealthHandler()
	healthHandler.SetupRoutes(publicGroup)
	
	versionHandler := NewVersionHandler()
	versionHandler.SetupRoutes(publicGroup)
}

// setupAuthRoutes sets up authentication-related routes
func (r *Router) setupAuthRoutes() {
	authGroup := r.router.Group("/api/v1/auth")
	{
		// Login endpoint
		authGroup.POST("/login", r.auth.LoginHandler())
		
		// Logout endpoint (requires authentication)
		authGroup.POST("/logout", r.auth.AuthMiddleware(), r.auth.LogoutHandler())
	}
}

// setupProtectedRoutes sets up routes that require authentication
func (r *Router) setupProtectedRoutes() {
	protectedGroup := r.router.Group("/api/v1")
	protectedGroup.Use(r.auth.AuthMiddleware())
	
	exampleHandler := NewExampleHandler(r.config)
	exampleHandler.SetupRoutes(protectedGroup)
	
	// Add key-value handler
	keyValueHandler := handlers.NewKeyValueHandler(r.database)
	keyValueHandler.SetupRoutes(protectedGroup)
}

// GetRouter returns the gin router
func (r *Router) GetRouter() *gin.Engine {
	return r.router
}

// Run starts the server
func (r *Router) Run(address string) error {
	return r.router.Run(address)
}