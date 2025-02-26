package commands

import (
	"context"
	"os"
	"strings"

	"github.com/Ali-Farhadnia/goshell/internal/service/shell"
)

// EchoCommand implements the echo command
type EchoCommand struct {
	sessionRepo shell.SessionRepository
}

// New creates a new echo command
func NewEchoCommand(sessionRepo shell.SessionRepository) *EchoCommand {
	return &EchoCommand{
		sessionRepo: sessionRepo,
	}
}

// Name returns the command name
func (e *EchoCommand) Name() string {
	return "echo"
}

// MaxArguments returns the maximum number of arguments allowed for the Command.
func (e *EchoCommand) MaxArguments() int {
	return -1 // Unlimited arguments
}

// Execute runs the command
func (e *EchoCommand) Execute(ctx context.Context, args []string) (string, error) {
	var output strings.Builder

	for _, arg := range args {
		if strings.HasPrefix(arg, "'") && strings.HasSuffix(arg, "'") {
			// Literal string, remove quotes
			output.WriteString(strings.Trim(arg, "'"))
		} else {
			// Handle environment variables and other strings
			var expandedArg strings.Builder
			var inEnvVar bool
			var envVarName strings.Builder

			for _, char := range arg {
				if char == '$' && !inEnvVar {
					inEnvVar = true
					envVarName.Reset()
				} else if inEnvVar && (('a' <= char && char <= 'z') || ('A' <= char && char <= 'Z') || ('0' <= char && char <= '9') || char == '_') {
					envVarName.WriteRune(char)
				} else if inEnvVar {
					envVar := os.Getenv(envVarName.String())
					expandedArg.WriteString(envVar)
					expandedArg.WriteRune(char)
					inEnvVar = false
					envVarName.Reset()
				} else {
					expandedArg.WriteRune(char)
				}
			}
			if inEnvVar {
				expandedArg.WriteString(os.Getenv(envVarName.String()))
			}

			output.WriteString(expandedArg.String())
		}
		output.WriteString(" ")
	}

	return strings.TrimSpace(output.String()), nil
}

// Help returns the help text
func (e *EchoCommand) Help() string {
	return "echo [args...] - Prints the provided arguments to the output, supports environment variables and multiple expressions"
}
