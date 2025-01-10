package config

import (
	c "github.com/BwezB/Wikno-backend/pkg/configs"
	e "github.com/BwezB/Wikno-backend/pkg/errors"
	g "github.com/BwezB/Wikno-backend/pkg/graph"
	h "github.com/BwezB/Wikno-backend/pkg/health"
	l "github.com/BwezB/Wikno-backend/pkg/log"

	"github.com/BwezB/Wikno-backend/internal/auth/api"
	"github.com/BwezB/Wikno-backend/internal/auth/db"
	"github.com/BwezB/Wikno-backend/internal/auth/service"
	"github.com/go-playground/validator/v10"
)

type AuthConfig struct {
	c.Common `yaml:",inline"`
	Server   api.ServerConfig
	Database db.DatabaseConfig
	Logger   l.LoggerConfig
	Health   h.HealthServiceConfig
	Service  service.ServiceConfig
	Graph    g.GraphConfig
}

func New(validator *validator.Validate) (*AuthConfig, error) {
	authConfig := &AuthConfig{}
	if err := c.LoadValidatedConfig(authConfig, validator); err != nil {
		return nil, e.Wrap("Failed to load config", err)
	}
	return authConfig, nil
}

func (a *AuthConfig) SetDefaults() {
	a.Common.SetDefaults()
	a.Server.SetDefaults()
	a.Database.SetDefaults()
	a.Logger.SetDefaults()
	a.Health.SetDefaults()
	a.Service.SetDefaults()
	a.Graph.SetDefaults()
}

func (a *AuthConfig) AddFromEnv() {
	a.Common.AddFromEnv()
	a.Server.AddFromEnv()
	a.Database.AddFromEnv()
	a.Logger.AddFromEnv()
	a.Health.AddFromEnv()
	a.Service.AddFromEnv()
	a.Graph.AddFromEnv()
}

func (a *AuthConfig) AddFromFlags() {
	a.Common.AddFromFlags()
	a.Server.AddFromFlags()
	a.Database.AddFromFlags()
	a.Logger.AddFromFlags()
	a.Health.AddFromFlags()
	a.Service.AddFromFlags()
	a.Graph.AddFromFlags()
}
