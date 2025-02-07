package server

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/MetsysEht/setuProject/internal/boot"
	"github.com/MetsysEht/setuProject/internal/config"
	"github.com/MetsysEht/setuProject/pkg/logger"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcprometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

type Server struct {
	config         config.NetworkInterfaces
	internalServer *http.Server
	grpcServer     *grpc.Server
	httpServer     *http.Server
}

type (
	RegisterGrpcHandlers func(server *grpc.Server) error
	RegisterHttpHandlers func(mux *runtime.ServeMux) error
)

func NewServer(config config.NetworkInterfaces, grpcHandlers RegisterGrpcHandlers, httpHandlers RegisterHttpHandlers, interceptors ...grpc.UnaryServerInterceptor) (*Server, error) {

	grpcprometheus.EnableHandlingTimeHistogram(func(opts *prometheus.HistogramOpts) {
		opts.Name = "grpc_server_handled_duration_seconds"
		opts.Buckets = prometheus.ExponentialBuckets(0.002, 2, 19)
	})

	grpcServer, err := newGRPCServer(grpcHandlers, interceptors...)
	if err != nil {
		return nil, err
	}

	grpcprometheus.Register(grpcServer)

	httpServer, err := newHttpServer(httpHandlers)
	if err != nil {
		return nil, err
	}

	internalServer, err := newInternalServer()
	if err != nil {
		return nil, err
	}
	return &Server{
		config:         config,
		internalServer: internalServer,
		grpcServer:     grpcServer,
		httpServer:     httpServer,
	}, nil
}

func newGRPCServer(r RegisterGrpcHandlers, interceptors ...grpc.UnaryServerInterceptor) (*grpc.Server, error) {
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(interceptors...)))
	err := r(grpcServer)
	if err != nil {
		return nil, err
	}
	return grpcServer, nil
}

func newHttpServer(r RegisterHttpHandlers) (*http.Server, error) {
	mux := runtime.NewServeMux(getServeMuxOptions()...)
	if r != nil {
		err := r(mux)
		if err != nil {
			return nil, err
		}
	}
	handler := cors.AllowAll().Handler(mux)
	server := http.Server{Handler: handler}
	return &server, nil
}

func newInternalServer() (*http.Server, error) {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
	mux.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	mux.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	mux.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	mux.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
	server := http.Server{Handler: mux}

	return &server, nil
}

func (s *Server) StartGrpcServer() {
	listener, err := net.Listen("tcp", s.config.GrpcServerAddress)
	if err != nil {
		panic(err)
	}

	err = s.grpcServer.Serve(listener)
	if err != nil {
		panic(err)
	}
}

func (s *Server) StartMetricsServer() {
	listener, err := net.Listen("tcp", s.config.InternalServerAddress)
	if err != nil {
		panic(err)
	}

	err = s.internalServer.Serve(listener)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}

func (s *Server) StartHttpServer() {
	listener, err := net.Listen("tcp", s.config.HttpServerAddress)
	if err != nil {
		panic(err)
	}

	err = s.httpServer.Serve(listener)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}

func (s *Server) Start() error {
	// Start gRPC server.
	go s.StartGrpcServer()

	// Start internal HTTP server. Used for exposing prometheus metrics.
	go s.StartMetricsServer()

	// Start HTTP server for gRPC gateway.
	go s.StartHttpServer()

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	logger.L.Info("Server Shutdown")
	time.Sleep(time.Duration(boot.Config.App.ShutdownDelay) * time.Second)
	ctx, cancel := context.WithTimeout(ctx, time.Duration(boot.Config.App.ShutdownTimeout)*time.Second)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)

	// gracefully shutdown http server for gateway
	g.Go(func() error {
		return s.httpServer.Shutdown(ctx)
	})

	// gracefully shutdown internal server
	g.Go(func() error {
		return s.internalServer.Shutdown(ctx)
	})

	// gracefully shutdown grpc-server
	g.Go(func() error {
		return gracefullyShutdownGrpcServerWithTimeout(ctx, s.grpcServer)
	})

	return g.Wait()
}

func gracefullyShutdownGrpcServerWithTimeout(ctx context.Context, server *grpc.Server) error {
	stopped := make(chan bool)
	go func() {
		server.GracefulStop()
		close(stopped)
	}()

	select {
	case <-ctx.Done():
		server.Stop()
		return ctx.Err()
	case <-stopped:
		return nil
	}
}

func getServeMuxOptions() []runtime.ServeMuxOption {
	return []runtime.ServeMuxOption{}
}
