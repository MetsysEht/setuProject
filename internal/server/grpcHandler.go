package server

import (
	"context"

	"github.com/MetsysEht/setuProject/internal/boot"
	setuGateway2 "github.com/MetsysEht/setuProject/internal/gateway/setuGateway"
	"github.com/MetsysEht/setuProject/internal/health"
	"github.com/MetsysEht/setuProject/internal/kycVerification"
	healthv1 "github.com/MetsysEht/setuProject/rpc/health"
	kycverificationv1 "github.com/MetsysEht/setuProject/rpc/kycVerification"
	"google.golang.org/grpc"
)

func GrpcHandlerFunc(_ context.Context) func(server *grpc.Server) error {

	healthService := health.NewService()
	setuGateway := setuGateway2.NewGateway(boot.Config.SetuGatewayService)
	kycVerifRepo := kycVerification.NewRepository(boot.DB)
	kycVerifManager := kycVerification.NewManager(kycVerifRepo, setuGateway)

	return func(server *grpc.Server) error {
		healthv1.RegisterHealthServiceServer(server, health.NewServer(*healthService))
		kycverificationv1.RegisterKYCVerificationServiceServer(server, kycVerification.NewServer(kycVerifManager))
		return nil
	}
}
