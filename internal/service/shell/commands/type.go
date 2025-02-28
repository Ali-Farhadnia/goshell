package commands

import (
	"context"
	"fmt"
	"io"

	"github.com/Ali-Farhadnia/goshell/internal/service/shell"
	"github.com/Ali-Farhadnia/goshell/pkg/execpath"
)

type TypeCommand struct {
	cmdRepo shell.CommandRepository
	path    string
}

func NewTypeCommand(
	cmdRepo shell.CommandRepository,
	path string,
) *TypeCommand {
	return &TypeCommand{
		cmdRepo: cmdRepo,
		path:    path,
	}
}

func (t *TypeCommand) Name() string {
	return "type"
}

func (t *TypeCommand) MaxArguments() int {
	return 1
}

// Execute runs the command
func (t *TypeCommand) Execute(ctx context.Context, args []string, inputReader io.Reader, outputWriter, errorOutputWriter io.Writer) error {
	if len(args) == 0 {
		_, err := fmt.Fprintf(errorOutputWriter, "usage: type <command>\n")
		return err
	}

	cmdName := args[0]

	// Check if it's a shell builtin
	if _, err := t.cmdRepo.Get(cmdName); err == nil {
		_, err = fmt.Fprintf(outputWriter, "%s is a shell builtin\n", cmdName)
		return err
	}

	// Check if it's an executable in $PATH
	cmdPath, err := execpath.FindExecutable(cmdName, t.path)
	if err != nil {
		_, err = fmt.Fprintf(errorOutputWriter, "%v\n", err)
		return err
	}

	_, err = fmt.Fprintf(outputWriter, "%s is %s\n", cmdName, cmdPath)
	if err != nil {
		_, err = fmt.Fprintf(errorOutputWriter, "error writing output: %v\n", err)
		return err
	}

	return nil
}

func (t *TypeCommand) Help() string {
	return "type <command> - Identifies if the command is a shell builtin or an external executable"
}
