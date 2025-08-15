// Package main is the entry point for the gRPC client application.
// It handles the initialization and execution of the gRPC client.
// This is only for testing purpose
package main

import (
	"context"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"github.com/mrityunjoydey/go-grpc/src/common/constant"
	pb "github.com/mrityunjoydey/go-grpc/rpc"
)

//nolint:funlen
func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("did not connect: %v", err)
		return
	}

	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("failed to close connection: %v", err)
		}
	}()

	c := pb.NewGreeterClient(conn)

	// --- Test Case 1: No request ID in header ---
	log.Println("--- Sending request without X-Request-ID header ---")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: "world"})
	if err != nil {
		log.Printf("could not greet: %v", err)
		return
	}

	log.Printf("Greeting: %s", r.GetMessage())

	// --- Test Case 2: With request ID in header ---
	log.Println("\n--- Sending request with X-Request-ID header ---")

	ctx2, cancel2 := context.WithTimeout(context.Background(), time.Second)
	defer cancel2()

	md := metadata.New(map[string]string{string(constant.RequestIDHeader): "client-generated-id-123"})
	ctx2 = metadata.NewOutgoingContext(ctx2, md)

	r2, err := c.SayHello(ctx2, &pb.HelloRequest{Name: "world"})
	if err != nil {
		log.Printf("could not greet: %v", err)
		return
	}

	log.Printf("Greeting: %s", r2.GetMessage())

	// --- Test Case 3: StreamGreetings ---
	log.Println("\n--- Calling StreamGreetings ---")

	ctx3, cancel3 := context.WithTimeout(context.Background(), time.Second)
	defer cancel3()

	stream, err := c.StreamGreetings(ctx3, &pb.HelloRequest{Name: "streaming world"})
	if err != nil {
		log.Printf("could not start stream: %v", err)
		return
	}

	for {
		feature, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Printf("%v.StreamGreetings(_) = _, %v", c, err)
			return
		}

		log.Printf("Streamed Greeting: %s", feature.GetMessage())
	}
}
