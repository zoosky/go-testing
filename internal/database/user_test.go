package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGetUser tests the GetUser method
func TestGetUser(t *testing.T) {
	// Setup - create a repository and add a user
	repo := NewUserRepository()
	user := &User{
		Username: "testuser",
		Email:    "test@example.com",
	}
	
	err := repo.CreateUser(user)
	assert.NoError(t, err)
	assert.NotEqual(t, 0, user.ID, "User ID should be assigned")
	
	// Test - retrieve the user
	retrievedUser, err := repo.GetUser(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, retrievedUser.ID)
	assert.Equal(t, user.Username, retrievedUser.Username)
	assert.Equal(t, user.Email, retrievedUser.Email)
	
	// Test - try to get a non-existent user
	_, err = repo.GetUser(999)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

// TestCreateUser tests the CreateUser method
func TestCreateUser(t *testing.T) {
	repo := NewUserRepository()
	
	// Create multiple users and verify IDs are assigned sequentially
	for i := 1; i <= 3; i++ {
		user := &User{
			Username: "user",
			Email:    "user@example.com",
		}
		
		err := repo.CreateUser(user)
		assert.NoError(t, err)
		assert.Equal(t, i, user.ID)
	}
	
	// Verify we can find all users
	users, err := repo.ListUsers()
	assert.NoError(t, err)
	assert.Len(t, users, 3)
}

// TestUpdateUser tests the UpdateUser method
func TestUpdateUser(t *testing.T) {
	repo := NewUserRepository()
	
	// Create a user
	user := &User{
		Username: "original",
		Email:    "original@example.com",
	}
	
	err := repo.CreateUser(user)
	assert.NoError(t, err)
	
	// Update the user
	user.Username = "updated"
	user.Email = "updated@example.com"
	
	err = repo.UpdateUser(user)
	assert.NoError(t, err)
	
	// Verify the update
	retrievedUser, err := repo.GetUser(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, "updated", retrievedUser.Username)
	assert.Equal(t, "updated@example.com", retrievedUser.Email)
	
	// Try to update non-existent user
	nonExistentUser := &User{
		ID:       999,
		Username: "nonexistent",
		Email:    "nonexistent@example.com",
	}
	
	err = repo.UpdateUser(nonExistentUser)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

// TestDeleteUser tests the DeleteUser method
func TestDeleteUser(t *testing.T) {
	repo := NewUserRepository()
	
	// Create a user
	user := &User{
		Username: "delete_me",
		Email:    "delete@example.com",
	}
	
	err := repo.CreateUser(user)
	assert.NoError(t, err)
	
	// Verify the user exists
	_, err = repo.GetUser(user.ID)
	assert.NoError(t, err)
	
	// Delete the user
	err = repo.DeleteUser(user.ID)
	assert.NoError(t, err)
	
	// Verify the user no longer exists
	_, err = repo.GetUser(user.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
	
	// Try to delete a non-existent user
	err = repo.DeleteUser(999)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

// TestListUsers tests the ListUsers method
func TestListUsers(t *testing.T) {
	repo := NewUserRepository()
	
	// Initially, no users
	users, err := repo.ListUsers()
	assert.NoError(t, err)
	assert.Empty(t, users)
	
	// Add some users
	userCount := 5
	for i := 0; i < userCount; i++ {
		user := &User{
			Username: "user",
			Email:    "user@example.com",
		}
		err := repo.CreateUser(user)
		assert.NoError(t, err)
	}
	
	// Verify all users are listed
	users, err = repo.ListUsers()
	assert.NoError(t, err)
	assert.Len(t, users, userCount)
}