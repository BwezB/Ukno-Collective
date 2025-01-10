package service

import (
	"context"

	"github.com/BwezB/Wikno-backend/internal/auth/db"
	"github.com/BwezB/Wikno-backend/internal/auth/model"

	e "github.com/BwezB/Wikno-backend/pkg/errors"
	g "github.com/BwezB/Wikno-backend/pkg/graph"
	r "github.com/BwezB/Wikno-backend/pkg/requestid"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	db     *db.Database
	graph  *g.GraphService
	config ServiceConfig
	id     string // My id for calling other services
	email  string // My email for calling other services
	token  string // My token for calling other services
}

func NewAuthService(database *db.Database, graph *g.GraphService, config ServiceConfig) (*AuthService, error) {
	// Hash the password for auth user, so it is not stored in plain text
	hashedPassword, err := hashPassword(config.password)
	if err != nil {
		return nil, e.Wrap("RegisterUser failed", err)
	}

	// Create the user in the DB if it does not exist
	user, err := database.CreateUser(r.WithRequestID(context.Background(), "0"), &model.AuthRequest{
		Email:    config.email,
		Password: config.password,
	}, hashedPassword)
	if err != nil && !e.Is(err, db.ErrDuplicateEntry) {
		return nil, e.Wrap("RegisterUser failed", err)
	}

	// JWT will be created with the first request

	authService := &AuthService{
		db:     database,
		config: config,
		graph:  graph,
		id:     user.ID,
		email:  user.Email,
	}

	return authService, nil
}

func (s *AuthService) RegisterUser(ctx context.Context, req *model.AuthRequest) (*model.AuthResponse, error) {
	// Hash the password, so it is not stored in plain text
	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return nil, e.Wrap("RegisterUser failed", err)
	}

	// Create the user in the DB
	user, err := s.db.CreateUser(ctx, req, hashedPassword)
	if err != nil {
		return nil, e.Wrap("RegisterUser failed", err)
	}

	// Create user in GRAPH SERVICE
	myToken, err := s.getJWToken()
	if err != nil {
		return nil, e.Wrap("Couldnt get my JWT token", err)
	}
	err = s.graph.CreateUser(user.ID, myToken)
	if err != nil {
		return nil, e.Wrap("RegisterUser failed", err)
	}

	// Create the jwt token
	token, err := generateJWT(user.ID, user.Email, s.config.jwtSecret, s.config.jwtExpiry)
	if err != nil {
		return nil, e.Wrap("RegisterUser failed", err)
	}

	response := model.AuthResponse{
		User:  *user,
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
		User:  *user,
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

// Helper functions

func (s *AuthService) getJWToken() (string, error) {
	if s.token == "" {
		// Create the jwt token
		token, err := generateJWT(s.id, s.email, s.config.jwtSecret, s.config.jwtExpiry)
		if err != nil {
			return "", e.Wrap("Couldnt get jwt token", err)
		}
		s.token = token
	}

	return s.token, nil
}

// HELPER PASSWORDFUNCTIONS

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
