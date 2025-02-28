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
	"github.com/Ali-Farhadnia/goshell/pkg/execpath/inputprocessor"
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
	shellSVC := shell.NewService(historySVC, sessionRepo, cmdRepo, shell.NewSystemCommand(sessionRepo, os.Getenv("PATH")), os.Getenv("PATH"))

	// register commands

	// exit
	shellSVC.RegisterCommand(commands.NewExitCommand(nil, nil))
	// echo
	shellSVC.RegisterCommand(commands.NewEchoCommand())
	// cat
	shellSVC.RegisterCommand(commands.NewCatCommand(sessionRepo))
	// type
	shellSVC.RegisterCommand(commands.NewTypeCommand(cmdRepo, os.Getenv("PATH")))
	// pwd
	shellSVC.RegisterCommand(commands.NewPWDCommand(sessionRepo))
	// login
	shellSVC.RegisterCommand(commands.NewLoginCommand(userSVC, sessionRepo))
	// adduser
	shellSVC.RegisterCommand(commands.NewAddUserCommand(userSVC))
	// logout
	shellSVC.RegisterCommand(commands.NewLogoutCommand(sessionRepo))
	// ls
	shellSVC.RegisterCommand(commands.NewLSCommand(sessionRepo))
	// cd
	shellSVC.RegisterCommand(commands.NewCDCommand(sessionRepo))
	// history
	shellSVC.RegisterCommand(commands.NewHistoryCommand(historySVC, sessionRepo))
	// help
	shellSVC.RegisterCommand(commands.NewHelpCommand(cmdRepo))
	// users
	shellSVC.RegisterCommand(commands.NewUsersCommand(userSVC))

	curDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// create guest user
	sessionRepo.SetSession(shell.Session{
		User:       nil,
		WorkingDir: curDir,
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
		if session.User != nil {
			fmt.Printf("%s:$ ", session.User.Username)
		} else {
			fmt.Print("$ ")

		}

		input, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println("\nExiting...")
				return nil
			}
			return fmt.Errorf("error reading input: %w", err)
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		// Parse input line into arguments (handles quotes and escaping)
		args, err := inputprocessor.ParseArguments(input)
		if err != nil {
			fmt.Println("error:", err)
			continue
		}

		// Handle redirections and get clean arguments
		inputReader, outputWriter, errorOutputWriter, cleanArgs, cleanup := inputprocessor.ProcessRedirections(args, session.WorkingDir)
		defer cleanup()

		if len(cleanArgs) == 0 {
			continue
		}

		commandName, commandArgs := cleanArgs[0], cleanArgs[1:]
		err = a.shellSVC.ExecuteCommand(ctx, commandName, commandArgs, inputReader, outputWriter, errorOutputWriter)
		if err != nil {
			fmt.Fprintln(errorOutputWriter, "error:", err)
		}
	}
}
