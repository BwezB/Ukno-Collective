package auth

import (
	"context"
	"time"
	
	l "github.com/BwezB/Wikno-backend/pkg/log"
	r "github.com/BwezB/Wikno-backend/pkg/requestid"
	e "github.com/BwezB/Wikno-backend/pkg/errors"
	h "github.com/BwezB/Wikno-backend/pkg/health"

	pb "github.com/BwezB/Wikno-backend/api/proto/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// CONTEXT KEYS
type contextKey string

// userIDKey is the key used to store the user ID in the context
const userIDKey contextKey = "user_id"
// userEmailKey is the key used to store the user email in the context
const userEmailKey contextKey = "user_email"
// authorizationKey is the key used to store the authorization token in the context
const authorizationKey = "authorization"


// GET/SET CONTEXT VALUES

// User id

// WithUserID adds a user ID to the context
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// GetUserID returns the user ID from the context
func GetUserID(ctx context.Context) string {
	if userID, ok := ctx.Value(userIDKey).(string); ok {
		return userID
	}
	return ""
}

// User email

// WithUserEmail adds a user email to the context
func WithUserEmail(ctx context.Context, userEmail string) context.Context {
	return context.WithValue(ctx, userEmailKey, userEmail)
}

// GetUserEmail returns the user email from the context
func GetUserEmail(ctx context.Context) string {
	if userEmail, ok := ctx.Value(userEmailKey).(string); ok {
		return userEmail
	}
	return ""
}

// Authentication token

// WithAuthorizationToken adds an authorization token to the OUTGOING context TO BE SENT OVER GRPC
func WithAuthorizationToken(ctx context.Context, token string) context.Context {
    md := metadata.New(map[string]string{
        authorizationKey: token,
    })
    return metadata.NewOutgoingContext(ctx, md)
}

// GetAuthorizationToken returns the authorization token from the context
func GetAuthorizationToken(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	tokens := md.Get(authorizationKey)
	if len(tokens) == 0 {
		return ""
	}
	return tokens[0]
}

// GRPC INTERCEPTOR

type AuthService struct {
	authClient 	pb.AuthServiceClient
	authHealthClient grpc_health_v1.HealthClient
}

func NewAuthService(config AuthConfig) (*AuthService, error) {
	l.Debug("Connecting to auth service", l.String("address", config.GetAddress()))
	conn, err := grpc.Dial(config.GetAddress(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, e.New("Failed to connect to auth service", e.ErrConnectionFailed, err)
	}

	l.Info("Connected to auth service", l.String("address", config.GetAddress()))

	return &AuthService{
		authClient: pb.NewAuthServiceClient(conn),
		authHealthClient: grpc_health_v1.NewHealthClient(conn),
	}, nil
}

// Health checks

func (s *AuthService) HealthCheck(ctx context.Context) *h.HealthStatus {
	healthResponse, err := s.authHealthClient.Check(ctx, &grpc_health_v1.HealthCheckRequest{})
	if err != nil || healthResponse.Status != grpc_health_v1.HealthCheckResponse_SERVING {
		return &h.HealthStatus{
			Healthy: false,
			Err:     e.New("auth service health check failed", e.ErrHealthCheckFailed, err),
			Time:    time.Now(),
		}
	}
	return &h.HealthStatus{
		Healthy: true,
		Time:    time.Now(),
	}
}

// UnaryAuthInterceptor returns a new unary interceptor that performs token validation
func UnaryAuthInterceptor(authService *AuthService) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if info.FullMethod == "/grpc.health.v1.Health/Check" {
            return handler(ctx, req)
        }

		// Extract token from metadata
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			l.Warn("Missing metadata", l.String("request_id", r.GetRequestID(ctx)))
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		tokens := md.Get(authorizationKey)
		if len(tokens) == 0 {
			l.Warn("Missing authorization token", l.String("request_id", r.GetRequestID(ctx)))
			return nil, status.Error(codes.Unauthenticated, "missing authorization token")
		}

		// Verify token with auth service
		resp, err := authService.authClient.VerifyToken(ctx, &pb.VerifyTokenRequest{
			Token: tokens[0],
		})
		if err != nil {
			l.Warn("Token verification failed",
				l.String("request_id", r.GetRequestID(ctx)),
				l.ErrField(err))
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		// Add user info to context
		ctx = WithUserID(ctx, resp.UserId)
		ctx = WithUserEmail(ctx, resp.Email)

		// Call the handler
		return handler(ctx, req)
	}
}


// GRPC DIAL OPTIONS

func WithToken(token string) grpc.DialOption {
    return grpc.WithUnaryInterceptor(func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
        md := metadata.Pairs("authorization", token)
        ctx = metadata.NewOutgoingContext(ctx, md)
        return invoker(ctx, method, req, reply, cc, opts...)
    })
}