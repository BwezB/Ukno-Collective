package log

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Global variables
var (
	zapLogger      *zap.Logger
	zapSugarLogger *zap.SugaredLogger
	//config LoggerConfig
)

// InitLogger initializes the logger so it can be used.
func InitLogger(conf LoggerConfig) error {
	// Set config
	//config = conf

	// Create appropreate config
	var zapConfig zap.Config
	switch conf.LoggerEnvironment {
	case "development":
		zapConfig = zap.NewDevelopmentConfig()
	case "production":
		zapConfig = zap.NewProductionConfig()
	default:
		return fmt.Errorf("invalid logger environment: %s", conf.LoggerEnvironment)
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
		return fmt.Errorf("invalid logger level: %s", conf.Level)
	}

	// General zapConfig settings
	zapConfig.Encoding = conf.Encoding

	// Specific settings for development and production
	if conf.LoggerEnvironment == "development" { // development
		zapConfig.Development = true
		zapConfig.DisableStacktrace = true
		zapConfig.DisableCaller = false

		timeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("15:04:05.000"))
		}

		zapConfig.EncoderConfig = zapcore.EncoderConfig{
			MessageKey:     "M",
			LevelKey: 	    "L",
			TimeKey:        "T",
			//NameKey: 	    "N",
			CallerKey:  	"C",
			//FunctionKey:  	"F",
			//StacktraceKey:  "S",
			EncodeLevel:    zapcore.CapitalColorLevelEncoder,
			EncodeTime:     timeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,

		}
	
	} else if conf.LoggerEnvironment == "production" { // production
		zapConfig.DisableStacktrace = false
		zapConfig.DisableCaller = true

		zapConfig.EncoderConfig = zapcore.EncoderConfig{
			MessageKey:     "M",
			LevelKey: 	    "L",
			TimeKey:        "T",
			NameKey: 	    "N",
			CallerKey:  	"C",
			FunctionKey:  	"F",
			StacktraceKey:  "S",
			EncodeLevel:    zapcore.CapitalColorLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}
	} else {
		return fmt.Errorf("invalid logger environment: %s", conf.LoggerEnvironment)
	}

	// Create logger
	var err error
	zapLogger, err = zapConfig.Build(zap.AddCallerSkip(1)) // Skip 1 because of the wrapper function
	if err != nil {
		return err
	}

	zapSugarLogger = zapLogger.Sugar()

	return nil
}

// DebugFunc logs the entry and EXIT!! of a function.
// CALL IT LIKE THIS IN THE BEGINING OF FUNCTION:
// defer log.DebugFunc()()
func DebugFunc(funcName string, args ...zap.Field) func() {
	funcName = "[ " + funcName + " ] "
	zapLogger.Debug(funcName + "IN", args...)
	start := time.Now()

	return func() {
		duration := time.Since(start)
		zapLogger.Debug(funcName + "OUT", zap.Duration("duration", duration))
	}
}

// TODO:
// 3. preverjanje zdravja
// 4. skaliranje
// 5. sporocilni sistem
// 6. jwt tokeni
// ostala funkcionalnost
