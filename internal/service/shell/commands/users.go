package commands

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/Ali-Farhadnia/goshell/internal/service/user"
)

// UsersCommand implements the users command
type UsersCommand struct {
	userSVC *user.Service
}

// NewUsersCommand creates a new users command
func NewUsersCommand(userSVC *user.Service) *UsersCommand {
	return &UsersCommand{
		userSVC: userSVC,
	}
}

// Name returns the command name
func (c *UsersCommand) Name() string {
	return "users"
}

// MaxArguments returns the maximum number of arguments allowed for the Command.
func (c *UsersCommand) MaxArguments() int {
	return 0
}

// Execute runs the command
func (c *UsersCommand) Execute(ctx context.Context, args []string, inputReader io.Reader, outputWriter, errorOutputWriter io.Writer) error {
	users, err := c.userSVC.ListUsers()
	if err != nil {
		_, err = fmt.Fprintf(errorOutputWriter, "%v\n", err)
		return err
	}

	var result strings.Builder
	_, err = result.WriteString("Registered users:\n")
	if err != nil {
		return err
	}

	_, err = result.WriteString("----------------\n")
	if err != nil {
		return err
	}

	for _, user := range users {
		lastLogin := "Never"
		if user.LastLogin != nil {
			lastLogin = user.LastLogin.Format(time.RFC822)
		}

		_, err = result.WriteString(fmt.Sprintf("%-15s Last login: %s\n", user.Username, lastLogin))
		if err != nil {
			return err
		}
	}

	_, err = fmt.Fprintf(outputWriter, "%s", result.String())
	if err != nil {
		_, err = fmt.Fprintf(errorOutputWriter, "error writing output: %v\n", err)
		return err
	}

	return nil
}

// Help returns the help text
func (c *UsersCommand) Help() string {
	return "users - List all registered users"
}
