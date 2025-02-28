package commands_test

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/Ali-Farhadnia/goshell/internal/service/shell/commands"
	"github.com/Ali-Farhadnia/goshell/internal/service/shell/repository"
	"github.com/stretchr/testify/assert"
)

func TestTypeCommand_Execute(t *testing.T) {
	ctx := context.Background()

	// Create a temporary executable file for testing
	tempExec, err := os.CreateTemp("", "testexec")
	assert.NoError(t, err)
	defer os.Remove(tempExec.Name())

	// Make the file executable
	err = os.Chmod(tempExec.Name(), 0755)
	assert.NoError(t, err)

	// Add the directory containing the executable to the PATH environment variable
	path := os.Getenv("PATH")
	path = filepath.Dir(tempExec.Name()) + string(os.PathListSeparator) + path

	cases := []struct {
		name           string
		args           []string
		setupRepo      func(repo *repository.CommandRepositoryMock)
		expectedOutput string
		expectedError  string
	}{
		{
			name: "success - shell builtin",
			args: []string{"builtin"},
			setupRepo: func(repo *repository.CommandRepositoryMock) {
				repo.On("Get", "builtin").Return(&MockCommand{}, nil).Once()
			},
			expectedOutput: "builtin is a shell builtin\n",
			expectedError:  "",
		},
		{
			name: "success - external executable",
			args: []string{filepath.Base(tempExec.Name())},
			setupRepo: func(repo *repository.CommandRepositoryMock) {
				repo.On("Get", filepath.Base(tempExec.Name())).Return(&MockCommand{}, errors.New("not found")).Once()
			},
			expectedOutput: filepath.Base(tempExec.Name()) + " is " + tempExec.Name() + "\n",
			expectedError:  "",
		},
		{
			name: "failure - command not found",
			args: []string{"nonexistent"},
			setupRepo: func(repo *repository.CommandRepositoryMock) {
				repo.On("Get", "nonexistent").Return(&MockCommand{}, errors.New("not found")).Once()
			},
			expectedOutput: "",
			expectedError:  "command not found: nonexistent\n",
		},
		{
			name:           "failure - missing argument",
			args:           []string{},
			setupRepo:      func(repo *repository.CommandRepositoryMock) {},
			expectedOutput: "",
			expectedError:  "usage: type <command>\n",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(repository.CommandRepositoryMock)
			tc.setupRepo(mockRepo)

			var outputBuffer bytes.Buffer
			var errorBuffer bytes.Buffer

			cmd := commands.NewTypeCommand(mockRepo, path)
			err := cmd.Execute(ctx, tc.args, nil, &outputBuffer, &errorBuffer)

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedOutput, outputBuffer.String())
			assert.Equal(t, tc.expectedError, errorBuffer.String())
			mockRepo.AssertExpectations(t)
		})
	}
}
