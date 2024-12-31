package config

import (
	"github.com/BwezB/Wikno-backend/pkg/common/configs"
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
	configs.SetEnvValue(&d.Host, envDBHost)
	configs.SetEnvValue(&d.Port, envDBPort)
	configs.SetEnvValue(&d.User, envDBUser)
	configs.SetEnvValue(&d.Password, envDBPassword)
	configs.SetEnvValue(&d.DBName, envDBName)
}

func (d *Database) AddFromFlags() {
	configs.SetFlagValue(&d.Host, flagDatabaseHost)
	configs.SetFlagValue(&d.Port, flagDatabasePort)
	configs.SetFlagValue(&d.User, flagDatabaseUser)
	configs.SetFlagValue(&d.Password, flagDatabasePassword)
	configs.SetFlagValue(&d.DBName, flagDatabaseName)
}

func (d *Database) GetAddress() string {
	return d.Host + ":" + d.Port
}

func (d *Database) GetDSN() string {
	return "host=" + d.Host + " port=" + d.Port + " user=" + d.User + " password=" + d.Password + " dbname=" + d.DBName + " sslmode=disable"
}
