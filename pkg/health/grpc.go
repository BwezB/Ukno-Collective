package health

import (
    "context"

    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/health/grpc_health_v1"
    "google.golang.org/grpc/status"
	"google.golang.org/grpc"
)

type GRPCHealthServer struct {
	grpc_health_v1.UnimplementedHealthServer
	healthService *HealthService
}

func NewGRPCHealthServer(hs *HealthService) *GRPCHealthServer {
	healthServer := &GRPCHealthServer{
        healthService: hs,
    }
	return healthServer
}

func (hs *GRPCHealthServer) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	select {
	case <-ctx.Done():
		return nil, status.Error(codes.Canceled, "context canceled")
	default:
		if !hs.healthService.IsHealthy() {
            return &grpc_health_v1.HealthCheckResponse{
                Status: grpc_health_v1.HealthCheckResponse_NOT_SERVING,
            }, nil
        }
        return &grpc_health_v1.HealthCheckResponse{
            Status: grpc_health_v1.HealthCheckResponse_SERVING,
        }, nil
	}
}

func (hs *GRPCHealthServer) Watch(req *grpc_health_v1.HealthCheckRequest, srv grpc_health_v1.Health_WatchServer) error {
	return status.Error(codes.Unimplemented, "Watch is not implemented")
}

func (hs *GRPCHealthServer) Serve() {
	// Start the health checking service in a new goroutine
	go func () {
		hs.healthService.Start()
	}()
}

func (hs *GRPCHealthServer) Shutdown() {
	hs.healthService.Stop()
}

// HELPER FUNCTIONS

func RegisterHealthServer(grpcServer *grpc.Server, healthServer *GRPCHealthServer) {
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
}
