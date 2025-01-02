package errors

import (
	"fmt"
	"errors"
)

// Error struct
type AppError struct {
	Code string
	Msg  string
	Err  error
}

func (e *AppError) Error() string {
	if e.Err == nil {
		return e.Msg
	}
	return fmt.Sprintf("%s: %s", e.Msg, e.Err.Error())
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func (e *AppError) Is(target error) bool {
	t, ok := target.(*AppError)
	if !ok {
		return false
	}
	return e.Code == t.Code
}


// FUNCTIONS FOR CREATING ERRORS

// NewErrorType returns a new error type with the given code and a default message.
func NewErrorType(code, default_msg string) *AppError {
	return &AppError{
		Code: code,
		Msg:  default_msg,
		Err:  nil,
	}
}

// New returns a new checkable error.
// External errors should be wrapped with this function.
// If message is empty, the default message of the error type is used.
func New(msg string, errorType *AppError, err error) error {
	if errorType == nil {
		errorType = NewErrorType("UNKNOWN_ERROR", "An unknown error occurred")
	}
	if msg == "" {
		msg = errorType.Msg
	}
	return &AppError{
		Code: errorType.Code,
		Msg:  msg,
		Err:  err,
	}
}

// Wrap adds context to existing errors.
func Wrap(msg string, err error) error {
	return &AppError{
		Msg: msg,
		Err: err,
	}
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func As(err error, target interface{}) bool {
	return errors.As(err, target)
}