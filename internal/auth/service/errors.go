package service

import (
	"github.com/BwezB/Wikno-backend/pkg/errors"
)

var (
	ErrUserExists = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
)