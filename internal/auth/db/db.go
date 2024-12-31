package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/BwezB/Wikno-backend/internal/auth/model"
	"github.com/BwezB/Wikno-backend/pkg/log"
)

type Database struct {
	*gorm.DB
}

func New(config *DatabaseConfig) (*Database, error) {
	defer log.DebugFunc("Address:", config.GetAddress(), "User:", config.User, "DBName:", config.DBName)()
	dsn := config.GetDSN()

	log.Info("Connecting to database with gorm")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Error("Gorm open failed. Could not connect to", config.DBName, "at", config.GetAddress())
		return nil, ErrDatabaseConnectionFailed
	}

	return &Database{db}, nil
}

func (db *Database) AutoMigrate() error {
	log.Info("Auto migrating database")
	err := db.DB.AutoMigrate(&model.User{})
	if err != nil {
		return log.Errore(err, "Gorm auto migration failed. Could not migrate database")
	}
	return nil
}

func (db *Database) CreateUser(user *model.User) error {
	res := db.Create(user)
	if res.Error != nil {
		if res.Error.Error() == "UNIQUE constraint failed: users.email" {
			return ErrUserExists
		}
		return log.Errore(res.Error, "Gorm create failed. Could not create user")
	}
	return nil
}

func (db *Database) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	res := db.First(&user, "email = ?", email)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return nil, ErrUserNotFound
		}
		return nil, log.Errore(res.Error, "Gorm first failed. Could not get user by email")
	}
	return &user, nil
}
