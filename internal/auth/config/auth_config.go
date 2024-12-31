package config

import (
	"github.com/BwezB/Wikno-backend/pkg/common/configs"
	"github.com/BwezB/Wikno-backend/pkg/log"
)

type AuthConfig struct {
	configs.BaseConfig
	Server   Server
	Database Database
}

func New() (*AuthConfig, error) {
	defer log.DebugFunc()()
	authConfig := &AuthConfig{}
	if err := configs.LoadValidatedConfig(authConfig); err != nil {
		return nil, log.Errore(err, "Failed to load config - authConfig")
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
