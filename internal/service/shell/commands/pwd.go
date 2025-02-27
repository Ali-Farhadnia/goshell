package commands

import (
	"context"
	"fmt"

	"github.com/Ali-Farhadnia/goshell/internal/service/shell"
)

// PWDCommand implements the pwd command
type PWDCommand struct {
	sessionRepo shell.SessionRepository
}

// NewPWDCommand creates a new pwd command
func NewPWDCommand(sessionRepo shell.SessionRepository) *PWDCommand {
	return &PWDCommand{
		sessionRepo: sessionRepo,
	}
}

// Name returns the command name
func (c *PWDCommand) Name() string {
	return "pwd"
}

// MaxArguments returns the maximum number of arguments allowed for the Command.
func (c *PWDCommand) MaxArguments() int {
	return 0
}

// Execute runs the command
func (c *PWDCommand) Execute(ctx context.Context, args []string) (string, error) {
	if len(args) > 0 {
		return "", fmt.Errorf("pwd: too many arguments")
	}

	session, err := c.sessionRepo.GetSession()
	if err != nil {
		return "", err
	}

	return session.WorkingDir, nil
}

// Help returns the help text
func (c *PWDCommand) Help() string {
	return "pwd - Prints the current working directory"
}
