package commands

import (
	"context"
	"fmt"
	"io"
	"strconv"
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
func (c *HistoryCommand) Execute(ctx context.Context, args []string, inputReader io.Reader, outputWriter, errorOutputWriter io.Writer) error {
	session, err := c.sessionRepo.GetSession()
	if err != nil {
		_, err = fmt.Fprintf(errorOutputWriter, "error getting session: %v\n", err)
		return err
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
				_, err = fmt.Fprintf(errorOutputWriter, "error clearing history: %v\n", err)
				return err
			}

			_, err = fmt.Fprintln(outputWriter, "History cleared.")
			return err

		case "-n", "--limit":
			if len(args) < 2 {
				_, err = fmt.Fprintf(errorOutputWriter, "usage: history -n <limit>\n")
				return err
			}

			limit, err := strconv.Atoi(args[1])
			if err != nil {
				_, err = fmt.Fprintf(errorOutputWriter, "invalid limit: %s\n", args[1])
				return err
			}

			return c.showHistory(limit, outputWriter, errorOutputWriter)
		}
	}

	return c.showHistory(0, outputWriter, errorOutputWriter)
}

// showHistory returns command history as a string
func (c *HistoryCommand) showHistory(limit int, outputWriter, errorOutputWriter io.Writer) error {
	session, err := c.sessionRepo.GetSession()
	if err != nil {
		_, err = fmt.Fprintf(errorOutputWriter, "error getting session: %v\n", err)
		return err
	}

	var userID *int64
	if session.User != nil {
		userID = &session.User.ID
	}

	historyStats, err := c.historySVC.GetCommandHistoryStats(userID, limit)
	if err != nil {
		_, err = fmt.Fprintf(errorOutputWriter, "error retrieving history: %v\n", err)
		return err
	}

	w := tabwriter.NewWriter(outputWriter, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "Command\tCount")
	fmt.Fprintln(w, "---\t---")

	for _, stat := range historyStats {
		_, err = fmt.Fprintf(w, "%s\t%d\n", stat.Command, stat.Count)
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

// Help returns the help text
func (c *HistoryCommand) Help() string {
	return "history [clear|-n <limit>] - Show command history, clear history, or limit results"
}
