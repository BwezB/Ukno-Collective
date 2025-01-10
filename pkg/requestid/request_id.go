package requestid

import (
	"context"

	l "github.com/BwezB/Wikno-backend/pkg/log"

    "github.com/google/uuid"
    "google.golang.org/grpc"
)


// CONTEXT KEY
type contextKey string

// requestIDKey is the key used to store the request ID in the context
const requestIDKey contextKey = "request_id"

// GET/SET REQUEST ID

// WithRequestID adds a request ID to the context
func WithRequestID(ctx context.Context, requestID string) context.Context {
    return context.WithValue(ctx, requestIDKey, requestID)
}

// RequestID returns the request ID from the context
func GetRequestID(ctx context.Context) string {
    if requestID, ok := ctx.Value(requestIDKey).(string); ok {
        return requestID
    }
    l.Error("Request ID not found in context")
    return ""
}


// GRPC INTERCEPTOR



func UnaryRequestIDInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	requestID := uuid.New().String()
	ctx = WithRequestID(ctx, requestID)

	return handler(ctx, req)
}