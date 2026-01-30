package database

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Migration represents a database migration
type Migration struct {
	Name string
	Path string
}

// RunMigrations runs all database migrations
func (d *Database) RunMigrations() error {
	migrations, err := d.getMigrations()
	if err != nil {
		return fmt.Errorf("failed to get migrations: %w", err)
	}

	// Sort migrations by name
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Name < migrations[j].Name
	})

	// Run each migration
	for _, migration := range migrations {
		if err := d.runMigration(migration); err != nil {
			return fmt.Errorf("failed to run migration %s: %w", migration.Name, err)
		}
	}

	return nil
}

// getMigrations gets all migration files
func (d *Database) getMigrations() ([]Migration, error) {
	migrationsDir := "internal/database/migrations"
	
	// Read migration files
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read migrations directory: %w", err)
	}

	var migrations []Migration
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		
		if filepath.Ext(file.Name()) != ".sql" {
			continue
		}
		
		migrations = append(migrations, Migration{
			Name: file.Name(),
			Path: filepath.Join(migrationsDir, file.Name()),
		})
	}

	return migrations, nil
}

// runMigration runs a single migration
func (d *Database) runMigration(migration Migration) error {
	log.Printf("Running migration: %s", migration.Name)

	// Read migration file
	content, err := os.ReadFile(migration.Path)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	// Split into individual statements
	statements := strings.Split(string(content), ";")
	
	// Execute each statement
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		
		if _, err := d.db.Exec(stmt); err != nil {
			return fmt.Errorf("failed to execute statement '%s': %w", stmt, err)
		}
	}

	log.Printf("Migration completed: %s", migration.Name)
	return nil
}

// InitDatabaseWithMigrations initializes the database with migrations
func (d *Database) InitDatabaseWithMigrations() error {
	// Run migrations
	if err := d.RunMigrations(); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	// Create additional tables if needed
	if err := d.CreateKeyValueTable(); err != nil {
		return fmt.Errorf("failed to create key-value table: %w", err)
	}

	return nil
}