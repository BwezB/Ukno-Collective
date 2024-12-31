package log

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Global variables
var (
	l *zap.Logger
	s *zap.SugaredLogger
	config LoggerConfig
)

// InitLogger initializes the logger so it can be used.
func InitLogger(conf LoggerConfig) error {
	// Set config
	config = conf
	// Create appropreate config
	var zapConfig zap.Config
	switch conf.LoggerEnvironment {
	case "development":
		zapConfig = zap.NewDevelopmentConfig()
	case "production":
		zapConfig = zap.NewProductionConfig()
	default:
		return fmt.Errorf("Invalid logger environment: %s", conf.LoggerEnvironment)
	}
	
	// Set the logger level
	switch conf.Level {
	case "debug":
		zapConfig.Level.SetLevel(zapcore.DebugLevel)
	case "info":
		zapConfig.Level.SetLevel(zapcore.InfoLevel)
	case "warn":
		zapConfig.Level.SetLevel(zapcore.WarnLevel)
	case "error":
		zapConfig.Level.SetLevel(zapcore.ErrorLevel)
	default:
		return fmt.Errorf("Invalid logger level: %s", conf.Level)
	}

	zapConfig.Encoding = conf.Encoding
	zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	zapConfig.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	zapConfig.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	zapConfig.EncoderConfig.FunctionKey = "function"

	// Create logger
	var err error
	l, err = zapConfig.Build()
	if err != nil {
		return err
	}

	s = l.Sugar()

	return nil
}

// LOGGING FUNCTIONS
func Debug(msg string, args zap.Field) {
	l.Debug(msg, args)
}
func Debugf(msg string, args ...interface{}) {
	s.Debugf(msg, args...)
}

func Info(msg string, args zap.Field) {
	l.Info(msg, args)
}
func Infof(msg string, args ...interface{}) {
	s.Infof(msg, args...)
}

func Warn(msg string, args zap.Field) {
	l.Warn(msg, args)
}
func Warnf(msg string, args ...interface{}) {
	s.Warnf(msg, args...)
}

func Error(msg string, args zap.Field) {
	l.Error(msg, args)
}
func Errorf(msg string, args ...interface{}) {
	s.Errorf(msg, args...)
}

// OTHER FUNCTIONS
// Sync flushes any buffered log entries
func Sync() {
	l.Sync()
}







// TODO:
// 1. nared error handling
// 2. dodej error handling in logging v app
// 3. preverjanje zdravja
// 4. skaliranje
// 5. sporocilni sistem
// 6. jwt tokeni
// ostala funkcionalnost
