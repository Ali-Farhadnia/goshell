package user

import (
	"fmt"
)

// UserRepository defines operations on users
type UserRepository interface {
	FindUserByUsername(username string) (*User, error)
	CreateUser(user *User) error
	UpdateUser(user *User) error
	ListUsers() ([]User, error)
	UpdateLastLogin(userID int64) error
}

// Service provides high-level functionality for the shell
type Service struct {
	userRepo UserRepository
}

// New creates a new Service instance
func New(userRepo UserRepository) *Service {
	return &Service{
		userRepo: userRepo,
	}
}

// FindUser finds a user by username
func (s *Service) FindUser(username string) (*User, error) {
	return s.userRepo.FindUserByUsername(username)
}

// CreateUser creates a new user
func (s *Service) CreateUser(username string) (*User, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.FindUserByUsername(username)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("user already exists: %s", username)
	}

	// Create user in database
	user := &User{
		Username: username,
	}

	if err := s.userRepo.CreateUser(user); err != nil {
		return nil, fmt.Errorf("failed to create user in database: %w", err)
	}

	return user, nil
}

// LoginUser logs in a user and updates last login time
func (s *Service) LoginUser(username string) (*User, error) {
	user, err := s.userRepo.FindUserByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("error on finding user: %w", err)
	}

	// Update last login time
	if err := s.userRepo.UpdateLastLogin(user.ID); err != nil {
		return nil, fmt.Errorf("failed to update last login time: %w", err)
	}

	return user, nil
}

// ListUsers lists all users
func (s *Service) ListUsers() ([]User, error) {
	return s.userRepo.ListUsers()
}
