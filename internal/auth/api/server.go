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
	m "github.com/BwezB/Wikno-backend/pkg/metrics"

	pb "github.com/BwezB/Wikno-backend/api/proto/auth"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedAuthServiceServer // Embed the generated server interface
	GrpcServer                        *grpc.Server
	netListener                       net.Listener
	service                           *service.AuthService

	metricsServer *m.MetricsServer
	healthServer  *h.GRPCHealthServer
}

func NewServer(service *service.AuthService,
			   healthService *h.HealthService,
			   metrics *m.MetricsService,
			   config ServerConfig) (*Server, error) {

	// Create a new server
	l.Debug("Creating new server")
	server := &Server{
		service: service,
	}

	// Set up the health server
	l.Debug("Creating health server")
	healthServer := h.NewGRPCHealthServer(healthService)
	server.healthServer = healthServer
	

	// Set up the metrics server
	l.Debug("Creating metrics server")
	metricsServer := m.NewMetricsServer(metrics, config.Metrics)
	server.metricsServer = metricsServer

	// Set up the gRPC server
	l.Debug("Creating gprc server")
	server.GrpcServer = grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			UnaryRequestIDInterceptor,
			MetricsInterceptor(metricsServer.MetricsService),
		),
	)
	pb.RegisterAuthServiceServer(server.GrpcServer, server) // Register auth service server
	h.RegisterHealthServer(server.GrpcServer, healthServer) // Register the health server

	// Set up the listener
	l.Debug("Creating net listener", l.String("address", config.GetAddress()))

	lis, err := net.Listen("tcp", config.GetAddress())
	if err != nil {
		return nil, e.Wrap("failed to create net listener", err)
	}
	server.netListener = lis

	return server, nil
}

func (s *Server) Serve() {
	// Start the health checks
	l.Debug("Starting health checks")
	s.healthServer.Serve()

	// Start the metrics server
	l.Debug("Starting metrics server")
	s.metricsServer.Serve()

	// Start the gRPC server
	l.Info("Starting gRPC server", l.String("address", s.netListener.Addr().String()))
	go func() {
		err := s.GrpcServer.Serve(s.netListener)
		if err != nil {
			l.Error("gRPC server error", l.ErrField(err))
		}
	}()
}

func (s *Server) Shutdown(ctx context.Context) error {
	l.Debug("Stopping health checks")
	s.healthServer.Shutdown()

	l.Debug("Stopping metrics server")
	err := s.metricsServer.Shutdown(ctx)
	if err != nil {
		return e.Wrap("failed to shutdown metrics server", err)
	}

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
	return &h.HealthStatus{
		Healthy: true,
		Time:    time.Now(),
	}
}
