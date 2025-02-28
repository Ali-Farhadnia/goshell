package repository

import (
	"errors"
	"time"

	"github.com/Ali-Farhadnia/goshell/internal/database"
	"github.com/Ali-Farhadnia/goshell/internal/service/user"
	"gorm.io/gorm"
)

// repository implements the Repository interface
type Repository struct {
	db *database.DB
}

// New creates a new Repository
func New(db *database.DB) *Repository {
	return &Repository{db: db}
}

// FindUserByUsername finds a user by username
func (r *Repository) FindUserByUsername(username string) (*user.User, error) {
	var usr user.User
	result := r.db.Where("username = ?", username).First(&usr)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, user.ErrUserNotFound
		}

		return nil, result.Error
	}

	return &usr, nil
}

// CreateUser creates a new user
func (r *Repository) CreateUser(user *user.User) error {
	return r.db.Create(user).Error
}

// UpdateUser updates an existing user
func (r *Repository) UpdateUser(user *user.User) error {
	return r.db.Save(user).Error
}

// ListUsers lists all users
func (r *Repository) ListUsers() ([]user.User, error) {
	var users []user.User
	result := r.db.Find(&users)
	return users, result.Error
}

// UpdateLastLogin updates a user's last login time
func (r *Repository) UpdateLastLogin(userID int64) error {
	now := user.User{LastLogin: new(time.Time)}
	*now.LastLogin = time.Now()

	return r.db.Model(&user.User{}).
		Where("id = ?", userID).
		Update("last_login", now.LastLogin).Error
}
