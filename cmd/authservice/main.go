// main for authservice
package main

import (
	"github.com/BwezB/Wikno-backend/internal/auth/api"
	"github.com/BwezB/Wikno-backend/internal/auth/config"
	"github.com/BwezB/Wikno-backend/internal/auth/db"
	"github.com/BwezB/Wikno-backend/internal/auth/service"

	l "github.com/BwezB/Wikno-backend/pkg/log"
	"log" // Using log before logger is initialized
)

func main() {
	// Get configuration
	config, err := config.New()
	if err != nil {
		log.Fatalf("Could not get configuration: %v", err) // Using log, as logger is not yet initialized
	}

	// Set up logging
	l.InitLogger(config.Logger)

	// Connect to the database
	database, err := db.New(&config.Database)
	if err != nil {
		l.Fatal("Could not connect to database:", l.ErrField(err))
	}

	if err := database.AutoMigrate(); err != nil {
		l.Fatal("Could not migrate database:", l.ErrField(err))
	}

	// Set up the service
	authService := service.New(database)
	server, err := api.NewServer(authService, &config.Server)
	if err != nil {
		l.Fatal("Could not create server:", l.ErrField(err))
	}

	// Start the server
	if err := server.Serve(); err != nil {
		l.Fatal("Could not start server:", l.ErrField(err))
	}
}
