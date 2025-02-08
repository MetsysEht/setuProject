package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/MetsysEht/setuProject/internal/boot"
	"github.com/MetsysEht/setuProject/internal/server"
	"github.com/MetsysEht/setuProject/pkg/logger"
	grpcprometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
)

func main() {
	// Initialize context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	boot.Initialize()
	s, err := server.NewServer(
		boot.Config.App.Interfaces.Service,
		server.GrpcHandlerFunc(ctx),
		server.HttpHandlerFunc(ctx),
		getInterceptors(ctx)...,
	)
	if err != nil {
		logger.L.Fatalf("failed to create new server: %v", err)
	}

	err = s.Start()
	if err != nil {
		logger.L.Fatalf("failed to start server: %v", err)
	}
	logger.L.Infof("Server up and running...")
	// accept graceful shutdowns when quit via SIGINT (Ctrl+C) or SIGTERM.
	// SIGKILL, SIGQUIT will not be caught.
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)

	// Block until signal is received.
	<-c
	err = s.Stop(ctx)
	if err != nil {
		panic(err)
	}
}

func getInterceptors(_ context.Context) []grpc.UnaryServerInterceptor {
	return []grpc.UnaryServerInterceptor{
		grpcprometheus.UnaryServerInterceptor,
	}
}
