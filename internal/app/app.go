package app

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/Ali-Farhadnia/goshell/internal/config"
	"github.com/Ali-Farhadnia/goshell/internal/database"
	"github.com/Ali-Farhadnia/goshell/internal/service/history"
	historyRepository "github.com/Ali-Farhadnia/goshell/internal/service/history/repository"
	"github.com/Ali-Farhadnia/goshell/internal/service/shell"
	"github.com/Ali-Farhadnia/goshell/internal/service/shell/commands"
	shellRepository "github.com/Ali-Farhadnia/goshell/internal/service/shell/repository"
	"github.com/Ali-Farhadnia/goshell/internal/service/user"
	userRepository "github.com/Ali-Farhadnia/goshell/internal/service/user/repository"
)

// Shell is the main shell application
type App struct {
	shellSVC    *shell.Service
	sessionRepo shell.SessionRepository
}

// NewShell creates and initializes a new shell
func New(cfg *config.Config) (*App, error) {
	// Initialize database
	db, err := database.New(&cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Initialize repository and service
	usrRepo := userRepository.New(db)
	sessionRepo := shellRepository.NewSessionRepository()
	historyRepo := historyRepository.New(db)
	guestHisotryCache := historyRepository.NewInMemory()
	cmdRepo := shellRepository.NewInMemoryCommandRepository()

	userSVC := user.New(usrRepo)
	historySVC := history.New(historyRepo, guestHisotryCache, -1)
	shellSVC := shell.NewService(historySVC, sessionRepo, cmdRepo)

	// register commands

	// exit
	shellSVC.RegisterCommand(commands.NewExitCommand(nil))
	// echo
	shellSVC.RegisterCommand(commands.NewEchoCommand(sessionRepo))
	// cat
	shellSVC.RegisterCommand(commands.NewCatCommand(sessionRepo))
	// type
	shellSVC.RegisterCommand(commands.NewTypeCommand(cmdRepo))
	// adduser
	shellSVC.RegisterCommand(commands.NewAddUserCommand(userSVC))
	// ls
	shellSVC.RegisterCommand(commands.NewLSCommand(sessionRepo))
	// cd
	shellSVC.RegisterCommand(commands.NewCDCommand(sessionRepo))
	// history
	shellSVC.RegisterCommand(commands.NewHistoryCommand(historySVC, sessionRepo))
	// help
	shellSVC.RegisterCommand(commands.NewHelpCommand(cmdRepo))

	// create guest user
	sessionRepo.SetSession(shell.Session{
		User:       nil,
		WorkingDir: "/",
	})

	return &App{
		shellSVC:    shellSVC,
		sessionRepo: sessionRepo,
	}, nil
}

func (a *App) Run() error {
	reader := bufio.NewReader(os.Stdin)
	ctx := context.Background()

	for {
		session, err := a.sessionRepo.GetSession()
		if err != nil {
			return err
		}
		var userName string
		if session.User != nil {
			userName = session.User.Username
		}

		// todo: check if the sign diffrent for loged in user
		fmt.Printf("%s> ", userName)
		input, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println("\nExiting...")
				return nil
			}
			return fmt.Errorf("error reading input: %w", err)
		}

		input = strings.TrimSpace(input)
		args := strings.Fields(input)
		if len(args) == 0 {
			continue
		}

		commandName, commandArgs := args[0], args[1:]
		result, err := a.shellSVC.ExecuteCommand(ctx, commandName, commandArgs)
		if err != nil {
			fmt.Println("error:", err)
		}

		fmt.Println(result)
	}
}
