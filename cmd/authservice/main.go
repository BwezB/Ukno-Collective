// main for authservice
package main

import (
	"context"
	"log" // Using log before logger is initialized
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BwezB/Wikno-backend/internal/auth/api"
	"github.com/BwezB/Wikno-backend/internal/auth/config"
	"github.com/BwezB/Wikno-backend/internal/auth/db"
	"github.com/BwezB/Wikno-backend/internal/auth/service"

	"github.com/go-playground/validator/v10"

	g "github.com/BwezB/Wikno-backend/pkg/graph"
	h "github.com/BwezB/Wikno-backend/pkg/health"
	l "github.com/BwezB/Wikno-backend/pkg/log"
	m "github.com/BwezB/Wikno-backend/pkg/metrics"
)

func main() {
	// SETUP

	validator := validator.New()
	// Get configuration
	config, err := config.New(validator)
	if err != nil {
		log.Fatalf("Could not get configuration: %v", err) // Using log, as logger is not yet initialized
	}

	// Set up logging
	l.InitLogger(config.Logger)

	// Create the health checks service to add checks to
	healthService := h.NewHealthService(config.Health)

	// Create the database
	database, err := db.New(config.Database)
	if err != nil {
		l.Fatal("Could not connect to database:", l.ErrField(err))
	}
	healthService.AddCheck(database) // Health check the database
	// Drop tables if needed
	if config.Database.DropTables {
		if config.Environment == "production" {
			l.Error("Cannot drop tables in production!")
		} else {
			if err := database.DropTables(); err != nil {
				l.Fatal("Could not drop tables:", l.ErrField(err))
			}
		}
	}
	// Auto migrate the database
	if err := database.AutoMigrate(); err != nil {
		l.Fatal("Could not migrate database:", l.ErrField(err))
	}

	// Create the graph service
	graphService, err := g.NewGraphService(config.Graph)
	if err != nil {
		l.Fatal("Could not create graph service:", l.ErrField(err))
	}

	// Create the service
	authService, err := service.NewAuthService(database, graphService, config.Service)
	if err != nil {
		l.Fatal("Could not create service:", l.ErrField(err))
	}

	// Create the metrics
	metrics := m.NewMetrics("authservice")

	// Create the server
	server, err := api.NewServer(authService, healthService, metrics, validator, config.Server)
	if err != nil {
		l.Fatal("Could not create server:", l.ErrField(err))
	}

	// START
	// Start server and metrics server
	server.Serve()

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

// TODO:
// 6. ostali servisi
// 7. sporocilni sistem
