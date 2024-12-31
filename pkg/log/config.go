package log

import (
	"github.com/BwezB/Wikno-backend/pkg/configs"
)

type LoggerConfig struct {
	LoggerEnvironment string `yaml:"environment" validate:"required,oneof=development production"`
	Level string `yaml:"level" validate:"required,oneof=debug info warning error"`
	Encoding string `yaml:"encoding" validate:"required,oneof=json console"`
}

func (l *LoggerConfig) SetDefaults() {
	l.LoggerEnvironment = defaultLoggerEnvironment
	l.Level = defaultLogLevel
	l.Encoding = defaultLogEncoding
}

func (l *LoggerConfig) AddFromEnv() {
	configs.SetEnvValue(&l.LoggerEnvironment, envLoggerEnvironment)
	configs.SetEnvValue(&l.Level, envLogLevel)
	configs.SetEnvValue(&l.Encoding, envLogEncoding)
}

func (l *LoggerConfig) AddFromFlags() {
	configs.SetFlagValue(&l.LoggerEnvironment, flagLoggerEnvironment)
	configs.SetFlagValue(&l.Level, flagLogLevel)
	configs.SetFlagValue(&l.Encoding, flagLogEncoding)
}

// DEFAULTS
const (
	// defaultLogEncoding is the default value for the log encoding.
	defaultLoggerEnvironment = "development"
	// defaultLogLevel is the default value for the log level.
	defaultLogLevel = "debug"
	// defaultLogEncoding is the default value for the log encoding.
	defaultLogEncoding = "console"
)

// ENV
const (
	// envLoggerEnvironment is the default value for the logger environment.
	envLoggerEnvironment = "LOGGER_ENVIRONMENT"
	// envLogLevel is the environment variable for the log level.
	envLogLevel = "LOG_LEVEL"
	// envLogEncoding is the environment variable for the log encoding.
	envLogEncoding = "LOG_ENCODING"
)

// FLAGS
var (
	// flagLoggerEnvironment is the flag for the logger environment.
	flagLoggerEnvironment = configs.NewFlag("logger_environment", "", "Logger Environment")
	// flagLogLevel is the flag for the log level.
	flagLogLevel = configs.NewFlag("log_level", "", "Log Level")
	// flagLogEncoding is the flag for the log encoding.
	flagLogEncoding = configs.NewFlag("log_encoding", "", "Log Encoding")
)