package commands

import (
	"context"
	"fmt"

	"github.com/Ali-Farhadnia/goshell/internal/service/shell"
	"github.com/Ali-Farhadnia/goshell/pkg/execpath"
)

type TypeCommand struct {
	cmdRepo shell.CommandRepository
}

func NewTypeCommand(cmdRepo shell.CommandRepository) *TypeCommand {
	return &TypeCommand{cmdRepo: cmdRepo}
}

func (t *TypeCommand) Name() string {
	return "type"
}

func (t *TypeCommand) MaxArguments() int {
	return 1
}

func (t *TypeCommand) Execute(ctx context.Context, args []string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("usage: type <command>")
	}

	cmdName := args[0]

	// Check if it's a shell builtin
	if _, err := t.cmdRepo.Get(cmdName); err == nil {
		return fmt.Sprintf("%s is a shell builtin", cmdName), nil
	}

	// Check if it's an executable in $PATH
	cmdPath, err := execpath.FindExecutable(cmdName)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s is %s", cmdName, cmdPath), nil
}

func (t *TypeCommand) Help() string {
	return "type <command> - Identifies if the command is a shell builtin or an external executable"
}
