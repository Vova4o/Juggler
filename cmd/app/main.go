package main

import (
	"log"
	"os"

	"juggler/internal/app"
	"juggler/internal/config"
)

func main() {
	cfg, err := config.LoadFromArgs(os.Args)
	if err != nil {
		config.PrintUsage()
		log.Fatalf("Configuration error: %v", err)
	}

	if err := cfg.Validate(); err != nil {
		log.Fatalf("Configuration validation error: %v", err)
	}

	application := app.NewApp(cfg)
	if err := application.Run(); err != nil {
		log.Fatalf("Error starting application: %v", err)
	}
}
