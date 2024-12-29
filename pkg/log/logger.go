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

// formatMessage formats the message with the given format and values
func formatMessage(format string, v ...interface{}) string {
    if len(v) == 0 {
        return format
    }
    return fmt.Sprintf(format, v...)
}


// LEVELS
// Debug logs debug information
func Debug(format string, v ...interface{}) {
    if shouldLog(levelDebugVal) {
        msg := formatMessage(format, v...)
        logMsgWithLevel(levelDebugVal, msg)
    }
}

// Info logs general information
func Info(format string, v ...interface{}) {
    if shouldLog(levelInfoVal) {
        msg := formatMessage(format, v...)
        logMsgWithLevel(levelInfoVal, msg)
    }
}

// Warning logs potential issues
func Warning(format string, v ...interface{}) {
    if shouldLog(levelWarningVal) {
        msg := formatMessage(format, v...)
        logMsgWithLevel(levelWarningVal, msg)
    }
}

// Error logs errors that didnt stop the application
func Error(format string, v ...interface{}) {
    if shouldLog(levelErrorVal) {
        msg := formatMessage(format, v...)
        logMsgWithLevel(levelErrorVal, msg)
    }
}

// ErrorErr logs errors that didnt stop the application, and returns the error
func ErrorErr(err error, format string, v ...interface{}) error{
    msg := formatMessage(format, v...)
    if shouldLog(levelErrorVal) {
        logMsgWithLevel(levelErrorVal, msg)
    }

    return e.Error(err, msg)
}

// Fatal logs the issue and exits.
func Fatal(format string, v ...interface{}) {
    msg := formatMessage(format, v...)
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
