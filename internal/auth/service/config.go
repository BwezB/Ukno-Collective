package service

import (
	"time"
	c "github.com/BwezB/Wikno-backend/pkg/configs"
)

type ServiceConfig struct {
	jwtSecret 	string
	jwtExpiry   time.Duration
	email 	 	string // The email the authservice uses to send requests to the graph service
	password    string // The password the authservice uses to send requests to the graph service
}

func (sc *ServiceConfig) SetDefaults() {
	// Left out jwtSecretKey for security reasons
	sc.jwtExpiry = 24 * time.Hour
	sc.email = "authservice@wikno.com"
	// Left out password for security reasons
}

func (sc *ServiceConfig) AddFromEnv() {
	c.SetEnvValue(&sc.jwtSecret, "JWT_SECRET")
	c.SetEnvValue(&sc.jwtExpiry, "JWT_EXPIRY")
	c.SetEnvValue(&sc.email, "AUTH_EMAIL")
	c.SetEnvValue(&sc.password, "AUTH_PASSWORD")
}

var (
	flagJWTSecret = c.NewFlag("jwt-secret", "", "Secret key for JWT")
	flagJWTExpiry = c.NewFlag("jwt-expiry", "", "Expiry time for JWT")
	flagEmail = c.NewFlag("auth-email", "", "Email for the auth service")
	flagPassword = c.NewFlag("auth-password", "", "Password for the auth service")
)
func (sc *ServiceConfig) AddFromFlags() {
	c.SetFlagValue(&sc.jwtSecret, flagJWTSecret)
	c.SetFlagValue(&sc.jwtExpiry, flagJWTExpiry)
	c.SetFlagValue(&sc.email, flagEmail)
	c.SetFlagValue(&sc.password, flagPassword)
}
