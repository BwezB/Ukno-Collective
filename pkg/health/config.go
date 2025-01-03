package health

import (
	"time"
	c "github.com/BwezB/Wikno-backend/pkg/configs"
)

// HealthServiceConfig is a struct that represents the configuration for a HealthService
type HealthServiceConfig struct {
	HealthCheckInterval time.Duration
}


// DEFAULTS

func (hsc *HealthServiceConfig) SetDefaults() {
	hsc.HealthCheckInterval = 5 * time.Second
}


// ENV

func (hsc *HealthServiceConfig) AddFromEnv() {
	c.SetEnvValue(&hsc.HealthCheckInterval, "HEALTH_CHECK_INTERVAL")
}


// FLAGS

var (
	// flagHealthCheckInterval is the flag for the health check interval.
	flagHealthCheckInterval = c.NewFlag("health-check-interval", "", "Health check interval")
)
func (hsc *HealthServiceConfig) AddFromFlags() {
	c.SetFlagValue(&hsc.HealthCheckInterval, flagHealthCheckInterval)
}
