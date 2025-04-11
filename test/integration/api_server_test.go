// +build integration

package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	_ "go-testing/docs" // Import swagger docs
	"go-testing/internal/api"
	"go-testing/internal/calculator"
	"go-testing/internal/database"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	serverURL string
	client    *http.Client
)

// TestMain starts the server before running tests and shuts it down afterward
func TestMain(m *testing.M) {
	// Create server with real dependencies
	repo := database.NewUserRepository()
	calc := calculator.NewCalculator()
	server := api.NewServer(repo, calc)

	// Choose a random port to avoid conflicts
	serverURL = "http://localhost:8081"
	
	// Start server in a goroutine
	_, cancel := context.WithCancel(context.Background())
	go func() {
		http.ListenAndServe(":8081", server.Router())
	}()
	
	// Give the server a moment to start
	time.Sleep(100 * time.Millisecond)
	
	// Create a client with a timeout
	client = &http.Client{
		Timeout: 5 * time.Second,
	}
	
	// Run tests
	exitCode := m.Run()
	
	// Shut down server
	cancel()
	
	// Exit with the test exit code
	os.Exit(exitCode)
}

// TestServerHealth tests that the server is running and responding
func TestServerHealth(t *testing.T) {
	// Make a simple request to the server's users endpoint
	resp, err := client.Get(serverURL + "/users")
	require.NoError(t, err, "Server should be reachable")
	defer resp.Body.Close()
	
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Should get 200 OK from users endpoint")
}

// TestSwaggerEndpoint tests that the Swagger documentation is accessible
func TestSwaggerEndpoint(t *testing.T) {
	// Test Swagger UI endpoint
	t.Run("Swagger UI", func(t *testing.T) {
		resp, err := client.Get(serverURL + "/swagger/index.html")
		require.NoError(t, err, "Swagger UI should be reachable")
		defer resp.Body.Close()
		
		assert.Equal(t, http.StatusOK, resp.StatusCode, "Should get 200 OK from Swagger UI")
		
		// Check for HTML content
		contentType := resp.Header.Get("Content-Type")
		assert.Contains(t, contentType, "text/html", "Content-Type should be HTML")
	})
	
	// Test Swagger JSON endpoint
	t.Run("Swagger JSON", func(t *testing.T) {
		resp, err := client.Get(serverURL + "/swagger/doc.json")
		require.NoError(t, err, "Swagger JSON should be reachable")
		defer resp.Body.Close()
		
		assert.Equal(t, http.StatusOK, resp.StatusCode, "Should get 200 OK from Swagger JSON")
		
		// Check for JSON content
		contentType := resp.Header.Get("Content-Type")
		assert.Contains(t, contentType, "application/json", "Content-Type should be JSON")
		
		// Check that it contains basic Swagger structure
		var swaggerDoc map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&swaggerDoc)
		assert.NoError(t, err, "Should be valid JSON")
		
		// Check for basic Swagger structure
		assert.Contains(t, swaggerDoc, "swagger", "Should contain swagger version")
		assert.Contains(t, swaggerDoc, "info", "Should contain API info")
		assert.Contains(t, swaggerDoc, "paths", "Should contain API paths")
	})
}

// TestUserCRUD tests the full CRUD cycle for users
func TestUserCRUD(t *testing.T) {
	// Create a test user
	newUser := map[string]interface{}{
		"username": "testuser",
		"email":    "test@example.com",
	}
	
	var createdUser map[string]interface{}
	
	// CREATE
	t.Run("Create User", func(t *testing.T) {
		body, _ := json.Marshal(newUser)
		resp, err := client.Post(serverURL+"/users", "application/json", bytes.NewBuffer(body))
		require.NoError(t, err, "Should be able to create user")
		defer resp.Body.Close()
		
		assert.Equal(t, http.StatusCreated, resp.StatusCode, "Should return 201 Created")
		
		err = json.NewDecoder(resp.Body).Decode(&createdUser)
		assert.NoError(t, err, "Should decode user response")
		
		assert.NotNil(t, createdUser["id"], "User should have an ID")
		assert.Equal(t, newUser["username"], createdUser["username"], "Username should match")
		assert.Equal(t, newUser["email"], createdUser["email"], "Email should match")
	})
	
	// Make sure we have a user to work with for subsequent tests
	require.NotNil(t, createdUser, "User should be created before continuing tests")
	userID := int(createdUser["id"].(float64))
	
	// READ
	t.Run("Get User", func(t *testing.T) {
		resp, err := client.Get(fmt.Sprintf("%s/users/%d", serverURL, userID))
		require.NoError(t, err, "Should be able to get user")
		defer resp.Body.Close()
		
		assert.Equal(t, http.StatusOK, resp.StatusCode, "Should return 200 OK")
		
		var user map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&user)
		assert.NoError(t, err, "Should decode user response")
		
		assert.Equal(t, float64(userID), user["id"], "ID should match")
		assert.Equal(t, newUser["username"], user["username"], "Username should match")
		assert.Equal(t, newUser["email"], user["email"], "Email should match")
	})
	
	// LIST
	t.Run("List Users", func(t *testing.T) {
		resp, err := client.Get(serverURL + "/users")
		require.NoError(t, err, "Should be able to list users")
		defer resp.Body.Close()
		
		assert.Equal(t, http.StatusOK, resp.StatusCode, "Should return 200 OK")
		
		var users []map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&users)
		assert.NoError(t, err, "Should decode users response")
		
		assert.NotEmpty(t, users, "Users list should not be empty")
		
		// Find our created user in the list
		found := false
		for _, user := range users {
			if int(user["id"].(float64)) == userID {
				found = true
				assert.Equal(t, newUser["username"], user["username"], "Username should match")
				assert.Equal(t, newUser["email"], user["email"], "Email should match")
				break
			}
		}
		assert.True(t, found, "Created user should be found in list")
	})
	
	// UPDATE
	t.Run("Update User", func(t *testing.T) {
		updatedUser := map[string]interface{}{
			"id":       userID,
			"username": "updateduser",
			"email":    "updated@example.com",
		}
		
		body, _ := json.Marshal(updatedUser)
		
		req, _ := http.NewRequest("PUT", fmt.Sprintf("%s/users/%d", serverURL, userID), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		
		resp, err := client.Do(req)
		require.NoError(t, err, "Should be able to update user")
		defer resp.Body.Close()
		
		assert.Equal(t, http.StatusOK, resp.StatusCode, "Should return 200 OK")
		
		var user map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&user)
		assert.NoError(t, err, "Should decode user response")
		
		assert.Equal(t, float64(userID), user["id"], "ID should match")
		assert.Equal(t, updatedUser["username"], user["username"], "Username should be updated")
		assert.Equal(t, updatedUser["email"], user["email"], "Email should be updated")
		
		// Verify update with GET
		resp, err = client.Get(fmt.Sprintf("%s/users/%d", serverURL, userID))
		require.NoError(t, err, "Should be able to get updated user")
		defer resp.Body.Close()
		
		var updatedUserData map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&updatedUserData)
		assert.NoError(t, err, "Should decode user response")
		
		assert.Equal(t, updatedUser["username"], updatedUserData["username"], "Username should be updated")
		assert.Equal(t, updatedUser["email"], updatedUserData["email"], "Email should be updated")
	})
	
	// DELETE
	t.Run("Delete User", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/users/%d", serverURL, userID), nil)
		
		resp, err := client.Do(req)
		require.NoError(t, err, "Should be able to delete user")
		defer resp.Body.Close()
		
		assert.Equal(t, http.StatusNoContent, resp.StatusCode, "Should return 204 No Content")
		
		// Verify deletion with GET
		resp, err = client.Get(fmt.Sprintf("%s/users/%d", serverURL, userID))
		require.NoError(t, err, "Should be able to attempt get on deleted user")
		defer resp.Body.Close()
		
		assert.Equal(t, http.StatusNotFound, resp.StatusCode, "Should return 404 Not Found for deleted user")
	})
}

// TestCalculatorOperations tests all calculator API endpoints
func TestCalculatorOperations(t *testing.T) {
	tests := []struct {
		name         string
		endpoint     string
		a, b         float64
		expectedCode int
		expectedVal  float64
		expectError  bool
	}{
		{"Add", "/calculator/add", 5, 3, http.StatusOK, 8, false},
		{"Subtract", "/calculator/subtract", 5, 3, http.StatusOK, 2, false},
		{"Multiply", "/calculator/multiply", 5, 3, http.StatusOK, 15, false},
		{"Divide", "/calculator/divide", 6, 3, http.StatusOK, 2, false},
		{"Divide by zero", "/calculator/divide", 5, 0, http.StatusBadRequest, 0, true},
	}
	
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("%s%s?a=%v&b=%v", serverURL, tc.endpoint, tc.a, tc.b)
			
			resp, err := client.Get(url)
			require.NoError(t, err, "API request should not fail")
			defer resp.Body.Close()
			
			assert.Equal(t, tc.expectedCode, resp.StatusCode, "Status code should match expected")
			
			if !tc.expectError {
				var result map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&result)
				assert.NoError(t, err, "Should decode result")
				
				assert.Equal(t, tc.expectedVal, result["result"], "Result should match expected")
			}
		})
	}
}

// TestMissingEndpoint tests handling of a non-existent endpoint
func TestMissingEndpoint(t *testing.T) {
	resp, err := client.Get(serverURL + "/non-existent")
	require.NoError(t, err, "Request to non-existent endpoint should not error")
	defer resp.Body.Close()
	
	assert.Equal(t, http.StatusNotFound, resp.StatusCode, "Should return 404 Not Found")
}

// TestBadRequest tests handling of a bad request
func TestBadRequest(t *testing.T) {
	// Send invalid JSON to the create user endpoint
	resp, err := client.Post(
		serverURL+"/users", 
		"application/json",
		bytes.NewBufferString("{invalid json}"),
	)
	require.NoError(t, err, "Request with invalid JSON should not error")
	defer resp.Body.Close()
	
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Should return 400 Bad Request")
}