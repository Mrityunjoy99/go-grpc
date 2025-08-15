package constant

// RequestHeader is the header key in gRPC metadata.
type RequestHeader string

const (
	// RequestIDHeader is the header key for the request ID in gRPC metadata.
	RequestIDHeader RequestHeader = "X-Request-ID"
)
