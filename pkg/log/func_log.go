package log

import (
	"runtime"
	"strings"
	"time"
)

// DebugFunc logs the entry and EXIT!! of a function.
// CALL IT LIKE THIS IN THE BEGINING OF FUNCTION:
// defer log.DebugFunc()()
func DebugFunc(v ...interface{}) func() {
    if !shouldLog(levelDebugVal) { // Skip logging if log level higher than debug
        return func() {}
    }

    pc, _, _, _ := runtime.Caller(1)  // Skip 1 frame to get the caller
    nameParts := strings.Split(runtime.FuncForPC(pc).Name(), "/")
    funcName := nameParts[len(nameParts)-1]

    msg_in := "[Function] IN " + funcName
    if len(v) != 0 {
        msg_in += " with arguments: " + formatMessage(v...)
    }

    logMsgWithLevel(levelDebugVal, msg_in)
    start := time.Now()
    
    return func() {
        duration := time.Since(start)
        logMsgWithLevel(levelDebugVal, "[Function] OUT " + funcName + " after " + duration.String())
    }
}
