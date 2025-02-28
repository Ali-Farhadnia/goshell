package user

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

// UserRepository defines operations on users
type UserRepository interface {
	FindUserByUsername(username string) (User, error)
	CreateUser(user User) error
	UpdateUser(user User) error
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
func (s *Service) FindUser(username string) (User, error) {
	return s.userRepo.FindUserByUsername(username)
}

// CreateUser creates a new user
func (s *Service) CreateUser(username, password string) (User, error) {
	// Check if user already exists
	_, err := s.userRepo.FindUserByUsername(username)

	if err != nil {
		if !errors.Is(err, ErrUserNotFound) {
			return User{}, err
		}
	} else {
		return User{}, fmt.Errorf("user already exists: %s", username)
	}

	var passwordHash *string
	if password != "" {
		hash, err := hashPassword(password)
		if err != nil {
			return User{}, fmt.Errorf("failed to hash password: %w", err)
		}
		passwordHash = &hash
	}

	// Create user in database
	user := User{
		Username:     username,
		PasswordHash: passwordHash,
	}

	if err := s.userRepo.CreateUser(user); err != nil {
		return User{}, fmt.Errorf("failed to create user in database: %w", err)
	}

	return user, nil
}

// LoginUser logs in a user and updates last login time
func (s *Service) LoginUser(username, password string) (User, error) {
	user, err := s.userRepo.FindUserByUsername(username)
	if err != nil {
		return User{}, fmt.Errorf("error on finding user: %w", err)
	}

	// Verify password if it is set
	if user.PasswordHash != nil {
		if err := verifyPassword(*user.PasswordHash, password); err != nil {
			return User{}, fmt.Errorf("invalid password")
		}
	}

	// Update last login time
	if err := s.userRepo.UpdateLastLogin(user.ID); err != nil {
		return User{}, fmt.Errorf("failed to update last login time: %w", err)
	}

	return user, nil
}

// ListUsers lists all users
func (s *Service) ListUsers() ([]User, error) {
	return s.userRepo.ListUsers()
}

// hashPassword hashes a password using bcrypt
func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// verifyPassword checks if the provided password matches the hashed password
func verifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
