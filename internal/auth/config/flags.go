package config
import (
	"github.com/BwezB/Wikno-backend/pkg/common/configs"
)

// DATABASE CONFIG
var (
	flagDatabaseHost     = configs.NewFlag("database_host", "", "Database Host")
	flagDatabasePort     = configs.NewFlag("database_port", "", "Database Port")
	flagDatabaseUser     = configs.NewFlag("database_user", "", "Database User")
	flagDatabasePassword = configs.NewFlag("database_password", "", "Database Password")
	flagDatabaseName     = configs.NewFlag("database_name", "", "Database Name")
)

// SERVER CONFIG
var (
	flagServerHost = configs.NewFlag("server_host", "", "Server Host")
	flagServerPort = configs.NewFlag("server_port", "", "Server Port")
)
