package ls

import (
	"context"
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
func New(sessionRepo shell.SessionRepository) *LSCommand {
	return &LSCommand{
		sessionRepo: sessionRepo,
	}
}

// Name returns the command name
func (l *LSCommand) Name() string {
	return "ls"
}

// Execute runs the command
func (l *LSCommand) Execute(ctx context.Context, args []string) (string, error) {
	session, err := l.sessionRepo.GetSession()
	if err != nil {
		return "", err
	}

	dirPath := session.WorkingDir
	if len(args) > 0 {
		dirPath = filepath.Join(session.WorkingDir, args[0])
	}

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return "", err
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

	return strings.Join(output, "\n"), nil
}

// Help returns the help text
func (l *LSCommand) Help() string {
	return "ls [dir] - Lists the contents of the specified directory (or current directory if none specified)"
}
