package metrics

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	l "github.com/BwezB/Wikno-backend/pkg/log"

)


// SERVER

type MetricsServer struct {
	server         *http.Server
	MetricsService *MetricsService
}

func NewMetricsServer(metrics *MetricsService, config MetricsServerConfig) *MetricsServer {
	mux := http.NewServeMux()

	// Register metrics handler with the provided registry
	handler := promhttp.HandlerFor(metrics.registry, promhttp.HandlerOpts{
		EnableOpenMetrics: true,
		Registry:          metrics.registry,
	})

	mux.Handle(config.Path, handler)

	// Create the server
	server := &http.Server{
		Addr:    config.GetAddress(),
		Handler: mux,
	}

	return &MetricsServer{
		server:         server,
		MetricsService: metrics,
	}
}

// Start starts the metrics server in a new goroutine
func (s *MetricsServer) Serve() error {
	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// Log error but don't return it since we're in a goroutine
			l.Error("metrics server error", l.ErrField(err))
		}
	}()

	return nil
}

// Stop gracefully shuts down the metrics server
func (s *MetricsServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

// GetRegistry returns the prometheus registry used by the server
func (s *MetricsServer) GetRegistry() *prometheus.Registry {
	return s.MetricsService.registry
}


// INTERCEPTOR

// MetricsInterceptor creates a new unary interceptor for collecting metrics
func MetricsInterceptor(metrics *MetricsService) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		method := extractMethodName(info.FullMethod)
		start_time := time.Now()
		// Pre-request metrics:
		metrics.InFlightGauge.WithLabelValues(extractMethodName(method)).Inc() // Increment in-flight gauge
		metrics.RequestCounter.WithLabelValues(method, Started).Inc()        // Increment request counter

		// Handle the request
		resp, err := handler(ctx, req)

		duration := time.Since(start_time).Seconds()
		// Post-request metrics:
		metrics.InFlightGauge.WithLabelValues(extractMethodName(method)).Dec() // Decrement in-flight gauge
		metrics.RequestDuration.WithLabelValues(method).Observe(duration)      // Observe request duration
		if err != nil {
			errStatus, _ := status.FromError(err)
			errCode := errStatus.Code().String()

			metrics.RequestCounter.WithLabelValues(method, Failed).Inc() // Increment failed request counter
			metrics.ErrorCounter.WithLabelValues(method, errCode).Inc()    // Increment error counter
		} else {
			metrics.RequestCounter.WithLabelValues(method, Completed).Inc() // Increment completed request counter
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