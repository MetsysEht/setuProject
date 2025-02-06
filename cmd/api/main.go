package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MetsysEht/setuProject/internal/boot"
	"github.com/MetsysEht/setuProject/internal/server"
	"github.com/MetsysEht/setuProject/pkg/logger"
)

func main() {
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	boot.Initialize()
	server.Initialize()
	srv := &http.Server{
		Addr:    ":8080",
		Handler: server.S.Handler(),
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.L.Fatalf("listen: %s\n", err)
		}
	}()
	logger.L.Infof("Server running at port :8080")
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.L.Infof("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.L.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		logger.L.Infof("timeout of 5 seconds.")
	}
	logger.L.Infof("Server exiting")

}
