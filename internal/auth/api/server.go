package api

import (
	"context"
	"github.com/BwezB/Wikno-backend/internal/auth/service"
	"github.com/BwezB/Wikno-backend/internal/auth/model"
	pb "github.com/BwezB/Wikno-backend/api/proto/auth"
)

type AuthServer struct {
	pb.UnimplementedAuthServiceServer // Embed the generated server interface
	service *service.AuthService
}

func NewAuthServer(service *service.AuthService) *AuthServer {
	return &AuthServer{service: service}
}

func (s *AuthServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
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
		Email: response.User.Email,
	}

	return &res, nil
}

func (s *AuthServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
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
		Email: response.User.Email,
	}

	return &res, nil
}
