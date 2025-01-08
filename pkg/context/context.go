package context

import (
	"context"
	"time"
    
)


// CONTEXT KEYS

// contextKey is the type that defines the key for the context
type contextKey string


// WRAPPER FUNCTIONS

func WithTimeout(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
    return context.WithTimeout(ctx, timeout)
}

func Background() context.Context {
    return context.Background()
}
