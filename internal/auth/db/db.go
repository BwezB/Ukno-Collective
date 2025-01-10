package db

import (
	"context"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/BwezB/Wikno-backend/internal/auth/model"

	r "github.com/BwezB/Wikno-backend/pkg/requestid"
	e "github.com/BwezB/Wikno-backend/pkg/errors"
	h "github.com/BwezB/Wikno-backend/pkg/health"
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

func (db *Database) AutoMigrate() error {
	l.Debug("Auto migrating database")

	err := db.DB.AutoMigrate(&model.User{})
	if err != nil {
		return e.New("Auto migration failed", ErrInternal, err)
	}

	l.Info("Auto migration successful")
	return nil
}

func (db *Database) DropTables() error {
	l.Debug("Resetting tables")
	err := db.Migrator().DropTable(&model.User{})
	if err != nil {
		return e.New("Failed to reset tables", ErrInternal, err)
	}
	l.Info("Tables reset")
	return nil
}


// GET/SET METHODS

// CreateUser needs to get a hashed password!
func (db *Database) CreateUser(ctx context.Context, req *model.AuthRequest, hashedPassword string) (*model.User, error) {
	l.Debug("Creating user",
		l.String("email", req.Email),
		l.String("request_id", r.GetRequestID(ctx)))
	
	// Create the user object that will be stored in the DB
	user := &model.User{
		Email:    req.Email,
		Password: hashedPassword,
	}
	
	res := db.WithContext(ctx).Create(user)
	if res.Error != nil {
		return nil, TranslateDatabaseError(res.Error)
	}

	l.Info("Created user", l.String("email", req.Email), l.String("id", user.ID))
	return user, nil
}

func (db *Database) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	l.Debug("Getting user by email",
		l.String("email", email),
		l.String("request_id", r.GetRequestID(ctx)))

	var user model.User
	res := db.WithContext(ctx).First(&user, "email = ?", email)
	if res.Error != nil {
		return nil, TranslateDatabaseError(res.Error)
	}
	return &user, nil
}

func (db *Database) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	l.Debug("Getting user by id",
		l.String("id", id),
		l.String("request_id", r.GetRequestID(ctx)))

	var user model.User
	res := db.WithContext(ctx).First(&user, "id = ?", id)
	if res.Error != nil {
		return nil, TranslateDatabaseError(res.Error)
	}
	return &user, nil
}


// HEALTH CHECK

// HealthCheck checks the health of the database
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
