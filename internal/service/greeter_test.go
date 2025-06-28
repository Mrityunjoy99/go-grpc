package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mrityunjoydey/go-grpc/pkg/logger"
	pb "github.com/mrityunjoydey/go-grpc/rpc"
)

func TestGreeterService_SayHello(t *testing.T) {
	logger, _ := logger.NewZapLogger("test")
	s := NewGreeterService(logger)

	req := &pb.HelloRequest{Name: "World"}
	res, err := s.SayHello(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "Hello, World", res.Message)
}
