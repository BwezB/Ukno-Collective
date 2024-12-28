package configuration
// Flags for the common configuration package

var (
	// flagEnvironment is the flag for the environment.
	flagEnvironment = NewFlag("environment", "", "Environment")
	// flagConfigFilePath is the flag for the config file path.
	flagConfigFilePath = NewFlag("config-file-path", "", "Config file path")
)
