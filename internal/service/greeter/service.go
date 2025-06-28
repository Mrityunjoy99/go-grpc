package greeter

import (
	"context"
	"fmt"
	"io"

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

// StreamGreetings implements the StreamGreetings RPC method for server-side streaming.
func (s *GreeterService) StreamGreetings(req *pb.HelloRequest, stream pb.Greeter_StreamGreetingsServer) error {
	s.logger.Info("StreamGreetings request received", zap.String("name", req.GetName()))
	for i := 0; i < 5; i++ {
		response := &pb.HelloReply{
			Message: fmt.Sprintf("Hello, %s! (Greeting #%d)", req.GetName(), i+1),
		}
		if err := stream.Send(response); err != nil {
			// s.logger.Error("Failed to send greeting", zap.Error(err))
			return err
		}
	}
	return nil
}

// GreetManyTimes implements the GreetManyTimes RPC method for client-side streaming.
func (s *GreeterService) GreetManyTimes(stream pb.Greeter_GreetManyTimesServer) error {
	s.logger.Info("GreetManyTimes request received")
	var names []string
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			message := fmt.Sprintf("Hello, %v!", names)
			return stream.SendAndClose(&pb.HelloReply{Message: message})
		}
		if err != nil {
			s.logger.Error("Failed to receive request", zap.Error(err))
			return err
		}
		s.logger.Info("Received name", zap.String("name", req.GetName()))
		names = append(names, req.GetName())
	}
}

// Chat implements the Chat RPC method for bi-directional streaming.
func (s *GreeterService) Chat(stream pb.Greeter_ChatServer) error {
	s.logger.Info("Chat session started")
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			s.logger.Error("Failed to receive message", zap.Error(err))
			return err
		}

		s.logger.Info("Received message", zap.String("name", req.GetName()))
		response := &pb.HelloReply{
			Message: "Hello, " + req.GetName(),
		}

		if err := stream.Send(response); err != nil {
			s.logger.Error("Failed to send message", zap.Error(err))
			return err
		}
	}
}
