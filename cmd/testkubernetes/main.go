package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	// Connect to your service
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Create health client
	healthClient := grpc_health_v1.NewHealthClient(conn)

	// Check health every 5 seconds
	for {
		// Perform health check
		resp, err := healthClient.Check(context.Background(), &grpc_health_v1.HealthCheckRequest{})
		if err != nil {
			log.Printf("❌ Health check failed: %v", err)
		} else {
			switch resp.Status {
			case grpc_health_v1.HealthCheckResponse_SERVING:
				log.Printf("✅ Service is healthy")
			case grpc_health_v1.HealthCheckResponse_NOT_SERVING:
				log.Printf("❌ Service is unhealthy")
			default:
				log.Printf("❓ Service status unknown")
			}
		}

		time.Sleep(5 * time.Second)
	}
}