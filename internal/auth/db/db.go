package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/BwezB/Wikno-backend/internal/auth/config"
	"github.com/BwezB/Wikno-backend/internal/auth/model"
)

type Database struct {
	*gorm.DB
}

func New(config *config.Database) (*Database, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		config.Host, config.User, config.Password, config.DBName, config.Port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &Database{db}, nil
}

func (db *Database) AutoMigrate() error {
	return db.DB.AutoMigrate(&model.User{})
}

func (db *Database) CreateUser(user *model.User) error {
	res := db.Create(user)
	return res.Error
}

func (db *Database) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	res := db.First(&user, "email = ?", email)
	return &user, res.Error
}
