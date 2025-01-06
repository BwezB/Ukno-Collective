package api

import (
	"context"
	"strings"
	"time"

	c "github.com/BwezB/Wikno-backend/pkg/context"
	m "github.com/BwezB/Wikno-backend/pkg/metrics"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// CONTEXT KEYS

func UnaryRequestIDInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	requestID := uuid.New().String()
	ctx = c.WithRequestID(ctx, requestID)

	return handler(ctx, req)
}

// METRICS

// MetricsInterceptor creates a new unary interceptor for collecting metrics
func MetricsInterceptor(metrics *m.MetricsService) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		method := extractMethodName(info.FullMethod)
		start_time := time.Now()
		// Pre-request metrics:
		metrics.InFlightGauge.WithLabelValues(extractMethodName(method)).Inc() // Increment in-flight gauge
		metrics.RequestCounter.WithLabelValues(method, m.Started).Inc()        // Increment request counter

		// Handle the request
		resp, err := handler(ctx, req)

		duration := time.Since(start_time).Seconds()
		// Post-request metrics:
		metrics.InFlightGauge.WithLabelValues(extractMethodName(method)).Dec() // Decrement in-flight gauge
		metrics.RequestDuration.WithLabelValues(method).Observe(duration)      // Observe request duration
		if err != nil {
			errStatus, _ := status.FromError(err)
			errCode := errStatus.Code().String()

			metrics.RequestCounter.WithLabelValues(method, m.Failed).Inc() // Increment failed request counter
			metrics.ErrorCounter.WithLabelValues(method, errCode).Inc()    // Increment error counter
		} else {
			metrics.RequestCounter.WithLabelValues(method, m.Completed).Inc() // Increment completed request counter
		}

		return resp, err
	}
}

// extractMethodName removes the service prefix from the full method name
// e.g. "/auth.AuthService/Login" becomes "Login"
func extractMethodName(fullMethod string) string {
	parts := strings.Split(fullMethod, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return "unknown"
}
