package graph

import (
	"context"
	"time"

	e "github.com/BwezB/Wikno-backend/pkg/errors"
	l "github.com/BwezB/Wikno-backend/pkg/log"
	a "github.com/BwezB/Wikno-backend/pkg/auth"
	h "github.com/BwezB/Wikno-backend/pkg/health"

	pb "github.com/BwezB/Wikno-backend/api/proto/graph"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GraphService struct {
	graphClient pb.GraphServiceClient
}

func NewGraphService(config GraphConfig) (*GraphService, error) {
	l.Debug("Connecting to graph service", l.String("address", config.GetAddress()))
	conn, err := grpc.Dial(config.GetAddress(), grpc.WithTransportCredentials(
		insecure.NewCredentials()),
	)
	if err != nil {
		return nil, e.New("Failed to connect to graph service", e.ErrConnectionFailed, err)
	}

	l.Info("Connected to graph service", l.String("address", config.GetAddress()))

	return &GraphService{
		graphClient: pb.NewGraphServiceClient(conn),
	}, nil
}

func (gs *GraphService) CreateUser(id, token string) error {
	l.Debug("Creating user in graph service", l.String("id", id))

	ctx := a.WithAuthorizationToken(context.Background(), token)

	_, err := gs.graphClient.CreateUser(ctx, &pb.UserRequest{
		Id: id,
	})
	if err != nil {
		return e.Wrap("CreateUser failed", err)
	}

	l.Info("Created user in graph service", l.String("id", id))

	return nil
}

// HealthCheck

func (gs *GraphService) HealthCheck(ctx context.Context) *h.HealthStatus {
	_, err := gs.graphClient.Ping(ctx, &pb.PingRequest{})
	if err != nil {
		return &h.HealthStatus{
			Healthy: false,
			Err: err,
			Time: time.Now(),
		}
	}
	return &h.HealthStatus{
		Healthy: true,
		Time: time.Now(),
	}
}
