package log

import (
	"fmt"
	"log"
	"runtime"
	"strings"
	"time"

	e "github.com/BwezB/Wikno-backend/pkg/errors"
)

var counter int = 0
func init() {
	counter++
	fmt.Println("Logger init", counter)
	log.SetFlags(0) // Remove timestamp from log messages
}

// logWithLevel has to be called from a specific log function, that is called by the app!
func logMsgWithLevel(level Level, msg string) {
	// Get caller info
	_, file, line, _ := runtime.Caller(2)

	// Format timestamp
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	// Final log format: timestamp | level | file:line | message
	log.Printf("%s | %s | %s\t [%s:%d]\n", timestamp, levelNames[level], msg, file, line)
}

func formatMessage(v ...interface{}) string {
    parts := make([]string, len(v))
    for i, item := range v {
        parts[i] = fmt.Sprint(item)
    }
    return strings.Join(parts, " ")
}

func formatMessagef(format string, v ...interface{}) string {
	return fmt.Sprintf(format, v...) // added \n so next message is in new line
}

// LEVELS
// Debug logs debug information, joined with spaces.
func Debug(v ...interface{}) {
	if shouldLog(levelDebugVal) {
		msg := formatMessage(v...)
		logMsgWithLevel(levelDebugVal, msg)
	}
}
// Debug logs debug information, formatted.
func Debugf(format string, v ...interface{}) {
	if shouldLog(levelDebugVal) {
		msg := formatMessagef(format, v...)
		logMsgWithLevel(levelDebugVal, msg)
	}
}

// Info logs general information, joined with spaces.
func Info(v ...interface{}) {
	if shouldLog(levelInfoVal) {
		msg := formatMessage(v...)
		logMsgWithLevel(levelInfoVal, msg)
	}
}
// Info logs general information
func Infof(format string, v ...interface{}) {
	if shouldLog(levelInfoVal) {
		msg := formatMessagef(format, v...)
		logMsgWithLevel(levelInfoVal, msg)
	}
}

// Warning logs potential issues, joined with spaces.
func Warning(v ...interface{}) {
	if shouldLog(levelWarningVal) {
		msg := formatMessage(v...)
		logMsgWithLevel(levelWarningVal, msg)
	}
}
// Warning logs potential issues
func Warningf(format string, v ...interface{}) {
	if shouldLog(levelWarningVal) {
		msg := formatMessagef(format, v...)
		logMsgWithLevel(levelWarningVal, msg)
	}
}

// Error logs errors that didnt stop the application
func Error(v ...interface{}) {
	if shouldLog(levelErrorVal) {
		msg := formatMessage(v...)
		logMsgWithLevel(levelErrorVal, msg)
	}
}
// Error logs errors that didnt stop the application
func Errorf(format string, v ...interface{}) {
	if shouldLog(levelErrorVal) {
		msg := formatMessagef(format, v...)
		logMsgWithLevel(levelErrorVal, msg)
	}
}

// Errore logs errors that didnt stop the application, and returns the error
func Errore(err error, v ...interface{}) error {
	msg := formatMessage(v...)
	if shouldLog(levelErrorVal) {
		logMsgWithLevel(levelErrorVal, msg)
	}
	msg = strings.ToLower(msg)
	return e.Error(err, msg)
}
// Errore logs errors that didnt stop the application, and returns the error
func Errorfe(err error, format string, v ...interface{}) error {
	msg := formatMessagef(format, v...)
	if shouldLog(levelErrorVal) {
		logMsgWithLevel(levelErrorVal, msg)
	}
	msg = strings.ToLower(msg)
	return e.Error(err, msg)
}

// Fatal logs the issue and exits.
func Fatal(v ...interface{}) {
	msg := formatMessage(v...)
	logMsgWithLevel(levelFatalVal, msg)
	log.Fatal("Application terminated due to fatal error") // Terminate the application
}
// Fatal logs the issue and exits.
func Fatalf(format string, v ...interface{}) {
	msg := formatMessagef(format, v...)
	logMsgWithLevel(levelFatalVal, msg)
	log.Fatal("Application terminated due to fatal error") // Terminate the application
}

// TODO:
// 1. nared error handling
// 2. dodej error handling in logging v app
// 3. preverjanje zdravja
// 4. skaliranje
// 5. sporocilni sistem
// 6. jwt tokeni
// ostala funkcionalnost
