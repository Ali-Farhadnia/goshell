package shell

import (
	"context"
	"fmt"
	"strings"

	"github.com/Ali-Farhadnia/goshell/internal/service/history"
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
	Execute(ctx context.Context, args []string) (string, error)

	// Help returns the help text for the command
	Help() string
}

type Service struct {
	historySVC  *history.Service
	sessionRepo SessionRepository
	commandRepo CommandRepository
}

func NewService(
	historySVC *history.Service,
	sessionRepo SessionRepository,
	commandRepo CommandRepository,
) *Service {
	return &Service{
		historySVC:  historySVC,
		sessionRepo: sessionRepo,
		commandRepo: commandRepo,
	}
}

func (s *Service) RegisterCommand(cmd Command) error {
	return s.commandRepo.Register(cmd)

}

func (s *Service) ExecuteCommand(ctx context.Context, cmdName string, args []string) (string, error) {
	cmd, err := s.GetCommand(cmdName)
	if err != nil {
		return "", err
	}

	if len(args) > cmd.MaxArguments() && cmd.MaxArguments() != -1 {
		return "", fmt.Errorf("too many arguments for command '%s'", cmdName)
	}

	session, err := s.sessionRepo.GetSession()
	if err != nil {
		return "", err
	}

	var userID *int64
	if session.User != nil {
		userID = &session.User.ID
	}

	history := fmt.Sprintf("%s %s", cmdName, strings.Join(args, " "))
	err = s.historySVC.SaveCommandHistory(userID, history)
	if err != nil {
		return "", err
	}

	switch {
	case len(args) > 0 && args[0] == "--help":
		return cmd.Help(), nil
	default:
		result, err := cmd.Execute(ctx, args)
		if err != nil {
			return "", err
		}
		return result, nil
	}
}

func (s *Service) GetCommand(name string) (Command, error) {
	cmd, err := s.commandRepo.Get(name)
	if err != nil {
		return nil, err
	}

	return cmd, nil
}

func (s *Service) Help() (string, error) {
	commands, err := s.commandRepo.List()
	if err != nil {
		return "", err
	}

	var result strings.Builder
	for _, cmd := range commands {
		result.WriteString(fmt.Sprintf("%s\n", cmd.Help()))
	}

	return result.String(), nil
}
