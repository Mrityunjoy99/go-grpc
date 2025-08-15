package greeter_test

import (
	"context"
	"fmt"
	"io"
	"testing"

	"github.com/mrityunjoydey/go-grpc/src/service/greeter"
	"github.com/mrityunjoydey/go-grpc/pkg/logger"
	pb "github.com/mrityunjoydey/go-grpc/rpc"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

// mockGreeterServerStream is a mock implementation of the stream interfaces.
type mockGreeterServerStream struct {
	ctx         context.Context
	requests    []*pb.HelloRequest
	sentReplies []*pb.HelloReply
	recvIndex   int
	finalReply  *pb.HelloReply
}

// Implement grpc.ServerStream
func (m *mockGreeterServerStream) SetHeader(md metadata.MD) error  { return nil }
func (m *mockGreeterServerStream) SendHeader(md metadata.MD) error { return nil }
func (m *mockGreeterServerStream) SetTrailer(md metadata.MD)       {}
func (m *mockGreeterServerStream) Context() context.Context        { return m.ctx }
func (m *mockGreeterServerStream) SendMsg(msg interface{}) error   { return nil }
func (m *mockGreeterServerStream) RecvMsg(msg interface{}) error   { return nil }

// Implement methods for specific stream types
func (m *mockGreeterServerStream) Send(res *pb.HelloReply) error {
	m.sentReplies = append(m.sentReplies, res)
	return nil
}

func (m *mockGreeterServerStream) Recv() (*pb.HelloRequest, error) {
	if m.recvIndex >= len(m.requests) {
		return nil, io.EOF
	}

	req := m.requests[m.recvIndex]
	m.recvIndex++

	return req, nil
}

func (m *mockGreeterServerStream) SendAndClose(res *pb.HelloReply) error {
	m.finalReply = res
	return nil
}

func TestGreeterService_SayHello(t *testing.T) {
	logger, _ := logger.NewZapLogger("test", false)
	s := greeter.NewService(logger)

	req := &pb.HelloRequest{Name: "World"}
	res, err := s.SayHello(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "Hello, World", res.Message)
}

func TestGreeterService_StreamGreetings(t *testing.T) {
	logger, _ := logger.NewZapLogger("test", false)
	s := greeter.NewService(logger)

	req := &pb.HelloRequest{Name: "Streamer"}
	stream := &mockGreeterServerStream{ctx: context.Background()}

	err := s.StreamGreetings(req, stream)

	assert.NoError(t, err)
	assert.Len(t, stream.sentReplies, 5)

	for i, msg := range stream.sentReplies {
		expected := fmt.Sprintf("Hello, Streamer! (Greeting #%d)", i+1)
		assert.Equal(t, expected, msg.Message)
	}
}

func TestGreeterService_GreetManyTimes(t *testing.T) {
	logger, _ := logger.NewZapLogger("test", false)
	s := greeter.NewService(logger)

	requests := []*pb.HelloRequest{
		{Name: "Alice"},
		{Name: "Bob"},
		{Name: "Charlie"},
	}
	stream := &mockGreeterServerStream{ctx: context.Background(), requests: requests}

	err := s.GreetManyTimes(stream)

	assert.NoError(t, err)
	assert.NotNil(t, stream.finalReply)
	assert.Equal(t, "Hello, [Alice Bob Charlie]!", stream.finalReply.Message)
}

func TestGreeterService_Chat(t *testing.T) {
	logger, _ := logger.NewZapLogger("test", false)
	s := greeter.NewService(logger)

	requests := []*pb.HelloRequest{
		{Name: "Dave"},
		{Name: "Eve"},
		{Name: "Frank"},
	}
	stream := &mockGreeterServerStream{ctx: context.Background(), requests: requests}

	err := s.Chat(stream)

	assert.NoError(t, err)
	assert.Len(t, stream.sentReplies, len(requests))

	for i, reply := range stream.sentReplies {
		expected := "Hello, " + requests[i].GetName()
		assert.Equal(t, expected, reply.Message)
	}
}
