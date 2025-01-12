// tests/integration_test.go
package tests

import (
    "context"
    "testing"
    "time"

    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
    "google.golang.org/grpc/status"
    "google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"

    auth "github.com/BwezB/Wikno-backend/api/proto/auth"
    graph "github.com/BwezB/Wikno-backend/api/proto/graph"
)

const (
    auth_host = "192.168.49.2"
    auth_port = "50051"
    graph_host = "192.168.49.2"
    graph_port = "50052"
)

type testClients struct {
    authClient  auth.AuthServiceClient
    graphClient graph.GraphServiceClient
    ctx         context.Context
    cancel      context.CancelFunc
}

func setupClients(t *testing.T) *testClients {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    
    // Connect to auth service
    authConn, err := grpc.Dial(auth_host + ":" + auth_port, grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        t.Fatalf("Could not connect to auth service: %v", err)
    }
    
    // Connect to graph service
    graphConn, err := grpc.Dial(graph_host + ":" + graph_port, grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        t.Fatalf("Could not connect to graph service: %v", err)
    }
    
    return &testClients{
        authClient:  auth.NewAuthServiceClient(authConn),
        graphClient: graph.NewGraphServiceClient(graphConn),
        ctx:         ctx,
        cancel:      cancel,
    }
}

// Test Auth Service

func TestAuthService(t *testing.T) {
    clients := setupClients(t)
    defer clients.cancel()

    // Test successful registration
    t.Run("Register Success", func(t *testing.T) {
        resp, err := clients.authClient.Register(clients.ctx, &auth.AuthRequest{
            Email:    "test@example.com",
            Password: "testpassword123",
        })
        if err != nil {
            t.Fatalf("Registration failed: %v", err)
        }
        if resp.Token == "" {
            t.Error("Expected token in response")
        }
    })

    // Test duplicate registration
    t.Run("Register Duplicate", func(t *testing.T) {
        _, err := clients.authClient.Register(clients.ctx, &auth.AuthRequest{
            Email:    "test@example.com",
            Password: "testpassword123",
        })
        if status.Code(err) != codes.AlreadyExists {
            t.Errorf("Expected AlreadyExists error, got: %v", err)
        }
    })

    // Test successful login
    t.Run("Login Success", func(t *testing.T) {
        resp, err := clients.authClient.Login(clients.ctx, &auth.AuthRequest{
            Email:    "test@example.com",
            Password: "testpassword123",
        })
        if err != nil {
            t.Fatalf("Login failed: %v", err)
        }
        if resp.Token == "" {
            t.Error("Expected token in response")
        }
    })

    // Test invalid login
    t.Run("Login Invalid", func(t *testing.T) {
        _, err := clients.authClient.Login(clients.ctx, &auth.AuthRequest{
            Email:    "test@example.com",
            Password: "wrongpassword",
        })
        if status.Code(err) != codes.Unauthenticated {
            t.Errorf("Expected Unauthenticated error, got: %v", err)
        }
    })

    // Test token verification
    t.Run("Token Verification", func(t *testing.T) {
        // First login to get token
        loginResp, _ := clients.authClient.Login(clients.ctx, &auth.AuthRequest{
            Email:    "test@example.com",
            Password: "testpassword123",
        })

        // Verify token
        verifyResp, err := clients.authClient.VerifyToken(clients.ctx, &auth.VerifyTokenRequest{
            Token: loginResp.Token,
        })
        if err != nil {
            t.Fatalf("Token verification failed: %v", err)
        }
        if verifyResp.Email != "test@example.com" {
            t.Errorf("Expected email test@example.com, got %s", verifyResp.Email)
        }
    })

    // Test invalid token verification
    t.Run("Invalid Token Verification", func(t *testing.T) {
        _, err := clients.authClient.VerifyToken(clients.ctx, &auth.VerifyTokenRequest{
            Token: "invalid.token.here",
        })
        if status.Code(err) != codes.Unauthenticated {
            t.Errorf("Expected Unauthenticated error, got: %v", err)
        }
    })
}

// Test Graph Service

func TestGraphService(t *testing.T) {
    clients := setupClients(t)
    defer clients.cancel()

    // Get authenticated context
    authCtx, _ := getAuthenticatedContext(t, clients)

    // Test context without authentication
    t.Run("Create No Token", func(t *testing.T) {
        _, err := clients.graphClient.CreateEntity(clients.ctx, &graph.EntityRequest{ // Uses the context without token
            Name:       "Test Entity",
            Definition: "Test Definition",
        })
        if status.Code(err) != codes.Unauthenticated {
            t.Errorf("Expected Unauthenticated error, got: %v", err)
        }
    })

    // Test creating entity with auth
	// ID of Test Entity
    var entityID string
    t.Run("Create Entity", func(t *testing.T) {
        entity, err := clients.graphClient.CreateEntity(authCtx, &graph.EntityRequest{
            Name:       "Test Entity",
            Definition: "Test Definition",
        })
        if err != nil {
            t.Fatalf("Entity creation failed: %v", err)
        }
        entityID = entity.EntityId
    })

    // Test getting user data
    t.Run("Get User Data", func(t *testing.T) {
        userData, err := clients.graphClient.GetUserData(authCtx, &graph.Empty{})
        if err != nil {
            t.Fatalf("Getting user data failed: %v", err)
        }
        if len(userData.Entities) == 0 {
            t.Error("Expected at least one entity")
        }
    })

	// Test finding entity by Name
	t.Run("Find Entity By Name", func(t *testing.T) {
		entity, err := clients.graphClient.FindEntities(authCtx, &graph.SearchRequest{
			Name: "Test Entity",
		})
		if err != nil {
			t.Fatalf("Finding entity by name failed: %v", err)
		}
		for _, e := range entity.Entities {
			if e.EntityId == entityID {
				return
			}
		}
		t.Error("Should have found Test Entity, that was created earlier")
	})

    // Test updating entity
    t.Run("Update Entity", func(t *testing.T) {
        _, err := clients.graphClient.UpdateEntity(authCtx, &graph.EntityRequest{
            Id:         entityID,
            Name:       "Updated Entity",
            Definition: "Updated Definition",
        })
        if err != nil {
            t.Fatalf("Entity update failed: %v", err)
        }
    })

    // Test creating connection type
	var connectionTypeID string
    t.Run("Create Connection Type", func(t *testing.T) {
        connType, err := clients.graphClient.CreateConnectionType(authCtx, &graph.ConnectionTypeRequest{
            Name:       "Test Connection",
            Definition: "Test Connection Definition",
        })
        if err != nil {
            t.Fatalf("Connection type creation failed: %v", err)
        }
        if connType.ConnectionTypeId == "" {
            t.Error("Expected connection type ID")
        }
		connectionTypeID = connType.ConnectionTypeId
    })

	// Test finding connection type by Name
	t.Run("Find Connection Type By Name", func(t *testing.T) {
		connType, err := clients.graphClient.FindConnectionTypes(authCtx, &graph.SearchRequest{
			Name: "Test Connection",
		})
		if err != nil {
			t.Fatalf("Finding connection type by name failed: %v", err)
		}
		for _, c := range connType.ConnectionTypes {
			if c.ConnectionTypeId == connectionTypeID {
				return
			}
		}
		t.Error("Should have found Test Connection, that was created earlier")
	})

    // Test creating property type
	var propertyTypeID string
    t.Run("Create Property Type", func(t *testing.T) {
        propType, err := clients.graphClient.CreatePropertyType(authCtx, &graph.PropertyTypeRequest{
            Name:       "Test Property",
            Definition: "Test Property Definition",
            ValueType:  "string",
        })
        if err != nil {
            t.Fatalf("Property type creation failed: %v", err)
        }
        if propType.PropertyTypeId == "" {
            t.Error("Expected property type ID")
        }
		propertyTypeID = propType.PropertyTypeId
    })

	// Test finding property type by Name
	t.Run("Find Property Type By Name", func(t *testing.T) {
		propType, err := clients.graphClient.FindPropertyTypes(authCtx, &graph.SearchRequest{
			Name: "Test Property",
		})
		if err != nil {
			t.Fatalf("Finding property type by name failed: %v", err)
		}
		for _, p := range propType.PropertyTypes {
			if p.PropertyTypeId == propertyTypeID {
				return
			}
		}
		t.Error("Should have found Test Property, that was created earlier")
	})
}

// Helper function to get authenticated context
func getAuthenticatedContext(t *testing.T, clients *testClients) (context.Context, string) {
    resp, err := clients.authClient.Login(clients.ctx, &auth.AuthRequest{
        Email:    "test@example.com",
        Password: "testpassword123",
    })
    if err != nil {
        t.Fatalf("Could not login for authenticated context: %v", err)
    }

    // Create context with token
    md := metadata.NewOutgoingContext(clients.ctx, metadata.Pairs("authorization", resp.Token))
    return md, resp.Token
}