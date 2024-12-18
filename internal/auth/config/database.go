package config

import (
	"flag"
	"fmt"
)


// Defult vales for the database
const (
	defaultDBHost = "localhost"
	defaultDBPort = "5432"
	defaultDBUser = "postgres"
	defaultDBName = "auth_db"
)

// Environment variables for the database
const (
	envDBHost = "DB_HOST"
	envDBPort = "DB_PORT"
	envDBUser = "DB_USER"
	envDBPassword = "DB_PASSWORD"
	envDBName = "DB_NAME"
)

// Flags for the database
var (
	flagDBHost = flag.String("db-host", "", "Database host")
	flagDBPort = flag.String("db-port", "", "Database port")
	flagDBUser = flag.String("db-user", "", "Database user")
	flagDBName = flag.String("db-name", "", "Database name")
)


// DatabaseConfig holds all the configuration for the database
type DatabaseConfig struct {
	Host	 string `yaml:"host" validate:"required" json:"-"`
	Port     string	`yaml:"port" validate:"required,number" json:"-"`
	User     string `yaml:"user" validate:"required" json:"-"`
	Password string `yaml:"password" validate:"required" json:"-"`
	DBName   string `yaml:"name" validate:"required" json:"-"`
}

func NewDatabaseConfig(fileDBConfig *DatabaseConfig) (*DatabaseConfig, error) {
	if !flag.Parsed() {
		flag.Parse() // Parse the flags if they have not been parsed
	}

	return &DatabaseConfig{
		Host:     getConfigValue(*flagDBHost, envDBHost, fileDBConfig.Host,  defaultDBHost),
		Port:     getConfigValue(*flagDBPort, envDBPort, fileDBConfig.Port, defaultDBPort),
		User:     getConfigValue(*flagDBUser, envDBUser, fileDBConfig.User, defaultDBUser),
		Password: getConfigValue("", envDBPassword, fileDBConfig.Password, ""), // Raises error if password isnt in env or file
		DBName:   getConfigValue(*flagDBName, envDBName, fileDBConfig.DBName, defaultDBName),
	}, nil
}

func (c DatabaseConfig) String() string {
    return fmt.Sprintf("DatabaseConfig{Host: %s, Port: %s, User: %s, DBName: %s, Password: ***}", 
        c.Host, c.Port, c.User, c.DBName)
}
