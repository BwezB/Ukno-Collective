package log

import (
	"time"

	"go.uber.org/zap"
)

// DebugFunc logs the entry and EXIT!! of a function.
// CALL IT LIKE THIS IN THE BEGINING OF FUNCTION:
// defer log.DebugFunc()()
func DebugFunc(args ...zap.Field) func() {

    l.Debug("Function IN", args...)
    start := time.Now()
    
    return func() {
        duration := time.Since(start)
        l.Debug("Function OUT", zap.Duration("duration", duration))
    }
}
