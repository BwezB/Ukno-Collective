package metrics

import (
	"context"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	l "github.com/BwezB/Wikno-backend/pkg/log"
)

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
