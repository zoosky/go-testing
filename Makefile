.PHONY: build run test test-integration benchmark clean swagger help

# Default target
all: build

# Build the application
build:
	@echo "Building application..."
	@mkdir -p bin
	@go build -o bin/server ./cmd/server
	@echo "Build complete. Binary is located at bin/server"

# Run the application
run: build
	@echo "Running server..."
	@./bin/server

# Run the application in development mode
run-dev:
	@echo "Running server in development mode..."
	@go run cmd/server/main.go

# Run unit tests
test:
	@echo "Running unit tests..."
	@go test ./...

# Run unit tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -cover ./...
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run integration tests
test-integration:
	@echo "Running integration tests..."
	@go test -tags=integration ./test/integration

# Run all tests (unit and integration)
test-all: test test-integration

# Run benchmarks
benchmark:
	@echo "Running benchmarks..."
	@echo "Benchmarking calculator package..."
	@go test -bench=. -benchmem ./pkg/calculator
	@echo "\nBenchmarking database package..."
	@go test -bench=. -benchmem ./internal/database
	@echo "\nBenchmarking API package..."
	@go test -bench=. -benchmem ./internal/api
	@echo "\nBenchmark complete."

# Generate Swagger documentation
swagger:
	@echo "Generating Swagger documentation..."
	@# Check for swag command and install if not found
	@if ! command -v swag > /dev/null; then \
		echo "Installing swag..."; \
		go install github.com/swaggo/swag/cmd/swag@latest; \
	fi
	@# Try using swag from path, then fall back to GOPATH if needed
	@SWAG_CMD="swag"; \
	if ! command -v $$SWAG_CMD > /dev/null; then \
		SWAG_CMD="$(shell go env GOPATH)/bin/swag"; \
	fi; \
	$$SWAG_CMD init -g internal/api/server.go -o docs; \
	if [ $$? -ne 0 ]; then \
		echo "Error generating Swagger documentation"; \
		exit 1; \
	fi
	@echo "Swagger documentation generated in docs/ directory"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@echo "Clean complete"

# Show help
help:
	@echo "Available targets:"
	@echo "  build          - Build the application"
	@echo "  run            - Build and run the application"
	@echo "  run-dev        - Run the application in development mode (without build)"
	@echo "  test           - Run unit tests"
	@echo "  test-coverage  - Run unit tests with coverage"
	@echo "  test-integration - Run integration tests"
	@echo "  test-all       - Run all tests (unit and integration)"
	@echo "  benchmark      - Run performance benchmarks"
	@echo "  swagger        - Generate Swagger documentation"
	@echo "  clean          - Clean build artifacts"
	@echo "  help           - Show this help"