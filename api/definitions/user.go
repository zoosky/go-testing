package definitions

// User represents the user resource in the API
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// UserCreateRequest represents the request body for creating a user
type UserCreateRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

// UserUpdateRequest represents the request body for updating a user
type UserUpdateRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

// UserResponse represents a user response
type UserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// UsersResponse represents a list of users
type UsersResponse struct {
	Users []UserResponse `json:"users"`
}

// ErrorResponse represents an API error
type ErrorResponse struct {
	Error string `json:"error"`
}