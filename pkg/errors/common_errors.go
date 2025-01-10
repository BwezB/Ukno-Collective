package errors

var (
	ErrInvalidFunctionArgument = NewErrorType("INVALID_FUNCTION_ARGUMENT", "An invalid function argument was provided")
	ErrInternal = NewErrorType("INTERNAL_ERROR", "An internal error occurred")
	ErrConnectionFailed = NewErrorType("CONNECTION_FAILED", "Failed to connect to the server")
	ErrNotImplemented = NewErrorType("NOT_IMPLEMENTED", "This function is not implemented")
	ErrCancelled = NewErrorType("CANCELLED", "The operation was cancelled")
	ErrInvalidRequest = NewErrorType("INVALID_REQUEST", "Invalid request")
	ErrHealthCheckFailed = NewErrorType("HEALTH_CHECK_FAILED", "Health check failed")
)