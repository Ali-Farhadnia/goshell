package commands

import (
	"context"
	"fmt"
	"io"

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
func (c *PWDCommand) Execute(ctx context.Context, args []string, inputReader io.Reader, outputWriter, errorOutputWriter io.Writer) error {
	session, err := c.sessionRepo.GetSession()
	if err != nil {
		_, err := fmt.Fprintf(errorOutputWriter, "session error: %v\n", err)
		return err
	}

	dirPath := session.WorkingDir
	if len(args) > 0 {
		_, err := fmt.Fprintf(errorOutputWriter, "pwd: too many arguments\n")
		return err
	}

	_, err = fmt.Fprintf(outputWriter, "%s\n", dirPath)
	if err != nil {
		_, err := fmt.Fprintf(errorOutputWriter, "error writing output: %v\n", err)
		return err
	}

	return nil
}

// Help returns the help text
func (c *PWDCommand) Help() string {
	return "pwd - Prints the current working directory"
}
