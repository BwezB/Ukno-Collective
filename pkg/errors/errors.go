package errors

import (
	"fmt"
	"errors"
)

// Error struct
type appError struct {
	Code string
	Msg  string
	Err  error
}

func (e *appError) Error() string {
	if e.Err == nil {
		return e.Msg
	}
	return fmt.Sprintf("%s: %s", e.Msg, e.Err.Error())
}

func (e *appError) Unwrap() error {
	return e.Err
}

func (e *appError) Is(target error) bool {
	t, ok := target.(*appError)
	if !ok {
		return false
	}
	return e.Code == t.Code
}


// FUNCTIONS FOR CREATING ERRORS

// NewErrorType returns a new error type with the given code and a default message.
func NewErrorType(code, default_msg string) *appError {
	return &appError{
		Code: code,
		Msg:  default_msg,
		Err:  nil,
	}
}

// New returns a new checkable error.
// External errors should be wrapped with this function.
// If message is empty, the default message of the error type is used.
func New(msg string, errorType *appError, err error) error {
	if errorType == nil {
		errorType = NewErrorType("UNKNOWN_ERROR", "An unknown error occurred")
	}
	if msg == "" {
		msg = errorType.Msg
	}
	return &appError{
		Code: errorType.Code,
		Msg:  msg,
		Err:  err,
	}
}

// Wrap adds context to existing errors.
func Wrap(msg string, err error) error {
	return &appError{
		Msg: msg,
		Err: err,
	}
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}
