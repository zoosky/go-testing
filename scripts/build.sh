#!/bin/bash

# Create bin directory if it doesn't exist
mkdir -p bin

# Build the server
echo "Building server..."
go build -o bin/server ./cmd/server

# Check if the build was successful
if [ $? -eq 0 ]; then
    echo "Build successful. Binary is located at bin/server"
else
    echo "Build failed."
    exit 1
fi

# Generate Swagger documentation
echo "Generating Swagger documentation..."
go install github.com/swaggo/swag/cmd/swag@latest

# Try using swag from path, then fall back to GOPATH if needed
SWAG_CMD="swag"
if ! command -v $SWAG_CMD > /dev/null; then
    SWAG_CMD="$(go env GOPATH)/bin/swag"
fi
$SWAG_CMD init -g internal/api/server.go -o docs
if [ $? -ne 0 ]; then
    echo "Error generating Swagger documentation"
    exit 1
fi

echo "Build completed successfully."