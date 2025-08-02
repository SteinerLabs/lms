#!/bin/bash

# Exit on error
set -e

# Change to the directory containing the proto file
cd "$(dirname "$0")/.."

# Create the proto/gen directory if it doesn't exist
mkdir -p proto/gen

# Generate Go code from the proto file
protoc --go_out=proto/gen --go_opt=paths=source_relative \
       --go-grpc_out=proto/gen --go-grpc_opt=paths=source_relative \
       proto/auth.proto

echo "Proto code generation completed successfully!"