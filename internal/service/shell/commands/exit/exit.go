package exit

import (
	"context"
	"fmt"
	"os"
	"strconv"
)

// ExitCommand implements the exit command
type ExitCommand struct {
	onExit func(int)
}

// New creates a new exit command
func New(onExit func(int)) *ExitCommand {
	return &ExitCommand{onExit: onExit}
}

// Name returns the command name
func (c *ExitCommand) Name() string {
	return "exit"
}

// Execute runs the command
func (c *ExitCommand) Execute(ctx context.Context, args []string) (string, error) {
	exitCode := 0 // Default exit code

	if len(args) > 0 {
		code, err := strconv.Atoi(args[0])
		if err != nil {
			return "Invalid exit code: " + args[0], err
		}
		exitCode = code
	}

	if c.onExit != nil {
		c.onExit(exitCode)
	}

	fmt.Printf("exit status %d\n", exitCode)

	os.Exit(exitCode)

	return "", nil
}

// Help returns the help text
func (c *ExitCommand) Help() string {
	return "exit [code] - Exit the shell with optional exit code"
}
