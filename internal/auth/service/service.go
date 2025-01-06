package service

import (
	"context"
	"time"

	"github.com/BwezB/Wikno-backend/internal/auth/db"
	"github.com/BwezB/Wikno-backend/internal/auth/model"

	c "github.com/BwezB/Wikno-backend/pkg/context"
	e "github.com/BwezB/Wikno-backend/pkg/errors"
	h "github.com/BwezB/Wikno-backend/pkg/health"
	l "github.com/BwezB/Wikno-backend/pkg/log"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	db *db.Database
}

func New(database *db.Database) *AuthService {
	return &AuthService{db: database}
}

func (s *AuthService) RegisterUser(ctx context.Context, req *model.RegisterRequest) (*model.RegisterResponse, error) {
	l.Debug("Registering user",
		l.String("email", req.Email),
		l.String("request_id", c.GetRequestID(ctx)))
	
	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return nil, e.Wrap("RegisterUser failed", err)
	}

	user := model.User{
		Email:    req.Email,
		Password: hashedPassword,
	}

	if err := s.db.CreateUser(ctx, &user); err != nil {
		return nil, e.Wrap("RegisterUser failed", err)
	}
	response := model.RegisterResponse{User: user}

	return &response, nil

}

func (s *AuthService) LoginUser(ctx context.Context, req *model.LoginRequest) (*model.LoginResponse, error) {
	l.Debug("Logging in user",
		l.String("email", req.Email),
		l.String("request_id", c.GetRequestID(ctx)))
	
	user, err := s.db.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, e.Wrap("LoginUser failed", err)
	}

	if err := comparePasswords(user.Password, req.Password); err != nil {
		return nil, e.Wrap("LoginUser failed", err)
	}

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


// HEALTH CHECK

func (s *AuthService) HealthCheck(ctx context.Context) *h.HealthStatus {
	// If the service is responding, it will respond...
	_, err := hashPassword("health_check")
	if err != nil {
		return &h.HealthStatus{
			Healthy: false, 
			Time: time.Now(),
			Err: err,
		}
	}
	return &h.HealthStatus{Healthy: true}
}
