package middleware

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/mrityunjoydey/go-grpc/internal/common/constant"
)

// UnaryRequestIDInterceptor returns a new unary server interceptor that adds a request ID to the context.
func UnaryRequestIDInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		var requestID string
		if ok {
			header := md.Get(string(constant.RequestIDHeader))
			if len(header) > 0 {
				requestID = header[0]
			}
		}

		if requestID == "" {
			requestID = uuid.New().String()
		}

		ctx = context.WithValue(ctx, constant.RequestIDKey, requestID)

		return handler(ctx, req)
	}
}

// StreamRequestIDInterceptor returns a new stream server interceptor that adds a request ID to the context.
func StreamRequestIDInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := ss.Context()
		md, ok := metadata.FromIncomingContext(ctx)
		var requestID string
		if ok {
			header := md.Get(string(constant.RequestIDHeader))
			if len(header) > 0 {
				requestID = header[0]
			}
		}

		if requestID == "" {
			requestID = uuid.New().String()
		}

		newCtx := context.WithValue(ctx, constant.RequestIDKey, requestID)

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
