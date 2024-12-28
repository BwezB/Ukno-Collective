package config

import (
	"github.com/BwezB/Wikno-backend/pkg/common/baseconfig"
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
	baseconfig.SetEnvValue(&d.Host, envDBHost)
	baseconfig.SetEnvValue(&d.Port, envDBPort)
	baseconfig.SetEnvValue(&d.User, envDBUser)
	baseconfig.SetEnvValue(&d.Password, envDBPassword)
	baseconfig.SetEnvValue(&d.DBName, envDBName)
}

func (d *Database) AddFromFlags() {
	baseconfig.SetFlagValue(&d.Host, flagDatabaseHost)
	baseconfig.SetFlagValue(&d.Port, flagDatabasePort)
	baseconfig.SetFlagValue(&d.User, flagDatabaseUser)
	baseconfig.SetFlagValue(&d.Password, flagDatabasePassword)
	baseconfig.SetFlagValue(&d.DBName, flagDatabaseName)
}
