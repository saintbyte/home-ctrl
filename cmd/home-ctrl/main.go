package main

import (
	"flag"
	"fmt"
	"github.com/saintbyte/home-ctrl/internal/migrations"
	"log"
	"os"

	"github.com/saintbyte/home-ctrl/internal/app"
)

func main() {
	// Parse command line flags
	daemonMode := flag.Bool("daemon", true, "Run as a daemon with signal handling")
	flag.Parse()

	// Check if we should run migrations
	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		migrations.RunMigrations()
		return
	}

	fmt.Println("Starting home-ctrl application...")

	// Initialize application
	a, err := app.NewApp()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Run application in appropriate mode
	if *daemonMode {
		fmt.Println("Running in daemon mode...")
		if err := a.RunAsDaemon(); err != nil {
			log.Fatalf("Daemon failed: %v", err)
		}
	} else {
		if err := a.Run(); err != nil {
			log.Fatalf("Application failed: %v", err)
		}
	}
}
