# Go Testing Example Project

This project demonstrates various testing techniques and best practices in Go, including:

- Unit testing with table-driven tests
- Integration testing
- Benchmarking
- Mocking and dependency injection
- HTTP API testing

## Project Structure

The project follows the standard Go project layout:

```
go-testing/
├── cmd/          # Command applications
│   └── server/   # Main server application
├── internal/     # Private code
│   ├── api/      # API server implementation
│   ├── calculator/ # Internal calculator implementation
│   └── database/ # Database implementation
├── pkg/          # Public library code
│   └── calculator/ # Public calculator package
├── api/          # API definitions
│   └── definitions/ # API type definitions
├── configs/      # Configuration files
├── test/         # Additional test data/scripts
│   └── integration/ # Integration tests
└── scripts/      # Build/automation scripts
```

## Building and Running the Application

### Using Make

The project includes a Makefile with several useful targets:

```bash
# Build the application
make build

# Build and run the application
make run

# Run the application in development mode
make run-dev

# Run unit tests
make test

# Run unit tests with coverage
make test-coverage

# Run integration tests
make test-integration

# Run all tests (unit and integration)
make test-all

# Generate Swagger documentation
make swagger

# Clean build artifacts
make clean

# Show help with all available targets
make help
```

### Building Manually

You can use the provided build script:

```bash
./scripts/build.sh
```

Or build manually:

```bash
mkdir -p bin
go build -o bin/server ./cmd/server
```

### Running Manually

Start the built server:

```bash
./bin/server
```

Or run directly without building:

```bash
go run cmd/server/main.go
```

## Running Tests

You can use the provided test script:

```bash
./scripts/run_tests.sh       # Run unit tests, coverage, and benchmarks
./scripts/run_tests.sh -i    # Also run integration tests
```

### Manual Test Commands

#### Unit Tests

Run all unit tests:

```bash
go test ./...
```

Run tests with coverage:

```bash
go test -cover ./...
```

Generate detailed coverage report:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Integration Tests

Run integration tests (which are tagged with `// +build integration`):

```bash
go test -tags=integration ./test/integration
```

### Benchmarks

Run all benchmarks with the Make target:

```bash
make benchmark
```

Run specific benchmarks:

```bash
# Run calculator benchmarks
go test -bench=. -benchmem ./pkg/calculator

# Run database benchmarks
go test -bench=. -benchmem ./internal/database

# Run API benchmarks
go test -bench=. -benchmem ./internal/api
```

Run benchmarks with additional options:

```bash
# Run specific benchmark functions (regexp matching)
go test -bench=BenchmarkAdd ./pkg/calculator

# Show more detailed memory allocation statistics
go test -bench=. -benchmem ./pkg/calculator

# Run benchmarks with custom count (10 iterations)
go test -bench=. -count=10 ./pkg/calculator

# Compare benchmark results with benchstat (needs benchstat installed)
# First, run benchmarks and save results
go test -bench=. -count=10 ./pkg/calculator > old.txt
# Make changes, then run again
go test -bench=. -count=10 ./pkg/calculator > new.txt
# Compare results
benchstat old.txt new.txt
```

### Other Test Flags

- Use `-short` to skip long-running tests: `go test -short ./...`
- Use `-race` to detect race conditions: `go test -race ./...`
- Use `-v` for verbose output: `go test -v ./...`

## Testing Best Practices Demonstrated

This project demonstrates the following testing best practices:

1. **Project Structure**: Following the standard Go project layout
2. **Table-Driven Tests**: Using slices of test cases for comprehensive coverage
3. **Test Helpers**: With `t.Helper()` for better error reporting
4. **Subtests**: Using `t.Run()` for better organization
5. **Mocks**: Using interfaces and the testify/mock package
6. **HTTP Testing**: Using httptest for API testing
7. **Benchmarking**: Proper benchmark structure with `b.N` and `b.ResetTimer()`
8. **Integration Testing**: Separate integration tests with build tags
9. **Test Main**: Using TestMain for setup/teardown

## API Endpoints

The server provides the following endpoints:

### User Endpoints

- `GET /users`: List all users
- `GET /users/{id}`: Get a user by ID
- `POST /users`: Create a new user
- `PUT /users/{id}`: Update a user
- `DELETE /users/{id}`: Delete a user

### Calculator Endpoints

- `GET /calculator/add?a=5&b=3`: Add two numbers
- `GET /calculator/subtract?a=5&b=3`: Subtract b from a
- `GET /calculator/multiply?a=5&b=3`: Multiply two numbers
- `GET /calculator/divide?a=6&b=3`: Divide a by b

### API Documentation

The API is documented using Swagger/OpenAPI. When the server is running, you can access the interactive API documentation at:

```
http://localhost:8080/swagger/index.html
```

This provides an interactive UI to explore the API endpoints, see required parameters, and even test the API directly from your browser.