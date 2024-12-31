// main for authservice
package main

import (
	"github.com/BwezB/Wikno-backend/internal/auth/api"
	"github.com/BwezB/Wikno-backend/internal/auth/config"
	"github.com/BwezB/Wikno-backend/internal/auth/db"
	"github.com/BwezB/Wikno-backend/internal/auth/service"
	"github.com/BwezB/Wikno-backend/pkg/log"
)

func main() {
	// Get configuration
	config, err := config.New()
	if err != nil {
		log.Fatal("Could not get configuration:", err)
	}

	// Set up logging
	log.SetLevel(config.Logger.Level)

	// Connect to the database
	database, err := db.New(&config.Database)
	if err != nil {
		log.Fatal("Could not connect to database:", err)
	}

	if err := database.AutoMigrate(); err != nil {
		log.Fatal("Could not migrate database:", err)
	}

	// Set up the service
	authService := service.New(database)
	server, err := api.NewServer(authService, &config.Server)
	if err != nil {
		log.Fatal("Could not create server:", err)
	}

	// Start the server
	if err := server.Serve(); err != nil {
		log.Fatal("Server serving faliure:", err)
	}
}
