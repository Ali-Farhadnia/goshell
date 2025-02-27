package commands

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Ali-Farhadnia/goshell/internal/service/shell"
)

// LSCommand implements the ls command
type LSCommand struct {
	sessionRepo shell.SessionRepository
}

// New creates a new ls command
func NewLSCommand(sessionRepo shell.SessionRepository) *LSCommand {
	return &LSCommand{
		sessionRepo: sessionRepo,
	}
}

// Name returns the command name
func (l *LSCommand) Name() string {
	return "ls"
}

// MaxArguments returns the maximum number of arguments allowed for the Command.
func (c *LSCommand) MaxArguments() int {
	return 1
}

// Execute runs the command
func (l *LSCommand) Execute(ctx context.Context, args []string, inputReader io.Reader, outputWriter, errorOutputWriter io.Writer) error {
	session, err := l.sessionRepo.GetSession()
	if err != nil {
		_, err := fmt.Fprintf(errorOutputWriter, "session error: %v\n", err)
		return err
	}

	dirPath := session.WorkingDir
	if len(args) > 0 {
		dirPath = filepath.Join(session.WorkingDir, args[0])
	}

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		_, err := fmt.Fprintf(errorOutputWriter, "dir error: %v\n", err)
		return err
	}

	var output []string
	for _, entry := range entries {
		if entry.IsDir() {
			output = append(output, entry.Name()+"/")
		} else {
			output = append(output, entry.Name())
		}
	}

	sort.Strings(output)
	_, err = fmt.Fprintf(outputWriter, "%s\n", strings.Join(output, "\n"))
	if err != nil {
		_, err := fmt.Fprintf(errorOutputWriter, "error writing output: %v\n", err)
		return err
	}

	return nil
}

// Help returns the help text
func (l *LSCommand) Help() string {
	return "ls [dir] - Lists the contents of the specified directory (or current directory if none specified)"
}
