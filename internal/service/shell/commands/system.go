package commands

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Ali-Farhadnia/goshell/internal/service/shell"
	"github.com/Ali-Farhadnia/goshell/pkg/execpath"
)

// SystemCommand executes system commands
type SystemCommand struct {
	sessionRepo shell.SessionRepository
}

// NewSystemCommand creates a new system command handler
func NewSystemCommand(sessionRepo shell.SessionRepository) *SystemCommand {
	return &SystemCommand{
		sessionRepo: sessionRepo,
	}
}

// Name returns the command name (wildcard for any system executable)
func (c *SystemCommand) Name() string {
	return "system"
}

// MaxArguments returns the maximum number of arguments allowed (-1 for unlimited).
func (c *SystemCommand) MaxArguments() int {
	return -1
}

// Execute runs the system command
func (c *SystemCommand) Execute(ctx context.Context, args []string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("usage: <command> [args...]")
	}

	// Check if it's an executable in $PATH
	cmdPath, err := execpath.FindExecutable(args[0])
	if err != nil {
		return "", err
	}

	// Get the current working directory from session
	session, err := c.sessionRepo.GetSession()
	if err != nil {
		return "", err
	}

	// Prepare command execution
	cmd := exec.CommandContext(ctx, cmdPath, args[1:]...)
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
