package health

import (
	"context"
	"log"

	"github.com/MetsysEht/setuProject/internal/boot"
	"github.com/MetsysEht/setuProject/pkg/healthcheck"
)

// Service has the core healthcheck checking logic
type Service struct {
	HealthCheck healthcheck.HealthCheck
}

// Response wraps the response from healthcheck lib
type Response healthcheck.Response

// NewService uses the generic package healthcheck to add checks
// that are required for healthcheck checking this application
func NewService() *Service {
	h := healthcheck.New()
	h.AddReadinessCheck(
		"db", func() error {
			sqlDB, err := boot.DB.DB()
			if err != nil {
				log.Fatalf("Failed to get SQL database instance: %v", err)
			}
			return sqlDB.Ping()
		},
	)

	return &Service{HealthCheck: h}
}

// Ready executes the readiness checks and give back response
func (s *Service) Ready(ctx context.Context) Response {
	resp := s.HealthCheck.Ready(ctx)
	return Response(resp)
}

// Live executes the liveliness checks and give back response
func (s *Service) Live(ctx context.Context) Response {
	resp := s.HealthCheck.Live(ctx)
	return Response(resp)
}
