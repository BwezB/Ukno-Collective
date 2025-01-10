// Configuration precedence (highest to lowest):
// 1. Command line flags
// 2. Environment variables
// 3. Configuration file
// 4. Default values
package configs

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// CONFIGURABLE INTERFACE
type Configurable interface {
	SetDefaults()
	AddFromEnv()
	AddFromFlags()
}

// DEFINING BASE CONFIG
type Common struct {
	Environment string `yaml:"environment" validate:"required" json:"-"`
}

func New(validator *validator.Validate) (*Common, error) {
	baseConfig := &Common{}
	if err := LoadValidatedConfig(baseConfig, validator); err != nil {
		return nil, fmt.Errorf("failed to load base config: %w", err)
	}
	return baseConfig, nil
}

// DEFAULTS

// defaultConfigFilePath is the default value for the config file path.
const defaultConfigFilePath = "config.yaml"

// SetDefaults sets the default values for the base config.
func (c *Common) SetDefaults() {
	c.Environment = "production"
}

// ENVIRONMENT

// envConfigFilePath is the environment variable for the config file path.
const envConfigFilePath = "CONFIG_FILE_PATH"

func (c *Common) AddFromEnv() {
	SetEnvValue(&c.Environment, "ENVIRONMENT")
}

// FLAGS

var (
	// flagConfigFilePath is the flag for the config file path.
	flagConfigFilePath = NewFlag("config-file-path", "", "Config file path")

	flagEnvironment = NewFlag("environment", "", "Environment")
)

func (c *Common) AddFromFlags() {
	SetFlagValue(&c.Environment, flagEnvironment)
}
