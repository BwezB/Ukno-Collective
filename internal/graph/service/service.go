package service

import (
	"context"

	e "github.com/BwezB/Wikno-backend/pkg/errors"

	"github.com/BwezB/Wikno-backend/internal/graph/db"
	"github.com/BwezB/Wikno-backend/internal/graph/model"
)

type GraphService struct {
	db *db.Database
}

func NewService(database *db.Database) *GraphService {
	graphService := &GraphService{
		db: database,
	}
	return graphService
}

// DB PASSTHROUGHS

// CreateUser creates a new user
func (s *GraphService) CreateUser(ctx context.Context, req *model.UserRequest) error { // This should only be called by authservice
    err := s.db.CreateUser(ctx, req)
    if err != nil {
        return e.Wrap("CreateUser failed", err)
    }
    return nil
}

// GetUserData gets all user's entities and types
func (s *GraphService) GetUserData(ctx context.Context) (*model.UserDataResponse, error) {
	user, err := s.db.GetUserData(ctx)
	if err != nil {
		return nil, e.Wrap("GetUserData failed", err)
	}
	return user, nil
}

// CreateEntity creates a new entity or links to existing one
func (s *GraphService) CreateEntity(ctx context.Context, req *model.EntityRequest) (*model.UsersEntity, error) {
	usersEntity, err := s.db.CreateEntity(ctx, req)
	if err != nil {
		return nil, e.Wrap("CreateEntity failed", err)
	}
	return usersEntity, nil
}

// UpdateEntity updates an existing entity
func (s *GraphService) UpdateEntity(ctx context.Context, req *model.EntityRequest) error {
	if err := s.db.UpdateEntity(ctx, req); err != nil {
		return e.Wrap("UpdateEntity failed", err)
	}
	return nil
}

// FindEntities finds entities by name
func (s *GraphService) FindEntities(ctx context.Context, req *model.SearchRequest) ([]model.UsersEntity, error) {
	entities, err := s.db.FindEntitiesWithName(ctx, req)
	if err != nil {
		return nil, e.Wrap("FindEntities failed", err)
	}
	return entities, nil
}

// CreateConnectionType creates a new connection type or links to existing one
func (s *GraphService) CreateConnectionType(ctx context.Context, req *model.ConnectionTypeRequest) (*model.UsersConnectionType, error) {
	usersConnectionType, err := s.db.CreateConnectionType(ctx, req)
	if err != nil {
		return nil, e.Wrap("CreateConnectionType failed", err)
	}

	return usersConnectionType, nil
}

// FindConnectionTypes finds connection types by name
func (s *GraphService) FindConnectionTypes(ctx context.Context, req *model.SearchRequest) ([]model.UsersConnectionType, error) {
	types, err := s.db.FindConnectionTypesWithName(ctx, req)
	if err != nil {
		return nil, e.Wrap("FindConnectionTypes failed", err)
	}
	return types, nil
}

// CreatePropertyType creates a new property type or links to existing one
func (s *GraphService) CreatePropertyType(ctx context.Context, req *model.PropertyTypeRequest) (*model.PropertyTypeResponse, error) {
	usersPropertyType, err := s.db.CreatePropertyType(ctx, req)
	if err != nil {
		return nil, e.Wrap("CreatePropertyType failed", err)
	}
	return usersPropertyType, nil
}

// FindPropertyTypes finds property types by name
func (s *GraphService) FindPropertyTypes(ctx context.Context, req *model.SearchRequest) ([]model.PropertyTypeResponse, error) {
	types, err := s.db.FindPropertyTypesWithName(ctx, req)
	if err != nil {
		return nil, e.Wrap("FindPropertyTypes failed", err)
	}
	return types, nil
}
