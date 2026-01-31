package server

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/saintbyte/home-ctrl/internal/auth"
	"github.com/saintbyte/home-ctrl/internal/config"
	"github.com/saintbyte/home-ctrl/internal/database"
	"github.com/saintbyte/home-ctrl/internal/server/v1"
)

// Server represents the HTTP server
type Server struct {
	config   *config.Config
	auth     *auth.Auth
	v1Router *v1.Router
	router   *gin.Engine
}

// NewServer creates a new server instance
func NewServer(cfg *config.Config, authService *auth.Auth, db *database.Database) *Server {
	return &Server{
		config:   cfg,
		auth:     authService,
		v1Router: v1.NewRouter(cfg, authService, db),
		router:   gin.Default(),
	}
}

// SetupRoutes sets up the HTTP routes
func (s *Server) SetupRoutes() {
	// Configure CORS to allow all origins
	s.router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * 3600, // 12 hours
	}))

	// Setup v1 routes - use the main router instead of v1Router
	s.v1Router.SetupRoutesOn(s.router)

	// Health check endpoint (not versioned)
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Service is running",
		})
	})

	// Serve static files from public directory in root
	// Files from public will be accessible directly (e.g., /style.css, /script.js)
	// Use StaticFile for root to serve index.html
	s.router.StaticFile("/", "./public/index.html")

	// Serve other static files from public directory
	// This middleware checks if a file exists and serves it, otherwise passes to next handler
	s.router.Use(func(c *gin.Context) {
		path := c.Request.URL.Path
		// Skip root path, health check, and API routes
		if path == "/" || path == "/health" ||
			path == "/api" ||
			(len(path) > 5 && path[:5] == "/api/") {
			c.Next()
			return
		}

		// Check if file exists in public directory
		filePath := filepath.Join("./public", path)
		if info, err := os.Stat(filePath); err == nil && !info.IsDir() {
			c.File(filePath)
			c.Abort()
			return
		}

		c.Next()
	})

	// Fallback: serve index.html for all other routes (for SPA support)
	s.router.NoRoute(func(c *gin.Context) {
		// Skip health check and API routes
		if c.Request.URL.Path == "/health" ||
			c.Request.URL.Path == "/api" ||
			(len(c.Request.URL.Path) > 5 && c.Request.URL.Path[:5] == "/api/") {
			c.Next()
			return
		}
		c.File("./public/index.html")
	})
}

// Run starts the HTTP server
func (s *Server) Run() error {
	address := s.config.GetServerAddress()
	fmt.Printf("Starting server on %s\n", address)

	return s.router.Run(address)
}

// GetRouter returns the gin router (for testing)
func (s *Server) GetRouter() *gin.Engine {
	return s.v1Router.GetRouter()
}
