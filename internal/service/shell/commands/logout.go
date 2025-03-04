package commands

import (
	"context"
	"fmt"
	"io"

	"github.com/Ali-Farhadnia/goshell/internal/service/shell"
)

// LogoutCommand implements the logout command
type LogoutCommand struct {
	sessionRepo shell.SessionRepository
}

// NewLogoutCommand creates a new logout command
func NewLogoutCommand(sessionRepo shell.SessionRepository) *LogoutCommand {
	return &LogoutCommand{
		sessionRepo: sessionRepo,
	}
}

// Name returns the command name
func (c *LogoutCommand) Name() string {
	return "logout"
}

// MaxArguments returns the maximum number of arguments allowed for the Command.
func (c *LogoutCommand) MaxArguments() int {
	return 0
}

// Execute runs the command
func (c *LogoutCommand) Execute(ctx context.Context, args []string, inputReader io.Reader, outputWriter, errorOutputWriter io.Writer) error {
	session, err := c.sessionRepo.GetSession()
	if err != nil {
		_, err = fmt.Fprintf(errorOutputWriter, "session error: %v\n", err)
		return err
	}

	session.User = nil

	err = c.sessionRepo.SetSession(session)
	if err != nil {
		_, err = fmt.Fprintf(errorOutputWriter, "session save error: %v\n", err)
		return err
	}

	_, err = fmt.Fprintf(outputWriter, "Logged out.\n")
	if err != nil {
		_, err := fmt.Fprintf(errorOutputWriter, "error writing output: %v\n", err)
		return err
	}

	return nil
}

// Help returns the help text
func (c *LogoutCommand) Help() string {
	return "logout - Logout current user and return to guest"
}
