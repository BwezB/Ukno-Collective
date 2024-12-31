// Configuration precedence (highest to lowest):
// 1. Command line flags
// 2. Environment variables
// 3. Configuration file
// 4. Default values
package configs

import (
	"fmt"
)

// CONFIGURABLE INTERFACE
type Configurable interface {
	SetDefaults()
	AddFromEnv()
	AddFromFlags()
}

// DEFINING BASE CONFIG
type BaseConfig struct {
	Environment string `yaml:"environment" validate:"required" json:"-"`
}

func New() (*BaseConfig, error) {
	baseConfig := &BaseConfig{}
	if err := LoadValidatedConfig(baseConfig); err != nil {
		return nil, fmt.Errorf("failed to load base config: %w", err)
	}
	return baseConfig, nil
}

func (c *BaseConfig) SetDefaults() {
	c.Environment = defaultEnvironment
}

func (c *BaseConfig) AddFromEnv() {
	SetEnvValue(&c.Environment, envEnvironment)
}

func (c *BaseConfig) AddFromFlags() {
	ParseFlags()

	SetFlagValue(&c.Environment, flagEnvironment)
}

// DEFAULTS
const (
	// defaultEnvironment is the default value for the environment.
	defaultEnvironment = "production"
	// defaultConfigFilePath is the default value for the config file path.
	defaultConfigFilePath = "config.yaml"
)

// ENV
const (
	// envEnvironment is the environment variable for the environment.
	envEnvironment = "ENVIRONMENT"
	// envConfigFilePath is the environment variable for the config file path.
	envConfigFilePath = "CONFIG_FILE_PATH"
)

// FLAGS
var (
	// flagEnvironment is the flag for the environment.
	flagEnvironment = NewFlag("environment", "", "Environment")
	// flagConfigFilePath is the flag for the config file path.
	flagConfigFilePath = NewFlag("config-file-path", "", "Config file path")
)