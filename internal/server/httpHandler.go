package server

import (
	"context"

	"github.com/MetsysEht/setuProject/internal/boot"
	healthv1 "github.com/MetsysEht/setuProject/rpc/health"
	kyc_verificationv1 "github.com/MetsysEht/setuProject/rpc/kycVerification"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func HttpHandlerFunc(ctx context.Context) func(mux *runtime.ServeMux) error {

	return func(mux *runtime.ServeMux) error {
		optsWithoutTracing := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
		if err := healthv1.RegisterHealthServiceHandlerFromEndpoint(ctx, mux, boot.Config.App.Interfaces.Service.GrpcServerAddress, optsWithoutTracing); err != nil {
			return err
		}
		if err := kyc_verificationv1.RegisterKYCVerificationServiceHandlerFromEndpoint(ctx, mux, boot.Config.App.Interfaces.Service.GrpcServerAddress, optsWithoutTracing); err != nil {
			return err
		}
		return nil
	}
}
