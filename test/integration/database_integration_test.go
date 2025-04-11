// +build integration

package integration

import (
	"testing"

	"go-testing/internal/database"
	"github.com/stretchr/testify/assert"
)

// Note: The TestMain function has been moved to api_server_test.go to avoid multiple definitions

// TestRepositoryConcurrency tests concurrent operations on the repository
func TestRepositoryConcurrency(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	
	repo := database.NewUserRepository()
	
	// Create a user to work with
	user := &database.User{
		Username: "concurrent",
		Email:    "concurrent@example.com",
	}
	
	err := repo.CreateUser(user)
	assert.NoError(t, err)
	
	// Simulate concurrent reads
	t.Run("Concurrent reads", func(t *testing.T) {
		t.Parallel()
		
		for i := 0; i < 100; i++ {
			go func() {
				_, err := repo.GetUser(user.ID)
				assert.NoError(t, err)
			}()
		}
	})
	
	// Simulate concurrent writes (would be more meaningful with a real database)
	t.Run("Concurrent writes", func(t *testing.T) {
		t.Parallel()
		
		for i := 0; i < 10; i++ {
			go func(idx int) {
				newUser := &database.User{
					Username: "user",
					Email:    "user@example.com",
				}
				err := repo.CreateUser(newUser)
				assert.NoError(t, err)
			}(i)
		}
	})
}

// Additional integration tests would use real dependencies like databases
// For example:
//
// TestDatabaseConnection would test connecting to a real database
// TestTransactionRollback would test transaction rollback with a real database
// TestDatabaseReconnection would test reconnection behavior after connection failure