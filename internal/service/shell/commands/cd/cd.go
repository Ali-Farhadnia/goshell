package cd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Ali-Farhadnia/goshell/internal/service/shell"
)

// CDCommand implements the cd command
type CDCommand struct {
	sessionRepo shell.SessionRepository
}

// New creates a new cd command
func New(sessionRepo shell.SessionRepository) *CDCommand {
	return &CDCommand{
		sessionRepo: sessionRepo,
	}
}

// Name returns the command name
func (c *CDCommand) Name() string {
	return "cd"
}

// Execute runs the command
func (c *CDCommand) Execute(ctx context.Context, args []string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("usage: cd <dir>")
	}

	session, err := c.sessionRepo.GetSession()
	if err != nil {
		return "", err
	}

	newPath := filepath.Join(session.WorkingDir, args[0])

	info, err := os.Stat(newPath)
	if err != nil {
		return "", err
	}

	if !info.IsDir() {
		return "", fmt.Errorf("not a directory")
	}

	session.WorkingDir = newPath

	err = c.sessionRepo.SetSession(session)
	if err != nil {
		return "", err
	}

	return "", nil
}

// Help returns the help text
func (c *CDCommand) Help() string {
	return "cd [dir] - Changes the current working directory"
}
