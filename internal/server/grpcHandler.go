package server

import (
	"context"

	"github.com/MetsysEht/setuProject/internal/health"
	healthv1 "github.com/MetsysEht/setuProject/rpc/health"
	"google.golang.org/grpc"
)

func GrpcHandlerFunc(_ context.Context) func(server *grpc.Server) error {

	healthService := health.NewService()

	return func(server *grpc.Server) error {
		healthv1.RegisterHealthServiceServer(server, health.NewServer(*healthService))
		return nil
	}
}
