package migrations

import (
	"fmt"
	"github.com/saintbyte/home-ctrl/internal/database"
	"log"
)

// runMigrations runs database migrations
func RunMigrations() {
	fmt.Println("Running database migrations...")

	// Initialize database
	db, err := database.NewDatabase("data")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := db.InitDatabaseWithMigrations(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	fmt.Println("Database migrations completed successfully!")
}
