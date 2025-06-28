package service

import (
	"context"

	"github.com/mrityunjoydey/go-grpc/pkg/logger"
	pb "github.com/mrityunjoydey/go-grpc/rpc"
	"go.uber.org/zap"
)

// GreeterService implements the GreeterServer interface.
type GreeterService struct {
	pb.UnimplementedGreeterServer
	logger logger.Logger
}

// NewGreeterService creates a new GreeterService.
func NewGreeterService(logger logger.Logger) *GreeterService {
	return &GreeterService{logger: logger}
}

// SayHello implements the SayHello RPC method.
func (s *GreeterService) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	s.logger.Info("SayHello request received", zap.String("name", req.GetName()))
	return &pb.HelloReply{Message: "Hello, " + req.GetName()}, nil
}
