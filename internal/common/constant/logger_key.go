// Package constant provides constants for the application.
package constant

// LoggerKey is the key for the request ID in the logger context and logs.
type LoggerKey string

const (
	// RequestIDKey is the key for the request ID in the logger context and logs.
	RequestIDKey LoggerKey = "request_id"
)
