package shell

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/Ali-Farhadnia/goshell/pkg/execpath"
)

// SystemCommand executes system commands
type SystemCommand struct {
	sessionRepo SessionRepository
	path        string
}

// NewSystemCommand creates a new system command handler
func NewSystemCommand(
	sessionRepo SessionRepository,
	path string,
) *SystemCommand {
	return &SystemCommand{
		sessionRepo: sessionRepo,
		path:        path,
	}
}

// Execute runs the system command
func (c *SystemCommand) Execute(ctx context.Context, cmdName string, args []string, inputReader io.Reader, outputWriter, errorOutputWriter io.Writer) error {
	// Check if it's an executable in $PATH
	cmdPath, err := execpath.FindExecutable(cmdName, c.path)
	if err != nil {
		_, err := fmt.Fprintf(errorOutputWriter, "command not found: %s\n", cmdName)
		return err
	}

	// Get the current working directory from session
	session, err := c.sessionRepo.GetSession()
	if err != nil {
		_, err := fmt.Fprintf(errorOutputWriter, "error getting session: %v\n", err)
		return err
	}

	// Prepare command execution
	cmd := exec.CommandContext(ctx, cmdPath, args...)
	cmd.Dir = session.WorkingDir
	cmd.Env = os.Environ()

	// Set up input, output, and error streams
	cmd.Stdin = inputReader
	cmd.Stdout = outputWriter
	cmd.Stderr = errorOutputWriter

	// Execute command
	err = cmd.Run()
	if err != nil {
		_, err := fmt.Fprintf(errorOutputWriter, "command execution failed: %v\n", err)
		return err
	}

	return nil
}

// Help returns the help text
func (c *SystemCommand) Help() string {
	return "<command> [args...] - Executes a system command"
}
