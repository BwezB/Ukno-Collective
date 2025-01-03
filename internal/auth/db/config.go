package db

import (
	"time"

	c "github.com/BwezB/Wikno-backend/pkg/configs"
)

type DatabaseConfig struct {
	// DATABASE CONFIG
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

	// CONNECTION POOL
	// MaxOpenConns is the maximum number of open connections to the database
	MaxOpenConns int `yaml:"max_open_conns" validate:"number,min=1"`
	// MaxIdleConns is the maximum number of idle connections to the database
	MaxIdleConns int `yaml:"max_idle_conns" validate:"number,min=0"`
	// ConnMaxLifetime is the maximum lifetime of a connection to the database
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime" validate:"number,min=1"`
}

// DEFAULTS

func (d *DatabaseConfig) SetDefaults() {
	d.Host = "localhost"
	d.Port = "5432"
	d.User = "postgres"
	// Password is intentionally left blank
	d.DBName = "auth_db"

	d.MaxOpenConns = 10
	d.MaxIdleConns = 5
	d.ConnMaxLifetime = 5 * time.Minute
}

// ENVIRONMENT VARIABLES

func (d *DatabaseConfig) AddFromEnv() {
	c.SetEnvValue(&d.Host, "DB_HOST")
	c.SetEnvValue(&d.Port, "DB_PORT")
	c.SetEnvValue(&d.User, "DB_USER")
	c.SetEnvValue(&d.Password, "DB_PASSWORD")
	c.SetEnvValue(&d.DBName, "DB_NAME")

	c.SetEnvValue(&d.MaxOpenConns, "DB_MAX_OPEN_CONNS")
	c.SetEnvValue(&d.MaxIdleConns, "DB_MAX_IDLE_CONNS")
	c.SetEnvValue(&d.ConnMaxLifetime, "DB_CONN_MAX_LIFETIME")
}

// FLAGS

var (
	flagDatabaseHost     = c.NewFlag("database_host", "", "Database Host")
	flagDatabasePort     = c.NewFlag("database_port", "", "Database Port")
	flagDatabaseUser     = c.NewFlag("database_user", "", "Database User")
	flagDatabasePassword = c.NewFlag("database_password", "", "Database Password")
	flagDatabaseName     = c.NewFlag("database_name", "", "Database Name")

	flagDatabaseMaxOpenConns     = c.NewFlag("database_max_open_conns", "", "Database Max Open Connections")
	flagDatabaseMaxIdleConns     = c.NewFlag("database_max_idle_conns", "", "Database Max Idle Connections")
	flagDatabaseConnMaxLifetime  = c.NewFlag("database_conn_max_lifetime", "", "Database Connection Max Lifetime")
)

func (d *DatabaseConfig) AddFromFlags() {
	c.SetFlagValue(&d.Host, flagDatabaseHost)
	c.SetFlagValue(&d.Port, flagDatabasePort)
	c.SetFlagValue(&d.User, flagDatabaseUser)
	c.SetFlagValue(&d.Password, flagDatabasePassword)
	c.SetFlagValue(&d.DBName, flagDatabaseName)

	c.SetFlagValue(&d.MaxOpenConns, flagDatabaseMaxOpenConns)
	c.SetFlagValue(&d.MaxIdleConns, flagDatabaseMaxIdleConns)
	c.SetFlagValue(&d.ConnMaxLifetime, flagDatabaseConnMaxLifetime)
}

// HELPER FUNCTIONS

func (d *DatabaseConfig) GetAddress() string {
	return d.Host + ":" + d.Port
}

func (d *DatabaseConfig) GetDSN() string {
	return "host=" + d.Host + " port=" + d.Port + " user=" + d.User + " password=" + d.Password + " dbname=" + d.DBName + " sslmode=disable"
}
