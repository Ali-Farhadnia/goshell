package commands

import (
	"context"
	"fmt"
	"io"
	"sort"
	"text/tabwriter"

	"github.com/Ali-Farhadnia/goshell/internal/service/shell"
)

type HelpCommand struct {
	cmdRepo shell.CommandRepository
}

func NewHelpCommand(cmdRepo shell.CommandRepository) *HelpCommand {
	return &HelpCommand{cmdRepo: cmdRepo}
}

func (h *HelpCommand) Name() string {
	return "help"
}

func (h *HelpCommand) MaxArguments() int {
	return 0
}

func (h *HelpCommand) Execute(ctx context.Context, args []string, inputReader io.Reader, outputWriter, errorOutputWriter io.Writer) error {
	commands, err := h.cmdRepo.List()
	if err != nil {
		_, err = fmt.Fprintf(errorOutputWriter, "error listing commands: %v\n", err)
		return err
	}

	sort.Slice(commands, func(i, j int) bool {
		return commands[i].Name() < commands[j].Name()
	})

	w := tabwriter.NewWriter(outputWriter, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "Command\tDescription")
	fmt.Fprintln(w, "---\t---")

	for _, cmd := range commands {
		_, err = fmt.Fprintf(w, "%s\t%s\n", cmd.Name(), cmd.Help())
		if err != nil {
			return err
		}
	}

	err = w.Flush()
	if err != nil {
		_, err := fmt.Fprintf(errorOutputWriter, "error flushing tab writer: %v\n", err)
		return err
	}

	return nil
}

func (h *HelpCommand) Help() string {
	return "help - Displays available commands and their usage in a formatted table."
}
