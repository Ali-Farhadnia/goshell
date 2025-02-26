package adduser

import (
	"context"
	"fmt"

	"github.com/Ali-Farhadnia/goshell/internal/service/user"
)

// AddUserCommand implements the adduser command
type AddUserCommand struct {
	userSVC *user.Service
}

// New creates a new adduser command
func New(userSVC *user.Service) *AddUserCommand {
	return &AddUserCommand{
		userSVC: userSVC,
	}
}

// Name returns the command name
func (c *AddUserCommand) Name() string {
	return "adduser"
}

// MaxArguments returns the maximum number of arguments allowed for the Command.
func (c *AddUserCommand) MaxArguments() int {
	return 2
}

// Execute runs the command
func (c *AddUserCommand) Execute(ctx context.Context, args []string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("usage: adduser <username>")
	}

	// todo: add password
	username := args[0]
	user, err := c.userSVC.CreateUser(username)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("User added: %s\n", user.Username), nil
}

// Help returns the help text
func (c *AddUserCommand) Help() string {
	return "adduser <username> - Add a new user"
}
