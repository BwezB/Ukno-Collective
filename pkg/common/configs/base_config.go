// Configuration precedence (highest to lowest):
// 1. Command line flags
// 2. Environment variables
// 3. Configuration file
// 4. Default values
package configs

import "github.com/BwezB/Wikno-backend/pkg/log"

// CONFIGURABLE INTERFACE
type Configurable interface {
	SetDefaults()
	AddFromEnv()
	AddFromFlags()
}

// DEFINING BASE CONFIG
type BaseConfig struct {
	Environment string `yaml:"environment" validate:"required" json:"-"`
	Logger      Logger `yaml:"logger" validate:"required"`
}

func New() (*BaseConfig, error) {
	baseConfig := &BaseConfig{}
	if err := LoadValidatedConfig(baseConfig); err != nil {
		return nil, log.Errore(err, "Failed to load config - baseConfig")
	}
	return baseConfig, nil
}

func (c *BaseConfig) SetDefaults() {
	c.Environment = defaultEnvironment
	c.Logger.SetDefaults()
}

func (c *BaseConfig) AddFromEnv() {
	SetEnvValue(&c.Environment, envEnvironment)
	c.Logger.AddFromEnv()
}

func (c *BaseConfig) AddFromFlags() {
	ParseFlags()

	SetFlagValue(&c.Environment, flagEnvironment)
	c.Logger.AddFromFlags()
}
