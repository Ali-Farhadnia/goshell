package commands

import (
	"context"
	"fmt"
	"sort"
	"strings"
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

func (h *HelpCommand) Execute(ctx context.Context, args []string) (string, error) {
	commands, err := h.cmdRepo.List()
	if err != nil {
		return "", err
	}

	sort.Slice(commands, func(i, j int) bool {
		return commands[i].Name() < commands[j].Name()
	})

	var result strings.Builder
	w := tabwriter.NewWriter(&result, 0, 0, 3, ' ', 0) // Adjust spacing here
	fmt.Fprintln(w, "Command\tDescription")
	fmt.Fprintln(w, "---\t---")

	for _, cmd := range commands {
		fmt.Fprintf(w, "%s\t%s\n", cmd.Name(), cmd.Help())
	}

	w.Flush() // Important! Flush the writer

	return result.String(), nil
}

func (h *HelpCommand) Help() string {
	return "help - Displays available commands and their usage in a formatted table."
}
