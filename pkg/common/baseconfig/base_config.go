// Configuration precedence (highest to lowest):
// 1. Command line flags
// 2. Environment variables
// 3. Configuration file
// 4. Default values
package baseconfig

import "fmt"

// CONFIGURABLE INTERFACE
type Configurable interface {
	SetDefaults()
	AddFromEnv()
	AddFromFlags()
}

// DEFINING BASE CONFIG
type BaseConfig struct {
	Environment string `yaml:"environment" validate:"required" json:"-"`
	Logger	    Logger `yaml:"logger" validate:"required"`
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
