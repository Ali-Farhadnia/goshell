package commands_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"

	"github.com/Ali-Farhadnia/goshell/internal/service/shell"
	"github.com/Ali-Farhadnia/goshell/internal/service/shell/commands"
	"github.com/Ali-Farhadnia/goshell/internal/service/shell/repository"
	"github.com/stretchr/testify/assert"
)

type MockCommand struct {
	NameVal string
	HelpVal string
}

func (m *MockCommand) Name() string {
	return m.NameVal
}

func (m *MockCommand) Help() string {
	return m.HelpVal
}

func (m *MockCommand) MaxArguments() int {
	return 0
}

func (m *MockCommand) Execute(ctx context.Context, args []string, inputReader io.Reader, outputWriter, errorOutputWriter io.Writer) error {
	return nil
}

func TestHelpCommand_Execute(t *testing.T) {
	ctx := context.Background()

	cases := []struct {
		name           string
		setupRepo      func(repo *repository.CommandRepositoryMock)
		expectedOutput string
		expectedError  string
	}{
		{
			name: "success - list commands",
			setupRepo: func(repo *repository.CommandRepositoryMock) {
				repo.On("List").Return([]shell.Command{
					&MockCommand{NameVal: "cmd1", HelpVal: "Help for cmd1"},
					&MockCommand{NameVal: "cmd2", HelpVal: "Help for cmd2"},
				}, nil).Once()
			},
			expectedOutput: "Command   Description\n---       ---\ncmd1      Help for cmd1\ncmd2      Help for cmd2\n",
			expectedError:  "",
		},
		{
			name: "success - sorted list",
			setupRepo: func(repo *repository.CommandRepositoryMock) {
				repo.On("List").Return([]shell.Command{
					&MockCommand{NameVal: "cmd2", HelpVal: "Help for cmd2"},
					&MockCommand{NameVal: "cmd1", HelpVal: "Help for cmd1"},
				}, nil).Once()
			},
			expectedOutput: "Command   Description\n---       ---\ncmd1      Help for cmd1\ncmd2      Help for cmd2\n",
			expectedError:  "",
		},
		{
			name: "failure - repo error",
			setupRepo: func(repo *repository.CommandRepositoryMock) {
				repo.On("List").Return([]shell.Command{}, errors.New("repo error")).Once()
			},
			expectedOutput: "",
			expectedError:  "error listing commands: repo error\n",
		},
		{
			name: "success - empty command list",
			setupRepo: func(repo *repository.CommandRepositoryMock) {
				repo.On("List").Return([]shell.Command{}, nil).Once()
			},
			expectedOutput: "Command   Description\n---       ---\n",
			expectedError:  "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(repository.CommandRepositoryMock)
			tc.setupRepo(mockRepo)

			var outputBuffer bytes.Buffer
			var errorBuffer bytes.Buffer

			cmd := commands.NewHelpCommand(mockRepo)
			err := cmd.Execute(ctx, []string{}, nil, &outputBuffer, &errorBuffer)

			assert.NoError(t, err)

			assert.Equal(t, tc.expectedOutput, outputBuffer.String())
			assert.Equal(t, tc.expectedError, errorBuffer.String())
			mockRepo.AssertExpectations(t)
		})
	}
}
