package commands

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/Ali-Farhadnia/goshell/internal/service/history"
	"github.com/Ali-Farhadnia/goshell/internal/service/shell"
)

// HistoryCommand implements the history command
type HistoryCommand struct {
	historySVC  *history.Service
	sessionRepo shell.SessionRepository
}

// New creates a new history command
func NewHistoryCommand(
	historySVC *history.Service,
	sessionRepo shell.SessionRepository,
) *HistoryCommand {
	return &HistoryCommand{
		historySVC:  historySVC,
		sessionRepo: sessionRepo,
	}
}

// Name returns the command name
func (c *HistoryCommand) Name() string {
	return "history"
}

// MaxArguments returns the maximum number of arguments allowed for the Command.
func (c *HistoryCommand) MaxArguments() int {
	return 2
}

// Execute runs the command
func (c *HistoryCommand) Execute(ctx context.Context, args []string) (string, error) {
	session, err := c.sessionRepo.GetSession()
	if err != nil {
		return "", err
	}

	var userID *int64
	if session.User != nil {
		userID = &session.User.ID
	}

	// Parse flags
	if len(args) > 0 {
		switch args[0] {
		case "clear":
			err := c.historySVC.ClearCommandHistory(userID)
			if err != nil {
				return "", fmt.Errorf("error clearing history: %w", err)
			}

			return "History cleared.", nil

		case "-n", "--limit":
			if len(args) < 2 {
				return "", fmt.Errorf("usage: history -n <limit>")
			}

			limit, err := strconv.Atoi(args[1])
			if err != nil {
				return "", fmt.Errorf("invalid limit: %s", args[1])
			}

			return c.showHistory(limit)
		}
	}

	return c.showHistory(0)
}

// showHistory returns command history as a string
func (c *HistoryCommand) showHistory(limit int) (string, error) {
	session, err := c.sessionRepo.GetSession()
	if err != nil {
		return "", err
	}

	var userID *int64
	if session.User != nil {
		userID = &session.User.ID
	}

	historyStats, err := c.historySVC.GetCommandHistoryStats(userID, limit)
	if err != nil {
		return fmt.Sprintf("Error retrieving history: %v\n", err), fmt.Errorf("")
	}

	var result strings.Builder
	w := tabwriter.NewWriter(&result, 0, 0, 3, ' ', 0) // Adjust spacing here

	fmt.Fprintln(w, "Command\tCount")
	fmt.Fprintln(w, "---\t---")

	for _, stat := range historyStats {
		fmt.Fprintf(w, "%s\t%d\n", stat.Command, stat.Count)
	}

	w.Flush()

	return result.String(), nil
}

// Help returns the help text
func (c *HistoryCommand) Help() string {
	return "history [clear|-n <limit>] - Show command history, clear history, or limit results"
}
