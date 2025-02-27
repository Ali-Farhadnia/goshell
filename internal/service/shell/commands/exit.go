package commands

import (
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
)

// ExitCommand implements the exit command
type ExitCommand struct {
	onExit func(int)
}

// New creates a new exit command
func NewExitCommand(onExit func(int)) *ExitCommand {
	return &ExitCommand{onExit: onExit}
}

// Name returns the command name
func (c *ExitCommand) Name() string {
	return "exit"
}

// MaxArguments returns the maximum number of arguments allowed for the Command.
func (c *ExitCommand) MaxArguments() int {
	return 1
}

// Execute runs the command
func (c *ExitCommand) Execute(ctx context.Context, args []string, inputReader io.Reader, outputWriter, errorOutputWriter io.Writer) error {
	exitCode := 0 // Default exit code

	if len(args) > 0 {
		code, err := strconv.Atoi(args[0])
		if err != nil {
			_, err = fmt.Fprintf(errorOutputWriter, "Invalid exit code: %s\n", args[0])
			return err
		}
		exitCode = code
	}

	if c.onExit != nil {
		c.onExit(exitCode)
	}

	_, err := fmt.Fprintf(outputWriter, "exit status %d\n", exitCode)
	if err != nil {
		_, err = fmt.Fprintf(errorOutputWriter, "error writing output: %v\n", err)
		return err
	}

	os.Exit(exitCode)

	return nil
}

// Help returns the help text
func (c *ExitCommand) Help() string {
	return "exit [code] - Exit the shell with optional exit code"
}
