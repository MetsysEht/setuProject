package health

import (
	"context"
	"fmt"

	"github.com/MetsysEht/setuProject/pkg/healthcheck"
	healthv1 "github.com/MetsysEht/setuProject/rpc/health"
)

// Server has methods implementing of server rpc.
type Server struct {
	healthv1.HealthServiceServer
	check Service
}

// NewServer returns a server.
func NewServer(check Service) *Server {
	return &Server{check: check}
}

// LivenessCheck returns service's serving status.
func (s *Server) LivenessCheck(ctx context.Context, _ *healthv1.LivenessRequest) (*healthv1.LivenessResponse, error) {
	response := s.check.Live(ctx)
	if !response.Status {
		resp := &healthv1.LivenessResponse{
			Status:       healthv1.ServingStatus_SERVING_STATUS_NOT_SERVING,
			StatusChecks: toRPCStatusChecks(response.StatusChecks),
		}
		return nil, fmt.Errorf("liveness check failed, resp: %+v", resp)
	}
	return &healthv1.LivenessResponse{
		Status:       healthv1.ServingStatus_SERVING_STATUS_SERVING,
		StatusChecks: toRPCStatusChecks(response.StatusChecks),
	}, nil
}

// ReadinessCheck is the rpc method implemented
func (s *Server) ReadinessCheck(ctx context.Context, _ *healthv1.ReadinessRequest) (*healthv1.ReadinessResponse, error) {
	response := s.check.Ready(ctx)
	if !response.Status {
		resp := &healthv1.ReadinessResponse{
			Status:       healthv1.ServingStatus_SERVING_STATUS_NOT_SERVING,
			StatusChecks: toRPCStatusChecks(response.StatusChecks),
		}
		return nil, fmt.Errorf("readiness check failed, resp: %+v", resp)
	}

	return &healthv1.ReadinessResponse{
		Status:       healthv1.ServingStatus_SERVING_STATUS_SERVING,
		StatusChecks: toRPCStatusChecks(response.StatusChecks),
	}, nil
}

func toRPCStatusChecks(checks []healthcheck.StatusCheck) []*healthv1.StatusCheck {
	var statusChecks []*healthv1.StatusCheck
	for _, check := range checks {
		if check.Error == nil {
			statusChecks = append(
				statusChecks, &healthv1.StatusCheck{
					Name:   check.Name,
					Status: "ok",
				},
			)
			continue
		}
		statusChecks = append(
			statusChecks, &healthv1.StatusCheck{
				Name:   check.Name,
				Status: fmt.Sprintf("error: %v", *check.Error),
			},
		)
	}

	return statusChecks
}
