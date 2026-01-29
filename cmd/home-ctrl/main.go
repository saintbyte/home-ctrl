package main

import (
	"fmt"
	"log"

	"github.com/saintbyte/home-ctrl/internal/app"
)

func main() {
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