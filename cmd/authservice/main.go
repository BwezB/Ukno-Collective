// main for authservice
package main

import (
	"context"
	"os"
	"os/signal"
	"log" // Using log before logger is initialized
	"syscall"
	"time"

	"github.com/BwezB/Wikno-backend/internal/auth/api"
	"github.com/BwezB/Wikno-backend/internal/auth/config"
	"github.com/BwezB/Wikno-backend/internal/auth/db"
	"github.com/BwezB/Wikno-backend/internal/auth/service"

	h "github.com/BwezB/Wikno-backend/pkg/health"
	l "github.com/BwezB/Wikno-backend/pkg/log"
)

func main() {
	// SETUP

	// Get configuration
	config, err := config.New()
	if err != nil {
		log.Fatalf("Could not get configuration: %v", err) // Using log, as logger is not yet initialized
	}

	// Set up logging
	l.InitLogger(config.Logger)

	// Create the health checks service
	healthService := h.NewHealthService(config.Health)

	// Connect to the database
	database, err := db.New(&config.Database)
	if err != nil {
		l.Fatal("Could not connect to database:", l.ErrField(err))
	}
	healthService.AddCheck(database) // Health check the database

	if err := database.AutoMigrate(); err != nil {
		l.Fatal("Could not migrate database:", l.ErrField(err))
	}

	// Set up the service
	authService := service.New(database)
	healthService.AddCheck(authService) // Health check the service

	// Set up the server
	server, err := api.NewServer(authService, healthService, &config.Server)
	if err != nil {
		l.Fatal("Could not create server:", l.ErrField(err))
	}
	healthService.AddCheck(server) // Health check the server


	// START

	go func() {
		if err := server.Serve(); err != nil {
			l.Fatal("Could not start server:", l.ErrField(err))
		}
	}()

	
	// SHUTDOWN
	// Create the stop channel
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	// Wait for a stop signal
	<-stop
	l.Info("Shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		l.Fatal("Could not shutdown server:", l.ErrField(err))
	}
}
