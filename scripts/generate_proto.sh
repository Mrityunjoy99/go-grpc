PROTO_DIR=proto
MODULE_NAME=github.com/mrityunjoydey/go-grpc

echo "Generating gRPC code..."
protoc --experimental_allow_proto3_optional --go_out=. --go_opt=module="$MODULE_NAME" --go-grpc_out=. --go-grpc_opt=module="$MODULE_NAME" "$PROTO_DIR"/**/*.proto
# --experimental_allow_proto3_optional is required for proto3 optional fields