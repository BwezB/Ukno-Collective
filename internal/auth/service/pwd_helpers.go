package service

import (
	"github.com/BwezB/Wikno-backend/pkg/log"
	"golang.org/x/crypto/bcrypt"
)

func (s *AuthService) hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", log.Errore(err, "Password hashing failed")
	}

	return string(hashedPassword), nil
}

func (s *AuthService) comparePasswords(hashedPassword string, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return log.Errore(err, "Password comparison failed")
	}

	return nil
}