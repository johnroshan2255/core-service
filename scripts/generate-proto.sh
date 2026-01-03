#!/bin/bash

# Script to generate Go code from proto files

set -e

echo "Generating proto files..."

# Create output directory if it doesn't exist
mkdir -p proto/notification/v1

# Generate notification proto
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/notification/v1/notification.proto

echo "Proto files generated successfully!"

