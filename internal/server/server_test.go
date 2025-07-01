package server

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/mrityunjoydey/go-grpc/pkg/logger"
	pb "github.com/mrityunjoydey/go-grpc/rpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

func newBufconnListener() *bufconn.Listener {
	return bufconn.Listen(bufSize)
}

func TestNew(t *testing.T) {
	logger, err := logger.NewZapLogger("test", false)
	require.NoError(t, err)

	srv := New("8080", logger)
	assert.NotNil(t, srv)
}

func TestServer_StartStop(t *testing.T) {
	logger, err := logger.NewZapLogger("test", false)
	require.NoError(t, err)

	// Use a bufconn listener to avoid using a real port
	bufListener := newBufconnListener()

	srv := New("", logger) // Port is not used with bufconn

	// Start the server in a separate goroutine
	go func() {
		if err := srv.serve(bufListener); err != nil {
			// This error is expected after the listener is closed.
			if err != grpc.ErrServerStopped {
				t.Logf("server error: %v", err)
			}
		}
	}()

	// Create a client connection to the server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Use a custom resolver scheme for bufconn
	resolver.SetDefaultScheme("passthrough")

	conn, err := grpc.NewClient("passthrough:///bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return bufListener.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)

	defer func() {
		if err := conn.Close(); err != nil {
			t.Logf("failed to close connection: %v", err)
		}
	}()

	// Check the health of the Greeter service
	healthClient := grpc_health_v1.NewHealthClient(conn)
	resp, err := healthClient.Check(ctx, &grpc_health_v1.HealthCheckRequest{
		Service: pb.Greeter_ServiceDesc.ServiceName,
	})

	require.NoError(t, err)
	assert.Equal(t, grpc_health_v1.HealthCheckResponse_SERVING, resp.Status)

	// Stop the server
	srv.Stop()

	// After stopping the server, the health check should fail with an "Unavailable" error.
	_, err = healthClient.Check(ctx, &grpc_health_v1.HealthCheckRequest{
		Service: pb.Greeter_ServiceDesc.ServiceName,
	})

	require.Error(t, err)
	st, ok := status.FromError(err)
	require.True(t, ok, "error should be a gRPC status error")
	assert.Equal(t, codes.Unavailable, st.Code(), "expected status code to be Unavailable")
}
