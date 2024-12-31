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
    regResp, err := client.Register(ctx, &pb.RegisterRequest{
        Email:    "test@example.com",
        Password: "testpassword123",
    })
    if err != nil {
        log.Printf("Registration failed: %v", err)
    } else {
        log.Printf("Registration successful: %+v", regResp)
    }

    // Test Login
    loginResp, err := client.Login(ctx, &pb.LoginRequest{
        Email:    "test@example.com",
        Password: "testpassword123",
    })
    if err != nil {
        log.Printf("Login failed: %v", err)
    } else {
        log.Printf("Login successful: %+v", loginResp)
    }

    // Test Invalid Login
    invalidLoginResp, err := client.Login(ctx, &pb.LoginRequest{
        Email:    "test@example.com",
        Password: "wrongpassword",
    })
    if err != nil {
        log.Printf("Invalid login correctly failed: %v", err)
    } else {
        log.Printf("Unexpected successful login: %+v", invalidLoginResp)
    }
}