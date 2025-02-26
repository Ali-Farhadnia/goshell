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
	"github.com/Ali-Farhadnia/goshell/internal/service/shell/commands/adduser"
	"github.com/Ali-Farhadnia/goshell/internal/service/shell/commands/cd"
	"github.com/Ali-Farhadnia/goshell/internal/service/shell/commands/echo"
	"github.com/Ali-Farhadnia/goshell/internal/service/shell/commands/exit"
	historyCommand "github.com/Ali-Farhadnia/goshell/internal/service/shell/commands/history"
	"github.com/Ali-Farhadnia/goshell/internal/service/shell/commands/ls"
	sessionRepository "github.com/Ali-Farhadnia/goshell/internal/service/shell/repository"
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
	sessionRepo := sessionRepository.NewSessionRepository()
	historyRepo := historyRepository.New(db)
	guestHisotryCache := historyRepository.NewInMemory()

	userSVC := user.New(usrRepo)
	historySVC := history.New(historyRepo, guestHisotryCache, -1)
	shellSVC := shell.NewService(historySVC, sessionRepo)

	// register commands

	// exit
	shellSVC.RegisterCommand(exit.New(nil))
	// echo
	shellSVC.RegisterCommand(echo.New(sessionRepo))
	// adduser
	shellSVC.RegisterCommand(adduser.New(userSVC))
	// ls
	shellSVC.RegisterCommand(ls.New(sessionRepo))
	// cd
	shellSVC.RegisterCommand(cd.New(sessionRepo))
	// history
	shellSVC.RegisterCommand(historyCommand.New(historySVC, sessionRepo))

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

		var result string

		switch commandName {
		case "help":
			result, err = a.shellSVC.Help()
		default:
			result, err = a.shellSVC.ExecuteCommand(ctx, commandName, commandArgs)
		}

		if err != nil {
			fmt.Println("error:", err)
		}

		fmt.Println(result)
	}
}
