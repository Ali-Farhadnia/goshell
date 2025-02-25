package repository

import (
	"sort"
	"sync"
	"time"

	"github.com/Ali-Farhadnia/goshell/internal/service/history"
)

// InMemoryRepository implements the Repository interface with in-memory storage
type InMemoryRepository struct {
	mu          sync.RWMutex
	commands    map[int64][]history.CommandHistory // map[userID][]CommandHistory
	lastID      int64
	commandMeta map[int64]map[string]time.Time // map[userID][command]lastExecutionTime
}

// NewInMemory creates a new in-memory Repository
func NewInMemory() *InMemoryRepository {
	return &InMemoryRepository{
		commands:    make(map[int64][]history.CommandHistory),
		commandMeta: make(map[int64]map[string]time.Time),
	}
}

// SaveCommand saves a command to history
func (r *InMemoryRepository) SaveCommand(command *history.CommandHistory) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Auto-increment ID if not provided
	if command.ID == 0 {
		r.lastID++
		command.ID = r.lastID
	}

	// Set timestamps
	now := time.Now()
	if command.CreatedAt.IsZero() {
		command.CreatedAt = now
	}
	command.UpdatedAt = now

	// Store command
	r.commands[command.UserID] = append(r.commands[command.UserID], *command)

	// Update metadata for command stats sorting
	if _, exists := r.commandMeta[command.UserID]; !exists {
		r.commandMeta[command.UserID] = make(map[string]time.Time)
	}
	r.commandMeta[command.UserID][command.Command] = command.CreatedAt

	return nil
}

// GetUserHistory gets command history for a user
func (r *InMemoryRepository) GetUserHistory(userID int64, limit int) ([]history.CommandHistory, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	commands, exists := r.commands[userID]
	if !exists {
		return []history.CommandHistory{}, nil
	}

	// Sort commands by creation time (newest first)
	sort.Slice(commands, func(i, j int) bool {
		return commands[i].CreatedAt.After(commands[j].CreatedAt)
	})

	// Apply limit if needed
	result := commands
	if limit > 0 && limit < len(commands) {
		result = commands[:limit]
	}

	return result, nil
}

// ClearUserHistory clears command history for a user
func (r *InMemoryRepository) ClearUserHistory(userID int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Clear commands
	delete(r.commands, userID)
	delete(r.commandMeta, userID)

	return nil
}

// GetUserCommandStats returns command execution counts for a user, sorted by creation time
func (r *InMemoryRepository) GetUserCommandStats(userID int64) ([]history.CommandStats, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	commands, exists := r.commands[userID]
	if !exists {
		return []history.CommandStats{}, nil
	}

	// Count commands
	commandCounts := make(map[string]int64)
	for _, cmd := range commands {
		commandCounts[cmd.Command]++
	}

	// Create stats slice
	stats := make([]history.CommandStats, 0, len(commandCounts))
	for cmd, count := range commandCounts {
		stats = append(stats, history.CommandStats{
			Command: cmd,
			Count:   count,
		})
	}

	// Sort by most recent creation time
	meta := r.commandMeta[userID]
	sort.Slice(stats, func(i, j int) bool {
		timeI := meta[stats[i].Command]
		timeJ := meta[stats[j].Command]
		return timeI.After(timeJ)
	})

	return stats, nil
}
