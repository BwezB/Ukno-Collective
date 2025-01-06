package config

import (
	c "github.com/BwezB/Wikno-backend/pkg/configs"
	e "github.com/BwezB/Wikno-backend/pkg/errors"
	h "github.com/BwezB/Wikno-backend/pkg/health"
	l "github.com/BwezB/Wikno-backend/pkg/log"

	"github.com/go-playground/validator/v10"
	"github.com/BwezB/Wikno-backend/internal/auth/api"
	"github.com/BwezB/Wikno-backend/internal/auth/db"
	"github.com/BwezB/Wikno-backend/internal/auth/service"
)

type AuthConfig struct {
	c.BaseConfig
	Server   api.ServerConfig
	Database db.DatabaseConfig
	Logger   l.LoggerConfig
	Health   h.HealthServiceConfig
	Service  service.ServiceConfig
}

func New(validator *validator.Validate) (*AuthConfig, error) {
	authConfig := &AuthConfig{}
	if err := c.LoadValidatedConfig(authConfig, validator); err != nil {
		return nil, e.Wrap("Failed to load config", err)
	}
	return authConfig, nil
}

func (a *AuthConfig) SetDefaults() {
	a.BaseConfig.SetDefaults()
	a.Server.SetDefaults()
	a.Database.SetDefaults()
	a.Logger.SetDefaults()
	a.Health.SetDefaults()
	a.Service.SetDefaults()
}

func (a *AuthConfig) AddFromEnv() {
	a.BaseConfig.AddFromEnv()
	a.Server.AddFromEnv()
	a.Database.AddFromEnv()
	a.Logger.AddFromEnv()
	a.Health.AddFromEnv()
	a.Service.AddFromEnv()
}

func (a *AuthConfig) AddFromFlags() {
	a.BaseConfig.AddFromFlags()
	a.Server.AddFromFlags()
	a.Database.AddFromFlags()
	a.Logger.AddFromFlags()
	a.Health.AddFromFlags()
	a.Service.AddFromFlags()
}
