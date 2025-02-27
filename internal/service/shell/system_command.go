package shell

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"strings"

	"github.com/Ali-Farhadnia/goshell/pkg/execpath"
)

// SystemCommand executes system commands
type SystemCommand struct {
	sessionRepo SessionRepository
}

// NewSystemCommand creates a new system command handler
func NewSystemCommand(sessionRepo SessionRepository) *SystemCommand {
	return &SystemCommand{
		sessionRepo: sessionRepo,
	}
}

// Execute runs the system command
func (c *SystemCommand) Execute(ctx context.Context, cmdName string, args []string) (string, error) {
	// Check if it's an executable in $PATH
	cmdPath, err := execpath.FindExecutable(cmdName)
	if err != nil {
		return "", err
	}

	// Get the current working directory from session
	session, err := c.sessionRepo.GetSession()
	if err != nil {
		return "", err
	}

	// Prepare command execution
	cmd := exec.CommandContext(ctx, cmdPath, args...)
	cmd.Dir = session.WorkingDir
	cmd.Env = os.Environ()

	// Capture output
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Execute command
	err = cmd.Run()
	if err != nil {
		return stderr.String(), err
	}

	return strings.TrimSpace(stdout.String()), nil
}

// Help returns the help text
func (c *SystemCommand) Help() string {
	return "<command> [args...] - Executes a system command"
}
