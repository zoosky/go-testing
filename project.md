# Go Project Testing Best Practices

Here are the proven best practices for setting up testing in Go projects:

## Typical Project structure

    myproject/
    ├── cmd/          # Command applications
    ├── internal/     # Private code
    ├── pkg/          # Public library code
    ├── api/          # API definitions
    ├── configs/      # Configuration files
    ├── test/         # Additional test data/scripts
    └── scripts/      # Build/automation scripts

## Unit Testing

1. **Follow the Standard Library Pattern**:

- Create `*_test.go` files next to the code they test
- Use the standard `testing` package
- Name test functions as `TestXxx` where `Xxx` describes what's being tested

5. **Table-Driven Tests**:

- Use slices of test cases for comprehensive coverage
- Each test case includes inputs and expected outputs
- Allows for easy addition of new test cases

9. **Test Helpers**:

- Create helper functions for common setup/teardown
- Use `t.Helper()` to improve error reporting
- Consider subtests with `t.Run()` for better organization

13. **Testable Code Structure**:

- Use dependency injection for easier mocking
- Keep functions small and focused
- Avoid global state that makes testing difficult

17. **Mock External Dependencies**:

- Use interfaces for external services
- Create test doubles (mocks, stubs, fakes) that implement these interfaces
- Consider using `github.com/golang/mock/gomock` or `github.com/stretchr/testify/mock`

## Integration Testing

1. **Test Directory Structure**:

- Keep integration tests separate from unit tests
- Use `*_integration_test.go` naming convention
- Use build tags like `// +build integration` to selectively run tests

5. **Test Real Dependencies**:

- Use Docker containers for databases, message queues, etc.
- Consider `testcontainers-go` for managing test infrastructure
- Use environment variables to configure test connections

9. **Test Setup/Teardown**:

- Create robust setup/teardown functions using `TestMain(m *testing.M)`
- Ensure tests clean up after themselves
- Consider parallel test execution with `t.Parallel()`

13. **API Testing**:

- Use `net/http/httptest` for testing HTTP handlers
- Test against your actual router config
- Consider using `github.com/gavv/httpexpect` for fluent API testing

## Benchmark Testing

1. **Naming and Organization**:

- Name benchmark functions as `BenchmarkXxx`
- Create separate benchmark files when appropriate
- Group related benchmarks together

5. **Benchmark Structure**:

- Reset the timer if performing setup: `b.ResetTimer()`
- Run operations in a loop using `b.N`: `for i := 0; i < b.N; i++ {}`
- Avoid allocation in the timing loop by pre-allocating

9. **Memory Benchmarks**:

- Use `b.ReportAllocs()` to track memory allocations
- Consider `testing.AllocsPerRun()` for detailed allocation tracking
- Use `go test -benchmem` to include memory stats

13. **Comparison Benchmarks**:

- Use tools like `benchstat` to compare benchmark results
- Maintain a set of baseline benchmarks for regression detection
- Run benchmarks on CI for consistent environments

## Tools and Utilities

1. **Coverage Analysis**:

- Use `go test -cover` for basic coverage reporting
- Consider `go test -coverprofile=coverage.out` and `go tool cover -html=coverage.out`
- Set coverage thresholds in CI

5. **Go Test Flags**:

- `-short` flag for skipping long-running tests
- `-race` for detecting race conditions
- `-timeout` for setting test timeouts

9. **Test Frameworks and Assertion Libraries**:

- `github.com/stretchr/testify` for assertions and mocking
- `github.com/onsi/ginkgo` and `github.com/onsi/gomega` for BDD-style tests
- Keep in mind Go's philosophy of simplicity before adding dependencies

# Go Testing Best Practices

Go's standard library provides robust testing capabilities through the `testing` package. Here are the proven best practices for setting up different types of tests in a Go project:

## Project Structure

Go projects typically follow a standard layout:
 ``myproject/
 ├── cmd/          # Command applications
 ├── internal/     # Private code
 ├── pkg/          # Public library code
 ├── api/          # API definitions
 ├── configs/      # Configuration files
 ├── test/         # Additional test data/scripts
 └── scripts/      # Build/automation scripts

## Unit Testing

1. **File Naming**: Place tests in the same package as the code, with `_test.go` suffix
 - Example: `user.go` → `user_test.go`
3. **Function Naming**: Prefix test functions with `Test` followed by the function name
 - Example: `func TestCreateUser(t *testing.T) {...}`
5. **Table-Driven Tests**: Use slices of test cases for comprehensive coverage
  go
  ``func TestAdd(t *testing.T) {
      tests := []struct {
          name     string
          a, b     int
          expected int
      }{
          {"positive", 2, 3, 5},
          {"negative", -1, -2, -3},
          {"mixed", -1, 5, 4},
      }
      
      for _, tc := range tests {
          t.Run(tc.name, func(t *testing.T) {
              got := Add(tc.a, tc.b)
              if got != tc.expected {
                  t.Errorf("Add(%d, %d) = %d; want %d", tc.a, tc.b, got, tc.expected)
              }
          })
      }
  }
8. **Test Helpers**: Create helper functions for common setup/teardown
  go
  ``func setupTest(t *testing.T) (*Database, func()) {
      db := NewTestDatabase()
      return db, func() {
          db.Close()
      }
  }
11. **Subtests**: Use `t.Run()` to group related tests and enable selective test runs
  go
  ``func TestUser(t *testing.T) {
      t.Run("Create", func(t*testing.T) {...})
      t.Run("Update", func(t *testing.T) {...})
  }
14. **Testify**: Consider using `github.com/stretchr/testify` for more expressive assertions
  go
  ``import "github.com/stretchr/testify/assert"

  func TestSomething(t *testing.T) {
      result := Calculate()
      assert.Equal(t, 42, result)
  }

## Integration Testing

1. **Separate Package**: Use package `package_test` to test from an external perspective
  go
  ``// user_test.go
  package user_test // Not package user

  import (
      "testing"
      "myproject/user"
  )
4. **Test Flags**: Add build tags to separate slow integration tests
  go
  ``// +build integration

  package mypackage_test
7. **Docker Compose**: Use containers for external dependencies
  go
  ``func TestDatabaseIntegration(t *testing.T) {
      if testing.Short() {
          t.Skip("skipping integration test")
      }
      // Test with real database
  }
10. **Test Main**: Use `TestMain` for setup/teardown of test environment
  go
  ``func TestMain(m *testing.M) {
      // Setup test environment
      setup()
      code := m.Run()
      // Teardown
      teardown()
      os.Exit(code)
  }
13. **Interfaces**: Design with interfaces to make mocking easier for unit tests while using real implementations for integration tests

## Benchmark Testing

1. **Function Naming**: Prefix with `Benchmark` followed by what's being measured
  go
  ``func BenchmarkCalculate(b *testing.B) {
      for i := 0; i < b.N; i++ {
          Calculate()
      }
  }
4. **Reset Timer**: Reset the timer if setup work is needed
  go
  ``func BenchmarkComplexOperation(b *testing.B) {
      data := prepareData()
      b.ResetTimer() // Reset timer after setup
      for i := 0; i < b.N; i++ {
          ProcessData(data)
      }
  }
7. **Sub-benchmarks**: Use `b.Run()` for comparative benchmarking
  go
  ``func BenchmarkAlgorithms(b *testing.B) {
      b.Run("algorithm1", func(b*testing.B) {
          for i := 0; i < b.N; i++ {
              Algorithm1()
          }
      })
      b.Run("algorithm2", func(b *testing.B) {
          for i := 0; i < b.N; i++ {
              Algorithm2()
          }
      })
  }
10. **Memory Allocation**: Track memory allocations
  go
  ``func BenchmarkMemoryUsage(b *testing.B) {
      b.ReportAllocs() // Report memory allocations
      for i := 0; i < b.N; i++ {
          GenerateData()
      }
  }

## Running Tests

- Run all tests: `go test ./...`
- Run specific tests: `go test -run=TestUserCreate`
- Run integration tests: `go test -tags=integration ./...`
- Run benchmarks: `go test -bench=. -benchmem ./...`
- Generate coverage report: `go test -cover ./...` or `go test -coverprofile=coverage.out ./...`
