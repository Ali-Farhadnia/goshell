package commands_test

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/Ali-Farhadnia/goshell/internal/service/shell"
	"github.com/Ali-Farhadnia/goshell/internal/service/shell/commands"
	"github.com/Ali-Farhadnia/goshell/internal/service/shell/repository"
	"github.com/stretchr/testify/assert"
)

func TestLSCommand_Execute(t *testing.T) {
	ctx := context.Background()

	curDir, err := os.Getwd()
	assert.NoError(t, err)

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp(curDir, "testls")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create files and directories within the temporary directory
	os.Mkdir(filepath.Join(tempDir, "dir1"), 0755)
	os.Mkdir(filepath.Join(tempDir, "dir2"), 0755)
	os.WriteFile(filepath.Join(tempDir, "file1.txt"), []byte("file1"), 0644)
	os.WriteFile(filepath.Join(tempDir, "file2.txt"), []byte("file2"), 0644)

	cases := []struct {
		name           string
		args           []string
		setupSession   func(repo *repository.SessionRepositoryMock)
		expectedOutput string
		expectedError  string
	}{
		{
			name: "success - list current directory",
			setupSession: func(repo *repository.SessionRepositoryMock) {
				repo.On("GetSession").Return(shell.Session{WorkingDir: tempDir}, nil).Once()
			},
			expectedOutput: "dir1/\ndir2/\nfile1.txt\nfile2.txt\n",
			expectedError:  "",
		},
		{
			name: "success - list subdirectory",
			args: []string{"dir1"},
			setupSession: func(repo *repository.SessionRepositoryMock) {
				repo.On("GetSession").Return(shell.Session{WorkingDir: tempDir}, nil).Once()
			},
			expectedOutput: "\n", // dir1 is empty
			expectedError:  "",
		},
		{
			name: "failure - session error",
			setupSession: func(repo *repository.SessionRepositoryMock) {
				repo.On("GetSession").Return(shell.Session{}, errors.New("session error")).Once()
			},
			expectedOutput: "",
			expectedError:  "session error: session error\n",
		},
		{
			name: "failure - directory does not exist",
			args: []string{"nonexistent"},
			setupSession: func(repo *repository.SessionRepositoryMock) {
				repo.On("GetSession").Return(shell.Session{WorkingDir: tempDir}, nil).Once()
			},
			expectedOutput: "",
			expectedError:  "dir error: open " + filepath.Join(tempDir, "nonexistent") + ": no such file or directory\n",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockSessionRepo := new(repository.SessionRepositoryMock)
			tc.setupSession(mockSessionRepo)

			var outputBuffer bytes.Buffer
			var errorBuffer bytes.Buffer

			cmd := commands.NewLSCommand(mockSessionRepo)
			err := cmd.Execute(ctx, tc.args, nil, &outputBuffer, &errorBuffer)

			assert.NoError(t, err)

			// Normalize output for consistent testing
			expectedOutput := tc.expectedOutput
			if runtime.GOOS == "windows" {
				expectedOutput = strings.ReplaceAll(expectedOutput, "/", "\\")
			}

			assert.Equal(t, expectedOutput, outputBuffer.String())
			assert.Equal(t, tc.expectedError, errorBuffer.String())
			mockSessionRepo.AssertExpectations(t)
		})
	}
}
