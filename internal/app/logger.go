package app

import (
	"log/slog"
	"os"
)

// NewLogger creates a new slog logger with appropriate configuration
func NewLogger() *slog.Logger {
	// Create a handler with JSON format and appropriate level
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	})

	return slog.New(handler)
}

// GetLogger returns a default logger instance
var logger = NewLogger()

// Log is a convenience variable for the default logger
var Log = logger
