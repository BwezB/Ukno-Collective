package log

import (
	"time"
)

func EnterFunc(name string) func() {
    Debug("Entering function: %s", name)
    start := time.Now()
    
    return func() {
        duration := time.Since(start)
        Debug("Exiting function: %s (took: %v)", name, duration)
    }
}