package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/saintbyte/home-ctrl/internal/config"
)

// Daemon represents a daemonized application
type Daemon struct {
	app          *App
	configPath   string
	signalChan   chan os.Signal
	shutdownChan chan struct{}
	ctx          context.Context
	cancel       context.CancelFunc
}

// NewDaemon creates a new daemon instance
func NewDaemon() (*Daemon, error) {
	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())

	return &Daemon{
		signalChan:   make(chan os.Signal, 1),
		shutdownChan: make(chan struct{}),
		ctx:          ctx,
		cancel:       cancel,
	}, nil
}

// SetupSignalHandling sets up signal handling for the daemon
func (d *Daemon) SetupSignalHandling() {
	signal.Notify(d.signalChan,
		syscall.SIGINT,  // Ctrl+C
		syscall.SIGTERM, // Termination signal
		syscall.SIGHUP,  // Hangup signal (for config reload)
	)
}

// ReloadConfig reloads the configuration
func (d *Daemon) ReloadConfig() error {
	if d.configPath == "" {
		return fmt.Errorf("config path not set")
	}

	log.Printf("Reloading configuration from %s", d.configPath)

	// Reload configuration
	cfg, err := config.LoadConfig(d.configPath)
	if err != nil {
		return fmt.Errorf("failed to reload config: %w", err)
	}

	// Update app configuration
	d.app.config = cfg

	// Reinitialize authentication with new config
	for username, password := range cfg.Auth.Users {
		d.app.auth.AddUser(username, password)
	}

	log.Printf("Configuration reloaded successfully")
	return nil
}

// RunDaemon runs the application as a daemon
func (d *Daemon) RunDaemon() error {
	// Load initial configuration
	configPath := os.Getenv("HOME_CTRL_CONFIG")
	if configPath == "" {
		configPath = "config.yaml"
	}
	d.configPath = configPath

	// Create app instance
	app, err := NewApp()
	if err != nil {
		return fmt.Errorf("failed to create app: %w", err)
	}
	d.app = app

	// Setup signal handling
	d.SetupSignalHandling()

	// Start server in a goroutine
	go func() {
		if err := app.Run(); err != nil {
			log.Printf("Server error: %v", err)
			d.cancel()
		}
	}()

	log.Printf("Daemon started, PID: %d", os.Getpid())

	// Main signal handling loop
	for {
		select {
		case sig := <-d.signalChan:
			d.handleSignal(sig)
		case <-d.shutdownChan:
			log.Printf("Shutdown complete")
			return nil
		case <-d.ctx.Done():
			log.Printf("Context cancelled, shutting down")
			return nil
		}
	}
}

// handleSignal handles incoming signals
func (d *Daemon) handleSignal(sig os.Signal) {
	switch sig {
	case syscall.SIGHUP:
		log.Printf("Received SIGHUP, reloading configuration")
		if err := d.ReloadConfig(); err != nil {
			log.Printf("Failed to reload config: %v", err)
		}

	case syscall.SIGINT, syscall.SIGTERM:
		log.Printf("Received %v, shutting down gracefully", sig)
		d.shutdown()

	default:
		log.Printf("Received unexpected signal: %v", sig)
	}
}

// shutdown performs graceful shutdown
func (d *Daemon) shutdown() {
	// Cancel context to signal shutdown
	d.cancel()

	// Close database connection
	if d.app != nil && d.app.db != nil {
		if err := d.app.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}

	// Give some time for cleanup
	time.Sleep(1 * time.Second)

	// Signal shutdown complete
	close(d.shutdownChan)
	log.Printf("Daemon stopped")
}
