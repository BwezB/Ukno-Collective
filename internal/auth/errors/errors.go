package errors

import (
	"errors"
)

// Database caused errors
var (
	ErrUserNotFound = errors.New("user not found")
	ErrUserExists   = errors.New("user already exists")
	ErrInvalidPassword = errors.New("invalid password")
	ErrFailedToConnectToDatabase = errors.New("failed to connect to database")
)
