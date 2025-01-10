package service

import (
	"time"
	c "github.com/BwezB/Wikno-backend/pkg/configs"
	g "github.com/BwezB/Wikno-backend/pkg/graph"
)

type ServiceConfig struct {
	jwtSecret 	string
	jwtExpiry   time.Duration
	email 	 	string // The email the authservice uses to send requests to the graph service
	password    string // The password the authservice uses to send requests to the graph service
	graph       g.GraphConfig
}

func (sc *ServiceConfig) SetDefaults() {
	// Left out jwtSecretKey for security reasons
	sc.jwtExpiry = 24 * time.Hour
	sc.email = "authservice@wikno.com"
	// Left out password for security reasons
	sc.graph.SetDefaults()
}

func (sc *ServiceConfig) AddFromEnv() {
	c.SetEnvValue(&sc.jwtSecret, "JWT_SECRET")
	c.SetEnvValue(&sc.jwtExpiry, "JWT_EXPIRY")
	c.SetEnvValue(&sc.email, "AUTH_SERVICE_EMAIL")
	c.SetEnvValue(&sc.password, "AUTH_SERVICE_PASSWORD")
	sc.graph.AddFromEnv()
}

var (
	flagJWTSecret = c.NewFlag("jwt_secret", "", "Secret key for JWT")
	flagJWTExpiry = c.NewFlag("jwt_expiry", "", "Expiry time for JWT")
	flagEmail = c.NewFlag("email", "", "Email for the auth service")
	flagPassword = c.NewFlag("password", "", "Password for the auth service")
)
func (sc *ServiceConfig) AddFromFlags() {
	c.SetFlagValue(&sc.jwtSecret, flagJWTSecret)
	c.SetFlagValue(&sc.jwtExpiry, flagJWTExpiry)
	c.SetFlagValue(&sc.email, flagEmail)
	c.SetFlagValue(&sc.password, flagPassword)
	sc.graph.AddFromFlags()
}
