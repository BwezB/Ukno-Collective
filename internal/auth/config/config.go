// Package config holds all the configuration for the auth service.
// It follows the order of flag, environment variable, default value (defined here).
package config

import (
	"flag"
	"fmt"
	"os"
	"errors"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

// Defult vales for config
const (
	defaultConfigPath = "cmd/authservice/config.yaml"

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
	Server      ServerConfig   `yaml:"server" validate:"required"`
}

func New() (*Config, error) {
	if !flag.Parsed() {
		flag.Parse() // Parse the flags if they have not been parsed
	}

	// Load the config from file if available
	configFileLoaded := true
	configFilePath := getConfigValue(*flagConfigPath, envConfigPath, "", defaultConfigPath)
	fileConfig, err := loadConfigFromFile(configFilePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fileConfig = &Config{} // If file does not exist, create a new empty config
			configFileLoaded = false
		} else {
			return nil, fmt.Errorf("could not load config from file: %w", err)
		}
	}

	// CREATING CONFIG
	config := Config{}

	// Set the global config values
	config.Environment = getConfigValue(*flagEnvironment, envEnvironment, fileConfig.Environment, defaultEnvironment)

	// Set the database config
	dbConfig, err := NewDatabaseConfig(&fileConfig.Database)
	if err != nil {
		return nil, err
	}
	config.Database = *dbConfig

	// Set the server config
	serverConfig, err := NewServerConfig(&fileConfig.Server)
	if err != nil {
		return nil, err
	}
	config.Server = *serverConfig

	// VALIDATE CONFIG
	// Validate the structs
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("config validation error: %w", err)
	}

	// Check if config file is in production
	if config.Environment == "production" && configFileLoaded {
		return nil, errors.New("config file is not allowed in production environment")
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
