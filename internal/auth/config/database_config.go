package config

import (
	"github.com/BwezB/Wikno-backend/pkg/common/configuration"
)

type Database struct {
	// Host is the address of the database server
	Host string `yaml:"host" validate:"required,hostname"`
	// Port is the port of the database server
	Port string `yaml:"port" validate:"required,number,min=1,max=65535"`
	// User is the username to connect to the database
	User string `yaml:"user" validate:"required"`
	// Password is the password to connect to the database
	Password string `yaml:"password" validate:"required" json:"-"`
	// DBName is the name of the database to connect to
	DBName string `yaml:"dbname" validate:"required"`
}

func (d *Database) SetDefaults() {
	d.Host = "localhost"
	d.Port = "5432"
	d.User = "postgres"
	// Password is intentionally left blank
	d.DBName = "auth_db"
}

func (d *Database) AddFromEnv() {
	configuration.SetEnvValue(&d.Host, "DATABASE_HOST")
	configuration.SetEnvValue(&d.Port, "DATABASE_PORT")
	configuration.SetEnvValue(&d.User, "DATABASE_USER")
	configuration.SetEnvValue(&d.Password, "DATABASE_PASSWORD")
	configuration.SetEnvValue(&d.DBName, "DATABASE_NAME")
}

var (
	flagDatabaseHost     = configuration.NewFlag("database_host", "", "Database Host")
	flagDatabasePort     = configuration.NewFlag("database_port", "", "Database Port")
	flagDatabaseUser     = configuration.NewFlag("database_user", "", "Database User")
	flagDatabasePassword = configuration.NewFlag("database_password", "", "Database Password")
	flagDatabaseName     = configuration.NewFlag("database_name", "", "Database Name")
)

func (d *Database) AddFromFlags() {
	configuration.SetFlagValue(&d.Host, flagDatabaseHost)
	configuration.SetFlagValue(&d.Port, flagDatabasePort)
	configuration.SetFlagValue(&d.User, flagDatabaseUser)
	configuration.SetFlagValue(&d.Password, flagDatabasePassword)
	configuration.SetFlagValue(&d.DBName, flagDatabaseName)
}
