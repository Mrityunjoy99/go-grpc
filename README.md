# Go gRPC Boilerplate Microservice

This is a production-ready gRPC microservice boilerplate built with Go, featuring Zap for logging, Viper for configuration, and a complete suite of development and CI tools.

## Features

- **Language:** Go (v1.24)
- **Framework:** gRPC
- **Logging:** Uber's Zap logger (production config)
- **Configuration:** Viper with YAML and environment variable support
- **Linting:** `golangci-lint` with a curated ruleset
- **Testing:** `testify` for assertions and mocking
- **Containerization:** Multi-stage Dockerfile and Docker Compose
- **CI/CD:** GitHub Actions for automated linting and testing
- **Tooling:** Makefile for common development tasks

## Implemented RPCs

The `Greeter` service now includes the following methods:

-   **Unary RPC**: `SayHello(HelloRequest) returns (HelloReply)`
    -   A simple request-response method.
-   **Server Streaming RPC**: `StreamGreetings(HelloRequest) returns (stream HelloReply)`
    -   The client sends a single request, and the server returns a stream of responses.
-   **Client Streaming RPC**: `GreetManyTimes(stream HelloRequest) returns (HelloReply)`
    -   The client sends a stream of requests, and the server returns a single response.
-   **Bi-directional Streaming RPC**: `Chat(stream HelloRequest) returns (stream HelloReply)`
    -   Both the client and server can send a stream of messages to each other.

## Configuration

The application uses a flexible configuration system with the following features:

- **Type-safe configuration** using Go structs
- **Environment variable support** with automatic binding
- **YAML configuration** file support
- **Default values** using struct tags
- **Validation** using go-playground/validator

### Configuration Structure

The main configuration is defined in `internal/common/config/config.go`:

```go
type Config struct {
    Server ServerConfig
}

type ServerConfig struct {
    Port string
}
```

### Environment Variables

Configuration can be overridden using environment variables. The variable names should be in uppercase with dots (`.`) replaced by underscores (`_`). For example:

- `SERVER_PORT=8080` will override the server port

### Usage in Code

To use the configuration in your code:

```go
import "github.com/mrityunjoydey/go-grpc/internal/common/config"
import config_pkg "github.com/mrityunjoydey/go-grpc/pkg/config"

// Initialize config
cfg := &config.Config{}
cfg, err := config_pkg.LoadConfig(cfg)
if err != nil {
    // handle error
}

// Use the config
port := cfg.Server.Port
```

### Default Values

Default values can be set using struct tags from the `default` package. For example:

```go
type ServerConfig struct {
    Port string `default:"50051"`
}
```

### Validation

The configuration is automatically validated using `go-playground/validator`. Add validation tags to your config structs:

```go
type ServerConfig struct {
    Port string `validate:"required,numeric"`
}
```

## Project Structure

```
.
├── .github/workflows/ci.yml   # GitHub Actions CI pipeline
├── .gitignore                 # Standard Go gitignore
├── .golangci.yml              # Linter configuration
├── Dockerfile                 # Multi-stage Dockerfile
├── Makefile                   # Development commands
├── cmd/server/main.go         # Application entry point
├── configs/config.yaml        # Default configuration
├── docker-compose.yml         # Docker Compose setup
├── go.mod                     # Go module definition
├── go.sum                     # Go module checksums
├── grpc-service.postman_collection.json # Postman collection for testing
├── internal                   # Internal application logic
│   ├── common                 # Common utilities
│   ├── server                 # gRPC server setup
│   └── service                # Service implementations
├── pkg                        # Shared packages
├── proto                      # Protocol Buffers definitions
└── rpc                        # Generated gRPC code
├── proto/greeter.proto        # Protocol Buffers service definition
└── rpc/                       # Generated gRPC code
```

## Setup

### Prerequisites

- Go v1.24 or later
- Docker and Docker Compose
- `protoc` compiler
- `golangci-lint`

### Installation

1.  **Clone the repository:**

    ```sh
    git clone https://github.com/mrityunjoydey/go-grpc.git
    cd go-grpc
    ```

2.  **Install Go dependencies:**

    ```sh
    go mod tidy
    ```

3.  **Generate gRPC code:**

    You can generate the gRPC code from the `.proto` file using the Makefile:

    ```sh
    make proto
    ```

## Usage

### Running Locally

Before running the service for the first time, you need to generate the gRPC code:

```sh
make proto
```

Then you can run the gRPC server:

```sh
make run
```

The server will start on `localhost:50051` by default.

### Running with Docker

To build and run the service using Docker Compose:

```sh
docker-compose up --build
```

## Postman

A Postman collection is included at `grpc-service.postman_collection.json`.

For the best experience, it is recommended to import the `proto/greeter.proto` file directly into Postman. This enables full support for gRPC service discovery and request building.

1.  Open Postman and create a new gRPC request.
2.  Select the `greeter.proto` file as the service definition.
3.  Postman will automatically detect the `Greeter` service and the `SayHello` method.

## Testing

To run the unit tests:

```sh
make test
```

## Linting

To run the linter and check for code quality issues:

```sh
make lint
```

## CI/CD

The project includes a GitHub Actions workflow defined in `.github/workflows/ci.yml`. This pipeline automatically runs the linter and tests on every push and pull request to the `main` branch.
