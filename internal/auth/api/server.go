package api

import (
	"context"
	"net"
	"time"

	"github.com/BwezB/Wikno-backend/internal/auth/model"
	"github.com/BwezB/Wikno-backend/internal/auth/service"

	c "github.com/BwezB/Wikno-backend/pkg/context"
	e "github.com/BwezB/Wikno-backend/pkg/errors"
	h "github.com/BwezB/Wikno-backend/pkg/health"
	l "github.com/BwezB/Wikno-backend/pkg/log"

	pb "github.com/BwezB/Wikno-backend/api/proto/auth"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedAuthServiceServer // Embed the generated server interface
	GrpcServer                        *grpc.Server
	netListener                       net.Listener
	service                           *service.AuthService
	healthService				      *h.HealthService
}

func NewServer(service *service.AuthService, healthService *h.HealthService, config *ServerConfig) (*Server, error) {
	// VALIDATE INPUTS
	if service == nil {
		return nil, e.Wrap("service cannot be nil", ErrInvalidFunctionArgument)
	}
	if healthService == nil {
		return nil, e.Wrap("healthService cannot be nil", ErrInvalidFunctionArgument)
	}
	if config == nil {
		return nil, e.Wrap("config cannot be nil", ErrInvalidFunctionArgument)
	}

	// BUSINESS LOGIC
	// Create the grpc server
	l.Debug("Creating gprc server")
	server := &Server{
		service: service,
		healthService: healthService,
	}

	// Set up the gRPC server
	server.GrpcServer = grpc.NewServer(
		grpc.UnaryInterceptor(UnaryRequestIDInterceptor),
	)

	// Create the health check server
	grpcHealthServer := h.NewGRPCHealthServer(healthService)

	// Register the servers
	pb.RegisterAuthServiceServer(server.GrpcServer, server)
	h.RegisterHealthServer(server.GrpcServer, grpcHealthServer)

	// Set up the listener
	l.Debug("Creating net listener", l.String("address", config.GetAddress()))

	lis, err := net.Listen("tcp", config.GetAddress())
	if err != nil {
		return nil, e.Wrap("failed to create net listener", err)
	}
	server.netListener = lis

	return server, nil
}

func (s *Server) Serve() error {
	// Start the health checks
	l.Debug("Starting health checks")
	go func() {
		s.healthService.Start()
	}()

	// Start the gRPC server
	l.Info("Starting gRPC server", l.String("address", s.netListener.Addr().String()))
	return s.GrpcServer.Serve(s.netListener)
}

func (s *Server) Shutdown(ctx context.Context) error {
	l.Debug("Stopping health checks")
	s.healthService.Stop()

	l.Info("Shutting down gRPC server")
	s.GrpcServer.GracefulStop()
	return nil
}


// AUTH SERVICE FUNCTIONS

func (s *Server) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	l.Debug("Registering user", 
		l.String("email", req.GetEmail()), 
		l.String("request_id", c.GetRequestID(ctx)))

	request := model.RegisterRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	response, err := s.service.RegisterUser(ctx, &request)
	if err != nil {
		l.Warn("Failed to register user:", l.ErrField(err))
		return nil, translateToGrpcError(err)
	}

	res := pb.RegisterResponse{
		UserId: response.User.ID,
		Email:  response.User.Email,
	}

	l.Info("User registration successful",
		l.String("email", response.User.Email),
		l.String("id", response.User.ID),
		l.String("request_id", c.GetRequestID(ctx)))

	return &res, nil
}

func (s *Server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	l.Debug("Logging in user",
		l.String("email", req.GetEmail()),
		l.String("request_id", c.GetRequestID(ctx)))

	request := model.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	response, err := s.service.LoginUser(ctx, &request)
	if err != nil {
		l.Warn("Failed to login user:", l.ErrField(err))
		return nil, translateToGrpcError(err)
	}

	res := pb.LoginResponse{
		UserId: response.User.ID,
		Email:  response.User.Email,
	}

	l.Info("User login successful",
		l.String("email", response.User.Email),
		l.String("id", response.User.ID),
		l.String("request_id", c.GetRequestID(ctx)))
	
	return &res, nil
}


// HEALTH CHECK

func (s *Server) HealthCheck(ctx context.Context) *h.HealthStatus {
	// If the server is responding, it will respond...
	l.Debug("Health check successful")
	return &h.HealthStatus{
		Healthy: true,
		Time:    time.Now(),
	}
}
