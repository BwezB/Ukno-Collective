package db

import (
	"context"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/BwezB/Wikno-backend/internal/graph/model"

	c "github.com/BwezB/Wikno-backend/pkg/context"
	e "github.com/BwezB/Wikno-backend/pkg/errors"
	l "github.com/BwezB/Wikno-backend/pkg/log"
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

	// Set up join tables
	if err := setupJoinTable(db); err != nil {
		return nil, e.Wrap("Failed to setup join tables", err)
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

func setupJoinTable(db *gorm.DB) error {
	if err := db.SetupJoinTable(&model.User{}, "Entities", &model.UserEntity{}); err != nil {
		return e.New("Failed to setup join table for Entries", ErrInternal, err)
	}

	if err := db.SetupJoinTable(&model.User{}, "ConnectionTypes", &model.UserConnectionType{}); err != nil {
		return e.New("Failed to setup join table for ConnectionTypes", ErrInternal, err)
	}

	if err := db.SetupJoinTable(&model.User{}, "PropertyTypes", &model.UserPropertyType{}); err != nil {
		return e.New("Failed to setup join table for PropertyTypes", ErrInternal, err)
	}
	return nil
}

// Other database setup functions

func (db *Database) AutoMigrate() error {
	l.Debug("Auto migrating database")

	err := db.DB.AutoMigrate(
		&model.User{},
		&model.Entity{},
		&model.ConnectionType{},
		&model.PropertyType{},
		&model.UserEntity{},
		&model.UserConnectionType{},
		&model.UserPropertyType{},
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
		&model.User{},
		&model.Entity{},
		&model.ConnectionType{},
		&model.PropertyType{},
		&model.UserEntity{},
		&model.UserConnectionType{},
		&model.UserPropertyType{},
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
		l.String("request_id", c.GetRequestID(ctx)))

	// Create the user object that will be stored in the DB
	user := model.User{ID: req.ID}
	res := db.WithContext(ctx).Create(&user)
	if res.Error != nil {
		return TranslateDatabaseError(res.Error)
	}
	return nil
}

// GetUser gets the user by ID (from the context) and preloads the user's entities, connection types, and property types.
func (db *Database) GetUser(ctx context.Context) (*model.User, error) {
	// Get ID from context
	userID := c.GetUserID(ctx)
	if userID == "" {
		return nil, e.New("Failed to get user ID from context", ErrInternal, nil)
	}

	l.Debug("Getting user",
		l.String("user_id", userID),
		l.String("request_id", c.GetRequestID(ctx)))

	var user model.User
	res := db.WithContext(ctx).
		Preload("Entities").
		Preload("ConnectionTypes").
		Preload("PropertyTypes").
		First(&user, "id = ?", userID)
	if res.Error != nil {
		return nil, TranslateDatabaseError(res.Error)
	}

	return &user, nil
}

// CreateEntity creates a UserEntity and creates an Entity if one does not already exist.
func (db *Database) CreateEntity(ctx context.Context, req *model.EntityRequest) error {
	// Get ID from context
	userID := c.GetUserID(ctx)
	if userID == "" {
		return e.New("Failed to get user ID from context", ErrInternal, nil)
	}

	l.Debug("Creating entity",
		l.String("user_id", userID),
		l.String("name", req.Name),
		l.String("definition", req.Definition),
		l.String("entity_id", req.ID),
		l.String("request_id", c.GetRequestID(ctx)))

	// Start transaction
	tx := db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return TranslateDatabaseError(tx.Error)
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
			return TranslateDatabaseError(err)
		}
	} else {
		// Create new entity
		entity = model.Entity{}
		if err := tx.Create(&entity).Error; err != nil {
			tx.Rollback()
			return TranslateDatabaseError(err)
		}
	}

	// Create user-entity relationship
	userEntity := model.UserEntity{
		UserID:     userID,
		EntityID:   entity.ID,
		Name:       req.Name,
		Definition: req.Definition,
	}
	if err := tx.Create(&userEntity).Error; err != nil {
		tx.Rollback()
		return TranslateDatabaseError(err)
	}

	return TranslateDatabaseError(tx.Commit().Error) // Returns nil if error is nil
}

// UpdateEntity updates the UserEntity
func (db *Database) UpdateEntity(ctx context.Context, req *model.EntityRequest) error {
	// Get ID from context
	userID := c.GetUserID(ctx)
	if userID == "" {
		return e.New("Failed to get user ID from context", ErrInternal, nil)
	}

	l.Debug("Updating entity",
		l.String("user_id", userID),
		l.String("name", req.Name),
		l.String("definition", req.Definition),
		l.String("entity_id", req.ID),
		l.String("request_id", c.GetRequestID(ctx)))

	res := db.WithContext(ctx).Model(&model.UserEntity{}).
		Where("user_id = ? AND entity_id = ?", userID, req.ID).
		Updates(map[string]interface{}{
			"name":       req.Name,
			"definition": req.Definition,
		})

	if res.Error != nil {
		return TranslateDatabaseError(res.Error)
	}
	if res.RowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

// FindEntitiesWithName finds the Entity that has the given name written in the UserEntity table.
func (db *Database) FindEntitiesWithName(ctx context.Context, req *model.SearchRequest) ([]model.UserEntity, error) {
    l.Debug("Finding entities with name",
        l.String("name", req.Name),
        l.String("request_id", c.GetRequestID(ctx)))

    var userEntities []model.UserEntity
    res := db.WithContext(ctx).
        Raw(`SELECT DISTINCT ON (entity_id) *
             FROM user_entities 
             WHERE name = ?`, req.Name).
        Scan(&userEntities)
    
    if res.Error != nil {
        return nil, TranslateDatabaseError(res.Error)
    }

    return userEntities, nil
}

// CreateConnectionType creates a UserConnectionType and creates a ConnectionType if one does not already exist.
func (db *Database) CreateConnectionType(ctx context.Context, req *model.ConnectionTypeRequest) error {
	// Get ID from context
	userID := c.GetUserID(ctx)
	if userID == "" {
		return e.New("Failed to get user ID from context", ErrInternal, nil)
	}

	l.Debug("Creating connection type",
		l.String("user_id", userID),
		l.String("name", req.Name),
		l.String("definition", req.Definition),
		l.String("connection_type_id", req.ID),
		l.String("request_id", c.GetRequestID(ctx)))

	// Start transaction
	tx := db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return TranslateDatabaseError(tx.Error)
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
			return TranslateDatabaseError(err)
		}
	} else {
		// Create new connection type
		connectionType = model.ConnectionType{}
		if err := tx.Create(&connectionType).Error; err != nil {
			tx.Rollback()
			return TranslateDatabaseError(err)
		}
	}

	// Create user-connection type relationship
	userConnectionType := model.UserConnectionType{
		UserID:           userID,
		ConnectionTypeID: connectionType.ID,
		Name:             req.Name,
		Definition:       req.Definition,
	}
	if err := tx.Create(&userConnectionType).Error; err != nil {
		tx.Rollback()
		return TranslateDatabaseError(err)
	}

	return TranslateDatabaseError(tx.Commit().Error) // Returns nil if error is nil
}

// FindConnectionTypesWithName finds the ConnectionType that has the given name written in the UserConnectionType table.
func (db *Database) FindConnectionTypesWithName(ctx context.Context, req *model.SearchRequest) ([]model.UserConnectionType, error) {
	l.Debug("Finding connection types with name",
		l.String("name", req.Name),
		l.String("request_id", c.GetRequestID(ctx)))

	var userConnectionTypes []model.UserConnectionType
	res := db.WithContext(ctx).
		Raw(`SELECT DISTINCT ON (connection_type_id) *
			 FROM user_connection_types
			 WHERE name = ?`, req.Name).
		Scan(&userConnectionTypes)

	if res.Error != nil {
		return nil, TranslateDatabaseError(res.Error)
	}
	return userConnectionTypes, nil
}

// CreatePropertyType creates a UserPropertyType and creates a PropertyType if one does not already exist.
func (db *Database) CreatePropertyType(ctx context.Context, req *model.PropertyTypeRequest) error {
	// Get ID from context
	userID := c.GetUserID(ctx)
	if userID == "" {
		return e.New("Failed to get user ID from context", ErrInternal, nil)
	}

	l.Debug("Creating property type",
		l.String("user_id", userID),
		l.String("name", req.Name),
		l.String("definition", req.Definition),
		l.String("property_type_id", req.ID),
		l.String("request_id", c.GetRequestID(ctx)))

	// Start transaction
	tx := db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return TranslateDatabaseError(tx.Error)
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
			return TranslateDatabaseError(err)
		}
	} else {
		// Create new property type
		propertyType = model.PropertyType{}
		if err := tx.Create(&propertyType).Error; err != nil {
			tx.Rollback()
			return TranslateDatabaseError(err)
		}
	}

	// Create user-property type relationship
	userPropertyType := model.UserPropertyType{
		UserID:         userID,
		PropertyTypeID: propertyType.ID,
		Name:           req.Name,
		Definition:     req.Definition,
	}
	if err := tx.Create(&userPropertyType).Error; err != nil {
		tx.Rollback()
		return TranslateDatabaseError(err)
	}

	return TranslateDatabaseError(tx.Commit().Error) // Returns nil if error is nil
}

// FindPropertyTypesWithName finds the PropertyType that has the given name written in the UserPropertyType table.
func (db *Database) FindPropertyTypesWithName(ctx context.Context, req *model.SearchRequest) ([]model.UserPropertyType, error) {
	l.Debug("Finding property types with name",
		l.String("name", req.Name),
		l.String("request_id", c.GetRequestID(ctx)))

	var userPropertyTypes []model.UserPropertyType
	res := db.WithContext(ctx).
		Raw(`SELECT DISTINCT ON (property_type_id) *
			 FROM user_property_types
			 WHERE name = ?`, req.Name).
		Scan(&userPropertyTypes)

	if res.Error != nil {
		return nil, TranslateDatabaseError(res.Error)
	}
	return userPropertyTypes, nil
}
