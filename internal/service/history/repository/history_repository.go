package repository

import (
	"github.com/Ali-Farhadnia/goshell/internal/database"
	"github.com/Ali-Farhadnia/goshell/internal/service/history"
)

// repository implements the Repository interface
type Repository struct {
	db *database.DB
}

// New creates a new Repository
func New(db *database.DB) *Repository {
	return &Repository{db: db}
}

// SaveCommand saves a command to history
func (r *Repository) SaveCommand(command *history.CommandHistory) error {
	return r.db.Create(command).Error
}

// GetUserHistory gets command history for a user
func (r *Repository) GetUserHistory(userID int64, limit int) ([]history.CommandHistory, error) {
	var history []history.CommandHistory

	query := r.db.Where("user_id = ?", userID).Order("created_at DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}

	result := query.Find(&history)
	return history, result.Error
}

// ClearUserHistory clears command history for a user
func (r *Repository) ClearUserHistory(userID int64) error {
	return r.db.Where("user_id = ?", userID).Delete(&history.CommandHistory{}).Error
}

// GetUserCommandStatsGORM returns command execution counts for a user, sorted by creation time
func (r *Repository) GetUserCommandStats(userID int64) ([]history.CommandStats, error) {
	var stats []history.CommandStats

	result := r.db.Model(&history.CommandHistory{}).
		Select("command, COUNT(*) as count").
		Where("user_id = ?", userID).
		Group("command").
		Order("MAX(created_at) DESC").
		Find(&stats)

	return stats, result.Error
}
