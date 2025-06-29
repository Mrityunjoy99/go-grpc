# Makefile for the Go gRPC Boilerplate project

.PHONY: all proto build run test lint clean help

# Variables
APP_NAME=grpc-server
CMD_PATH=cmd/server/main.go
BUILD_PATH=build
PROTO_DIR=proto
MODULE_NAME=github.com/mrityunjoydey/go-grpc

all: build

# Generate gRPC code from proto files
proto:
	echo "Generating gRPC code..."
	sh scripts/generate_proto.sh

# Build the application binary
build: proto
	echo "Building the application..."
	go build -o $(BUILD_PATH)/$(APP_NAME) $(CMD_PATH)

run-build:
	echo "Running the application..."
	$(BUILD_PATH)/$(APP_NAME)

# Run the application
run: 
	echo "Running the application..."
	go run $(CMD_PATH)

# Run the tests
test:
	echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	echo "To view the coverage report, run: go tool cover -html=coverage.out"

# Run the linter
lint:
	echo "Running linter..."
	golangci-lint run

# Clean the build artifacts
clean:
	echo "Cleaning up..."
	rm -f $(BUILD_PATH)/$(APP_NAME)

# Display help
help:
	echo "Available commands:"
	echo "  proto   - Generate gRPC code"
	echo "  build   - Build the application"
	echo "  run     - Run the application"
	echo "  test    - Run tests"
	echo "  test-coverage - Run tests and generate coverage report"
	echo "  lint    - Run linter"
	echo "  clean   - Clean build artifacts"
	echo "  help    - Display thitest-coverages help message"
