package commands

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Ali-Farhadnia/goshell/internal/service/shell"
)

// CatCommand implements the cat command
type CatCommand struct {
	sessionRepo shell.SessionRepository
}

// New creates a new cat command
func NewCatCommand(sessionRepo shell.SessionRepository) *CatCommand {
	return &CatCommand{
		sessionRepo: sessionRepo,
	}
}

// Name returns the command name
func (c *CatCommand) Name() string {
	return "cat"
}

// MaxArguments returns the maximum number of arguments allowed for the Command.
func (c *CatCommand) MaxArguments() int {
	return 1
}

// Execute runs the command
func (c *CatCommand) Execute(ctx context.Context, args []string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("usage: cat <filename>")
	}

	session, err := c.sessionRepo.GetSession()
	if err != nil {
		return "", err
	}

	filePath := args[0]
	if !filepath.IsAbs(filePath) {
		filePath = filepath.Join(session.WorkingDir, args[0])
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// Help returns the help text
func (c *CatCommand) Help() string {
	return "cat <filename> - Displays the content of the specified file"
}
