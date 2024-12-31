package db

import (
	"github.com/BwezB/Wikno-backend/pkg/errors"
)

var (
	ErrDatabaseConnectionFailed = errors.New("database connection failed")
	ErrDatabaseMigrationFailed  = errors.New("database migration failed")

	ErrUserNotFound = errors.New("user not found")
	ErrUserExists   = errors.New("user already exists")
)