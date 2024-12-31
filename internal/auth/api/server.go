package api

import (
	"context"
	"net"

	"github.com/BwezB/Wikno-backend/internal/auth/config"
	"github.com/BwezB/Wikno-backend/internal/auth/model"
	"github.com/BwezB/Wikno-backend/internal/auth/service"
	"github.com/BwezB/Wikno-backend/pkg/log"

	pb "github.com/BwezB/Wikno-backend/api/proto/auth"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedAuthServiceServer // Embed the generated server interface
	grpcServer                        *grpc.Server
	service                           *service.AuthService
	netListener                       net.Listener
}

func NewServer(service *service.AuthService, serverConfig *config.Server) (*Server, error) {
	defer log.DebugFunc("service:", service, "serverConfig:", serverConfig)()

	// VALIDATE INPUTS
	if service == nil {
		return nil, log.Errore(nil, "NewServer failed: service cannot be nil")
	}
	if serverConfig == nil {
		return nil, log.Errore(nil, "NewServer failed: serverConfig cannot be nil")
	}

	// BUSINESS LOGIC
	log.Info("Creating gprc server")
	
	server := &Server{service: service}
	// Set up the gRPC server
	server.grpcServer = grpc.NewServer()
	pb.RegisterAuthServiceServer(server.grpcServer, server)

	// Set up the listener
	log.Info("Creating net listener with address:", serverConfig.GetAddress())

	lis, err := net.Listen("tcp", serverConfig.GetAddress())
	if err != nil {
		return nil, log.Errore(err, "Failed to listen")
	}
	server.netListener = lis

	return server, nil
}

func (s *Server) Serve() error {
	log.Info("Starting gRPC server")
	return s.grpcServer.Serve(s.netListener)
}

func (s *Server) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	defer log.DebugFunc("email:", req.Email)()

	request := model.RegisterRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	response, err := s.service.RegisterUser(&request)
	if err != nil {
		log.Error("Failed to register user:", err)
		return nil, ErrUserExists // Return a custom error for the client
	}

	log.Debugf("User response: %v", response)

	res := pb.RegisterResponse{
		UserId: response.User.ID,
		Email:  response.User.Email,
	}
	return &res, nil
}

func (s *Server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	defer log.DebugFunc("email:", req.Email)()

	request := model.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	response, err := s.service.LoginUser(&request)
	if err != nil {
		if 
		log.Error("Failed to login user:", err)
		return nil, 
	}
	log.Debugf("User response: %v", response)

	res := pb.LoginResponse{
		UserId: response.User.ID,
		Email:  response.User.Email,
	}

	return &res, nil
}
