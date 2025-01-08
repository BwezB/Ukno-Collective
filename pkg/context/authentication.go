package context

import (
	"context"
)

// CONTEXT KEYS

// userIDKey is the key used to store the user ID in the context
const userIDKey contextKey = "user_id"
// userEmailKey is the key used to store the user email in the context
const userEmailKey contextKey = "user_email"


// GET/SET USER ID

// WithUserID adds a user ID to the context
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// GetUserID returns the user ID from the context
func GetUserID(ctx context.Context) string {
	if userID, ok := ctx.Value(userIDKey).(string); ok {
		return userID
	}
	return ""
}

// GET/SET USER EMAIL

// WithUserEmail adds a user email to the context
func WithUserEmail(ctx context.Context, userEmail string) context.Context {
	return context.WithValue(ctx, userEmailKey, userEmail)
}

// GetUserEmail returns the user email from the context
func GetUserEmail(ctx context.Context) string {
	if userEmail, ok := ctx.Value(userEmailKey).(string); ok {
		return userEmail
	}
	return ""
}


// GRPC INTERCEPTOR

// TODO