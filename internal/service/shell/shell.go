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

// Command defines the interface for all shell commands
type Command interface {
	// Name returns the command name
	Name() string

	// Execute runs the command with the given arguments
	Execute(ctx context.Context, args []string) (string, error)

	// Help returns the help text for the command
	Help() string
}

type Service struct {
	commands    map[string]Command
	historySVC  *history.Service
	sessionRepo SessionRepository
}

func NewService(
	historySVC *history.Service,
	sessionRepo SessionRepository,
) *Service {
	return &Service{
		commands:    make(map[string]Command),
		historySVC:  historySVC,
		sessionRepo: sessionRepo,
	}
}

func (s *Service) RegisterCommand(cmd Command) {
	s.commands[cmd.Name()] = cmd
}

func (s *Service) ExecuteCommand(ctx context.Context, cmdName string, args []string) (string, error) {
	cmd, err := s.GetCommand(cmdName)
	if err != nil {
		return "", err
	}

	session, err := s.sessionRepo.GetSession()
	if err != nil {
		return "", err
	}

	var userID *int64
	if session.User != nil {
		userID = &session.User.ID
	}

	err = s.historySVC.SaveCommandHistory(userID, fmt.Sprintf("%s %s", cmdName, strings.Join(args, " ")))
	if err != nil {
		return "", err
	}

	// help command
	if len(args) > 0 && args[0] == "--help" {
		return cmd.Help(), nil
	}

	result, err := cmd.Execute(ctx, args)
	if err != nil {
		return "", err
	}

	return result, nil
}

func (s *Service) GetCommand(name string) (Command, error) {
	cmd, ok := s.commands[name]
	if !ok {
		return nil, fmt.Errorf("command '%s' not found", name)
	}

	return cmd, nil
}

func (s *Service) Help() (string, error) {
	var result strings.Builder
	for _, cmd := range s.commands {
		result.WriteString(fmt.Sprintf("%s\n", cmd.Help()))
	}

	return result.String(), nil
}
