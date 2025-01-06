// main for testclient
package main

import (
    "context"
    "flag"
    "log"
    "time"

    pb "github.com/BwezB/Wikno-backend/api/proto/auth"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
    "google.golang.org/grpc/health/grpc_health_v1"
)

var (
    addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

func main() {
    flag.Parse()

    // Set up a connection to the server
    conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Fatalf("did not connect: %v", err)
    }
    defer conn.Close()
    
    client := pb.NewAuthServiceClient(conn)

    // Context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    
    // Test Registration
    regResp, err := client.Register(ctx, &pb.AuthRequest{
        Email:    "test@example.com",
        Password: "testpassword123",
    })
    if err != nil {
        log.Printf("Registration failed: %v", err)
    } else {
        log.Printf("Registration successful: %+v", regResp)
    }


    // Test Login
    loginResp, err := client.Login(ctx, &pb.AuthRequest{
        Email:    "test@example.com",
        Password: "testpassword123",
    })
    if err != nil {
        log.Printf("Login failed: %v", err)
    } else {
        log.Printf("Login successful: %+v", loginResp)
    }

    token := loginResp.GetToken()


    // Test Invalid Login
    invalidLoginResp, err := client.Login(ctx, &pb.AuthRequest{
        Email:    "test@example.com",
        Password: "wrongpassword",
    })
    if err != nil {
        log.Printf("Invalid login correctly failed: %v", err)
    } else {
        log.Printf("Unexpected successful login: %+v", invalidLoginResp)
    }


    // Test token validation
    validateResp, err := client.VerifyToken(ctx, &pb.VerifyTokenRequest{
        Token: token,
    })
    if err != nil {
        log.Printf("Token validation failed: %v", err)
    } else {
        log.Printf("Token validation successful: %+v", validateResp)
    }


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