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

func TestCDCommand_Execute(t *testing.T) {
	ctx := context.Background()

	// Mock session repository
	mockRepo := new(repository.SessionRepositoryMock)

	curDir, err := os.Getwd()
	assert.NoError(t, err)

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp(curDir, "testdir")
	assert.NoError(t, err)
	defer os.Remove(tempDir)

	// Create a file for testing
	tempFile, err := os.CreateTemp("", "testfile")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

	cases := []struct {
		name           string
		args           []string
		setupRepo      func()
		expectedOutput string
		expectedError  string
	}{
		{
			name: "success - change to valid directory (absolute path)",
			args: []string{tempDir},
			setupRepo: func() {
				mockRepo.On("GetSession").Return(shell.Session{WorkingDir: "/"}, nil).Once()
				mockRepo.On("SetSession", mock.MatchedBy(func(s shell.Session) bool {
					return s.WorkingDir == tempDir
				})).Return(nil).Once()
			},
			expectedOutput: "",
			expectedError:  "",
		},
		{
			name: "success - change to valid directory (relative path)",
			args: []string{filepath.Base(tempDir)},
			setupRepo: func() {
				mockRepo.On("GetSession").Return(shell.Session{WorkingDir: filepath.Dir(tempDir)}, nil).Once()
				mockRepo.On("SetSession", mock.MatchedBy(func(s shell.Session) bool {
					return s.WorkingDir == tempDir
				})).Return(nil).Once()
			},
			expectedOutput: "",
			expectedError:  "",
		},
		{
			name: "failure - directory does not exist",
			args: []string{"/nonexistentdir"},
			setupRepo: func() {
				mockRepo.On("GetSession").Return(shell.Session{WorkingDir: "/tmp"}, nil).Once()
			},
			expectedOutput: "",
			expectedError:  "no such file or directory\n",
		},
		{
			name: "failure - path is a file, not a directory",
			args: []string{tempFile.Name()},
			setupRepo: func() {
				mockRepo.On("GetSession").Return(shell.Session{WorkingDir: "/tmp"}, nil).Once()
			},
			expectedOutput: "",
			expectedError:  "no such file or directory",
		},
		{
			name: "failure - session retrieval error",
			args: []string{"somewhere"},
			setupRepo: func() {
				mockRepo.On("GetSession").Return(shell.Session{}, errors.New("session error")).Once()
			},
			expectedOutput: "",
			expectedError:  "session error",
		},
		{
			name: "failure - session update error",
			args: []string{tempDir},
			setupRepo: func() {
				mockRepo.On("GetSession").Return(shell.Session{WorkingDir: "/"}, nil).Once()
				mockRepo.On("SetSession", mock.Anything).Return(errors.New("session update error")).Once()
			},
			expectedOutput: "",
			expectedError:  "session update error",
		},
		{
			name:           "failure - no arguments",
			args:           []string{},
			setupRepo:      func() {},
			expectedOutput: "",
			expectedError:  "usage: cd <dir>",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo.Mock = mock.Mock{} // Reset mock
			tc.setupRepo()

			var outputBuffer bytes.Buffer
			var errorBuffer bytes.Buffer

			cmd := commands.NewCDCommand(mockRepo)
			err := cmd.Execute(ctx, tc.args, nil, &outputBuffer, &errorBuffer)

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedOutput, outputBuffer.String())
			assert.Contains(t, errorBuffer.String(), tc.expectedError)
			mockRepo.AssertExpectations(t)
		})
	}
}
