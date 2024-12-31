package db

import (
	c "github.com/BwezB/Wikno-backend/pkg/configs"
)

type DatabaseConfig struct {
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

func (d *DatabaseConfig) SetDefaults() {
	d.Host = defaultDBHost
	d.Port = defaultDBPort
	d.User = defaultDBUser
	// Password is intentionally left blank
	d.DBName = defaultDBName
}

func (d *DatabaseConfig) AddFromEnv() {
	c.SetEnvValue(&d.Host, envDBHost)
	c.SetEnvValue(&d.Port, envDBPort)
	c.SetEnvValue(&d.User, envDBUser)
	c.SetEnvValue(&d.Password, envDBPassword)
	c.SetEnvValue(&d.DBName, envDBName)
}

func (d *DatabaseConfig) AddFromFlags() {
	c.SetFlagValue(&d.Host, flagDatabaseHost)
	c.SetFlagValue(&d.Port, flagDatabasePort)
	c.SetFlagValue(&d.User, flagDatabaseUser)
	c.SetFlagValue(&d.Password, flagDatabasePassword)
	c.SetFlagValue(&d.DBName, flagDatabaseName)
}

func (d *DatabaseConfig) GetAddress() string {
	return d.Host + ":" + d.Port
}

func (d *DatabaseConfig) GetDSN() string {
	return "host=" + d.Host + " port=" + d.Port + " user=" + d.User + " password=" + d.Password + " dbname=" + d.DBName + " sslmode=disable"
}

// DEFAULTS
const (
	// defaultDBHost is the default host for the database.
	defaultDBHost = "localhost"
	// defaultPort is the default port for the database.
	defaultDBPort = "5432"
	// defaultUser is the default user for the database.
	defaultDBUser = "postgres"
	// No default password
	// defaultDBName is the default database name.
	defaultDBName = "auth_db"
)

// ENV
const (
	// envDBHost is the environment variable for the database host.
	envDBHost = "DB_HOST"
	// envDBPort is the environment variable for the database port.
	envDBPort = "DB_PORT"
	// envDBUser is the environment variable for the database user.
	envDBUser = "DB_USER"
	// envDBPassword is the environment variable for the database password.
	envDBPassword = "DB_PASSWORD"
	// envDBName is the environment variable for the database name.
	envDBName = "DB_NAME"
)

// FLAGS
var (
	flagDatabaseHost     = c.NewFlag("database_host", "", "Database Host")
	flagDatabasePort     = c.NewFlag("database_port", "", "Database Port")
	flagDatabaseUser     = c.NewFlag("database_user", "", "Database User")
	flagDatabasePassword = c.NewFlag("database_password", "", "Database Password")
	flagDatabaseName     = c.NewFlag("database_name", "", "Database Name")
)
