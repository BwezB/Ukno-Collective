package service

import (
	"github.com/BwezB/Wikno-backend/internal/auth/db"
	"github.com/BwezB/Wikno-backend/internal/auth/model"
	"github.com/BwezB/Wikno-backend/pkg/log"
)

type AuthService struct {
	db *db.Database
}

func New(database *db.Database) *AuthService {
	return &AuthService{db: database}
}

func (s *AuthService) RegisterUser(req *model.RegisterRequest) (*model.RegisterResponse, error) {
	defer log.DebugFunc("email:", req.Email)()

	hashedPassword, err := s.hashPassword(req.Password)
	if err != nil {
		return nil, log.Errore(err, "Password hashing for RegisterUser failed")
	}
	
	user := model.User{
		Email: req.Email,
		Password: hashedPassword,
	}

	if err := s.db.CreateUser(&user); err != nil {
		return nil, log.Errore(err, "CreateUser failed")
	}
	response := model.RegisterResponse{User: user}

	log.Info("User with id:", user.ID, "registered with email:", user.Email)

	return &response, nil

}

func (s *AuthService) LoginUser(req *model.LoginRequest) (*model.LoginResponse, error) {
	defer log.DebugFunc("email:", req.Email)()

	user, err := s.db.GetUserByEmail(req.Email)
	if err != nil {
		if err == db.ErrUserNotFound {
			log.Error("User with email ", req.Email, " not found: ", err) // Tko morm popravt stvari, da errorji niso kr usepovsod i think? Mejbi?
			return nil, ErrUserNotFound
		}
		return nil, log.Errore(err, "Couldnt get user by email")
	}

	if err := s.comparePasswords(user.Password, req.Password); err != nil {
		return nil, log.Errore(err, "Password comparison failed") 
	}

	log.Info("User with id:", user.ID, "logged in with email:", user.Email)

	response := model.LoginResponse{User: *user}
	return &response, nil
}
