package errors

import (
	"fmt"
)

// Error wraps an error with a message. If the error is nil, it returns nil.
func Error(err error, msg string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", msg, err)
}
