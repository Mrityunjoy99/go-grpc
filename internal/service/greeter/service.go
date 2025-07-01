// Package greeter implements the GreeterServer interface.
package greeter

import (
	"context"
	"fmt"
	"io"

	"github.com/mrityunjoydey/go-grpc/pkg/logger"
	pb "github.com/mrityunjoydey/go-grpc/rpc"
	"go.uber.org/zap"
)

// Service implements the GreeterServer interface.
type Service struct {
	pb.UnimplementedGreeterServer
	logger logger.Logger
}

// NewService creates a new Service.
func NewService(logger logger.Logger) *Service {
	return &Service{logger: logger}
}

// SayHello implements the SayHello RPC method.
func (s *Service) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	s.logger.WithContext(ctx).Info("SayHello request received", zap.String("name", req.GetName()))
	return &pb.HelloReply{Message: "Hello, " + req.GetName()}, nil
}

// StreamGreetings implements the StreamGreetings RPC method for server-side streaming.
func (s *Service) StreamGreetings(req *pb.HelloRequest, stream pb.Greeter_StreamGreetingsServer) error {
	s.logger.WithContext(stream.Context()).Info("StreamGreetings request received", zap.String("name", req.GetName()))

	for i := 0; i < 5; i++ {
		response := &pb.HelloReply{
			Message: fmt.Sprintf("Hello, %s! (Greeting #%d)", req.GetName(), i+1),
		}
		if err := stream.Send(response); err != nil {
			s.logger.WithContext(stream.Context()).Error("Failed to send greeting", zap.Error(err))
			return err
		}

		s.logger.WithContext(stream.Context()).Info("Sent greeting", zap.String("message", response.GetMessage()))
	}

	return nil
}

// GreetManyTimes implements the GreetManyTimes RPC method for client-side streaming.
func (s *Service) GreetManyTimes(stream pb.Greeter_GreetManyTimesServer) error {
	s.logger.WithContext(stream.Context()).Info("GreetManyTimes request received")

	var names []string

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			message := fmt.Sprintf("Hello, %v!", names)
			return stream.SendAndClose(&pb.HelloReply{Message: message})
		}

		if err != nil {
			s.logger.WithContext(stream.Context()).Error("Failed to receive request", zap.Error(err))
			return err
		}

		s.logger.WithContext(stream.Context()).Info("Received name", zap.String("name", req.GetName()))
		names = append(names, req.GetName())
	}
}

// Chat implements the Chat RPC method for bi-directional streaming.
func (s *Service) Chat(stream pb.Greeter_ChatServer) error {
	s.logger.WithContext(stream.Context()).Info("Chat session started")

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			s.logger.WithContext(stream.Context()).Error("Failed to receive message", zap.Error(err))
			return err
		}

		s.logger.WithContext(stream.Context()).Info("Received message", zap.String("name", req.GetName()))
		response := &pb.HelloReply{
			Message: "Hello, " + req.GetName(),
		}

		if err := stream.Send(response); err != nil {
			s.logger.WithContext(stream.Context()).Error("Failed to send message", zap.Error(err))
			return err
		}
	}
}
