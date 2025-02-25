package history

import (
	"time"
)

// HistoryRepository defines operations on command history
type HistoryRepository interface {
	SaveCommand(command *CommandHistory) error
	GetUserHistory(userID int64, limit int) ([]CommandHistory, error)
	GetUserCommandStats(userID int64) ([]CommandStats, error)
	ClearUserHistory(userID int64) error
}

// Service provides high-level functionality for the shell
type Service struct {
	historyRepo       HistoryRepository
	guestHistoryCache HistoryRepository
	guestID           int64
}

// New creates a new Service instance
func New(
	historyRepo HistoryRepository,
	guestHistoryCache HistoryRepository,
	guestID int64,
) *Service {
	return &Service{
		historyRepo:       historyRepo,
		guestHistoryCache: guestHistoryCache,
		guestID:           guestID,
	}
}

// SaveCommandHistory saves a command to history
func (s *Service) SaveCommandHistory(userID *int64, command string) error {
	history := &CommandHistory{
		Command:   command,
		CreatedAt: time.Now(),
	}

	if userID == nil {
		history.UserID = s.guestID

		return s.guestHistoryCache.SaveCommand(history)
	}

	history.UserID = *userID

	return s.historyRepo.SaveCommand(history)
}

// GetCommandHistory retrieves command history for a user
func (s *Service) GetCommandHistoryStats(userID *int64, limit int) ([]CommandStats, error) {
	if userID == nil {
		return s.guestHistoryCache.GetUserCommandStats(s.guestID)
	}

	return s.historyRepo.GetUserCommandStats(*userID)
}

// ClearCommandHistory clears command history for a user
func (s *Service) ClearCommandHistory(userID *int64) error {
	if userID == nil {
		return s.guestHistoryCache.ClearUserHistory(s.guestID)
	}

	return s.historyRepo.ClearUserHistory(*userID)
}
