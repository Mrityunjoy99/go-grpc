PROTO_DIR=proto
MODULE_NAME=github.com/mrityunjoydey/go-grpc

echo "Generating gRPC code..."
protoc --proto_path="$PROTO_DIR" --go_out=. --go_opt=module="$MODULE_NAME" --go-grpc_out=. --go-grpc_opt=module="$MODULE_NAME" "$PROTO_DIR"/*.proto