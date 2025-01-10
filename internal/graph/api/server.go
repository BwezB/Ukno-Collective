package api

import (
	"context"
	"net"

	"github.com/BwezB/Wikno-backend/internal/graph/model"
	"github.com/BwezB/Wikno-backend/internal/graph/service"
	"github.com/go-playground/validator/v10"

	a "github.com/BwezB/Wikno-backend/pkg/auth"
	r "github.com/BwezB/Wikno-backend/pkg/requestid"
	e "github.com/BwezB/Wikno-backend/pkg/errors"
	h "github.com/BwezB/Wikno-backend/pkg/health"
	l "github.com/BwezB/Wikno-backend/pkg/log"
	m "github.com/BwezB/Wikno-backend/pkg/metrics"

	pb "github.com/BwezB/Wikno-backend/api/proto/graph"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedGraphServiceServer
	GrpcServer    *grpc.Server
	netListener   net.Listener
	service       *service.GraphService
	validator     *validator.Validate
	metricsServer *m.MetricsServer
	healthServer  *h.GRPCHealthServer
}

func NewServer(service *service.GraphService,
	healthService *h.HealthService,
	metrics *m.MetricsService,
	authService *a.AuthService,
	validator *validator.Validate,
	config ServerConfig) (*Server, error) {

	l.Debug("Creating new server")
	server := &Server{
		service:   service,
		validator: validator,
	}

	l.Debug("Creating health server")
	healthServer := h.NewGRPCHealthServer(healthService)
	server.healthServer = healthServer

	l.Debug("Creating metrics server")
	metricsServer := m.NewMetricsServer(metrics, config.Metrics)
	server.metricsServer = metricsServer

	l.Debug("Creating grpc server")
	server.GrpcServer = grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			r.UnaryRequestIDInterceptor,
			m.MetricsInterceptor(metricsServer.MetricsService),
			a.UnaryAuthInterceptor(authService),
		),
	)
	pb.RegisterGraphServiceServer(server.GrpcServer, server)
	h.RegisterHealthServer(server.GrpcServer, healthServer)

	l.Debug("Creating net listener", l.String("address", config.GetAddress()))
	lis, err := net.Listen("tcp", config.GetAddress())
	if err != nil {
		return nil, e.Wrap("failed to create net listener", err)
	}
	server.netListener = lis

	return server, nil
}

func (s *Server) Serve() {
	l.Debug("Starting health checks")
	s.healthServer.Serve()

	l.Debug("Starting metrics server")
	s.metricsServer.Serve()

	l.Info("Starting gRPC server", l.String("address", s.netListener.Addr().String()))
	go func() {
		err := s.GrpcServer.Serve(s.netListener)
		if err != nil {
			l.Error("gRPC server error", l.ErrField(err))
		}
	}()
}

func (s *Server) Shutdown(ctx context.Context) error {
	l.Debug("Stopping health checks")
	s.healthServer.Shutdown()

	l.Debug("Stopping metrics server")
	err := s.metricsServer.Shutdown(ctx)
	if err != nil {
		return e.Wrap("failed to shutdown metrics server", err)
	}

	l.Info("Shutting down gRPC server")
	s.GrpcServer.GracefulStop()
	return nil
}


// GRAPH SERVICE METHODS

// Users

func (s *Server) CreateUser(ctx context.Context, req *pb.UserRequest) (*pb.Empty, error) {
	l.Debug("Creating user", 
		l.String("request_id", r.GetRequestID(ctx)),
		l.String("user_id", req.GetId()))
	
	// Translate request
	userReq := &model.UserRequest{
		ID: req.GetId(),
	}

	// Validate request
	if err := s.validator.Struct(userReq); err != nil {
		return nil, e.New("Request validation failed", ErrInvalidRequest, err)
	}

	// Create user
	err := s.service.CreateUser(ctx, userReq)
	if err != nil {
		l.Warn("Failed to create user:", l.ErrField(err))
		return nil, translateToGrpcError(err)
	}

	return &pb.Empty{}, nil
}

func (s *Server) GetUserData(ctx context.Context, _ *pb.Empty) (*pb.UserData, error) {
	l.Debug("Getting user data", l.String("request_id", r.GetRequestID(ctx)))

	// Get user data from service
	userData, err := s.service.GetUserData(ctx)
	if err != nil {
		l.Warn("Failed to get user data:", l.ErrField(err))
		return nil, translateToGrpcError(err)
	}

	// Validate the response
	if err := s.validator.Struct(userData); err != nil {
		return nil, e.New("Response validation failed", ErrInternal, err)
	}

	// Translate to protobuf response
	response := &pb.UserData{
		Entities:        translateEntitiesToProto(userData.Entities),
		ConnectionTypes: translateConnectionTypesToProto(userData.ConnectionTypes),
		PropertyTypes:   translatePropertyTypesToProto(userData.PropertyTypes),
	}

	return response, nil
}

// Entities

func (s *Server) CreateEntity(ctx context.Context, req *pb.EntityRequest) (*pb.UsersEntity, error) {
	l.Debug("Creating entity",
		l.String("name", req.GetName()),
		l.String("request_id", r.GetRequestID(ctx)))

	// Translate request
	entityReq := &model.EntityRequest{
		ID:         req.GetId(),
		Name:       req.GetName(),
		Definition: req.GetDefinition(),
	}

	// Validate request
	if err := s.validator.Struct(entityReq); err != nil {
		return nil, e.New("Request validation failed", ErrInvalidRequest, err)
	}

	// Create entity
	entity, err := s.service.CreateEntity(ctx, entityReq)
	if err != nil {
		l.Warn("Failed to create entity:", l.ErrField(err))
		return nil, translateToGrpcError(err)
	}

	// Validate response
	if err := s.validator.Struct(entity); err != nil {
		return nil, e.New("Response validation failed", ErrInternal, err)
	}

	// Translate to protobuf response
	response := translateEntityToProto(entity)

	return response, nil
}

func (s *Server) UpdateEntity(ctx context.Context, req *pb.EntityRequest) (*pb.Empty, error) {
	l.Debug("Updating entity",
		l.String("id", req.GetId()),
		l.String("name", req.GetName()),
		l.String("request_id", r.GetRequestID(ctx)))

	// Translate request
	entityReq := &model.EntityRequest{
		ID:         req.GetId(),
		Name:       req.GetName(),
		Definition: req.GetDefinition(),
	}

	// Validate request
	if err := s.validator.Struct(entityReq); err != nil {
		return nil, e.New("Request validation failed", ErrInvalidRequest, err)
	}

	// Update entity
	err := s.service.UpdateEntity(ctx, entityReq)
	if err != nil {
		l.Warn("Failed to update entity:", l.ErrField(err))
		return nil, translateToGrpcError(err)
	}

	return &pb.Empty{}, nil
}

func (s *Server) FindEntities(ctx context.Context, req *pb.SearchRequest) (*pb.EntitiesList, error) {
	l.Debug("Finding entities",
		l.String("name", req.GetName()),
		l.String("request_id", r.GetRequestID(ctx)))

	// Translate request
	searchReq := &model.SearchRequest{
		Name: req.GetName(),
	}

	// Validate request
	if err := s.validator.Struct(searchReq); err != nil {
		return nil, e.New("Request validation failed", ErrInvalidRequest, err)
	}

	// Find entities
	entities, err := s.service.FindEntities(ctx, searchReq)
	if err != nil {
		l.Warn("Failed to find entities:", l.ErrField(err))
		return nil, translateToGrpcError(err)
	}

	// Translate to protobuf response
	response := &pb.EntitiesList{
		Entities: translateEntitiesToProto(entities),
	}

	return response, nil
}

// Connection types

func (s *Server) CreateConnectionType(ctx context.Context, req *pb.ConnectionTypeRequest) (*pb.UsersConnectionType, error) {
	l.Debug("Creating connection type",
		l.String("name", req.GetName()),
		l.String("request_id", r.GetRequestID(ctx)))

	// Translate request
	connectionTypeReq := &model.ConnectionTypeRequest{
		ID:         req.GetId(),
		Name:       req.GetName(),
		Definition: req.GetDefinition(),
	}

	// Validate request
	if err := s.validator.Struct(connectionTypeReq); err != nil {
		return nil, e.New("Request validation failed", ErrInvalidRequest, err)
	}

	// Create connection type
	connectionType, err := s.service.CreateConnectionType(ctx, connectionTypeReq)
	if err != nil {
		l.Warn("Failed to create connection type:", l.ErrField(err))
		return nil, translateToGrpcError(err)
	}

	// Validate response
	if err := s.validator.Struct(connectionType); err != nil {
		return nil, e.New("Response validation failed", ErrInternal, err)
	}

	// Translate to protobuf response
	response := translateConnectionTypeToProto(connectionType)

	return response, nil
}

func (s *Server) FindConnectionTypes(ctx context.Context, req *pb.SearchRequest) (*pb.ConnectionTypesList, error) {
	l.Debug("Finding connection types",
		l.String("name", req.GetName()),
		l.String("request_id", r.GetRequestID(ctx)))

	// Translate request
	searchReq := &model.SearchRequest{
		Name: req.GetName(),
	}

	// Validate request
	if err := s.validator.Struct(searchReq); err != nil {
		return nil, e.New("Request validation failed", ErrInvalidRequest, err)
	}

	// Find connection types
	connectionTypes, err := s.service.FindConnectionTypes(ctx, searchReq)
	if err != nil {
		l.Warn("Failed to find connection types:", l.ErrField(err))
		return nil, translateToGrpcError(err)
	}

	// Translate to protobuf response
	response := &pb.ConnectionTypesList{
		ConnectionTypes: translateConnectionTypesToProto(connectionTypes),
	}

	return response, nil
}

// Property types

func (s *Server) CreatePropertyType(ctx context.Context, req *pb.PropertyTypeRequest) (*pb.UsersPropertyType, error) {
	l.Debug("Creating property type",
		l.String("name", req.GetName()),
		l.String("request_id", r.GetRequestID(ctx)))

	// Translate request
	propertyTypeReq := &model.PropertyTypeRequest{
		ID:         req.GetId(),
		Name:       req.GetName(),
		Definition: req.GetDefinition(),
		ValueType: 	req.GetValueType(),
	}

	// Validate request
	if err := s.validator.Struct(propertyTypeReq); err != nil {
		return nil, e.New("Request validation failed", ErrInvalidRequest, err)
	}

	// Create property type
	propertyType, err := s.service.CreatePropertyType(ctx, propertyTypeReq)
	if err != nil {
		l.Warn("Failed to create property type:", l.ErrField(err))
		return nil, translateToGrpcError(err)
	}

	// Validate response
	if err := s.validator.Struct(propertyType); err != nil {
		return nil, e.New("Response validation failed", ErrInternal, err)
	}

	// Translate to protobuf response
	response := translatePropertyTypeToProto(propertyType)

	return response, nil
}

func (s *Server) FindPropertyTypes(ctx context.Context, req *pb.SearchRequest) (*pb.PropertyTypesList, error) {
	l.Debug("Finding property types",
		l.String("name", req.GetName()),
		l.String("request_id", r.GetRequestID(ctx)))

	// Translate request
	searchReq := &model.SearchRequest{
		Name: req.GetName(),
	}

	// Validate request
	if err := s.validator.Struct(searchReq); err != nil {
		return nil, e.New("Request validation failed", ErrInvalidRequest, err)
	}

	// Find property types
	propertyTypes, err := s.service.FindPropertyTypes(ctx, searchReq)
	if err != nil {
		l.Warn("Failed to find property types:", l.ErrField(err))
		return nil, translateToGrpcError(err)
	}

	// Translate to protobuf response
	response := &pb.PropertyTypesList{
		PropertyTypes: translatePropertyTypesToProto(propertyTypes),
	}

	return response, nil
}

// HELPER FUNCTIONS

func translateEntityToProto(entity *model.UsersEntity) *pb.UsersEntity {
	return &pb.UsersEntity{
		Name:       entity.Name,
		Definition: entity.Definition,
		UserId:     entity.UserID,
		EntityId:   entity.EntityID,
	}
}

func translateEntitiesToProto(entities []model.UsersEntity) []*pb.UsersEntity {
	result := make([]*pb.UsersEntity, len(entities))
	for i, entity := range entities {
		result[i] = translateEntityToProto(&entity)
	}
	return result
}

func translateConnectionTypeToProto(connectionType *model.UsersConnectionType) *pb.UsersConnectionType {
	return &pb.UsersConnectionType{
		Name:           connectionType.Name,
		Definition:     connectionType.Definition,
		UserId:         connectionType.UserID,
		ConnectionTypeId: connectionType.ConnectionTypeID,
	}
}

func translateConnectionTypesToProto(connectionTypes []model.UsersConnectionType) []*pb.UsersConnectionType {
	result := make([]*pb.UsersConnectionType, len(connectionTypes))
	for i, connectionType := range connectionTypes {
		result[i] = translateConnectionTypeToProto(&connectionType)
	}
	return result
}

func translatePropertyTypeToProto(propertyType *model.PropertyTypeResponse) *pb.UsersPropertyType {
	return &pb.UsersPropertyType{
		Name:         propertyType.Name,
		Definition:   propertyType.Definition,
		UserId:       propertyType.UserID,
		PropertyTypeId: propertyType.PropertyTypeID,
		ValueType:    propertyType.ValueType,
	}
}

func translatePropertyTypesToProto(propertyTypes []model.PropertyTypeResponse) []*pb.UsersPropertyType {
	result := make([]*pb.UsersPropertyType, len(propertyTypes))
	for i, propertyType := range propertyTypes {
		result[i] = translatePropertyTypeToProto(&propertyType)
	}
	return result
}
