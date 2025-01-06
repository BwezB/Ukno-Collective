package service

import (
	"context"
	"time"

	"github.com/BwezB/Wikno-backend/internal/auth/db"
	"github.com/BwezB/Wikno-backend/internal/auth/model"
	
	e "github.com/BwezB/Wikno-backend/pkg/errors"
	h "github.com/BwezB/Wikno-backend/pkg/health"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	db *db.Database
	config ServiceConfig
	
}

func New(database *db.Database, config ServiceConfig) *AuthService {
	authService := &AuthService{
		db: database, 
		config: config,
	}
	return authService
}

func (s *AuthService) RegisterUser(ctx context.Context, req *model.AuthRequest) (*model.AuthResponse, error) {	
	// Hash the password, so it is not stored in plain text
	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return nil, e.Wrap("RegisterUser failed", err)
	}

	// Create the user object that will be stored in the DB
	user := model.User{
		Email:    req.Email,
		Password: hashedPassword,
	}

	// Create the user in the DB
	if err := s.db.CreateUser(ctx, &user); err != nil {
		return nil, e.Wrap("RegisterUser failed", err)
	}
	
	// Create the jwt token
	token, err := generateJWT(user.ID, user.Email, s.config.jwtSecret, s.config.jwtExpiry)
	if err != nil {
		return nil, e.Wrap("RegisterUser failed", err)
	}

	response := model.AuthResponse{
		User: user,
		Token: token,
	} 

	return &response, nil

}

func (s *AuthService) LoginUser(ctx context.Context, req *model.AuthRequest) (*model.AuthResponse, error) {
	// Get the user from the DB
	user, err := s.db.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, e.Wrap("LoginUser failed", err)
	}

	// Compare the passwords
	if err := comparePasswords(user.Password, req.Password); err != nil {
		return nil, e.Wrap("LoginUser failed", err)
	}

	// Create the jwt token
	token, err := generateJWT(user.ID, user.Email, s.config.jwtSecret, s.config.jwtExpiry)
	if err != nil {
		return nil, e.Wrap("LoginUser failed", err)
	}

	response := model.AuthResponse{
		User: *user,
		Token: token,
	}
	return &response, nil
}

func (s *AuthService) VerifyToken(ctx context.Context, req *model.VerifyTokenRequest) (*model.VerifyTokenResponse, error) {
	// Verify the token
	claims, err := verifyJWT(req.Token, s.config.jwtSecret)
	if err != nil {
		return nil, e.Wrap("VerifyToken failed", err)
	}

	// Get the user from the DB
	user, err := s.db.GetUserByID(ctx, claims.UserID)
	if err != nil {
		return nil, e.Wrap("VerifyToken failed", err)
	}

	response := model.VerifyTokenResponse{
		User: *user,
	}
	return &response, nil
}


// HELPER PASSWORD FUNCTIONS

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
