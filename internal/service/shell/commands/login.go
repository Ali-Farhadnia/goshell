package commands

import (
	"context"
	"fmt"
	"io"

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
func (c *LoginCommand) Execute(ctx context.Context, args []string, inputReader io.Reader, outputWriter, errorOutputWriter io.Writer) error {
	if len(args) < 1 {
		_, err := fmt.Fprintf(errorOutputWriter, "usage: login <username> [password]\n")
		return err
	}

	username := args[0]
	password := ""
	if len(args) > 1 {
		password = args[1]
	}

	user, err := c.userSVC.LoginUser(username, password)
	if err != nil {
		_, err = fmt.Fprintf(errorOutputWriter, "login failed: %v\n", err)
		return err
	}

	session, err := c.sessionRepo.GetSession()
	if err != nil {
		_, err = fmt.Fprintf(errorOutputWriter, "session error: %v\n", err)
		return err
	}

	session.User = user

	err = c.sessionRepo.SetSession(session)
	if err != nil {
		_, err = fmt.Fprintf(errorOutputWriter, "session save error: %v\n", err)
		return err
	}

	_, err = fmt.Fprintf(outputWriter, "Logged in as: %s\n", user.Username)
	if err != nil {
		_, err := fmt.Fprintf(errorOutputWriter, "error writing output: %v\n", err)
		return err
	}

	return nil
}

// Help returns the help text
func (c *LoginCommand) Help() string {
	return "login <username> [password] - Login as specified user. Password is optional."
}
