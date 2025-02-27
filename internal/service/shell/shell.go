package shell

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/Ali-Farhadnia/goshell/internal/service/history"
	"github.com/Ali-Farhadnia/goshell/pkg/execpath"
)

var (
	ErrCommandNotFound = errors.New("command not found")
)

type SessionRepository interface {
	GetSession() (Session, error)
	SetSession(s Session) error
}

type CommandRepository interface {
	Register(cmd Command) error
	Get(name string) (Command, error)
	List() ([]Command, error)
}

// Command defines the interface for all shell commands
type Command interface {
	// Name returns the command name
	Name() string

	// MaxArguments returns the maximum number of arguments allowed for the Command.
	MaxArguments() int

	// Execute runs the command with the given arguments
	Execute(ctx context.Context, args []string, inputReader io.Reader, outputWriter, errorOutputWriter io.Writer) error

	// Help returns the help text for the command
	Help() string
}

type Service struct {
	historySVC    *history.Service
	sessionRepo   SessionRepository
	commandRepo   CommandRepository
	systemCommand *SystemCommand
}

func NewService(
	historySVC *history.Service,
	sessionRepo SessionRepository,
	commandRepo CommandRepository,
	systemCommand *SystemCommand,
) *Service {
	return &Service{
		historySVC:    historySVC,
		sessionRepo:   sessionRepo,
		commandRepo:   commandRepo,
		systemCommand: systemCommand,
	}
}

func (s *Service) RegisterCommand(cmd Command) error {
	return s.commandRepo.Register(cmd)
}

// ExecuteCommand determines if a command is built-in or system-based and executes it.
func (s *Service) ExecuteCommand(ctx context.Context, cmdName string, args []string, inputReader io.Reader, outputWriter, errorOutputWriter io.Writer) error {
	// Check built-in command
	cmd, err := s.commandRepo.Get(cmdName)
	if err != nil {
		// If error is not "command not found", return immediately
		if !errors.Is(err, ErrCommandNotFound) {
			return err
		}

		// Check if it's a system command
		if _, err := execpath.FindExecutable(cmdName); err != nil {
			return err
		}

		return s.executeSystemCommand(ctx, cmdName, args, inputReader, outputWriter, errorOutputWriter)
	}

	return s.executeBuiltinCommand(ctx, cmd, args, inputReader, outputWriter, errorOutputWriter)
}

// executeBuiltinCommand runs a built-in command after performing necessary checks.
func (s *Service) executeBuiltinCommand(ctx context.Context, cmd Command, args []string, inputReader io.Reader, outputWriter, errorOutputWriter io.Writer) error {
	// Handle help flag
	if isHelpRequested(args) {
		_, err := fmt.Fprintf(outputWriter, "%s", cmd.Help())
		if err != nil {
			_, err = fmt.Fprintf(errorOutputWriter, "error writing output: %v\n", err)
			return err
		}

		return nil
	}

	// Validate argument count
	if err := validateArgs(cmd, args); err != nil {
		return err
	}

	// Get user ID
	userID, err := s.getUserID()
	if err != nil {
		return err
	}

	// Save command to history
	if err := s.saveCommandHistory(userID, cmd.Name(), args); err != nil {
		return err
	}

	return cmd.Execute(ctx, args, inputReader, outputWriter, errorOutputWriter)
}

// executeSystemCommand runs a system command.
func (s *Service) executeSystemCommand(ctx context.Context, cmdName string, args []string, inputReader io.Reader, outputWriter, errorOutputWriter io.Writer) error {
	// Handle help flag
	if isHelpRequested(args) {
		_, err := fmt.Fprintf(outputWriter, "%s", s.systemCommand.Help())
		if err != nil {
			_, err = fmt.Fprintf(errorOutputWriter, "error writing output: %v\n", err)
			return err
		}

		return nil
	}

	// Get user ID
	userID, err := s.getUserID()
	if err != nil {
		return err
	}

	// Save command to history
	if err := s.saveCommandHistory(userID, cmdName, args); err != nil {
		return err
	}

	return s.systemCommand.Execute(ctx, cmdName, args, inputReader, outputWriter, errorOutputWriter)
}

// Helper function to check if help is requested
func isHelpRequested(args []string) bool {
	return len(args) > 0 && args[0] == "--help"
}

// Helper function to validate argument count
func validateArgs(cmd Command, args []string) error {
	maxArgs := cmd.MaxArguments()
	if maxArgs != -1 && len(args) > maxArgs {
		return fmt.Errorf("too many arguments for command '%s'", cmd.Name())
	}
	return nil
}

// Retrieves the user ID from the session
func (s *Service) getUserID() (*int64, error) {
	session, err := s.sessionRepo.GetSession()
	if err != nil {
		return nil, err
	}

	if session.User != nil {
		return &session.User.ID, nil
	}
	return nil, nil
}

// Saves the command execution to history
func (s *Service) saveCommandHistory(userID *int64, cmdName string, args []string) error {
	history := fmt.Sprintf("%s %s", cmdName, strings.Join(args, " "))
	return s.historySVC.SaveCommandHistory(userID, history)
}
