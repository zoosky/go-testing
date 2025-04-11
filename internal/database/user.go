package database

import (
	"errors"
	"sync"
)

// User represents a user in the system
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// UserRepository interface defines methods for user data operations
type UserRepository interface {
	GetUser(id int) (*User, error)
	CreateUser(user *User) error
	UpdateUser(user *User) error
	DeleteUser(id int) error
	ListUsers() ([]*User, error)
}

// InMemoryUserRepository implements UserRepository with an in-memory storage
type InMemoryUserRepository struct {
	users map[int]*User
	mutex sync.RWMutex
	nextID int
}

// NewUserRepository creates a new InMemoryUserRepository
func NewUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users:  make(map[int]*User),
		mutex:  sync.RWMutex{},
		nextID: 1,
	}
}

// GetUser retrieves a user by ID
func (r *InMemoryUserRepository) GetUser(id int) (*User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	user, exists := r.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	
	return user, nil
}

// CreateUser adds a new user to the repository
func (r *InMemoryUserRepository) CreateUser(user *User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	// Assign a new ID
	user.ID = r.nextID
	r.nextID++
	
	// Store the user
	r.users[user.ID] = user
	
	return nil
}

// UpdateUser updates an existing user
func (r *InMemoryUserRepository) UpdateUser(user *User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	if _, exists := r.users[user.ID]; !exists {
		return errors.New("user not found")
	}
	
	r.users[user.ID] = user
	
	return nil
}

// DeleteUser removes a user from the repository
func (r *InMemoryUserRepository) DeleteUser(id int) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	if _, exists := r.users[id]; !exists {
		return errors.New("user not found")
	}
	
	delete(r.users, id)
	
	return nil
}

// ListUsers returns all users in the repository
func (r *InMemoryUserRepository) ListUsers() ([]*User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	users := make([]*User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}
	
	return users, nil
}