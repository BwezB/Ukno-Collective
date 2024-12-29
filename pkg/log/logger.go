package log

import (
    "fmt"
    "log"
    "runtime"
    "time"

    e "github.com/BwezB/Wikno-backend/pkg/errors"
)

// logWithLevel has to be called from a specific log function, that is called by the app!
func logMsgWithLevel(level Level, msg string) {
    // Get caller info
    _, file, line, _ := runtime.Caller(2)
    
    // Format timestamp
    timestamp := time.Now().Format("2006-01-02 15:04:05")
    
    // Final log format: timestamp | level | file:line | message
    log.Printf("%s | %-5s | %s:%d | %s\n", timestamp, levelNames[level], file, line, msg)
}

func logWithLevel(level Level, v ...interface{}) {
    msg := fmt.Sprint(v...)
    logMsgWithLevel(level, msg)
}

func logWithLevelf(level Level, format string, v ...interface{}) {
    msg := fmt.Sprintf(format, v...)
    logMsgWithLevel(level, msg)
}


// LEVELS

// Debugf logs debug information with formatting
func Debug(v ...interface{}) {
    if shouldLog(levelDebugVal) {
        logWithLevel(levelDebugVal, v...)
    }
}

// Debug logs debug information
func Debugf(format string, v ...interface{}) {
    if shouldLog(levelDebugVal) {
        logWithLevelf(levelDebugVal, format, v...)
    }
}


// Infof logs general information with formatting
func Info(v ...interface{}) {
    if shouldLog(levelInfoVal) {
        logWithLevel(levelInfoVal, v...)
    }
}

// Info logs general information
func Infof(format string, v ...interface{}) {
    if shouldLog(levelInfoVal) {
        logWithLevelf(levelInfoVal, format, v...)
    }
}


// Warningf logs potential issues with formatting
func Warning(v ...interface{}) {
    if shouldLog(levelWarningVal) {
        logWithLevel(levelWarningVal, v...)
    }
}

// Warning logs potential issues
func Warningf(format string, v ...interface{}) {
    if shouldLog(levelWarningVal) {
        logWithLevelf(levelWarningVal, format, v...)
    }
}


// Errorf logs errors that didnt stop the application with formatting
func Error(v ...interface{}) {
    if shouldLog(levelErrorVal) {
        logWithLevel(levelErrorVal, v...)
    }
}

// Error logs errors that didnt stop the application
func Errorf(format string, v ...interface{}) {
    if shouldLog(levelErrorVal) {
        logWithLevelf(levelErrorVal, format, v...)
    }
}

// ErrorErr logs errors and returns the error
func ErrorErr(err error, v ...interface{}) error{
    msg := fmt.Sprint(v...)
    if shouldLog(levelErrorVal) {
        logMsgWithLevel(levelErrorVal, msg)
    }
    return e.Error(err, msg)
}

// ErrorfErr logs errors with formatting and returns the error
func ErrorfErr(err error, format string, v ...interface{}) error{
    msg := fmt.Sprintf(format, v...)
    if shouldLog(levelErrorVal) {
        logMsgWithLevel(levelErrorVal, msg)
    }
    return e.Error(err, msg)
}


// Fatalf logs the issue and exits with formatting.
func Fatal(v ...interface{}) {
    logWithLevel(levelFatalVal, v...)
    log.Fatal("Application terminated due to fatal error")
}

// Fatal logs the issue and exits.
func Fatalf(format string, v ...interface{}) {
    logWithLevelf(levelFatalVal, format, v...)
    log.Fatal("Application terminated due to fatal error")
}

// TODO: 
// 1. nared error handling
// 2. dodej error handling in logging v app
// 3. preverjanje zdravja
// 4. skaliranje
// 5. sporocilni sistem
// 6. jwt tokeni
// ostala funkcionalnost
