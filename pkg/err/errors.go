package err

import (
	"fmt"
)

func Error(msg string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", msg, err)
}

// Zdj morm v vsak pkg dat errors.go in poupdejtat
