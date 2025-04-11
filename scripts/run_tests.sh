#!/bin/bash

# Run unit tests
echo "Running unit tests..."
go test ./...

# Run tests with coverage
echo -e "\nRunning tests with coverage..."
go test -cover ./...

# Generate detailed coverage report
echo -e "\nGenerating coverage report..."
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
echo "Coverage report generated: coverage.html"

# Run benchmarks
echo -e "\nRunning benchmarks..."
go test -bench=. -benchmem ./pkg/calculator

# Run integration tests if -i flag is provided
if [[ "$1" == "-i" ]]; then
  echo -e "\nRunning integration tests..."
  go test -tags=integration ./test/integration
fi

echo -e "\nAll tests completed!"