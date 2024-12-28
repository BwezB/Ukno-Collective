package config
import (
	"github.com/BwezB/Wikno-backend/pkg/common/baseconfig"
)

// DATABASE CONFIG
var (
	flagDatabaseHost     = baseconfig.NewFlag("database_host", "", "Database Host")
	flagDatabasePort     = baseconfig.NewFlag("database_port", "", "Database Port")
	flagDatabaseUser     = baseconfig.NewFlag("database_user", "", "Database User")
	flagDatabasePassword = baseconfig.NewFlag("database_password", "", "Database Password")
	flagDatabaseName     = baseconfig.NewFlag("database_name", "", "Database Name")
)

// SERVER CONFIG
var (
	flagServerHost = baseconfig.NewFlag("server_host", "", "Server Host")
	flagServerPort = baseconfig.NewFlag("server_port", "", "Server Port")
)
