package user

import (
	"time"

	"gorm.io/gorm"
)

// User represents a shell user
type User struct {
	ID        int64          `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Username  string     `gorm:"uniqueIndex;not null" json:"username"`
	LastLogin *time.Time `json:"last_login"`
}
