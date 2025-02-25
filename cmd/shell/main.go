package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Ali-Farhadnia/goshell/internal/app"
	"github.com/Ali-Farhadnia/goshell/internal/config"
)

func main() {
	// Parse command line flags
	configPath := flag.String("config", "config.yaml", "Path to configuration file")
	verbose := flag.Bool("verbose", false, "Enable verbose mode")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Override config with command line flags
	if *verbose {
		cfg.Shell.Verbose = true
	}

	// Initialize and run the shell
	app, err := app.New(cfg)
	if err != nil {
		fmt.Printf("Error initializing app: %v\n", err)
		os.Exit(1)
	}

	app.Run()
}
