package config

import (
	c "github.com/BwezB/Wikno-backend/pkg/configs"
	l "github.com/BwezB/Wikno-backend/pkg/log"
	e "github.com/BwezB/Wikno-backend/pkg/errors"

	"github.com/BwezB/Wikno-backend/internal/auth/db"
	"github.com/BwezB/Wikno-backend/internal/auth/api"
)

type AuthConfig struct {
	c.BaseConfig
	Server   api.ServerConfig
	Database db.DatabaseConfig
	Logger   l.LoggerConfig
}

func New() (*AuthConfig, error) {
	authConfig := &AuthConfig{}
	if err := c.LoadValidatedConfig(authConfig); err != nil {
		return nil, e.Wrap("Failed to load config", err)
	}
	return authConfig, nil
}

func (a *AuthConfig) SetDefaults() {
	a.BaseConfig.SetDefaults()
	a.Server.SetDefaults()
	a.Database.SetDefaults()
	a.Logger.SetDefaults()
}

func (a *AuthConfig) AddFromEnv() {
	a.BaseConfig.AddFromEnv()
	a.Server.AddFromEnv()
	a.Database.AddFromEnv()
	a.Logger.AddFromEnv()
}

func (a *AuthConfig) AddFromFlags() {
	a.BaseConfig.AddFromFlags()
	a.Server.AddFromFlags()
	a.Database.AddFromFlags()
	a.Logger.AddFromFlags()
}
