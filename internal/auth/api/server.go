package api

import (
	"context"
	"net"
	
	"github.com/BwezB/Wikno-backend/internal/auth/model"
	"github.com/BwezB/Wikno-backend/internal/auth/service"

	e "github.com/BwezB/Wikno-backend/pkg/errors"
	l "github.com/BwezB/Wikno-backend/pkg/log"

	pb "github.com/BwezB/Wikno-backend/api/proto/auth"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedAuthServiceServer // Embed the generated server interface
	grpcServer                        *grpc.Server
	service                           *service.AuthService
	netListener                       net.Listener
}

func NewServer(service *service.AuthService, config *ServerConfig) (*Server, error) {
	defer l.DebugFunc("NewServer (api)")() 

	// VALIDATE INPUTS
	if service == nil {
		return nil, e.Wrap("service cannot be nil", ErrInvalidFunctionArgument)
	}
	if config == nil {
		return nil, e.Wrap("config cannot be nil", ErrInvalidFunctionArgument)
	}

	// BUSINESS LOGIC
	l.Info("Creating gprc server")
	
	server := &Server{service: service}
	// Set up the gRPC server
	server.grpcServer = grpc.NewServer()
	pb.RegisterAuthServiceServer(server.grpcServer, server)

	// Set up the listener
	l.Info("Creating net listener", l.String("address", config.GetAddress()))

	lis, err := net.Listen("tcp", config.GetAddress())
	if err != nil {
		return nil, e.Wrap("failed to create net listener", err)
	}
	server.netListener = lis

	return server, nil
}

func (s *Server) Serve() error {
	l.Info("Starting gRPC server")
	return s.grpcServer.Serve(s.netListener)
}

func (s *Server) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	defer l.DebugFunc("Server.Register (api)", l.String("request email", req.GetEmail()))()

	request := model.RegisterRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	response, err := s.service.RegisterUser(&request)
	if err != nil {
		l.Warn("Failed to register user:", l.ErrField(err))
		return nil, err
	}

	l.Debug("User response:", l.String("user id", response.User.ID), l.String("email", response.User.Email))

	res := pb.RegisterResponse{
		UserId: response.User.ID,
		Email:  response.User.Email,
	}
	return &res, nil
}

func (s *Server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	defer l.DebugFunc("Server.Login (api)")()

	request := model.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	response, err := s.service.LoginUser(&request)
	if err != nil {
		l.Warn("Failed to login user:", l.ErrField(err))
		return nil, err
	}
	l.Debugf("User response: %v", response)

	res := pb.LoginResponse{
		UserId: response.User.ID,
		Email:  response.User.Email,
	}

	return &res, nil
}
