// Description: This package is used by all microservices to log errors and info messages.
package logging

import (
	"go.uber.org/zap"
)

type Logger struct {
	logger *zap.Logger
}

func New() *Logger {
	// Import config for logging and so on
}
