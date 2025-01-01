package service

import (
	"github.com/BwezB/Wikno-backend/internal/auth/db"
	"github.com/BwezB/Wikno-backend/internal/auth/model"

	e "github.com/BwezB/Wikno-backend/pkg/errors"
	l "github.com/BwezB/Wikno-backend/pkg/log"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	db *db.Database
}

func New(database *db.Database) *AuthService {
	return &AuthService{db: database}
}

func (s *AuthService) RegisterUser(req *model.RegisterRequest) (*model.RegisterResponse, error) {
	defer l.DebugFunc("RegisterUser (service)")()

	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return nil, e.Wrap("RegisterUser failed", err)
	}

	user := model.User{
		Email:    req.Email,
		Password: hashedPassword,
	}

	if err := s.db.CreateUser(&user); err != nil {
		return nil, e.Wrap("RegisterUser failed", err)
	}
	response := model.RegisterResponse{User: user}

	l.Info("User registration successful", 
		l.String("email", user.Email), 
		l.String("id", user.ID))

	return &response, nil

}

func (s *AuthService) LoginUser(req *model.LoginRequest) (*model.LoginResponse, error) {
	defer l.DebugFunc("LoginUser (service)")()

	user, err := s.db.GetUserByEmail(req.Email)
	if err != nil {
		return nil, e.Wrap("LoginUser failed", err)
	}

	if err := comparePasswords(user.Password, req.Password); err != nil {
		return nil, e.Wrap("LoginUser failed", err)
	}

	l.Info("User login successful", l.String("email", user.Email), l.String("id", user.ID))

	response := model.LoginResponse{User: *user}
	return &response, nil
}

// HELPER FUNCTIONS

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", translateBycriptError(err)
	}

	return string(hashedPassword), nil
}

func comparePasswords(hashedPassword string, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return translateBycriptError(err)
	}

	return nil
}
