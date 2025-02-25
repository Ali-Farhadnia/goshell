package exit

import (
	"context"
	"os"
)

// ExitCommand implements the exit command
type ExitCommand struct {
}

// New creates a new exit command
func New(onExit func(int)) *ExitCommand {
	return &ExitCommand{}
}

// Name returns the command name
func (c *ExitCommand) Name() string {
	return "exit"
}

// Execute runs the command
func (c *ExitCommand) Execute(ctx context.Context, args []string) (string, error) {
	os.Exit(0)

	return "", nil
}

// Help returns the help text
func (c *ExitCommand) Help() string {
	return "exit [code] - Exit the shell with optional exit code"
}
