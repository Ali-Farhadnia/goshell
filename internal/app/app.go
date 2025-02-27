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
	shellSVC := shell.NewService(historySVC, sessionRepo, cmdRepo, shell.NewSystemCommand(sessionRepo))

	// register commands

	// exit
	shellSVC.RegisterCommand(commands.NewExitCommand(nil))
	// echo
	shellSVC.RegisterCommand(commands.NewEchoCommand())
	// cat
	shellSVC.RegisterCommand(commands.NewCatCommand(sessionRepo))
	// type
	shellSVC.RegisterCommand(commands.NewTypeCommand(cmdRepo))
	// pwd
	shellSVC.RegisterCommand(commands.NewPWDCommand(sessionRepo))
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
		if input == "" {
			continue
		}

		// Parse input line into arguments (handles quotes and escaping)
		args, err := parseArguments(input)
		if err != nil {
			fmt.Println("error:", err)
			continue
		}

		// Handle redirections and get clean arguments
		inputReader, outputWriter, errorOutputWriter, cleanArgs, cleanup := processRedirections(args)
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

func parseArguments(input string) ([]string, error) {
	var args []string
	var current strings.Builder
	inQuotes := false
	escaped := false
	quoteChar := byte(0)

	for i := 0; i < len(input); i++ {
		c := input[i]

		if escaped {
			// Allow escaping only specific characters inside quotes
			if inQuotes && (c == '$' || c == '`' || c == '"' || c == '\\') {
				current.WriteByte(c)
			} else if !inQuotes {
				current.WriteByte(c) // keep escaped character
			} else {
				current.WriteByte('\\') // Keep backslash
				current.WriteByte(c)
			}
			escaped = false
			continue
		}

		if c == '\\' {
			escaped = true
			continue
		}

		if c == '"' {
			if inQuotes && quoteChar == '"' {
				inQuotes = false // End quote
			} else if !inQuotes {
				inQuotes = true
				quoteChar = '"'
			} else {
				current.WriteByte(c) // Inside different quote type, treat as normal character
			}
			continue
		}

		if (c == ' ' || c == '\t') && !inQuotes {
			// End
			if current.Len() > 0 {
				args = append(args, current.String())
				current.Reset()
			}
			continue
		}

		current.WriteByte(c)
	}

	// Append last argument if exists
	if current.Len() > 0 {
		args = append(args, current.String())
	}

	if inQuotes {
		return nil, fmt.Errorf("unterminated quote detected")
	}

	return args, nil
}

func processRedirections(args []string) (io.Reader, io.Writer, io.Writer, []string, func()) {
	inputReader := os.Stdin
	outputWriter := os.Stdout
	errorOutputWriter := os.Stderr
	var inputFile, outputFile, errorFile *os.File
	cleanArgs := []string{}

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case ">":
			if i+1 < len(args) {
				f, err := os.Create(args[i+1])
				if err != nil {
					fmt.Println("error:", err)
					continue
				}
				outputWriter = f
				outputFile = f
				i++
			}
		case ">>":
			if i+1 < len(args) {
				f, err := os.OpenFile(args[i+1], os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					fmt.Println("error:", err)
					continue
				}
				outputWriter = f
				outputFile = f
				i++
			}
		case "2>":
			if i+1 < len(args) {
				f, err := os.Create(args[i+1])
				if err != nil {
					fmt.Println("error:", err)
					continue
				}
				errorOutputWriter = f
				errorFile = f
				i++
			}
		case "2>>":
			if i+1 < len(args) {
				f, err := os.OpenFile(args[i+1], os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					fmt.Println("error:", err)
					continue
				}
				errorOutputWriter = f
				errorFile = f
				i++
			}
		case "<":
			if i+1 < len(args) {
				f, err := os.Open(args[i+1])
				if err != nil {
					fmt.Println("error:", err)
					continue
				}
				inputReader = f
				inputFile = f
				i++
			}
		default:
			cleanArgs = append(cleanArgs, args[i])
		}
	}

	// Cleanup function to close files
	cleanup := func() {
		if inputFile != nil {
			inputFile.Close()
		}
		if outputFile != nil {
			outputFile.Close()
		}
		if errorFile != nil {
			errorFile.Close()
		}
	}

	return inputReader, outputWriter, errorOutputWriter, cleanArgs, cleanup
}
