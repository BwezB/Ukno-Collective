// Description: This file contains the common errors used in the project.
package errors

import "errors"

var (
	ErrFlagsNotParsed = errors.New("flags must be parsed before creating configs")
)
