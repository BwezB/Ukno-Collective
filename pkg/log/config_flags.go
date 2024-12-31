package log

import (
	"github.com/BwezB/Wikno-backend/pkg/configs"
)

var (
	// flagLoggerEnvironment is the flag for the logger environment.
	flagLoggerEnvironment = configs.NewFlag("logger_environment", "", "Logger Environment")
	// flagLogLevel is the flag for the log level.
	flagLogLevel = configs.NewFlag("log_level", "", "Log Level")
)
