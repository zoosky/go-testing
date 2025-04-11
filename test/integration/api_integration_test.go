// +build integration

package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-testing/internal/calculator"
	"go-testing/internal/database"
	"go-testing/internal/api"

	"github.com/stretchr/testify/assert"
)

// TestFullAPIFlow tests the entire API flow with real dependencies
func TestFullAPIFlow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	
	// Setup real dependencies (not mocks)
	repo := database.NewUserRepository()
	calc := calculator.NewCalculator()
	server := api.NewServer(repo, calc)
	
	// Create a test server
	ts := httptest.NewServer(server.Router())
	defer ts.Close()
	
	// Create a new user
	t.Run("Create user", func(t *testing.T) {
		newUser := database.User{
			Username: "integration",
			Email:    "integration@example.com",
		}
		
		body, _ := json.Marshal(newUser)
		resp, err := http.Post(ts.URL+"/users", "application/json", bytes.NewBuffer(body))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		
		var createdUser database.User
		json.NewDecoder(resp.Body).Decode(&createdUser)
		resp.Body.Close()
		
		assert.NotEqual(t, 0, createdUser.ID)
		assert.Equal(t, newUser.Username, createdUser.Username)
		assert.Equal(t, newUser.Email, createdUser.Email)
		
		// List users
		t.Run("List users", func(t *testing.T) {
			resp, err := http.Get(ts.URL + "/users")
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			
			var users []*database.User
			json.NewDecoder(resp.Body).Decode(&users)
			resp.Body.Close()
			
			assert.NotEmpty(t, users)
			assert.Contains(t, extractUserIDs(users), createdUser.ID)
		})
		
		// Get user
		t.Run("Get user", func(t *testing.T) {
			resp, err := http.Get(fmt.Sprintf("%s/users/%d", ts.URL, createdUser.ID))
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			
			var user database.User
			json.NewDecoder(resp.Body).Decode(&user)
			resp.Body.Close()
			
			assert.Equal(t, createdUser.ID, user.ID)
			assert.Equal(t, createdUser.Username, user.Username)
			assert.Equal(t, createdUser.Email, user.Email)
		})
		
		// Update user
		t.Run("Update user", func(t *testing.T) {
			createdUser.Username = "updated"
			createdUser.Email = "updated@example.com"
			
			body, _ := json.Marshal(createdUser)
			
			req, _ := http.NewRequest("PUT", fmt.Sprintf("%s/users/%d", ts.URL, createdUser.ID), bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := http.DefaultClient.Do(req)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			
			var updatedUser database.User
			json.NewDecoder(resp.Body).Decode(&updatedUser)
			resp.Body.Close()
			
			assert.Equal(t, createdUser.ID, updatedUser.ID)
			assert.Equal(t, "updated", updatedUser.Username)
			assert.Equal(t, "updated@example.com", updatedUser.Email)
		})
		
		// Delete user
		t.Run("Delete user", func(t *testing.T) {
			req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/users/%d", ts.URL, createdUser.ID), nil)
			
			resp, err := http.DefaultClient.Do(req)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusNoContent, resp.StatusCode)
			resp.Body.Close()
			
			// Verify user is gone
			resp, err = http.Get(fmt.Sprintf("%s/users/%d", ts.URL, createdUser.ID))
			assert.NoError(t, err)
			assert.Equal(t, http.StatusNotFound, resp.StatusCode)
			resp.Body.Close()
		})
	})
	
	// Test calculator endpoints
	t.Run("Calculator operations", func(t *testing.T) {
		tests := []struct {
			name        string
			endpoint    string
			a, b        float64
			expected    float64
			expectError bool
		}{
			{"Add", "/calculator/add", 5, 3, 8, false},
			{"Subtract", "/calculator/subtract", 5, 3, 2, false},
			{"Multiply", "/calculator/multiply", 5, 3, 15, false},
			{"Divide", "/calculator/divide", 6, 3, 2, false},
			{"Divide by zero", "/calculator/divide", 5, 0, 0, true},
		}
		
		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				url := fmt.Sprintf("%s%s?a=%v&b=%v", ts.URL, tc.endpoint, tc.a, tc.b)
				
				resp, err := http.Get(url)
				assert.NoError(t, err)
				
				if tc.expectError {
					assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
				} else {
					assert.Equal(t, http.StatusOK, resp.StatusCode)
					
					var result map[string]float64
					json.NewDecoder(resp.Body).Decode(&result)
					resp.Body.Close()
					
					assert.Equal(t, tc.expected, result["result"])
				}
			})
		}
	})
}

// Helper function to extract user IDs from a slice of users
func extractUserIDs(users []*database.User) []int {
	ids := make([]int, len(users))
	for i, user := range users {
		ids[i] = user.ID
	}
	return ids
}