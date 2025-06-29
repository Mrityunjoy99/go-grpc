# Stage 1: Build the application
FROM golang:1.24-alpine AS builder

# Set the working directory
WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /grpc-server ./cmd/server

# Stage 2: Create the final image
FROM alpine:latest

# Install CA certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Set the working directory
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /grpc-server .

# Make the binary executable
RUN chmod +x /app/grpc-server

# Expose the gRPC port (documentation only, actual port mapping is done at runtime)
EXPOSE 50051

# Run the application
CMD ["/app/grpc-server"]
