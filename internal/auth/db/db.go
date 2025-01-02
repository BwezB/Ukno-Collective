package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/BwezB/Wikno-backend/internal/auth/model"

	e "github.com/BwezB/Wikno-backend/pkg/errors"
	l "github.com/BwezB/Wikno-backend/pkg/log"
)

type Database struct {
	*gorm.DB
}

func New(config *DatabaseConfig) (*Database, error) {
	defer l.DebugFunc("New (db)")()
	l.Info("Connecting to database with gorm",
		l.String("address", config.GetAddress()),
		l.String("user", config.User),
		l.String("dbname", config.DBName))

	dsn := config.GetDSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Disable gorm logging
	})
	if err != nil {
		return nil, e.New("Failed to connect to database", ErrDatabaseConnection, err)
	}

	return &Database{db}, nil
}

func (db *Database) AutoMigrate() error {
	l.Info("Auto migrating database")
	err := db.DB.AutoMigrate(&model.User{})
	if err != nil {
		return e.New("Auto migration failed", ErrInternal, err)
	}
	return nil
}

func (db *Database) CreateUser(user *model.User) error {
	res := db.Create(user)
	if res.Error != nil {
		return TranslateDatabaseError(res.Error)
	}
	return nil
}

func (db *Database) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	res := db.First(&user, "email = ?", email)
	if res.Error != nil {
		return nil, TranslateDatabaseError(res.Error)
	}
	return &user, nil
}
