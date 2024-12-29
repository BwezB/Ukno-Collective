package log

import (
    "strings"
    "fmt"
)

type Level int

const (
    LevelDebug = "DEBUG"
    LevelInfo  = "INFO"
    LevelWarning = "WARNING"
    LevelError = "ERROR"
    LevelFatal = "FATAL"
)

const (
    levelDebugVal Level = iota // 0
    levelInfoVal              // 1
    levelWarningVal           // 2
    levelErrorVal             // 3
    levelFatalVal            // 4
)

var levelValues = map[string]Level{
    LevelDebug: levelDebugVal,
    LevelInfo:  levelInfoVal,
    LevelWarning: levelWarningVal,
    LevelError: levelErrorVal,
    LevelFatal: levelFatalVal,
}

var levelNames = map[Level]string{
    levelDebugVal: LevelDebug,
    levelInfoVal:  LevelInfo,
    levelWarningVal: LevelWarning,
    levelErrorVal: LevelError,
    levelFatalVal: LevelFatal,
}



// FUNCTIONS FOR CONTROLLING LEVELS
// Default level is INFO
var currentLevel Level = levelDebugVal

// SetLevel allows changing the logging level
func SetLevel[L Level | string](level L) error{
    switch lvl := any(level).(type) {
    case Level:
        currentLevel = lvl
        return nil
    case string:
        lvlUp := strings.ToUpper(lvl)
        if lvlVal, ok := levelValues[lvlUp]; ok {
            currentLevel = lvlVal
            return nil
        }
    }

    // If we reach this point, the level is invalid
    Errorf("Invalid log level: %v", level)
    return fmt.Errorf("invalid log level: %v", level)
}

// shouldLog checks if we should log this level
func shouldLog(level Level) bool {
    return level >= currentLevel
}
