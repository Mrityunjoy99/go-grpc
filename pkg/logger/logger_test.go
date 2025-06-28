package logger

import (
	"context"
	"testing"

	"github.com/mrityunjoydey/go-grpc/internal/common/constant"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

// setupTestLogger creates a logger with an observer core to capture logs for testing.
func setupTestLogger() (*zapLogger, *observer.ObservedLogs) {
	core, logs := observer.New(zap.InfoLevel)
	logger := &zapLogger{logger: zap.New(core)}
	return logger, logs
}

func TestLogger_Info(t *testing.T) {
	logger, logs := setupTestLogger()
	logger.Info("test info", zap.String("key", "value"))

	assert.Equal(t, 1, logs.Len())
	entry := logs.All()[0]
	assert.Equal(t, zap.InfoLevel, entry.Level)
	assert.Equal(t, "test info", entry.Message)
	assert.ElementsMatch(t, []zapcore.Field{zap.String("key", "value")}, entry.Context)
}

func TestLogger_Debug(t *testing.T) {
	logger, logs := setupTestLogger()
	logger.Debug("test debug") // This should not be logged as the level is Info
	assert.Equal(t, 0, logs.Len())
}

func TestLogger_Warn(t *testing.T) {
	logger, logs := setupTestLogger()
	logger.Warn("test warn")
	assert.Equal(t, 1, logs.Len())
	assert.Equal(t, "test warn", logs.All()[0].Message)
}

func TestLogger_Error(t *testing.T) {
	logger, logs := setupTestLogger()
	logger.Error("test error")
	assert.Equal(t, 1, logs.Len())
	assert.Equal(t, "test error", logs.All()[0].Message)
}

func TestLogger_With(t *testing.T) {
	logger, logs := setupTestLogger()
	logger.With(zap.String("component", "test")).Info("message with field")

	assert.Equal(t, 1, logs.Len())
	entry := logs.All()[0]
	assert.Equal(t, "message with field", entry.Message)
	assert.ElementsMatch(t, []zapcore.Field{zap.String("component", "test")}, entry.Context)
}

func TestLogger_WithContext(t *testing.T) {
	t.Run("with request_id", func(t *testing.T) {
		logger, logs := setupTestLogger()
		ctx := context.WithValue(context.Background(), constant.RequestIDKey, "12345")
		contextualLogger := logger.WithContext(ctx)
		contextualLogger.Info("context message")

		assert.Equal(t, 1, logs.Len())
		entry := logs.All()[0]
		assert.Equal(t, "context message", entry.Message)
		assert.ElementsMatch(t, []zapcore.Field{zap.String("request_id", "12345")}, entry.Context)
	})

	t.Run("without request_id", func(t *testing.T) {
		logger, logs := setupTestLogger()
		ctx := context.Background()
		contextualLogger := logger.WithContext(ctx)
		contextualLogger.Info("no context message")

		assert.Equal(t, 1, logs.Len())
		entry := logs.All()[0]
		assert.Equal(t, "no context message", entry.Message)
		assert.Empty(t, entry.Context)
	})
}
