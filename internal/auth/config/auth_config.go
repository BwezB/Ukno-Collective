package config

import (
	"fmt"

	"github.com/BwezB/Wikno-backend/pkg/common/baseconfig"
)

type AuthConfig struct {
	baseconfig.BaseConfig
	Server   Server
	Database Database
}

func New() (*AuthConfig, error) {
	authConfig := &AuthConfig{}
	if err := baseconfig.LoadValidatedConfig(authConfig); err != nil {
		return nil, fmt.Errorf("failed to load auth config: %w", err)
	}
	return authConfig, nil
}

func (a *AuthConfig) SetDefaults() {
	a.BaseConfig.SetDefaults()
	a.Server.SetDefaults()
	a.Database.SetDefaults()
}

func (a *AuthConfig) AddFromEnv() {
	a.BaseConfig.AddFromEnv()
	a.Server.AddFromEnv()
	a.Database.AddFromEnv()
}

func (a *AuthConfig) AddFromFlags() {
	a.BaseConfig.AddFromFlags()
	a.Server.AddFromFlags()
	a.Database.AddFromFlags()
}
