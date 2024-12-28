package main

import (
	"github.com/BwezB/Wikno-backend/internal/auth/api"
	"github.com/BwezB/Wikno-backend/internal/auth/config"
	"github.com/BwezB/Wikno-backend/internal/auth/db"
	"github.com/BwezB/Wikno-backend/internal/auth/service"
	"github.com/BwezB/Wikno-backend/pkg/logger"
)

func main() {

	// Get configuration
	config, err := config.New()
	if err != nil {
		logger.Fatal("Could not get configuration:", err)
	}

	// Set up logging
	logger.SetLevel(config.Logger.Level)

	// Connect to the database
	database, err := db.New(&config.Database)
	if err != nil {
		logger.Fatal("Could not connect to database:", err)
	}

	if err := database.AutoMigrate(); err != nil {
		logger.Fatal("Could not migrate database:", err)
	}

	// Set up the service
	authService := service.New(database)
	server := api.NewServer(authService, &config.Server)

	// Start the server
	if err := server.Serve(); err != nil {
		logger.Fatal("Could not start server:", err)
	}
}
