package config

import (
	a "github.com/BwezB/Wikno-backend/pkg/auth"
	c "github.com/BwezB/Wikno-backend/pkg/configs"
	e "github.com/BwezB/Wikno-backend/pkg/errors"
	h "github.com/BwezB/Wikno-backend/pkg/health"
	l "github.com/BwezB/Wikno-backend/pkg/log"

	"github.com/BwezB/Wikno-backend/internal/graph/api"
	"github.com/BwezB/Wikno-backend/internal/graph/db"
	"github.com/go-playground/validator/v10"
)

type GraphConfig struct {
	c.Common `yaml:",inline"`
	Server   api.ServerConfig
	Database db.DatabaseConfig
	Logger   l.LoggerConfig
	Health   h.HealthServiceConfig
	Auth     a.AuthConfig
}

func New(validator *validator.Validate) (*GraphConfig, error) {
	graphConfig := &GraphConfig{}
	if err := c.LoadValidatedConfig(graphConfig, validator); err != nil {
		return nil, e.Wrap("Failed to load config", err)
	}
	return graphConfig, nil
}

func (a *GraphConfig) SetDefaults() {
	a.Common.SetDefaults()
	a.Server.SetDefaults()
	a.Database.SetDefaults()
	a.Logger.SetDefaults()
	a.Health.SetDefaults()
	a.Auth.SetDefaults()
}

func (a *GraphConfig) AddFromEnv() {
	a.Common.AddFromEnv()
	a.Server.AddFromEnv()
	a.Database.AddFromEnv()
	a.Logger.AddFromEnv()
	a.Health.AddFromEnv()
	a.Auth.AddFromEnv()
}

func (a *GraphConfig) AddFromFlags() {
	a.Common.AddFromFlags()
	a.Server.AddFromFlags()
	a.Database.AddFromFlags()
	a.Logger.AddFromFlags()
	a.Health.AddFromFlags()
	a.Auth.AddFromFlags()
}
