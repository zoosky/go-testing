package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"strconv"
	"testing"

	"go-testing/internal/calculator"
	"go-testing/internal/database"
)

// setupBenchServer creates a server for benchmarking
func setupBenchServer() *Server {
	repo := database.NewUserRepository()
	calc := calculator.NewCalculator()
	return NewServer(repo, calc)
}

// BenchmarkListUsers benchmarks the list users endpoint
func BenchmarkListUsers(b *testing.B) {
	// Create a repository and add test users
	repo := database.NewUserRepository()
	for i := 0; i < 100; i++ {
		user := &database.User{
			Username: "list" + strconv.Itoa(i),
			Email:    "list" + strconv.Itoa(i) + "@example.com",
		}
		repo.CreateUser(user)
	}
	
	// Create a server with the populated repository
	calc := calculator.NewCalculator()
	server := NewServer(repo, calc)
	handler := server.Router()
	
	// Reset the timer to exclude setup time
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/users", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
	}
}

// BenchmarkGetUser benchmarks the get user endpoint
func BenchmarkGetUser(b *testing.B) {
	// Create a repository and add a test user
	repo := database.NewUserRepository()
	user := &database.User{
		Username: "benchmark",
		Email:    "benchmark@example.com",
	}
	repo.CreateUser(user)
	
	// Create a server with the populated repository
	calc := calculator.NewCalculator()
	server := NewServer(repo, calc)
	handler := server.Router()
	
	// Reset the timer to exclude setup time
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", fmt.Sprintf("/users/%d", user.ID), nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
	}
}

// BenchmarkCreateUser benchmarks the create user endpoint
func BenchmarkCreateUser(b *testing.B) {
	server := setupBenchServer()
	handler := server.Router()
	
	// Reset the timer to exclude setup time
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		// Create a new user for each iteration
		newUser := database.User{
			Username: "create" + strconv.Itoa(i),
			Email:    "create" + strconv.Itoa(i) + "@example.com",
		}
		
		body, _ := json.Marshal(newUser)
		req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
	}
}

// BenchmarkCalculatorAdd benchmarks the add endpoint
func BenchmarkCalculatorAdd(b *testing.B) {
	server := setupBenchServer()
	handler := server.Router()
	
	// Reset the timer to exclude setup time
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/calculator/add?a=5&b=3", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
	}
}

// BenchmarkCalculatorMultipleOps benchmarks multiple calculator operations
func BenchmarkCalculatorMultipleOps(b *testing.B) {
	server := setupBenchServer()
	handler := server.Router()
	
	endpoints := []string{
		"/calculator/add?a=5&b=3",
		"/calculator/subtract?a=5&b=3",
		"/calculator/multiply?a=5&b=3",
		"/calculator/divide?a=6&b=3",
	}
	
	// Reset the timer to exclude setup time
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		// Cycle through endpoints
		endpoint := endpoints[i%len(endpoints)]
		req := httptest.NewRequest("GET", endpoint, nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
	}
}

// BenchmarkJsonSerialization benchmarks JSON serialization
func BenchmarkJsonSerialization(b *testing.B) {
	// Create a user to serialize
	user := &database.User{
		ID:       1,
		Username: "benchmark",
		Email:    "benchmark@example.com",
	}
	
	// Reset the timer to exclude setup time
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		_, _ = json.Marshal(user)
	}
}

// BenchmarkJsonDeserialization benchmarks JSON deserialization
func BenchmarkJsonDeserialization(b *testing.B) {
	// Create a user JSON to deserialize
	user := &database.User{
		ID:       1,
		Username: "benchmark",
		Email:    "benchmark@example.com",
	}
	userJSON, _ := json.Marshal(user)
	
	// Reset the timer to exclude setup time
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		var u database.User
		_ = json.Unmarshal(userJSON, &u)
	}
}