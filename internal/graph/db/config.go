package db

import (
	"time"
	"strconv"

	c "github.com/BwezB/Wikno-backend/pkg/configs"
)

type DatabaseConfig struct {
	// DATABASE CONFIG
	// Host is the address of the database server
	Host string `yaml:"host" validate:"required,hostname"`
	// Port is the port of the database server
	Port int `yaml:"port" validate:"required,min=1,max=65535"`
	// User is the username to connect to the database
	User string `yaml:"user" validate:"required"`
	// Password is the password to connect to the database
	Password string `yaml:"password" validate:"required" json:"-"`
	// DBName is the name of the database to connect to
	DBName string `yaml:"dbname" validate:"required"`


	// MaxOpenConns is the maximum number of open connections to the database
	MaxOpenConns int `yaml:"max_open_conns" validate:"number,min=1"`
	// MaxIdleConns is the maximum number of idle connections to the database
	MaxIdleConns int `yaml:"max_idle_conns" validate:"number,min=0"`
	// ConnMaxLifetime is the maximum lifetime of a connection to the database
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime" validate:"number,min=1"`

	// DropTables is a flag to drop tables on startup (DO NOT USE IN PRODUCTION)
	DropTables bool `yaml:"drop_tables" validate:"boolean"`
}

// DEFAULTS

func (d *DatabaseConfig) SetDefaults() {
	d.Host = "localhost"
	d.Port = 5432
	d.User = "postgres"
	// Password is intentionally left blank
	d.DBName = "graph_db"

	d.MaxOpenConns = 10
	d.MaxIdleConns = 5
	d.ConnMaxLifetime = 5 * time.Minute

	d.DropTables = false
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

	// DropTables should not be set from environment variables
}

// FLAGS

var (
	flagDatabaseHost     = c.NewFlag("db-host", "", "Database Host")
	flagDatabasePort     = c.NewFlag("db-port", "", "Database Port")
	flagDatabaseUser     = c.NewFlag("db-user", "", "Database User")
	flagDatabasePassword = c.NewFlag("db-password", "", "Database Password")
	flagDatabaseName     = c.NewFlag("db-name", "", "Database Name")

	flagDatabaseMaxOpenConns     = c.NewFlag("db-max-open-conns", "", "Database Max Open Connections")
	flagDatabaseMaxIdleConns     = c.NewFlag("db-max-idle-conns", "", "Database Max Idle Connections")
	flagDatabaseConnMaxLifetime  = c.NewFlag("db-conn-max-lifetime", "", "Database Connection Max Lifetime")

	flagDatabaseDropTables = c.NewFlag("db-drop-tables", "", "DROPS ALL TABLES! DO NOT USE IN PRODUCTION")
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

	c.SetFlagValue(&d.DropTables, flagDatabaseDropTables)
}

// HELPER FUNCTIONS

func (d *DatabaseConfig) GetAddress() string {
	return d.Host + ":" + strconv.Itoa(d.Port)
}

func (d *DatabaseConfig) GetDSN() string {
	return "host=" + d.Host + " port=" + strconv.Itoa(d.Port) + " user=" + d.User + " password=" + d.Password + " dbname=" + d.DBName + " sslmode=disable"
}
