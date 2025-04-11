package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-testing/internal/calculator"
	"go-testing/internal/database"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// setupTestServer creates a test server with mocked dependencies
func setupTestServer() (*Server, *database.MockUserRepository, *calculator.Calculator) {
	mockRepo := new(database.MockUserRepository)
	calc := calculator.NewCalculator()
	server := NewServer(mockRepo, calc)
	
	return server, mockRepo, calc
}

// TestListUsers tests the list users endpoint
func TestListUsers(t *testing.T) {
	server, mockRepo, _ := setupTestServer()
	
	// Mock data
	mockUsers := []*database.User{
		{ID: 1, Username: "user1", Email: "user1@example.com"},
		{ID: 2, Username: "user2", Email: "user2@example.com"},
	}
	
	// Setup mock expectations
	mockRepo.On("ListUsers").Return(mockUsers, nil)
	
	// Create a request
	req := httptest.NewRequest("GET", "/users", nil)
	rec := httptest.NewRecorder()
	
	// Serve the request
	server.Router().ServeHTTP(rec, req)
	
	// Assert response
	assert.Equal(t, http.StatusOK, rec.Code)
	
	// Verify the response body contains the expected users
	var responseUsers []*database.User
	err := json.NewDecoder(rec.Body).Decode(&responseUsers)
	assert.NoError(t, err)
	assert.Equal(t, len(mockUsers), len(responseUsers))
	
	// Verify the mock was called
	mockRepo.AssertExpectations(t)
}

// TestGetUser tests the get user endpoint
func TestGetUser(t *testing.T) {
	server, mockRepo, _ := setupTestServer()
	
	// Test cases
	tests := []struct {
		name           string
		userID         int
		mockUser       *database.User
		mockError      error
		expectedStatus int
	}{
		{
			name:           "Existing user",
			userID:         1,
			mockUser:       &database.User{ID: 1, Username: "user1", Email: "user1@example.com"},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Non-existent user",
			userID:         999,
			mockUser:       nil,
			mockError:      fmt.Errorf("user not found"),
			expectedStatus: http.StatusNotFound,
		},
	}
	
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup mock expectations for this test case
			mockRepo.On("GetUser", tc.userID).Return(tc.mockUser, tc.mockError).Once()
			
			// Create a request
			req := httptest.NewRequest("GET", fmt.Sprintf("/users/%d", tc.userID), nil)
			rec := httptest.NewRecorder()
			
			// Serve the request
			server.Router().ServeHTTP(rec, req)
			
			// Assert response status
			assert.Equal(t, tc.expectedStatus, rec.Code)
			
			// If expecting a success response, verify the user data
			if tc.expectedStatus == http.StatusOK {
				var user database.User
				err := json.NewDecoder(rec.Body).Decode(&user)
				assert.NoError(t, err)
				assert.Equal(t, tc.mockUser.ID, user.ID)
				assert.Equal(t, tc.mockUser.Username, user.Username)
				assert.Equal(t, tc.mockUser.Email, user.Email)
			}
		})
	}
	
	// Verify all mocks were called as expected
	mockRepo.AssertExpectations(t)
}

// TestCreateUser tests the create user endpoint
func TestCreateUser(t *testing.T) {
	server, mockRepo, _ := setupTestServer()
	
	// Test user data
	newUser := database.User{
		Username: "newuser",
		Email:    "newuser@example.com",
	}
	
	// After creation, user will have an ID
	createdUser := newUser
	createdUser.ID = 1
	
	// Setup mock expectations
	mockRepo.On("CreateUser", mock.MatchedBy(func(u *database.User) bool {
		return u.Username == newUser.Username && u.Email == newUser.Email
	})).Return(nil).Run(func(args mock.Arguments) {
		// Simulate ID assignment
		user := args.Get(0).(*database.User)
		user.ID = 1
	})
	
	// Create request with JSON body
	body, _ := json.Marshal(newUser)
	req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	
	// Serve the request
	server.Router().ServeHTTP(rec, req)
	
	// Assert response
	assert.Equal(t, http.StatusCreated, rec.Code)
	
	// Verify the response contains the created user with ID
	var responseUser database.User
	err := json.NewDecoder(rec.Body).Decode(&responseUser)
	assert.NoError(t, err)
	assert.Equal(t, createdUser.ID, responseUser.ID)
	assert.Equal(t, createdUser.Username, responseUser.Username)
	assert.Equal(t, createdUser.Email, responseUser.Email)
	
	// Verify the mock was called
	mockRepo.AssertExpectations(t)
}

// TestCalculatorEndpoints tests the calculator API endpoints
func TestCalculatorEndpoints(t *testing.T) {
	server, _, _ := setupTestServer()
	
	// Define test cases for each operation
	tests := []struct {
		name           string
		endpoint       string
		a, b           float64
		expectedStatus int
		expectedResult float64
		expectError    bool
	}{
		{"Add", "/calculator/add", 5, 3, http.StatusOK, 8, false},
		{"Subtract", "/calculator/subtract", 5, 3, http.StatusOK, 2, false},
		{"Multiply", "/calculator/multiply", 5, 3, http.StatusOK, 15, false},
		{"Divide", "/calculator/divide", 6, 3, http.StatusOK, 2, false},
		{"Divide by zero", "/calculator/divide", 5, 0, http.StatusBadRequest, 0, true},
		{"Missing parameters", "/calculator/add", 0, 0, http.StatusBadRequest, 0, true},
	}
	
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var url string
			
			if tc.name == "Missing parameters" {
				url = tc.endpoint
			} else {
				url = fmt.Sprintf("%s?a=%v&b=%v", tc.endpoint, tc.a, tc.b)
			}
			
			req := httptest.NewRequest("GET", url, nil)
			rec := httptest.NewRecorder()
			
			// Serve the request
			server.Router().ServeHTTP(rec, req)
			
			// Assert response status
			assert.Equal(t, tc.expectedStatus, rec.Code)
			
			// If expecting a success response, verify the result
			if !tc.expectError {
				var response map[string]float64
				err := json.NewDecoder(rec.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, response["result"])
			}
		})
	}
}