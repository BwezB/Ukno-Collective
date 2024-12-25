package configuration

import (
	"flag"
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

// SET UP VALIDATOR
var validate = validator.New()

// METHODS FOR SETTING VALUES

// SetEnvValue sets the value of the previous value to the environment value if the environment value is not empty.
func SetEnvValue(prevValue *string, envKey string) {
	if env := os.Getenv(envKey); env != "" {
		*prevValue = env
	}
}

// ParseFlags parses the flags if they have not been parsed.
func ParseFlags() {
	if !flag.Parsed() {
		flag.Parse()
	}
}

// NewFlag creates a new flag with the given name, value, and usage.
func NewFlag(name, value, usage string) *string {
	return flag.String(name, value, usage)
}

// SetFlagValue sets the value of the previous value to the flag value if the flag value is not empty.
func SetFlagValue(prevValue *string, flag *string) {
	if *flag != "" {
		*prevValue = *flag
	}
}

// METHODS FOR LOADING CONFIG

// LoadConfig loads the given config from the defaults<file<env<flags.
func LoadValidatedConfig[T Configurable](config T) error {
	config.SetDefaults()

	if err := getFileConfig(config); err != nil { // Sets just the fields that are in the file.
		return fmt.Errorf("failed to get config from file: %w", err)
	}
	config.AddFromEnv()
	config.AddFromFlags()

	if err := validate.Struct(config); err != nil {
		return fmt.Errorf("failed to validate config: %w", err)
	}

	return nil
}

// METHODS FOR GETTING CONFIG FROM FILE

// getFileConfig fills the given config from the yaml file.
func getFileConfig[T Configurable](config T) error {
	configFilePath := getConfigFilePath()
	configYamlFile, err := os.ReadFile(configFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // No config file, dont change the config
		}
		return fmt.Errorf("failed to read config file: %w", err)
	}

	err = yaml.Unmarshal(configYamlFile, config)
	if err != nil {
		return fmt.Errorf("failed to unmarshal config file: %w", err)
	}
	return nil
}

var flagConfigFilePath = flag.String("config-file-path", "", "Config file path")

func getConfigFilePath() string {
	configFilePath := "config.yaml"
	SetEnvValue(&configFilePath, "CONFIG_FILE_PATH")
	SetFlagValue(&configFilePath, flagConfigFilePath)
	return configFilePath
}
