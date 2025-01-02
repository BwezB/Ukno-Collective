package errors

var (
	ErrInvalidFunctionArgument = NewErrorType("INVALID_FUNCTION_ARGUMENT", "An invalid function argument was provided")
	ErrInternal = NewErrorType("INTERNAL_ERROR", "An internal error occurred")
)