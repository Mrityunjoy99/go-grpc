// Package middleware provides middleware functions for gRPC servers.
package middleware

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"

	"github.com/mrityunjoydey/go-grpc/internal/common/constant"
)

// UnaryRequestIDInterceptor returns a new unary server interceptor that adds a request ID to the context.
// It extracts the request ID from the 'x-request-id' metadata header if present, otherwise generates a new one.
// The request ID is then set in the context for logging and tracing purposes.
func UnaryRequestIDInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		_ *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)

		var requestID string

		// Try to get request ID from metadata
		if ok {
			// gRPC metadata keys are automatically lowercased
			values := md.Get(string(constant.RequestIDHeader))
			if len(values) > 0 && values[0] != "" {
				requestID = values[0]
			}
		}

		// Generate new request ID if not found in metadata
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Add request ID to outgoing metadata
		mdOut := metadata.Pairs(string(constant.RequestIDHeader), requestID)
		if err := grpc.SetHeader(ctx, mdOut); err != nil {
			// Log error but continue with the request
			grpclog.Errorf("Failed to set request ID in metadata: %v", err)
		}

		// Add request ID to context
		ctx = context.WithValue(ctx, constant.RequestIDKey, requestID)

		return handler(ctx, req)
	}
}

// StreamRequestIDInterceptor returns a new stream server interceptor that adds a request ID to the context.
// It extracts the request ID from the 'x-request-id' metadata header if present, otherwise generates a new one.
// The request ID is then set in the context for logging and tracing purposes.
func StreamRequestIDInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, _ *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := ss.Context()
		md, ok := metadata.FromIncomingContext(ctx)

		var requestID string

		// Try to get request ID from metadata
		if ok {
			// gRPC metadata keys are automatically lowercased
			values := md.Get(string(constant.RequestIDHeader))
			if len(values) > 0 && values[0] != "" {
				requestID = values[0]
			}
		}

		// Generate new request ID if not found in metadata
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Add request ID to outgoing metadata
		mdOut := metadata.Pairs(string(constant.RequestIDHeader), requestID)
		if err := ss.SetHeader(mdOut); err != nil {
			// Log error but continue with the request
			grpclog.Errorf("Failed to set request ID in stream metadata: %v", err)
		}

		// Create new context with request ID
		newCtx := context.WithValue(ctx, constant.RequestIDKey, requestID)

		// Wrap the server stream with our context
		wrapped := &wrappedStream{ServerStream: ss, newCtx: newCtx}

		return handler(srv, wrapped)
	}
}

// wrappedStream wraps a grpc.ServerStream and overrides its context.
type wrappedStream struct {
	grpc.ServerStream
	newCtx context.Context
}

func (w *wrappedStream) Context() context.Context {
	return w.newCtx
}
