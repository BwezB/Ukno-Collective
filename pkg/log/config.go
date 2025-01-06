package log

import (
	c "github.com/BwezB/Wikno-backend/pkg/configs"
)

type LoggerConfig struct {
	LoggerEnvironment string `yaml:"environment" validate:"required,oneof=development production"`
	Level             string `yaml:"level" validate:"required,oneof=debug info warning error"`
	Encoding          string `yaml:"encoding" validate:"required,oneof=json console"`
}


// DEFAULTS

func (l *LoggerConfig) SetDefaults() {
	l.LoggerEnvironment = "development"
	l.Level = "debug"
	l.Encoding = "console"
}


// ENV

func (l *LoggerConfig) AddFromEnv() {
	c.SetEnvValue(&l.LoggerEnvironment, "LOGGER_ENVIRONMENT")
	c.SetEnvValue(&l.Level, "LOG_LEVEL")
	c.SetEnvValue(&l.Encoding, "LOG_ENCODING")
}


// FLAGS

var (
	// flagLoggerEnvironment is the flag for the logger environment.
	flagLoggerEnvironment = c.NewFlag("logger_environment", "", "Logger Environment")
	// flagLogLevel is the flag for the log level.
	flagLogLevel = c.NewFlag("log_level", "", "Log Level")
	// flagLogEncoding is the flag for the log encoding.
	flagLogEncoding = c.NewFlag("log_encoding", "", "Log Encoding")
)
func (l *LoggerConfig) AddFromFlags() {
	c.SetFlagValue(&l.LoggerEnvironment, flagLoggerEnvironment)
	c.SetFlagValue(&l.Level, flagLogLevel)
	c.SetFlagValue(&l.Encoding, flagLogEncoding)
}

// FLAGS

