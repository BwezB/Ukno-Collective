package service

import (
	"github.com/BwezB/Wikno-backend/internal/auth/db"
	"github.com/BwezB/Wikno-backend/internal/auth/model"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	db *db.Database
}

func NewAuthService(database *db.Database) *AuthService {
	return &AuthService{db: database}
}

func (s *AuthService) RegisterUser(req *model.RegisterRequest) (*model.RegisterResponse, error) {
	hashedPassword, err := s.hashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	
	user := model.User{
		Email: req.Email,
		Password: hashedPassword,
	}

	if err := s.db.CreateUser(&user); err != nil {
		return nil, err
	}
	response := model.RegisterResponse{User: user}

	return &response, nil

}

func (s *AuthService) LoginUser(req *model.LoginRequest) (*model.LoginResponse, error) {
	user, err := s.db.GetUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	if err := s.comparePasswords(user.Password, req.Password); err != nil {
		return nil, err 
	}

	response := model.LoginResponse{User: *user}
	return &response, nil
}

func (s *AuthService) hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func (s *AuthService) comparePasswords(hashedPassword string, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return err
	}

	return nil
}
