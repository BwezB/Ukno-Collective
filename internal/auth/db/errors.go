// internal/auth/db/errors.go
package db

import (
	"strings"

	e "github.com/BwezB/Wikno-backend/pkg/errors"
	"gorm.io/gorm"
)

var (
    // Common errors
    ErrInternal           = e.ErrInternal
	// Common database error types
	ErrDatabaseConnection = e.NewErrorType("DB_CONNECTION_ERROR", "database connection error")
	ErrDuplicateEntry     = e.NewErrorType("DB_DUPLICATE_ENTRY", "resource already exists")
	ErrRecordNotFound     = e.NewErrorType("DB_NOT_FOUND", "resource not found")
)

// TranslateDatabaseError converts GORM and postgres errors into internal application errors
func TranslateDatabaseError(err error) error {
	if err == nil {
		return nil
	}

	// Check for GORM-specific errors
	if e.Is(err, gorm.ErrRecordNotFound) {
		return e.New("", ErrRecordNotFound, err)
	}

	// Check for Postgres-specific errors
	errMsg := err.Error()
	switch {
	case strings.Contains(errMsg, "connection refused"):
		return e.New("Database connection failed", ErrDatabaseConnection, err)

	case strings.Contains(errMsg, "duplicate key"):
		return e.New("Resource already exists", ErrDuplicateEntry, err)

	case strings.Contains(errMsg, "deadlock"):
		return e.New("Database transaction conflict", ErrInternal, err)

	case strings.Contains(errMsg, "foreign key"):
		return e.New("Database constraint violation", ErrInternal, err)
	}

	// Default case - wrap as internal error
	return e.New("Unexpected database error", ErrInternal, err)
}
