package commands_test

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/Ali-Farhadnia/goshell/internal/service/shell"
	"github.com/Ali-Farhadnia/goshell/internal/service/shell/commands"
	"github.com/Ali-Farhadnia/goshell/internal/service/shell/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCatCommand_Execute(t *testing.T) {
	ctx := context.Background()

	// Create a mock session repository
	mockRepo := new(repository.SessionRepositoryMock)

	// Create a temporary file for testing
	tempFile, err := os.CreateTemp("", "testfile")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

	content := "Hello, world!"
	_, err = tempFile.WriteString(content)
	assert.NoError(t, err)
	tempFile.Close()

	cases := []struct {
		name           string
		args           []string
		setupRepo      func()
		expectedOutput string
		expectedError  string
	}{
		{
			name: "success - absolute path",
			args: []string{tempFile.Name()},
			setupRepo: func() {
				mockRepo.On("GetSession").Return(shell.Session{WorkingDir: "/tmp"}, nil).Once()
			},
			expectedOutput: content,
			expectedError:  "",
		},
		{
			name: "success - relative path",
			args: []string{filepath.Base(tempFile.Name())},
			setupRepo: func() {
				mockRepo.On("GetSession").Return(shell.Session{WorkingDir: filepath.Dir(tempFile.Name())}, nil).Once()
			},
			expectedOutput: content,
			expectedError:  "",
		},
		{
			name: "failure - file not found",
			args: []string{"/nonexistent/file.txt"},
			setupRepo: func() {
				mockRepo.On("GetSession").Return(shell.Session{WorkingDir: "/tmp"}, nil).Once()
			},
			expectedOutput: "",
			expectedError:  "error reading file",
		},
		{
			name: "failure - session error",
			args: []string{"somefile.txt"},
			setupRepo: func() {
				mockRepo.On("GetSession").Return(shell.Session{}, errors.New("session error")).Once()
			},
			expectedOutput: "",
			expectedError:  "error getting session",
		},
		{
			name:           "failure - no arguments",
			args:           []string{},
			setupRepo:      func() {},
			expectedOutput: "",
			expectedError:  "usage: cat <filename>\n",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo.Mock = mock.Mock{} // Reset mock
			tc.setupRepo()

			var outputBuffer bytes.Buffer
			var errorBuffer bytes.Buffer

			cmd := commands.NewCatCommand(mockRepo)
			err := cmd.Execute(ctx, tc.args, nil, &outputBuffer, &errorBuffer)

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedOutput, outputBuffer.String())
			assert.Contains(t, errorBuffer.String(), tc.expectedError)
			mockRepo.AssertExpectations(t)
		})
	}
}
