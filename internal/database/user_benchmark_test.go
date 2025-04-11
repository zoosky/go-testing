package database

import (
	"strconv"
	"testing"
)

// BenchmarkCreateUser benchmarks the CreateUser method
func BenchmarkCreateUser(b *testing.B) {
	repo := NewUserRepository()
	
	// Reset the timer to exclude setup time
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		// Need to create a new user for each iteration to avoid ID conflicts
		user := &User{
			Username: "user" + strconv.Itoa(i),
			Email:    "user" + strconv.Itoa(i) + "@example.com",
		}
		_ = repo.CreateUser(user)
	}
}

// BenchmarkGetUser benchmarks the GetUser method
func BenchmarkGetUser(b *testing.B) {
	repo := NewUserRepository()
	
	// Create a user to get
	user := &User{
		Username: "benchmark",
		Email:    "benchmark@example.com",
	}
	repo.CreateUser(user)
	
	// Reset the timer to exclude setup time
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		_, _ = repo.GetUser(user.ID)
	}
}

// BenchmarkUpdateUser benchmarks the UpdateUser method
func BenchmarkUpdateUser(b *testing.B) {
	repo := NewUserRepository()
	
	// Create a user to update
	user := &User{
		Username: "benchmark",
		Email:    "benchmark@example.com",
	}
	repo.CreateUser(user)
	
	// Reset the timer to exclude setup time
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		// Modify the user and update
		user.Username = "updated" + strconv.Itoa(i)
		_ = repo.UpdateUser(user)
	}
}

// BenchmarkDeleteUser benchmarks the DeleteUser method
func BenchmarkDeleteUser(b *testing.B) {
	repo := NewUserRepository()
	
	// We need to create users just-in-time for deletion
	// because we can't delete the same user multiple times
	users := make([]*User, b.N)
	for i := 0; i < b.N; i++ {
		user := &User{
			Username: "delete" + strconv.Itoa(i),
			Email:    "delete" + strconv.Itoa(i) + "@example.com",
		}
		repo.CreateUser(user)
		users[i] = user
	}
	
	// Reset the timer to exclude setup time
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		_ = repo.DeleteUser(users[i].ID)
	}
}

// BenchmarkListUsers benchmarks the ListUsers method
func BenchmarkListUsers(b *testing.B) {
	repo := NewUserRepository()
	
	// Create some users to list
	for i := 0; i < 100; i++ {
		user := &User{
			Username: "list" + strconv.Itoa(i),
			Email:    "list" + strconv.Itoa(i) + "@example.com",
		}
		repo.CreateUser(user)
	}
	
	// Reset the timer to exclude setup time
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		_, _ = repo.ListUsers()
	}
}

// BenchmarkConcurrentReads benchmarks concurrent read operations
func BenchmarkConcurrentReads(b *testing.B) {
	repo := NewUserRepository()
	
	// Create users to read
	numUsers := 100
	userIDs := make([]int, numUsers)
	for i := 0; i < numUsers; i++ {
		user := &User{
			Username: "concurrent" + strconv.Itoa(i),
			Email:    "concurrent" + strconv.Itoa(i) + "@example.com",
		}
		repo.CreateUser(user)
		userIDs[i] = user.ID
	}
	
	// Reset the timer to exclude setup time
	b.ResetTimer()
	b.ReportAllocs()
	
	b.RunParallel(func(pb *testing.PB) {
		// Each goroutine will read random users
		i := 0
		for pb.Next() {
			id := userIDs[i%numUsers]
			_, _ = repo.GetUser(id)
			i++
		}
	})
}

// BenchmarkConcurrentWrites benchmarks concurrent write operations
func BenchmarkConcurrentWrites(b *testing.B) {
	repo := NewUserRepository()
	
	// Reset the timer to exclude setup time
	b.ResetTimer()
	b.ReportAllocs()
	
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			user := &User{
				Username: "parallel" + strconv.Itoa(i),
				Email:    "parallel" + strconv.Itoa(i) + "@example.com",
			}
			_ = repo.CreateUser(user)
			i++
		}
	})
}