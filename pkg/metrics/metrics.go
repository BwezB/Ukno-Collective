package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// EXPORTED CONSTANTS
const (
	Started   = "started"
	Completed = "completed"
	Failed    = "failed"
)

// METRICS

type MetricsService struct {
	registry *prometheus.Registry
	// RequestCounter is a counter for the total number of total gRPC requests, completed requests and failed requests. Use m.Started, m.Completed, m.Failed as status labels
	RequestCounter *prometheus.CounterVec
	// ErrorCounter is a counter for the gRPC errors. Use the gRPC status codes as status labels
	ErrorCounter *prometheus.CounterVec
	// InFlightGauge is a gauge for the number of gRPC requests in flight for each method
	InFlightGauge *prometheus.GaugeVec
	// RequestDuration is a histogram for the duration of gRPC requests for each method
	RequestDuration *prometheus.HistogramVec
}

func NewMetrics(namespace string) *MetricsService {
	// Create a new Metrics struct
	metrics := &MetricsService{
		registry: prometheus.NewRegistry(),
		RequestCounter: promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "grpc_requests_total",
			Help:      "Total number of gRPC requests",
		}, []string{"method", "status"}),
		ErrorCounter: promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "grpc_errors_total",
			Help:      "Total number of gRPC errors",
		}, []string{"method", "status"}),
		InFlightGauge: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "grpc_requests_in_flight",
			Help:      "Number of gRPC requests in flight",
		}, []string{"method"}),
		RequestDuration: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "grpc_request_duration_seconds",
			Help:      "Duration of gRPC requests in seconds",
			Buckets:   prometheus.DefBuckets,
		}, []string{"method"}),
	}

	// Register the metrics with the default registry
	metrics.Register(metrics.RequestCounter, metrics.ErrorCounter, metrics.RequestDuration)
	return metrics
}

// Register registers the given collectors with the metrics registry
func (m *MetricsService) Register(cs ...prometheus.Collector) {
	m.registry.MustRegister(cs...)
}

func (m *MetricsService) GetRegistry() *prometheus.Registry {
	return m.registry
}
