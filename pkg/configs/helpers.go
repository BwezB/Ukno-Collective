package configs

import (
	"flag"
	"fmt"
	"os"
	"time"
	"strconv"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

// SET UP VALIDATOR
var validate = validator.New()


// ENVIRONMENT VARIABLES

// SetEnvValue sets the value of the previous value to the environment value if the environment value is not empty.
func SetEnvValue[T *string | *int | *bool | *time.Duration] (prevValue T, envKey string) {
	env := os.Getenv(envKey)
	if env == "" {
		return // No env value, dont change the value
	}

	switch v := any(prevValue).(type) {
	case *string:
		*v = env
	case *int:
		parsed, err := strconv.Atoi(env)
		if err != nil {
			panic(fmt.Errorf("invalid int value for %s: %w", envKey, err))
		}
		*v = parsed
	case *time.Duration:
		parsed, err := time.ParseDuration(env)
		if err != nil {
			panic(fmt.Errorf("invalid duration value for %s: %w", envKey, err))
		}
		*v = parsed
	case *bool:
		parsed, err := strconv.ParseBool(env)
		if err != nil {
			panic(fmt.Errorf("invalid bool value for %s: %w", envKey, err))
		}
		*v = parsed
	default:
		panic(fmt.Errorf("unsupported type for SetEnvValue: %T", v))
	}
}


// FLAGS

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
func SetFlagValue[T *string | *int | *bool | *time.Duration] (prevValue T, flag *string) {
	if *flag == "" {
		return // No flag value, dont change the value
	}

	switch v := any(prevValue).(type) {
	case *string:
		*v = *flag
	case *int:
		parsed, err := strconv.Atoi(*flag)
		if err != nil {
			panic(fmt.Errorf("invalid int value for %s: %w", *flag, err))
		}
		*v = parsed
	case *time.Duration:
		parsed, err := time.ParseDuration(*flag)
		if err != nil {
			panic(fmt.Errorf("invalid duration value for %s: %w", *flag, err))
		}
		*v = parsed
	case *bool:
		parsed, err := strconv.ParseBool(*flag)
		if err != nil {
			panic(fmt.Errorf("invalid bool value for %s: %w", *flag, err))
		}
		*v = parsed
	default:
		panic(fmt.Errorf("unsupported type for SetFlagValue: %T", v))
	}
}


// METHODS FOR LOADING CONFIG

// LoadValidatedConfig loads the given (ANY) config from the defaults<file<env<flags.
func LoadValidatedConfig[T Configurable](config T) error {
	config.SetDefaults()

	if err := getFileConfig(config); err != nil { // Sets just the fields that are in the file.
		return fmt.Errorf("failed to get config from file: %w", err)
	}

	config.AddFromEnv()
	ParseFlags()
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

// getConfigFilePath gets the config file path from the flag, environment, or default.
func getConfigFilePath() string {
	ParseFlags()

	configFilePath := defaultConfigFilePath
	SetEnvValue(&configFilePath, envConfigFilePath)
	SetFlagValue(&configFilePath, flagConfigFilePath)
	return configFilePath
}
