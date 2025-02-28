package commands

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/Ali-Farhadnia/goshell/internal/service/shell"
)

// CatCommand implements the cat command
type CatCommand struct {
	sessionRepo shell.SessionRepository
}

// New creates a new cat command
func NewCatCommand(sessionRepo shell.SessionRepository) *CatCommand {
	return &CatCommand{
		sessionRepo: sessionRepo,
	}
}

// Name returns the command name
func (c *CatCommand) Name() string {
	return "cat"
}

// MaxArguments returns the maximum number of arguments allowed for the Command.
func (c *CatCommand) MaxArguments() int {
	return 1
}

// Execute runs the command
func (c *CatCommand) Execute(ctx context.Context, args []string, inputReader io.Reader, outputWriter, errorOutputWriter io.Writer) error {
	if len(args) == 0 {
		_, err := fmt.Fprintf(errorOutputWriter, "usage: cat <filename>\n")
		return err
	}

	session, err := c.sessionRepo.GetSession()
	if err != nil {
		_, err = fmt.Fprintf(errorOutputWriter, "error getting session: %v\n", err)
		return err
	}

	filePath := args[0]
	if !filepath.IsAbs(filePath) {
		filePath = filepath.Join(session.WorkingDir, args[0])
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		_, err = fmt.Fprintf(errorOutputWriter, "error reading file: %v\n", err)
		return err
	}

	_, err = fmt.Fprintf(outputWriter, "%s\n", string(data))
	if err != nil {
		_, err = fmt.Fprintf(errorOutputWriter, "error writing output: %v\n", err)
		return err
	}

	return nil
}

// Help returns the help text
func (c *CatCommand) Help() string {
	return "cat <filename> - Displays the content of the specified file"
}
