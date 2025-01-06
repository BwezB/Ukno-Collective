package service

import (
	"time"
	c "github.com/BwezB/Wikno-backend/pkg/configs"
)

type ServiceConfig struct {
	jwtSecret string
	jwtExpiry    time.Duration
}

func (sc *ServiceConfig) SetDefaults() {
	// Left out jwtSecretKey for security reasons
	sc.jwtExpiry = 24 * time.Hour
}

func (sc *ServiceConfig) AddFromEnv() {
	c.SetEnvValue(&sc.jwtSecret, "JWT_SECRET")
	c.SetEnvValue(&sc.jwtExpiry, "JWT_EXPIRY")
}

var (
	flagJWTSecret = c.NewFlag("jwt_secret", "", "Secret key for JWT")
	flagJWTExpiry = c.NewFlag("jwt_expiry", "", "Expiry time for JWT")
)
func (sc *ServiceConfig) AddFromFlags() {
	c.SetFlagValue(&sc.jwtSecret, flagJWTSecret)
	c.SetFlagValue(&sc.jwtExpiry, flagJWTExpiry)
}
