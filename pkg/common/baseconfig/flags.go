package baseconfig

// BASE CONFIG
var (
	// flagEnvironment is the flag for the environment.
	flagEnvironment = NewFlag("environment", "", "Environment")
	// flagConfigFilePath is the flag for the config file path.
	flagConfigFilePath = NewFlag("config-file-path", "", "Config file path")
)

// LOGGER
var (
	// flagLogLevel is the flag for the log level.
	flagLogLevel = NewFlag("log-level", "", "Log Level")
)