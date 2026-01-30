package main

import (
	"fmt"
	"log"
	"os"

	"github.com/saintbyte/home-ctrl/internal/app"
	"github.com/saintbyte/home-ctrl/internal/database"
)

func main() {
	// Check if we should run migrations
	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		runMigrations()
		return
	}

	fmt.Println("Starting home-ctrl application...")
	
	// Initialize application
	a, err := app.NewApp()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}
	
	// Run application
	if err := a.Run(); err != nil {
		log.Fatalf("Application failed: %v", err)
	}
}

// runMigrations runs database migrations
func runMigrations() {
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