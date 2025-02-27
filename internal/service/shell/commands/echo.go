package commands

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
)

// EchoCommand implements the echo command
type EchoCommand struct {
}

// New creates a new echo command
func NewEchoCommand() *EchoCommand {
	return &EchoCommand{}
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
func (e *EchoCommand) Execute(ctx context.Context, args []string, inputReader io.Reader, outputWriter, errorOutputWriter io.Writer) error {
	var output strings.Builder

	// If no arguments are provided, read from inputReader
	if len(args) == 0 {
		scanner := bufio.NewScanner(inputReader)
		for scanner.Scan() {
			_, err := output.WriteString(scanner.Text())
			if err != nil {
				return err
			}

			_, err = output.WriteString("\n")
			if err != nil {
				return err
			}
		}
		if err := scanner.Err(); err != nil {
			_, err = fmt.Fprintf(errorOutputWriter, "error reading input: %v\n", err)
			return err
		}
	} else {
		// Process arguments as before
		for _, arg := range args {
			if strings.HasPrefix(arg, "'") && strings.HasSuffix(arg, "'") {
				_, err := output.WriteString(strings.Trim(arg, "'"))
				if err != nil {
					return err
				}

			} else {
				var expandedArg strings.Builder
				var inEnvVar bool
				var envVarName strings.Builder

				for _, char := range arg {
					if char == '$' && !inEnvVar {
						inEnvVar = true
						envVarName.Reset()
					} else if inEnvVar && (('a' <= char && char <= 'z') || ('A' <= char && char <= 'Z') || ('0' <= char && char <= '9') || char == '_') {
						_, err := envVarName.WriteRune(char)
						if err != nil {
							return err
						}

					} else if inEnvVar {
						envVar := os.Getenv(envVarName.String())
						_, err := expandedArg.WriteString(envVar)
						if err != nil {
							return err
						}

						_, err = expandedArg.WriteRune(char)
						if err != nil {
							return err
						}

						inEnvVar = false
						envVarName.Reset()
					} else {
						_, err := expandedArg.WriteRune(char)
						if err != nil {
							return err
						}
					}
				}
				if inEnvVar {
					_, err := expandedArg.WriteString(os.Getenv(envVarName.String()))
					if err != nil {
						return err
					}
				}

				_, err := output.WriteString(expandedArg.String())
				if err != nil {
					return err
				}
			}

			_, err := output.WriteString(" ")
			if err != nil {
				return err
			}
		}
	}

	_, err := fmt.Fprintln(outputWriter, strings.TrimSpace(output.String()))
	if err != nil {
		_, err := fmt.Fprintf(errorOutputWriter, "error writing output: %v\n", err)
		return err
	}

	return nil
}

// Help returns the help text
func (e *EchoCommand) Help() string {
	return "echo [args...] - Prints the provided arguments to the output, supports environment variables and multiple expressions"
}
