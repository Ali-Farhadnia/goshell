package commands

import (
	"context"
	"fmt"

	"github.com/Ali-Farhadnia/goshell/internal/service/shell"
	"github.com/Ali-Farhadnia/goshell/internal/service/user"
)

// LoginCommand implements the login command
type LoginCommand struct {
	userSVC     *user.Service
	sessionRepo shell.SessionRepository
}

// NewLoginCommand creates a new login command
func NewLoginCommand(userSVC *user.Service, sessionRepo shell.SessionRepository) *LoginCommand {
	return &LoginCommand{
		userSVC:     userSVC,
		sessionRepo: sessionRepo,
	}
}

// Name returns the command name
func (c *LoginCommand) Name() string {
	return "login"
}

// MaxArguments returns the maximum number of arguments allowed for the Command.
func (c *LoginCommand) MaxArguments() int {
	return 2
}

// Execute runs the command
func (c *LoginCommand) Execute(ctx context.Context, args []string) (string, error) {
	if len(args) < 1 {
		return "", fmt.Errorf("usage: login <username>")
	}

	username := args[0]
	user, err := c.userSVC.LoginUser(username)
	if err != nil {
		return "", err
	}

	session, err := c.sessionRepo.GetSession()
	if err != nil {
		return "", err
	}

	session.User = user

	err = c.sessionRepo.SetSession(session)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Logged in as: %s\n", user.Username), nil
}

// Help returns the help text
func (c *LoginCommand) Help() string {
	return "login <username> - Login as specified user"
}
