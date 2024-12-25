package config

import (
	"fmt"
	"github.com/BwezB/Wikno-backend/pkg/common/configuration"
)

type AuthConfig struct {
	configuration.BaseConfig
	ServerConfig   ServerConfig
	DatabaseConfig DatabaseConfig
}

func New() (*AuthConfig, error) {
	authConfig := &AuthConfig{}
	if err := configuration.LoadValidatedConfig(authConfig); err != nil {
		return nil, fmt.Errorf("failed to load auth config: %w", err)
	}

	return authConfig, nil
}

func (a *AuthConfig) SetDefaults() {
	a.BaseConfig.SetDefaults()
	a.ServerConfig.SetDefaults()
	a.DatabaseConfig.SetDefaults()
}

func (a *AuthConfig) AddFromEnv() {
	a.BaseConfig.AddFromEnv()
	a.ServerConfig.AddFromEnv()
	a.DatabaseConfig.AddFromEnv()
}

func (a *AuthConfig) AddFromFlags() {
	a.BaseConfig.AddFromFlags()
	a.ServerConfig.AddFromFlags()
	a.DatabaseConfig.AddFromFlags()
}
