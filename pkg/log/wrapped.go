package log

import (
	"time"

	"go.uber.org/zap"
)

// LOGGING FUNCTIONS
// Debug: Detailed information for debugging
func Debug(msg string, args ...zap.Field) {
	zapLogger.Debug(msg, args...)
}
func Debugf(msg string, args ...interface{}) {
	zapSugarLogger.Debugf(msg, args...)
}

// Info: Notable expected events
func Info(msg string, args ...zap.Field) {
	zapLogger.Info(msg, args...)
}

// Warn: Handled issues
func Warn(msg string, args ...zap.Field) {
	zapLogger.Warn(msg, args...)
}

// Error: Unhandled issues
func Error(msg string, args ...zap.Field) {
	zapLogger.Error(msg, args...)
}

// Fatal: Unrecoverable issues
func Fatal(msg string, args ...zap.Field) {
	zapLogger.Sync()
	zapLogger.Fatal(msg, args...)
}

// String calls zap.String function
func String(key string, val string) zap.Field {
	return zap.String(key, val)
}

// Int calls zap.Int function
func Int(key string, val int) zap.Field {
	return zap.Int(key, val)
}

// Bool calls zap.Bool function
func Bool(key string, val bool) zap.Field {
	return zap.Bool(key, val)
}

// Duration calls zap.Duration function
func Duration(key string, val time.Duration) zap.Field {
	return zap.Duration(key, val)
}

// ErrField: Calls zap.Error function
func ErrField(err error) zap.Field {
	return zap.Error(err)
}

// Any calls zap.Any function
func Any(key string, val interface{}) zap.Field {
	return zap.Any(key, val)
}
