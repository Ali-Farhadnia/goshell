package logout

import (
	"context"

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

// Execute runs the command
func (c *LogoutCommand) Execute(ctx context.Context, args []string) (string, error) {
	session, err := c.sessionRepo.GetSession()
	if err != nil {
		return "", err
	}

	session.User = nil

	err = c.sessionRepo.SetSession(session)
	if err != nil {
		return "", err
	}

	return "Logged out.", nil
}

// Help returns the help text
func (c *LogoutCommand) Help() string {
	return "logout - Logout current user and return to guest"
}
