package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/mrityunjoydey/go-grpc/internal/common/config"
	"github.com/mrityunjoydey/go-grpc/internal/server"
	config_pkg "github.com/mrityunjoydey/go-grpc/pkg/config"
	"github.com/mrityunjoydey/go-grpc/pkg/logger"
)

func main() {
	cfg := &config.Config{}
	// Load configuration
	cfg, err := config_pkg.LoadConfig(cfg)
	if err != nil {
		panic("failed to load config: " + err.Error())
	}

	// Initialize logger
	log, err := logger.NewZapLogger("grpc-server", cfg.App.LogToFile)
	if err != nil {
		panic("failed to create logger: " + err.Error())
	}
	defer func() {
		if err := log.Flush(); err != nil {
			log.Error("failed to flush logs", zap.Error(err))
		}
	}()

	// Create and start server
	srv := server.New(cfg.Server.Port, log)

	// Graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := srv.Start(); err != nil {
			log.Fatal("gRPC server failed to start", zap.Error(err))
		}
	}()

	log.Info("gRPC server started")

	<-ctx.Done()

	log.Info("Shutting down gRPC server")
	srv.Stop()
}
