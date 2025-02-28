# GoShell

A feature-rich, POSIX-compliant shell implementation written in Go.

[![Go Report Card](https://goreportcard.com/badge/github.com/yourusername/goshell)](https://goreportcard.com/report/github.com/yourusername/goshell)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Features

GoShell is a powerful shell implementation with the following features:

- **Basic Command Support**: Built-in commands including `exit`, `echo`, `cat`, `type`, `pwd`, and `cd`.
- **System Command Execution**: Run any system executable.
- **User Management**: User registration, login and logout functionality.
- **Command History**: Persistent command history tracking for registered users.
- **I/O Redirection**: Support for input and output redirection (`>`, `>>`, `<`, `2>`, `2>>`).
- **Environment Variables**: Support for environment variable expansion.
- **Quote Handling**: Support for double quotes with proper escape character handling.
- **Error Handling**: Informative error messages for invalid commands and syntax.
- **Database Integration**: Persistence of user data and command history.

## Installation

### Prerequisites

- Go 1.20 or higher
- Docker (for running PostgreSQL)

### Setup

1. Clone the repository:

```bash
git clone https://github.com/Ali-Farhadnia/goshell.git
cd goshell
```

2. Create configuration file:

```bash
cp config.sample.yaml config.yaml
```

3. Edit `config.yaml` to match your desired configuration.

4. Start the PostgreSQL database:

```bash
make db-up
```

5. Run the shell:

```bash
make run
```

## Usage

### Basic Commands

```bash
# Print a message
$ echo Hello World
Hello World

# Display current directory
$ pwd
/home/user/goshell

# Change directory
$ cd /another/path

# View file contents
$ cat filename.txt

# Check command type
$ type echo
echo is a shell builtin

# Exit the shell
$ exit
```

### User Management

```bash
# Create a new user
$ adduser username password
user created successfully

# Login
$ login username password
username:$

# Logout
$ logout
```

### Command History

```bash
# View command history
$ history

# Clear command history
$ history clean
```

### I/O Redirection

```bash
# Redirect output to a file (overwrite)
$ echo Hello > output.txt

# Redirect output to a file (append)
$ echo World >> output.txt

# Redirect error output
$ unknown-command 2> error.txt
```

## Project Structure

The project is organized as follows:

```
.
├── cmd
│   └── shell
│       └── main.go
├── config.yaml
├── go.mod
├── go.sum
├── internal
│   ├── app
│   │   └── app.go
│   ├── config
│   │   └── config.go
│   ├── database
│   │   └── database.go
│   └── service
│       ├── history
│       │   ├── history.go
│       │   ├── history_test.go
│       │   ├── model.go
│       │   └── repository
│       │       ├── history_repository.go
│       │       ├── history_repository_mock.go
│       │       └── in_memory_history_repository.go
│       ├── shell
│       │   ├── commands
│       │   │   ├── adduser.go
│       │   │   ├── adduser_test.go
│       │   │   ├── cat.go
│       │   │   ├── cat_test.go
│       │   │   ├── cd.go
│       │   │   ├── cd_test.go
│       │   │   ├── echo.go
│       │   │   ├── echo_test.go
│       │   │   ├── exit.go
│       │   │   ├── exit_test.go
│       │   │   ├── help.go
│       │   │   ├── help_test.go
│       │   │   ├── history.go
│       │   │   ├── history_test.go
│       │   │   ├── login.go
│       │   │   ├── login_test.go
│       │   │   ├── logout.go
│       │   │   ├── logout_test.go
│       │   │   ├── ls.go
│       │   │   ├── ls_test.go
│       │   │   ├── model_test.go
│       │   │   ├── pwd.go
│       │   │   ├── pwd_test.go
│       │   │   ├── type.go
│       │   │   ├── type_test.go
│       │   │   ├── users.go
│       │   │   └── users_test.go
│       │   ├── model.go
│       │   ├── repository
│       │   │   ├── command_repo.go
│       │   │   ├── command_repo_mock.go
│       │   │   ├── session_repo.go
│       │   │   └── session_repo_mock.go
│       │   ├── shell.go
│       │   └── system_command.go
│       └── user
│           ├── model.go
│           ├── repository
│           │   ├── user_repository.go
│           │   └── user_repository_mock.go
│           └── user.go
├── makefile
├── pkg
│   ├── execpath
│   │   └── execpath.go
│   └── inputprocessor
│       ├── inputprocessor.go
│       └── inputprocessor_test.go
└── README.md
```

### Description of Key Components

- **`cmd/shell/main.go`**: The entry point of the application. This is where the shell application is initialized and executed.

- **`config.yaml`**: Configuration file for the application, containing settings and parameters.

- **`go.mod` & `go.sum`**: Go module files that manage dependencies for the project.

- **`internal/app/app.go`**: The main application logic, where the application is configured and started.

- **`internal/config/config.go`**: Handles the loading and parsing of the configuration file (`config.yaml`).

- **`internal/database/database.go`**: Contains database connection logic and database-related utilities.

- **`internal/service/history/`**: Manages the history of commands executed in the shell. Includes models, repositories, and business logic.

- **`internal/service/shell/`**: Contains the core shell functionality, including command definitions, repositories, and system commands.

  - **`commands/`**: Defines individual shell commands (e.g., `cd`, `ls`, `echo`) and their corresponding tests.
  
  - **`repository/`**: Handles data storage for commands and sessions, including mock implementations for testing.

- **`internal/service/user/`**: Manages user-related functionality, including user models and repositories.

- **`makefile`**: Contains build and automation commands for the project.

- **`pkg/execpath/execpath.go`**: Provides utilities for working with executable paths.

- **`pkg/inputprocessor/inputprocessor.go`**: Processes user input and prepares it for execution by the shell.

- **`README.md`**: This file, providing an overview of the project and its structure.

## Development

### Testing

Run the test suite:

```bash
make test
```

Run tests with coverage report:

```bash
make test-coverage
make coverage  # Open HTML coverage report
```

### Database Management

```bash
# Start PostgreSQL container
make db-up

# Log into PostgreSQL
make db-login

# Stop and remove PostgreSQL container
make clean-db
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
