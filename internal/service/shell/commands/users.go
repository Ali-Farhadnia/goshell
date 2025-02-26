package commands

import (
	"context"
	"fmt"
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
func (c *UsersCommand) Execute(ctx context.Context, args []string) (string, error) {
	users, err := c.userSVC.ListUsers()
	if err != nil {
		return "", err
	}

	var result strings.Builder
	result.WriteString("Registered users:\n")
	result.WriteString("----------------\n")
	for _, user := range users {
		lastLogin := "Never"
		if user.LastLogin != nil {
			lastLogin = user.LastLogin.Format(time.RFC822)
		}
		result.WriteString(fmt.Sprintf("%-15s Last login: %s\n", user.Username, lastLogin))
	}

	return result.String(), nil
}

// Help returns the help text
func (c *UsersCommand) Help() string {
	return "users - List all registered users"
}
