package healthcheck

import (
	"context"
	"net/http"
	"sync"
)

// HealthCheck is the interface for healthcheck checks, default impl is below
type HealthCheck interface {
	// AddReadinessCheck adds a check that indicates that this instance of the
	// application is currently unable to serve requests because of an upstream
	// or some transient failure. If a readiness check fails, this instance
	// should no longer receive requests, but should not be restarted or
	// destroyed.
	AddReadinessCheck(name string, check Check)

	// AddLivelinessCheck adds a check that indicates that this instance of the
	// application should be destroyed or restarted. A failed liveliness check
	// indicates that this instance is unhealthy, not some upstream dependency.
	// Every liveliness check is also included as a readiness check.
	AddLivelinessCheck(name string, check Check)

	// Ready runs the healthcheck on the readiness checks configured
	Ready(ctx context.Context) Response

	// Live runs the healthcheck on the liveliness checks configured
	Live(ctx context.Context) Response
}

// Response for the healthcheck check
type Response struct {
	// Status returns if the healthcheck was successful or not
	Status bool

	// Status per check is also returned to the clients
	StatusChecks []StatusCheck
}

// check implements HealthCheck
type check struct {
	http.ServeMux
	mutex            sync.RWMutex
	livelinessChecks map[string]Check
	readinessChecks  map[string]Check
}

// New creates a new basic Handler
func New() HealthCheck {
	return &check{
		livelinessChecks: make(map[string]Check),
		readinessChecks:  make(map[string]Check),
	}
}

// AddReadinessCheck registers new readiness checks
func (c *check) AddReadinessCheck(name string, check Check) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.readinessChecks[name] = check
}

// AddLivelinessCheck registers new liveliness checks
func (c *check) AddLivelinessCheck(name string, check Check) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.livelinessChecks[name] = check
}

// Ready runs all the readiness checks with liveliness checks combined
func (c *check) Ready(ctx context.Context) Response {
	return c.run(ctx, c.readinessChecks, c.livelinessChecks)
}

// Live runs all the liveliness checks
func (c *check) Live(ctx context.Context) Response {
	return c.run(ctx, c.livelinessChecks)
}

// pointer is a helper func
func pointer(s string) *string {
	return &s
}

// collect checks the main checks
func (c *check) collectChecks(checks map[string]Check) (bool, []StatusCheck) {

	c.mutex.RLock()
	defer c.mutex.RUnlock()

	status := true
	var statusChecks []StatusCheck

	for name, check := range checks {
		err := check()
		if err != nil {
			status = false
			statusChecks = append(
				statusChecks, StatusCheck{
					Name:  name,
					Error: pointer(err.Error()),
				},
			)
		} else {
			statusChecks = append(
				statusChecks, StatusCheck{
					Name:  name,
					Error: nil,
				},
			)
		}
	}

	return status, statusChecks
}

// run runs the checks
func (c *check) run(ctx context.Context, checks ...map[string]Check) Response {

	response := Response{Status: true, StatusChecks: []StatusCheck{}}

	select {
	case <-ctx.Done():
		return Response{Status: true}
	default:
		for _, checks := range checks {
			status, statusChecks := c.collectChecks(checks)
			if !status {
				response.Status = status
				response.StatusChecks = statusChecks
				return response // no need to check all, return fast
			}
			response.StatusChecks = append(
				response.StatusChecks,
				statusChecks...,
			)
		}
		return response
	}
}
