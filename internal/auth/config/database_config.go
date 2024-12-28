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
	d.Host = defaultDBHost
	d.Port = defaultDBPort
	d.User = defaultDBUser
	// Password is intentionally left blank
	d.DBName = defaultDBName
}

func (d *Database) AddFromEnv() {
	configuration.SetEnvValue(&d.Host, envDBHost)
	configuration.SetEnvValue(&d.Port, envDBPort)
	configuration.SetEnvValue(&d.User, envDBUser)
	configuration.SetEnvValue(&d.Password, envDBPassword)
	configuration.SetEnvValue(&d.DBName, envDBName)
}

func (d *Database) AddFromFlags() {
	configuration.SetFlagValue(&d.Host, flagDatabaseHost)
	configuration.SetFlagValue(&d.Port, flagDatabasePort)
	configuration.SetFlagValue(&d.User, flagDatabaseUser)
	configuration.SetFlagValue(&d.Password, flagDatabasePassword)
	configuration.SetFlagValue(&d.DBName, flagDatabaseName)
}
