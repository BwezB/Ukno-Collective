package api

import (
	"context"
	"net"

	"github.com/BwezB/Wikno-backend/internal/auth/config"
	"github.com/BwezB/Wikno-backend/internal/auth/model"
	"github.com/BwezB/Wikno-backend/internal/auth/service"
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

func NewServer(service *service.AuthService, serverConfig *config.Server) (*Server, error) {
	l.Info("Creating new server")
	server := &Server{service: service}

	// Set up the gRPC server
	server.grpcServer = grpc.NewServer()
	pb.RegisterAuthServiceServer(server.grpcServer, server)

	// Set up the listener
	lis, err := net.Listen("tcp", serverConfig.GetAddress())
	if err != nil {
		return nil, l.ErrorErr(err, "failed to listen")
	}
	server.netListener = lis
	
	return server, nil
}

func (s *Server) Serve() error {
	return s.grpcServer.Serve(s.netListener)
}

func (s *Server) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	request := model.RegisterRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	response, err := s.service.RegisterUser(&request)
	if err != nil {
		return nil, err
	}

	res := pb.RegisterResponse{
		UserId: response.User.ID,
		Email:  response.User.Email,
	}

	return &res, nil
}

func (s *Server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	request := model.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	response, err := s.service.LoginUser(&request)
	if err != nil {
		return nil, err
	}

	res := pb.LoginResponse{
		UserId: response.User.ID,
		Email:  response.User.Email,
	}

	return &res, nil
}
