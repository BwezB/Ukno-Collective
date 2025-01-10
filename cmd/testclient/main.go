package main

import (
    "context"
    "flag"
    "log"
    "time"
    //"os"

    auth "github.com/BwezB/Wikno-backend/api/proto/auth"
    graph "github.com/BwezB/Wikno-backend/api/proto/graph"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
    "google.golang.org/grpc/health/grpc_health_v1"
    "google.golang.org/grpc/metadata"
)

var (
    authAddr = flag.String("auth-addr", "localhost:50051", "the address of auth service")
    graphAddr = flag.String("graph-addr", "localhost:50052", "the address of graph service")
)

// withToken creates an interceptor that adds the token to all requests
func withToken(token string) grpc.DialOption {
    return grpc.WithUnaryInterceptor(func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
        md := metadata.Pairs("authorization", token)
        ctx = metadata.NewOutgoingContext(ctx, md)
        return invoker(ctx, method, req, reply, cc, opts...)
    })
}

func testAuthService(ctx context.Context, client auth.AuthServiceClient) (string, error) {
    log.Println("🔑 Testing Auth Service...")

    // Test Registration
    regResp, err := client.Register(ctx, &auth.AuthRequest{
        Email:    "test@example.com",
        Password: "testpassword123",
    })
    if err != nil {
        log.Printf("❌ Registration failed: %v", err)
    } else {
        log.Printf("✅ Registration successful: %+v", regResp)
    }

    // Test Login
    loginResp, err := client.Login(ctx, &auth.AuthRequest{
        Email:    "test@example.com",
        Password: "testpassword123",
    })
    if err != nil {
        log.Printf("❌ Login failed: %v", err)
        return "", err
    }
    log.Printf("✅ Login successful: %+v", loginResp)

    token := loginResp.GetToken()

    // Test Invalid Login
    _, err = client.Login(ctx, &auth.AuthRequest{
        Email:    "test@example.com",
        Password: "wrongpassword",
    })
    if err != nil {
        log.Printf("✅ Invalid login correctly failed: %v", err)
    } else {
        log.Printf("❌ Unexpected successful login with wrong password")
    }

    // Test token validation
    validateResp, err := client.VerifyToken(ctx, &auth.VerifyTokenRequest{
        Token: token,
    })
    if err != nil {
        log.Printf("❌ Token validation failed: %v", err)
        return "", err
    }
    log.Printf("✅ Token validation successful: %+v", validateResp)

    return token, nil
}

func testGraphService(ctx context.Context, client graph.GraphServiceClient) error {
    log.Println("📊 Testing Graph Service...")

    // Test entity operations
    entity, err := client.CreateEntity(ctx, &graph.EntityRequest{
        Name:       "Test Entity",
        Definition: "This is a test entity",
    })
    if err != nil {
        log.Printf("❌ Failed to create entity: %v", err)
        return err
    }
    log.Printf("✅ Successfully created entity: %+v", entity)

    // Test finding entities
    entities, err := client.FindEntities(ctx, &graph.SearchRequest{
        Name: "Test Entity",
    })
    if err != nil {
        log.Printf("❌ Failed to find entities: %v", err)
        return err
    }
    log.Printf("✅ Successfully found entities: %+v", entities)

    // Test connection type operations
    connType, err := client.CreateConnectionType(ctx, &graph.ConnectionTypeRequest{
        Name:       "Test Connection",
        Definition: "This is a test connection type",
    })
    if err != nil {
        log.Printf("❌ Failed to create connection type: %v", err)
        return err
    }
    log.Printf("✅ Successfully created connection type: %+v", connType)

    // Test property type operations
    propType, err := client.CreatePropertyType(ctx, &graph.PropertyTypeRequest{
        Name:       "Test Property",
        Definition: "This is a test property type",
        ValueType:  "string",
    })
    if err != nil {
        log.Printf("❌ Failed to create property type: %v", err)
        return err
    }
    log.Printf("✅ Successfully created property type: %+v", propType)

    // Get all user data
    userData, err := client.GetUserData(ctx, &graph.Empty{})
    if err != nil {
        log.Printf("❌ Failed to get user data: %v", err)
        return err
    }
    log.Printf("✅ Successfully retrieved user data: %+v", userData)

    return nil
}

func main() {
    // SETUP
    flag.Parse()

    // Set up connection to auth service
    authConn, err := grpc.Dial(*authAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Fatalf("did not connect to auth service: %v", err)
    }
    defer authConn.Close()
    
    authClient := auth.NewAuthServiceClient(authConn)
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()


    // TESTING

    // Test auth service and get token
    token, err := testAuthService(ctx, authClient)
    if err != nil {
        log.Fatalf("auth service test failed: %v", err)
    }

    // os.Exit(0)
    
    // Set up connection to graph service with token
    graphConn, err := grpc.Dial(*graphAddr, 
        grpc.WithTransportCredentials(insecure.NewCredentials()),
        withToken(token))
    if err != nil {
        log.Fatalf("did not connect to graph service: %v", err)
    }
    defer graphConn.Close()

    graphClient := graph.NewGraphServiceClient(graphConn)

    // Test graph service
    if err := testGraphService(ctx, graphClient); err != nil {
        log.Fatalf("graph service test failed: %v", err)
    }

    // Create health clients
    authHealthClient := grpc_health_v1.NewHealthClient(authConn)
    graphHealthClient := grpc_health_v1.NewHealthClient(graphConn)

    // Check health every 5 seconds
    for {
        // Check auth service health
        authResp, err := authHealthClient.Check(context.Background(), &grpc_health_v1.HealthCheckRequest{})
        if err != nil {
            log.Printf("❌ Auth service health check failed: %v", err)
        } else {
            log.Printf("✅ Auth service status: %v", authResp.Status)
        }

        // Check graph service health
        graphResp, err := graphHealthClient.Check(context.Background(), &grpc_health_v1.HealthCheckRequest{})
        if err != nil {
            log.Printf("❌ Graph service health check failed: %v", err)
        } else {
            log.Printf("✅ Graph service status: %v", graphResp.Status)
        }

        time.Sleep(5 * time.Second)
    }
}