package database

import (
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

// GetUser is a mocked method
func (m *MockUserRepository) GetUser(id int) (*User, error) {
	args := m.Called(id)
	
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	
	return args.Get(0).(*User), args.Error(1)
}

// CreateUser is a mocked method
func (m *MockUserRepository) CreateUser(user *User) error {
	args := m.Called(user)
	return args.Error(0)
}

// UpdateUser is a mocked method
func (m *MockUserRepository) UpdateUser(user *User) error {
	args := m.Called(user)
	return args.Error(0)
}

// DeleteUser is a mocked method
func (m *MockUserRepository) DeleteUser(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

// ListUsers is a mocked method
func (m *MockUserRepository) ListUsers() ([]*User, error) {
	args := m.Called()
	
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	
	return args.Get(0).([]*User), args.Error(1)
}