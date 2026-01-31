package app

import (
	"fmt"
	"os"

	"github.com/saintbyte/home-ctrl/internal/auth"
	"github.com/saintbyte/home-ctrl/internal/config"
	"github.com/saintbyte/home-ctrl/internal/database"
	"github.com/saintbyte/home-ctrl/internal/server"
)

// App represents the main application
type App struct {
	name    string
	version string
	config  *config.Config
	db      *database.Database
	auth    *auth.Auth
	server  *server.Server
}

// NewApp creates a new application instance
func NewApp() (*App, error) {
	// Load configuration
	configPath := os.Getenv("HOME_CTRL_CONFIG")
	if configPath == "" {
		configPath = "config.yaml"
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Initialize database
	db, err := database.NewDatabase(cfg.DataDir)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Initialize database with migrations
	if err := db.InitDatabaseWithMigrations(); err != nil {
		Log.Warn("failed to initialize database", "error", err)
	}

	// Initialize authentication
	authService := auth.NewAuth(cfg, db)

	// Add users from config
	for username, password := range cfg.Auth.Users {
		authService.AddUser(username, password)
	}

	// Create server with auth and database
	srv := server.NewServer(cfg, authService, db)
	srv.SetupRoutes()

	return &App{
		name:    "home-ctrl",
		version: "0.1.0",
		config:  cfg,
		db:      db,
		auth:    authService,
		server:  srv,
	}, nil
}

// Run starts the application
func (a *App) Run() error {
	fmt.Printf("Running %s v%s\n", a.name, a.version)
	fmt.Printf("Server listening on %s\n", a.config.GetServerAddress())
	fmt.Printf("Data directory: %s\n", a.config.DataDir)

	// Start HTTP server
	if err := a.server.Run(); err != nil {
		return fmt.Errorf("server failed: %w", err)
	}

	return nil
}

// Close closes the application resources
func (a *App) Close() error {
	if a.db != nil {
		return a.db.Close()
	}
	return nil
}

// RunAsDaemon runs the application as a daemon with signal handling
func (a *App) RunAsDaemon() error {
	daemon, err := NewDaemon()
	if err != nil {
		return fmt.Errorf("failed to create daemon: %w", err)
	}
	return daemon.RunDaemon()
}
