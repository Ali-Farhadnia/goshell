package commands

import (
	"context"
	"fmt"
	"io"

	"github.com/Ali-Farhadnia/goshell/internal/service/user"
)

// AddUserCommand implements the adduser command
type AddUserCommand struct {
	userSVC *user.Service
}

// New creates a new adduser command
func NewAddUserCommand(userSVC *user.Service) *AddUserCommand {
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
func (c *AddUserCommand) Execute(ctx context.Context, args []string, inputReader io.Reader, outputWriter, errorOutputWriter io.Writer) error {
	if len(args) == 0 {
		_, err := fmt.Fprintf(errorOutputWriter, "usage: adduser <username>\n")
		return err
	}

	// todo: add password
	username := args[0]
	user, err := c.userSVC.CreateUser(username)
	if err != nil {
		_, err = fmt.Fprintf(errorOutputWriter, "error creating user: %v\n", err)
		return err
	}

	_, err = fmt.Fprintf(outputWriter, "User added: %s\n", user.Username)
	if err != nil {
		_, err = fmt.Fprintf(errorOutputWriter, "error writing output: %v\n", err)
		return err
	}

	return nil
}

// Help returns the help text
func (c *AddUserCommand) Help() string {
	return "adduser <username> - Add a new user"
}
