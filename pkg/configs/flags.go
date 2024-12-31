package configs

// BASE CONFIG
var (
	// flagEnvironment is the flag for the environment.
	flagEnvironment = NewFlag("environment", "", "Environment")
	// flagConfigFilePath is the flag for the config file path.
	flagConfigFilePath = NewFlag("config-file-path", "", "Config file path")
)
