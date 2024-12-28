package config
import (
	"github.com/BwezB/Wikno-backend/pkg/common/configuration"
)

// DATABASE CONFIG
var (
	flagDatabaseHost     = configuration.NewFlag("database_host", "", "Database Host")
	flagDatabasePort     = configuration.NewFlag("database_port", "", "Database Port")
	flagDatabaseUser     = configuration.NewFlag("database_user", "", "Database User")
	flagDatabasePassword = configuration.NewFlag("database_password", "", "Database Password")
	flagDatabaseName     = configuration.NewFlag("database_name", "", "Database Name")
)

// SERVER CONFIG
var (
	flagServerHost = configuration.NewFlag("server_host", "", "Server Host")
	flagServerPort = configuration.NewFlag("server_port", "", "Server Port")
)
