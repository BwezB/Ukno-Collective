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
}

func (l *LoggerConfig) AddFromEnv() {
	configs.SetEnvValue(&l.Level, envLogLevel)
}

func (l *LoggerConfig) AddFromFlags() {
	configs.SetFlagValue(&l.Level, flagLogLevel)
}
