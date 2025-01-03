package errors

var (
	ErrInvalidFunctionArgument = NewErrorType("INVALID_FUNCTION_ARGUMENT", "An invalid function argument was provided")
	ErrInternal = NewErrorType("INTERNAL_ERROR", "An internal error occurred")
	ErrNotImplemented = NewErrorType("NOT_IMPLEMENTED", "This function is not implemented")
	ErrCancelled = NewErrorType("CANCELLED", "The operation was cancelled")
)