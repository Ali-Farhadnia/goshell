package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Shell struct {
	user    string
	history []string
}

func NewShell() *Shell {
	return &Shell{
		user:    "guest",
		history: []string{},
	}
}

func (s *Shell) run() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("%s$ ", s.user)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}
		s.history = append(s.history, input)
		args := strings.Fields(input)
		s.execute(args)
	}
}

func (s *Shell) execute(args []string) {
	switch args[0] {
	case "exit":
		s.exit(args)
	case "echo":
		s.echo(args[1:])
	case "cat":
		s.cat(args)
	case "type":
		s.cmdType(args)
	case "pwd":
		s.pwd()
	case "cd":
		s.cd(args)
	case "adduser":
		s.addUser(args)
	case "login":
		s.login(args)
	case "logout":
		s.logout()
	case "history":
		s.historyCmd(args)
	default:
		fmt.Println("Unknown command:", args[0])
	}
}

func (s *Shell) exit(args []string) {
	code := 0
	if len(args) > 1 {
		if val, err := strconv.Atoi(args[1]); err == nil {
			code = val
		}
	}
	os.Exit(code)
}

func (s *Shell) echo(args []string) {
	fmt.Println(strings.Join(args, " "))
}

func (s *Shell) cat(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: cat <filename>")
		return
	}
	data, err := os.ReadFile(args[1])
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Print(string(data))
}

func (s *Shell) cmdType(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: type <command>")
		return
	}
	cmd := args[1]
	_, err := exec.LookPath(cmd)
	if err != nil {
		fmt.Printf("%s is a shell built-in\n", cmd)
	} else {
		fmt.Printf("%s is an external command\n", cmd)
	}
}

func (s *Shell) pwd() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(dir)
}

func (s *Shell) cd(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: cd <directory>")
		return
	}
	if err := os.Chdir(args[1]); err != nil {
		fmt.Println("Error:", err)
	}
}

func (s *Shell) addUser(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: adduser <username>")
		return
	}
	s.user = args[1]
	fmt.Println("User added:", s.user)
}

func (s *Shell) login(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: login <username>")
		return
	}
	s.user = args[1]
	fmt.Println("Logged in as:", s.user)
}

func (s *Shell) logout() {
	s.user = "guest"
	fmt.Println("Logged out.")
}

func (s *Shell) historyCmd(args []string) {
	if len(args) > 1 && args[1] == "clean" {
		s.history = []string{}
		fmt.Println("History cleared.")
		return
	}
	for i, cmd := range s.history {
		fmt.Printf("%d: %s\n", i+1, cmd)
	}
}

func main() {
	shell := NewShell()
	shell.run()
}
