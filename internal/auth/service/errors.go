package service

import (
	e "github.com/BwezB/Wikno-backend/pkg/errors"

	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrInvalidPassword is returned when the password is invalid
	ErrInvalidPassword = e.NewErrorType("INVALID_PASSWORD", "Invalid password")
	ErrServiceInternal = e.NewErrorType("SERVICE_INTERNAL", "Service internal error")
)

// translateBycriptError translates a hash error to a service error
func translateBycriptError(err error) error {
	switch err {
	case bcrypt.ErrMismatchedHashAndPassword:
		return e.New("Password mismatch", ErrInvalidPassword, err)
	case bcrypt.ErrPasswordTooLong:
		return e.New("Password too long", ErrServiceInternal, err)
	case bcrypt.ErrHashTooShort:
		return e.New("Hash too short", ErrServiceInternal, err)
	default:
		return ErrServiceInternal
	}
}
