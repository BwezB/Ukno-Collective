package db

import (
	"context"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/BwezB/Wikno-backend/internal/graph/model"

	a "github.com/BwezB/Wikno-backend/pkg/auth"
	r "github.com/BwezB/Wikno-backend/pkg/requestid"
	e "github.com/BwezB/Wikno-backend/pkg/errors"
	l "github.com/BwezB/Wikno-backend/pkg/log"
	h "github.com/BwezB/Wikno-backend/pkg/health"
)

type Database struct {
	*gorm.DB
}

func New(config DatabaseConfig) (*Database, error) {
	l.Debug("Connecting to database with gorm",
		l.String("address", config.GetAddress()),
		l.String("user", config.User),
		l.String("dbname", config.DBName))

	// Connect to the database
	dsn := config.GetDSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Disable gorm logging
	})
	if err != nil {
		return nil, e.New("Failed to connect to database", ErrDatabaseConnection, err)
	}

	// Set up connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, e.New("Failed to get sql.DB from gorm.DB", ErrInternal, err)
	}

	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)

	l.Info("Connected to database",
		l.String("address", config.GetAddress()),
		l.String("user", config.User),
		l.String("dbname", config.DBName))

	return &Database{db}, nil
}

// Other database setup functions

func (db *Database) AutoMigrate() error {
	l.Debug("Auto migrating database")

	err := db.DB.AutoMigrate(
		&model.GraphUser{},
		&model.Entity{},
		&model.ConnectionType{},
		&model.PropertyType{},
		&model.UsersEntity{},
		&model.UsersConnectionType{},
		&model.UsersPropertyType{},
	)
	if err != nil {
		return e.New("Auto migration failed", ErrInternal, err)
	}

	l.Info("Auto migration successful")
	return nil
}

func (db *Database) DropTables() error {
	l.Debug("Resetting tables")
	err := db.Migrator().DropTable(
		&model.GraphUser{},
		&model.Entity{},
		&model.ConnectionType{},
		&model.PropertyType{},
		&model.UsersEntity{},
		&model.UsersConnectionType{},
		&model.UsersPropertyType{},
	)
	if err != nil {
		return e.New("Failed to reset tables", ErrInternal, err)
	}
	l.Info("Tables reset")
	return nil
}

// CRUD

// CreateUser creates a new user (with ID and email from the context)
func (db *Database) CreateUser(ctx context.Context, req *model.UserRequest) error {
	l.Debug("Creating user",
		l.String("user_id", req.ID),
		l.String("request_id", r.GetRequestID(ctx)))

	// Create the user object that will be stored in the DB
	user := model.GraphUser{ID: req.ID}
	res := db.WithContext(ctx).Create(&user)
	if res.Error != nil {
		return e.Wrap("Failed to create user", TranslateDatabaseError(res.Error))
	}

	l.Info("Created user", l.String("user_id", user.ID), l.String("request_id", r.GetRequestID(ctx)))

	return nil
}

// GetUserData gets the user by ID (from the context) and preloads the user's entities, connection types, and property types.
func (db *Database) GetUserData(ctx context.Context) (*model.UserDataResponse, error) {
	// Get ID from context
	userID := a.GetUserID(ctx)
	if userID == "" {
		return nil, e.New("Failed to get user ID from context", ErrInternal, nil)
	}

	l.Debug("Getting user",
		l.String("user_id", userID),
		l.String("request_id", r.GetRequestID(ctx)))

	var user model.GraphUser
	res := db.WithContext(ctx).
		Preload("UsersEntities").
		Preload("UsersConnectionTypes").
		Preload("UsersPropertyTypes").
		First(&user, "id = ?", userID)
	if res.Error != nil {
		return nil, e.Wrap("Failed to get user", TranslateDatabaseError(res.Error))
	}

	// Translate to response
	userData, err := db.translateUserToResponse(ctx, &user)
	if err != nil {
		return nil, e.Wrap("Failed to translate user to response", err)
	}

	return userData, nil
}

// CreateEntity creates a UserEntity and creates an Entity if one does not already exist.
func (db *Database) CreateEntity(ctx context.Context, req *model.EntityRequest) (*model.UsersEntity, error) {
	// Get ID from context
	userID := a.GetUserID(ctx)
	if userID == "" {
		return nil, e.New("Failed to get user ID from context", ErrInternal, nil)
	}

	l.Debug("Creating entity",
		l.String("user_id", userID),
		l.String("name", req.Name),
		l.String("definition", req.Definition),
		l.String("entity_id", req.ID),
		l.String("request_id", r.GetRequestID(ctx)))

	// Start transaction
	tx := db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, e.Wrap("Could not start transaction", TranslateDatabaseError(tx.Error))
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var entity model.Entity
	if req.ID != "" {
		// Check if entity exists
		if err := tx.First(&entity, "id = ?", req.ID).Error; err != nil {
			tx.Rollback() // Even if ErrRecordNotFound, rollback (entity must exist)
			return nil, e.Wrap("Could not find entity", TranslateDatabaseError(err))
		}
	} else {
		// Create new entity
		entity = model.Entity{}
		if err := tx.Create(&entity).Error; err != nil {
			tx.Rollback()
			return nil, e.Wrap("Could not create entity", TranslateDatabaseError(err))
		}
	}

	// Create user-entity relationship
	userEntity := model.UsersEntity{
		UserID:     userID,
		EntityID:   entity.ID,
		Name:       req.Name,
		Definition: req.Definition,
	}
	if err := tx.Create(&userEntity).Error; err != nil {
		tx.Rollback()
		return nil, e.Wrap("Could not create userEntity", TranslateDatabaseError(err))
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, e.Wrap("Could not commit", TranslateDatabaseError(err))
	}

	return &userEntity, nil
}

// UpdateEntity updates the UserEntity
func (db *Database) UpdateEntity(ctx context.Context, req *model.EntityRequest) error {
	// Get ID from context
	userID := a.GetUserID(ctx)
	if userID == "" {
		return e.New("Failed to get user ID from context", ErrInternal, nil)
	}

	l.Debug("Updating entity",
		l.String("user_id", userID),
		l.String("name", req.Name),
		l.String("definition", req.Definition),
		l.String("entity_id", req.ID),
		l.String("request_id", r.GetRequestID(ctx)))

	res := db.WithContext(ctx).Model(&model.UsersEntity{}).
		Where("user_id = ? AND entity_id = ?", userID, req.ID).
		Updates(map[string]interface{}{
			"name":       req.Name,
			"definition": req.Definition,
		})

	if res.Error != nil {
		return e.Wrap("Failed to update entity", TranslateDatabaseError(res.Error))
	}
	if res.RowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

// FindEntitiesWithName finds the Entity that has the given name written in the UserEntity table.
func (db *Database) FindEntitiesWithName(ctx context.Context, req *model.SearchRequest) ([]model.UsersEntity, error) {
	l.Debug("Finding entities with name",
		l.String("name", req.Name),
		l.String("request_id", r.GetRequestID(ctx)))

	var userEntities []model.UsersEntity
	res := db.WithContext(ctx).
		Raw(`SELECT DISTINCT ON (entity_id) *
             FROM users_entities
             WHERE name = ?`, req.Name).
		Scan(&userEntities)

	if res.Error != nil {
		return nil, e.Wrap("Failed to find entities with name", TranslateDatabaseError(res.Error))
	}

	return userEntities, nil
}

// CreateConnectionType creates a UserConnectionType and creates a ConnectionType if one does not already exist.
func (db *Database) CreateConnectionType(ctx context.Context, req *model.ConnectionTypeRequest) (*model.UsersConnectionType, error) {
	// Get ID from context
	userID := a.GetUserID(ctx)
	if userID == "" {
		return nil, e.New("Failed to get user ID from context", ErrInternal, nil)
	}

	l.Debug("Creating connection type",
		l.String("user_id", userID),
		l.String("name", req.Name),
		l.String("definition", req.Definition),
		l.String("connection_type_id", req.ID),
		l.String("request_id", r.GetRequestID(ctx)))

	// Start transaction
	tx := db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, e.Wrap("Could not start transaction", TranslateDatabaseError(tx.Error))
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var connectionType model.ConnectionType
	if req.ID != "" {
		// Check if connection type exists
		if err := tx.First(&connectionType, "id = ?", req.ID).Error; err != nil {
			tx.Rollback() // Even if ErrRecordNotFound, rollback (connection type must exist)
			return nil, e.Wrap("Could not find connection type", TranslateDatabaseError(err))
		}
	} else {
		// Create new connection type
		connectionType = model.ConnectionType{}
		if err := tx.Create(&connectionType).Error; err != nil {
			tx.Rollback()
			return nil, e.Wrap("Could not create connection type", TranslateDatabaseError(err))
		}
	}

	// Create user-connection type relationship
	userConnectionType := model.UsersConnectionType{
		UserID:           userID,
		ConnectionTypeID: connectionType.ID,
		Name:             req.Name,
		Definition:       req.Definition,
	}
	if err := tx.Create(&userConnectionType).Error; err != nil {
		tx.Rollback()
		return nil, e.Wrap("Could not create userConnectionType", TranslateDatabaseError(err))
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, e.Wrap("Could not commit", TranslateDatabaseError(err))
	}

	return &userConnectionType, nil
}

// FindConnectionTypesWithName finds the ConnectionType that has the given name written in the UserConnectionType table.
func (db *Database) FindConnectionTypesWithName(ctx context.Context, req *model.SearchRequest) ([]model.UsersConnectionType, error) {
	l.Debug("Finding connection types with name",
		l.String("name", req.Name),
		l.String("request_id", r.GetRequestID(ctx)))

	var userConnectionTypes []model.UsersConnectionType
	res := db.WithContext(ctx).
		Raw(`SELECT DISTINCT ON (connection_type_id) *
			 FROM users_connection_types
			 WHERE name = ?`, req.Name).
		Scan(&userConnectionTypes)

	if res.Error != nil {
		return nil, e.Wrap("Failed to find connection types with name", TranslateDatabaseError(res.Error))
	}
	return userConnectionTypes, nil
}

// CreatePropertyType creates a UserPropertyType and creates a PropertyType if one does not already exist.
func (db *Database) CreatePropertyType(ctx context.Context, req *model.PropertyTypeRequest) (*model.PropertyTypeResponse, error) {
	// Get ID from context
	userID := a.GetUserID(ctx)
	if userID == "" {
		return nil, e.New("Failed to get user ID from context", ErrInternal, nil)
	}

	l.Debug("Creating property type",
		l.String("user_id", userID),
		l.String("name", req.Name),
		l.String("definition", req.Definition),
		l.String("property_type_id", req.ID),
		l.String("request_id", r.GetRequestID(ctx)))

	// Start transaction
	tx := db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, e.Wrap("Could not start transaction", TranslateDatabaseError(tx.Error))
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var propertyType model.PropertyType
	if req.ID != "" {
		// Check if property type exists
		if err := tx.First(&propertyType, "id = ?", req.ID).Error; err != nil {
			tx.Rollback() // Even if ErrRecordNotFound, rollback (property type must exist)
			return nil, e.Wrap("Could not find property type", TranslateDatabaseError(err))
		}

		// Check if property type value type is the same
		if propertyType.ValueType != req.ValueType {
			tx.Rollback()
			return nil, e.New("Property type value type does not match", ErrInvalidRequest, nil)
		}
	} else {
		// Create new property type
		propertyType = model.PropertyType{
			ValueType: req.ValueType,
		}
		if err := tx.Create(&propertyType).Error; err != nil {
			tx.Rollback()
			return nil, e.Wrap("Could not create property type", TranslateDatabaseError(err))
		}
	}

	// Create user-property type relationship
	userPropertyType := model.UsersPropertyType{
		UserID:         userID,
		PropertyTypeID: propertyType.ID,
		Name:           req.Name,
		Definition:     req.Definition,
	}
	if err := tx.Create(&userPropertyType).Error; err != nil {
		tx.Rollback()
		return nil, e.Wrap("Could not create userPropertyType", TranslateDatabaseError(err))
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, e.Wrap("Could not commit", TranslateDatabaseError(err))
	}

	// Create response
	propertyTypeResponse := translatePropertyTypeToResponse(&propertyType, &userPropertyType)

	return propertyTypeResponse, nil
}

// FindPropertyTypesWithName finds the PropertyType that has the given name written in the UserPropertyType table.
func (db *Database) FindPropertyTypesWithName(ctx context.Context, req *model.SearchRequest) ([]model.PropertyTypeResponse, error) {
	l.Debug("Finding property types with name",
		l.String("name", req.Name),
		l.String("request_id", r.GetRequestID(ctx)))

	var userPropertyTypes []model.UsersPropertyType
	res := db.WithContext(ctx).
		Raw(`SELECT DISTINCT ON (property_type_id) *
			 FROM users_property_types
			 WHERE name = ?`, req.Name).
		Scan(&userPropertyTypes)
	if res.Error != nil {
		return nil, e.Wrap("Failed to find property types with name", TranslateDatabaseError(res.Error))
	}

	// Translate to response
	propertyTypes, err := db.translatePropertyTypesToResponse(ctx, userPropertyTypes)
	if err != nil {
		return nil, e.Wrap("Failed to translate property types to response", err)
	}

	return propertyTypes, nil
}


// HELPER FUNCTIONS

func (db *Database) translatePropertyTypesToResponse(ctx context.Context, usersPropertyTypes []model.UsersPropertyType) ([]model.PropertyTypeResponse, error) {
	propertyTypes := make([]model.PropertyTypeResponse, len(usersPropertyTypes))
	for _, userPropertyType := range usersPropertyTypes {
		// Get the property type so we can get the value type
		propertyType := model.PropertyType{}
		if err := db.WithContext(ctx).First(&propertyType, "id = ?", userPropertyType.PropertyTypeID).Error; err != nil {
			return nil, TranslateDatabaseError(err)
		}

		propertyTypes = append(propertyTypes, *translatePropertyTypeToResponse(&propertyType, &userPropertyType))
	}
	return propertyTypes, nil
}

func translatePropertyTypeToResponse(propertyType *model.PropertyType, userPropertyType *model.UsersPropertyType) *model.PropertyTypeResponse {
	return &model.PropertyTypeResponse{
		UserID:         userPropertyType.UserID,
		PropertyTypeID: userPropertyType.PropertyTypeID,
		Name:           userPropertyType.Name,
		Definition:     userPropertyType.Definition,
		ValueType:      propertyType.ValueType,
	}
}

func (db *Database) translateUserToResponse(ctx context.Context, user *model.GraphUser) (*model.UserDataResponse, error) {
	propertyTypes, err := db.translatePropertyTypesToResponse(ctx, user.UsersPropertyTypes)
	if err != nil {
		return nil, e.Wrap("Failed to translate property types to response", err)
	}

	return &model.UserDataResponse{
		ID:             user.ID,
		Entities:        user.UsersEntities,
		ConnectionTypes: user.UsersConnectionTypes,
		PropertyTypes:   propertyTypes,
	}, nil
}


// HEALTH CHECK

func (db *Database) HealthCheck(ctx context.Context) *h.HealthStatus {
	if err := db.WithContext(ctx).Exec("SELECT 1").Error; err != nil {
		return &h.HealthStatus{
			Healthy: false,
			Err:     e.New("health check gorm database connection failed", ErrDatabaseConnection, err),
			Time:    time.Now(),
		}
	}
	return &h.HealthStatus{
		Healthy: true,
		Time:    time.Now(),
	}
}