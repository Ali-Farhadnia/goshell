package repository

import (
	"errors"
	"sync"

	"github.com/Ali-Farhadnia/goshell/internal/service/shell"
)

// InMemoryCommandRepository is an in-memory implementation of CommandRepository.
type InMemoryCommandRepository struct {
	commands map[string]shell.Command
	mu       sync.RWMutex
}

// NewInMemoryCommandRepository creates a new instance of an in-memory command repository.
func NewInMemoryCommandRepository() *InMemoryCommandRepository {
	return &InMemoryCommandRepository{
		commands: make(map[string]shell.Command),
	}
}

// Register adds a new command to the repository.
func (r *InMemoryCommandRepository) Register(cmd shell.Command) error {
	if cmd == nil {
		return errors.New("command cannot be nil")
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	name := cmd.Name()
	if _, exists := r.commands[name]; exists {
		return errors.New("command already exists: " + name)
	}
	r.commands[name] = cmd
	return nil
}

// Get retrieves a command by name.
func (r *InMemoryCommandRepository) Get(name string) (shell.Command, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	cmd, exists := r.commands[name]
	if !exists {
		return nil, shell.ErrCommandNotFound
	}

	return cmd, nil
}

// List returns all registered commands.
func (r *InMemoryCommandRepository) List() ([]shell.Command, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if len(r.commands) == 0 {
		return nil, errors.New("no commands registered")
	}

	cmds := make([]shell.Command, 0, len(r.commands))
	for _, cmd := range r.commands {
		cmds = append(cmds, cmd)
	}
	return cmds, nil
}
