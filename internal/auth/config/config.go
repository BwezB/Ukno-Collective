// Package config holds all the configuration for the auth service.
// It follows the order of flag, environment variable, default value (defined here).
package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

// Defult vales for config
const (
	defaultConfigPath = "config.yaml"

	defaultEnvironment = "production"
)

// Environment variables for config
const (
	envConfigPath = "CONFIG_PATH"

	envEnvironment = "ENVIRONMENT"
)

// Flags for config
var (
	flagConfigPath = flag.String("config", "", "Path to config file")

	flagEnvironment = flag.String("env", "", "development, test or production")
)

// Set up the validator
var configValidate *validator.Validate

func init() {
	configValidate = validator.New()
}

// Config holds all the configuration for the auth service
type Config struct {
	Environment string         `yaml:"environment" validate:"required,oneof=development test production"`
	Database    DatabaseConfig `yaml:"database" validate:"required"`
}

func New() (*Config, error) {
	if !flag.Parsed() {
		flag.Parse() // Parse the flags if they have not been parsed
	}
	config := Config{} // New config

	// Set the environment
	config.Environment = getConfigValue(*flagEnvironment, envEnvironment, "", defaultEnvironment)

	// fileConfig (config from .yaml file) to override default values
	fileConfig := &Config{}
	if config.Environment != "production" {
		configFilePath := getConfigValue(*flagConfigPath, envConfigPath, "", defaultConfigPath)

		var err error
		fileConfig, err = loadConfigFromFile(configFilePath)
		if err != nil {
			return nil, fmt.Errorf("could not load config from file: %w", err)
		}
	}
	
	// Set the database config
	dbConfig, err := NewDatabaseConfig(fileConfig.Database)
	if err != nil {
		return nil, err
	}
	config.Database = *dbConfig

	// Validate the config
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("config validation error: %w", err)
	}

	return &config, nil
}

func loadConfigFromFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not load config file: %w", err)
	}

	config := Config{}
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("could not unmarshal config: %w", err)
	}

	return &config, nil
}

// getConfigValue returns the value of a flag, environment variable or default value (in that order)
func getConfigValue(flagValue, envKey, fileValue, defaultValue string) string {
	if flagValue != "" {
		return flagValue
	}

	if envValue := os.Getenv(envKey); envValue != "" {
		return envValue
	}

	if fileValue != "" {
		return fileValue
	}

	return defaultValue
}

func (c *Config) Validate() error {
	if err := configValidate.Struct(c); err != nil {
		return fmt.Errorf("config validation error: %w", err)
	}

	return nil
}
