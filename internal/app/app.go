package app

import (
	"fmt"
	"github.com/saintbyte/home-ctrl/internal/auth"
	"github.com/saintbyte/home-ctrl/internal/config"
	"github.com/saintbyte/home-ctrl/internal/database"
	"github.com/saintbyte/home-ctrl/internal/scheduler"
	"github.com/saintbyte/home-ctrl/internal/server"
	"github.com/saintbyte/home-ctrl/internal/weather"
	"os"
)

// App represents the main application
type App struct {
	name    string
	version string
	config  *config.Config
	db      *database.Database
	auth    *auth.Auth
	server  *server.Server
	sched   *scheduler.Scheduler
	weather *weather.YandexWeather
}

// NewApp creates a new application instance
func NewApp() (*App, error) {
	configPaths := []string{}
	if envPath := os.Getenv("HOME_CTRL_CONFIG"); envPath != "" {
		configPaths = []string{envPath}
	} else {
		configPaths = []string{
			"config.yaml",
			"/etc/home-ctrl/config.yaml",
		}
	}

	var cfg *config.Config
	var loadErr error
	for _, path := range configPaths {
		fmt.Println(path)
		cfg, loadErr = config.LoadConfig(path)
		if loadErr != nil {
			Log.Warn("failed loading configuration file:", "error=", loadErr)
		}
		if loadErr == nil && cfg != nil {
			break
		}
	}
	if cfg == nil {
		cfg = config.DefaultConfig()
	}

	// Initialize database
	db, err := database.NewDatabase(cfg.DataDir)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err, cfg.DataDir)
	}

	// Initialize database with migrations
	if err := db.InitDatabaseWithMigrations(); err != nil {
		Log.Warn("failed to initialize database", "error", cfg.DataDir)
	}

	// Initialize authentication
	authService := auth.NewAuth(cfg, db)

	// Add users from config
	for username, password := range cfg.Auth.Users {
		authService.AddUser(username, password)
	}

	// Create Yandex Weather service from the configured task
	var weatherService *weather.YandexWeather
	for _, task := range cfg.Tasks {
		if task.Command == "yandex_weather" {
			weatherService = weather.NewYandexWeather(task)
			break
		}
	}

	// Create scheduler
	sched := scheduler.NewScheduler(cfg, func(task config.Task) {
		if task.Command == "yandex_weather" && weatherService != nil {
			Log.Info("yandex weather: fetching")
			if err := weatherService.Fetch(); err != nil {
				Log.Warn("yandex weather: fetch failed", "error", err.Error())
			}
			return
		}
		Log.Info("task executed", "task", task.Name)
	})

	// Create server with auth and database
	srv := server.NewServer(cfg, authService, db, sched, weatherService)
	srv.SetupRoutes()

	return &App{
		name:    "home-ctrl",
		version: "0.1.0",
		config:  cfg,
		db:      db,
		auth:    authService,
		server:  srv,
		sched:   sched,
		weather: weatherService,
	}, nil
}

// Run starts the application
func (a *App) Run() error {
	fmt.Printf("Running %s v%s\n", a.name, a.version)
	fmt.Printf("Server listening on %s\n", a.config.GetServerAddress())
	fmt.Printf("Data directory: %s\n", a.config.DataDir)

	// Fetch weather once at startup so the API returns data immediately
	if a.weather != nil && a.weather.Enabled() {
		go func() {
			if err := a.weather.Fetch(); err != nil {
				Log.Warn("yandex weather: initial fetch failed", "error", err.Error())
			}
		}()
	}

	// Start scheduler
	a.sched.Start()

	// Start HTTP server
	if err := a.server.Run(); err != nil {
		return fmt.Errorf("server failed: %w", err)
	}

	return nil
}

// Close closes the application resources
func (a *App) Close() error {
	if a.sched != nil {
		a.sched.Stop()
	}
	if a.db != nil {
		return a.db.Close()
	}
	return nil
}

// GetScheduler returns the scheduler
func (a *App) GetScheduler() *scheduler.Scheduler {
	return a.sched
}

// RunAsDaemon runs the application as a daemon with signal handling
func (a *App) RunAsDaemon() error {
	daemon, err := NewDaemon()
	if err != nil {
		return fmt.Errorf("failed to create daemon: %w", err)
	}
	return daemon.RunDaemon()
}
