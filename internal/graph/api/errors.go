package api

import (
	"github.com/BwezB/Wikno-backend/internal/graph/db"
	e "github.com/BwezB/Wikno-backend/pkg/errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
    // ErrInvalidRequest is returned when the request is invalid
	ErrInvalidRequest = e.ErrInvalidRequest
    // ErrInternal is returned when an internal error occurs
    ErrInternal = e.ErrInternal
)

// translateIntoGrpcError translates an error into a gRPC error.
func translateToGrpcError(err error) error {

	if err == nil {
		return nil
	}

	// Get the code of the appropreate error

    var code codes.Code
	var message string
    switch {
    // Database errors
    case e.Is(err, db.ErrRecordNotFound):
        code = codes.NotFound
		message = "Resource not found"
    case e.Is(err, db.ErrDuplicateEntry):
        code = codes.AlreadyExists
		message = "Resource already exists"
    case e.Is(err, db.ErrDatabaseConnection):
        code = codes.Unavailable
		message = "Database connection error"

    // General errors
    case e.Is(err, e.ErrInvalidFunctionArgument):
        code = codes.InvalidArgument
		message = "Invalid function argument"
    case e.Is(err, e.ErrInternal):
        code = codes.Internal
		message = "Internal error"
    
    default:
        code = codes.Unknown
		message = "Unknown error"
    }

	return status.Error(code, message)
}
