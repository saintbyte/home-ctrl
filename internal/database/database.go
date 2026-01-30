package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// Database represents the SQLite database connection
type Database struct {
	db *sql.DB
}

// NewDatabase creates a new database instance
func NewDatabase(dataDir string) (*Database, error) {
	// Create data directory if it doesn't exist
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	// Database file path
	dbPath := filepath.Join(dataDir, "home-ctrl.db")

	// Open database connection
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Create tables
	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return &Database{db: db}, nil
}

// createTables creates the necessary database tables
func createTables(db *sql.DB) error {
	tables := []string{
		`CREATE TABLE IF NOT EXISTS api_keys (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			key TEXT UNIQUE NOT NULL,
			name TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			expires_at TIMESTAMP NULL
		)`,
		`CREATE TABLE IF NOT EXISTS sessions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			session_id TEXT UNIQUE NOT NULL,
			username TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			expires_at TIMESTAMP NOT NULL
		)`,
	}

	for _, table := range tables {
		if _, err := db.Exec(table); err != nil {
			return fmt.Errorf("failed to create table: %w", err)
		}
	}

	return nil
}

// Close closes the database connection
func (d *Database) Close() error {
	if d.db != nil {
		return d.db.Close()
	}
	return nil
}

// GetDB returns the underlying sql.DB instance
func (d *Database) GetDB() *sql.DB {
	return d.db
}

// InitDatabase initializes the database with default data
func (d *Database) InitDatabase() error {
	// Create key-value table
	if err := d.CreateKeyValueTable(); err != nil {
		return fmt.Errorf("failed to create key-value table: %w", err)
	}

	// Check if we have any API keys
	var count int
	if err := d.db.QueryRow("SELECT COUNT(*) FROM api_keys").Scan(&count); err != nil {
		return fmt.Errorf("failed to check api_keys count: %w", err)
	}

	// If no API keys, create a default one
	if count == 0 {
		defaultKey := "default-api-key-12345"
		_, err := d.db.Exec(
			"INSERT INTO api_keys (key, name) VALUES (?, ?)",
			defaultKey, "Default API Key",
		)
		if err != nil {
			return fmt.Errorf("failed to insert default API key: %w", err)
		}
		log.Printf("Created default API key: %s", defaultKey)
	}

	return nil
}