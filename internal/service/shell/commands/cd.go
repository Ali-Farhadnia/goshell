package commands

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/Ali-Farhadnia/goshell/internal/service/shell"
)

// CDCommand implements the cd command
type CDCommand struct {
	sessionRepo shell.SessionRepository
}

// New creates a new cd command
func NewCDCommand(sessionRepo shell.SessionRepository) *CDCommand {
	return &CDCommand{
		sessionRepo: sessionRepo,
	}
}

// Name returns the command name
func (c *CDCommand) Name() string {
	return "cd"
}

// MaxArguments returns the maximum number of arguments allowed for the Command.
func (c *CDCommand) MaxArguments() int {
	return 1
}

// Execute runs the command
func (c *CDCommand) Execute(ctx context.Context, args []string, inputReader io.Reader, outputWriter, errorOutputWriter io.Writer) error {
	if len(args) == 0 {
		_, err := fmt.Fprintf(errorOutputWriter, "usage: cd <dir>\n")
		return err
	}

	session, err := c.sessionRepo.GetSession()
	if err != nil {
		_, err = fmt.Fprintf(errorOutputWriter, "error getting session: %v\n", err)
		return err
	}

	newPath := filepath.Join(session.WorkingDir, args[0])

	info, err := os.Stat(newPath)
	if err != nil {
		_, err = fmt.Fprintf(errorOutputWriter, "error accessing path: %v\n", err)
		return err
	}

	if !info.IsDir() {
		_, err = fmt.Fprintf(errorOutputWriter, "not a directory\n")
		return err
	}

	session.WorkingDir = newPath

	err = c.sessionRepo.SetSession(session)
	if err != nil {
		_, err = fmt.Fprintf(errorOutputWriter, "error updating session: %v\n", err)
		return err
	}

	return nil
}

// Help returns the help text
func (c *CDCommand) Help() string {
	return "cd [dir] - Changes the current working directory"
}
