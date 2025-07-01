// Package main is the entry point of the application.
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/mrityunjoydey/go-grpc/internal/common/config"
	"github.com/mrityunjoydey/go-grpc/internal/common/constant"
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

	// Create a context with a request ID for lifecycle logs
	reqID := uuid.New().String()
	ctx := context.WithValue(context.Background(), constant.RequestIDKey, reqID)
	lifecycleLogger := log.WithContext(ctx)

	defer func() {
		if err := lifecycleLogger.Flush(); err != nil {
			lifecycleLogger.Error("failed to flush logs", zap.Error(err))
		}
	}()

	// Create and start server
	srv := server.New(cfg.Server.Port, log)

	// Graceful shutdown
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := srv.Start(); err != nil {
			lifecycleLogger.Fatal("gRPC server failed to start", zap.Error(err))
		}
	}()

	lifecycleLogger.Info("gRPC server started")

	<-ctx.Done()

	lifecycleLogger.Info("Shutting down gRPC server")
	srv.Stop()
}
