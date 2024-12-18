package main

import (
	"log"
	"net"

	"github.com/BwezB/Wikno-backend/internal/auth/api"
	"github.com/BwezB/Wikno-backend/internal/auth/db"
	"github.com/BwezB/Wikno-backend/internal/auth/config"
	"github.com/BwezB/Wikno-backend/internal/auth/service"
	"google.golang.org/grpc"
)

func main() {

	// Get configuration
	config, err := config.NewConfig()
	if err != nil {
		log.Fatal("Could not get configuration:", err)
	}

	database, err := db.NewConnection(&config.Database)
	if err != nil {
		log.Fatal("Could not connect to database:", err)
	}

	// Run migrations
	if err := database.AutoMigrate(); err != nil {
		log.Fatal("Could not migrate database:", err)
	}

	// Set up my stuff
	authService := service.NewAuthService(database)

	authServer := api.NewAuthServer(authService)

	// Set up the gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, &authServer)

	lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

	// Start the server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("Could not start server:", err)
	}


}
