package health

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	c "github.com/BwezB/Wikno-backend/pkg/context"
	l "github.com/BwezB/Wikno-backend/pkg/log"
)

// HealthStatus is a struct that represents the health status of a service
type HealthStatus struct {
	Healthy bool      `json:"healthy"`
	Err     error    `json:"error"`
	Time    time.Time `json:"time"`
}

// HealthCheckable structs can be checked for health with HealthCheck()
type HealthCheckable interface {
	HealthCheck(ctx context.Context) *HealthStatus
}

type HealthService struct {
	healthy     atomic.Bool
	checks  []HealthCheckable
	checksMu    sync.RWMutex
	illStatuses []HealthStatus
	statusesMu  sync.RWMutex
	ticker      *time.Ticker
	stopChan    chan struct{}
}

func NewHealthService(config HealthServiceConfig) *HealthService {
	hs := &HealthService{
		ticker:   time.NewTicker(config.HealthCheckInterval),
		stopChan: make(chan struct{}),
	}
	hs.healthy.Store(true)

	return hs
}

// AddCheck adds a HealthCheckable struct to the health service
func (hs *HealthService) AddCheck(check HealthCheckable) {
	hs.checksMu.Lock()
	defer hs.checksMu.Unlock()

	hs.checks = append(hs.checks, check)
}

// Start starts the health checking service in a new goroutine
func (hs *HealthService) Start() {
	for {
		select {
		case <-hs.ticker.C:
			hs.checkHealth()
		case <-hs.stopChan:
			hs.ticker.Stop()
			return
		}
	}
}

// Stop stops the health service
func (hs *HealthService) Stop() {
	close(hs.stopChan)
}

// IsHealthy returns the current health status of the service
func (hs *HealthService) IsHealthy() bool {
	return hs.healthy.Load()
}

// IllStatuses returns the health statuses of all checkables that are not healthy
func (hs *HealthService) IllStatuses() []HealthStatus {
	hs.statusesMu.RLock()
	defer hs.statusesMu.RUnlock()

	illStatuses := make([]HealthStatus, len(hs.illStatuses))
	copy(illStatuses, hs.illStatuses)

	return illStatuses
}


// HEHPER METHODS

// checkHealth checks the health of all checkables
// It updates healthy and illStatuses
func (hs *HealthService) checkHealth() {
	// Context so we can cancel the health check if it takes too long
	ctx, cancel := c.WithTimeout(c.Background(), 5*time.Second)
	defer cancel()

	hs.checksMu.RLock()
	checks := make([]HealthCheckable, len(hs.checks))
	copy(checks, hs.checks)
	hs.checksMu.RUnlock()

	// Reset illStates as we are about to check the health of all checkables
	hs.statusesMu.Lock()
	hs.illStatuses = hs.illStatuses[:0]
	hs.statusesMu.Unlock()

	allHealthy := true
	for _, check := range checks {
		status := check.HealthCheck(ctx)
		if !status.Healthy {
			allHealthy = false
			hs.statusesMu.Lock()
			hs.illStatuses = append(hs.illStatuses, *status)
			hs.statusesMu.Unlock()
			l.Warn("Health check failed", l.ErrField(status.Err))
		}
	}

	hs.healthy.Store(allHealthy)
}
