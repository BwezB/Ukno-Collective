package errors

import (
	"fmt"
)

func Error(err error, msg string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", msg, err)
}

// Zdj morm v vsak pkg dat errors.go in poupdejtat
