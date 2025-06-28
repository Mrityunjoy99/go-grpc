package server

import (
	"context"
	"fmt"
	"net"
	"runtime/debug"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"

	"github.com/mrityunjoydey/go-grpc/internal/service"
	"github.com/mrityunjoydey/go-grpc/pkg/logger"
	pb "github.com/mrityunjoydey/go-grpc/rpc"
)

// Server is the gRPC server.
type Server struct {
	logger     logger.Logger
	grpcServer *grpc.Server
	port       string
	healthSrv  *health.Server
}

// interceptorLogger adapts zap logger to the interceptor's logger interface.
func interceptorLogger(l logger.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		f := make([]zap.Field, 0, len(fields)/2)
		for i := 0; i < len(fields); i += 2 {
			key := fields[i]
			value := fields[i+1]

			switch v := value.(type) {
			case string:
				f = append(f, zap.String(key.(string), v))
			case int:
				f = append(f, zap.Int(key.(string), v))
			case bool:
				f = append(f, zap.Bool(key.(string), v))
			default:
				f = append(f, zap.Any(key.(string), v))
			}
		}

		logger := l.With(f...)

		switch lvl {
		case logging.LevelDebug:
			logger.Debug(msg)
		case logging.LevelInfo:
			logger.Info(msg)
		case logging.LevelWarn:
			logger.Warn(msg)
		case logging.LevelError:
			logger.Error(msg)
		default:
			panic(fmt.Sprintf("unknown level %v", lvl))
		}
	})
}

// New creates a new gRPC server.
func New(port string, logger logger.Logger) *Server {
	// Setup panic recovery handler
	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p any) (err error) {
			logger.Error("recovered from panic",
				zap.Any("panic", p),
				zap.String("stack", string(debug.Stack())),
			)
			return status.Errorf(codes.Internal, "internal server error")
		}),
	}

	// Create a new gRPC server with unary interceptors
	gs := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(interceptorLogger(logger)),
			recovery.UnaryServerInterceptor(recoveryOpts...),
		),
		grpc.ChainStreamInterceptor(
			logging.StreamServerInterceptor(interceptorLogger(logger)),
			recovery.StreamServerInterceptor(recoveryOpts...),
		),
	)

	// Register Greeter service
	greeterService := service.NewGreeterService(logger)
	pb.RegisterGreeterServer(gs, greeterService)

	// Register health check service
	healthSrv := health.NewServer()
	grpc_health_v1.RegisterHealthServer(gs, healthSrv)
	healthSrv.SetServingStatus(pb.Greeter_ServiceDesc.ServiceName, grpc_health_v1.HealthCheckResponse_SERVING)

	return &Server{
		logger:     logger,
		grpcServer: gs,
		port:       port,
		healthSrv:  healthSrv,
	}
}

// Start starts the gRPC server.
func (s *Server) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", s.port))
	if err != nil {
		s.logger.Fatal("Failed to listen", zap.Error(err))
		return err
	}

	s.logger.Info(fmt.Sprintf("gRPC server listening on port %s", s.port))
	if err := s.grpcServer.Serve(lis); err != nil {
		s.logger.Fatal("Failed to serve gRPC server", zap.Error(err))
		return err
	}

	return nil
}

// Stop gracefully stops the gRPC server.
func (s *Server) Stop() {
	s.logger.Info("Stopping gRPC server")
	// Set the health status to NOT_SERVING
	s.healthSrv.SetServingStatus(pb.Greeter_ServiceDesc.ServiceName, grpc_health_v1.HealthCheckResponse_NOT_SERVING)
	s.grpcServer.GracefulStop()
}
