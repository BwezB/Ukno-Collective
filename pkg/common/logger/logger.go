package logger

import (
    "fmt"
    "log"
    "runtime"
    "time"
)

// Levels as constants for consistency
const (
    levelDebug = "DEBUG"
    levelInfo  = "INFO"
    levelError = "ERROR"
	levelFatal = "FATAL"
)

// logWithLevel adds timestamp, level, and caller info
func logWithLevel(level, format string, v ...interface{}) {
    // Get caller info
    _, file, line, _ := runtime.Caller(2)
    
    // Format timestamp
    timestamp := time.Now().Format("2006-01-02 15:04:05")
    
    // Create message with all context
    message := fmt.Sprintf(format, v...)
    
    // Final log format: timestamp | level | file:line | message
    log.Printf("%s | %-5s | %s:%d | %s", timestamp, level, file, line, message)
}

// Exported functions
func Debug(format string, v ...interface{}) {
    logWithLevel(levelDebug, format, v...)
}

func Info(format string, v ...interface{}) {
    logWithLevel(levelInfo, format, v...)
}

func Error(format string, v ...interface{}) {
    logWithLevel(levelError, format, v...)
}

// Fatal logs and exits
func Fatal(format string, v ...interface{}) {
    logWithLevel(levelFatal, format, v...)
    log.Fatal("Application terminated due to fatal error")
}
// TODO: 
// 0. Vrz vn flag in env definitione v seperite file
// 1. nared error handling
// 2. dodej error handling in logging v app
// 3. preverjanje zdravja
// 4. skaliranje
// 5. sporocilni sistem
// 6. jwt tokeni
// ostala funkcionalnost
