package context

import (
	"context"
	"time"

	l "github.com/BwezB/Wikno-backend/pkg/log"
)


// CONTEXT KEYS
type contextKey string

// RequestIDKey is the key used to store the request ID in the context
const RequestIDKey contextKey = "request_id"


// REQUEST ID

// WithRequestID adds a request ID to the context
func WithRequestID(ctx context.Context, requestID string) context.Context {
    return context.WithValue(ctx, RequestIDKey, requestID)
}

// RequestID returns the request ID from the context
func GetRequestID(ctx context.Context) string {
    if requestID, ok := ctx.Value(RequestIDKey).(string); ok {
        return requestID
    }
    l.Error("Request ID not found in context")
    return ""
}


// WRAPPER FUNCTIONS

func WithTimeout(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
    return context.WithTimeout(ctx, timeout)
}

func Background() context.Context {
    return context.Background()
}