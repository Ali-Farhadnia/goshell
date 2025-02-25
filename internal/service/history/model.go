package history

import (
	"time"

	"gorm.io/gorm"
)

// CommandHistory represents a command executed by a user
type CommandHistory struct {
	ID        int64          `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Command string `gorm:"not null" json:"command"`
	UserID  int64  `gorm:"index;not null;references:users(id)" json:"user_id"`
}

// CommandStats represents aggregated command history
type CommandStats struct {
	Command string `json:"command"`
	Count   int64  `json:"count"`
}
