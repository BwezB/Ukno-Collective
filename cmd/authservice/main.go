package main

import (
	"log"

	"github.com/BwezB/Wikno-backend/internal/auth/api"
	"github.com/BwezB/Wikno-backend/internal/auth/config"
	"github.com/BwezB/Wikno-backend/internal/auth/db"
	"github.com/BwezB/Wikno-backend/internal/auth/service"
)

func main() {

	// Get configuration
	config, err := config.New()
	if err != nil {
		log.Fatal("Could not get configuration:", err)
	}

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
	server := api.NewServer(authService, &config.Server)

	// Start the server
	if err := server.Serve(); err != nil {
		log.Fatal("Could not start server:", err)
	}
}
