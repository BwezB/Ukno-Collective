package api

import (
    "context"


	c "github.com/BwezB/Wikno-backend/pkg/context"

    "github.com/google/uuid"
    "google.golang.org/grpc"
)


// CONTEXT KEYS

func UnaryRequestIDInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	requestID := uuid.New().String()
	ctx = c.WithRequestID(ctx, requestID)

	return handler(ctx, req)
}
