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
	cmd := commands.NewCDCommand(mockRepo)

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

	t.Run("success - change to valid directory (absolute path)", func(t *testing.T) {
		mockRepo.On("GetSession").Return(shell.Session{WorkingDir: "/"}, nil).Once()
		mockRepo.On("SetSession", mock.MatchedBy(func(s shell.Session) bool {
			return s.WorkingDir == tempDir
		})).Return(nil).Once()

		var outputBuffer bytes.Buffer
		var errorBuffer bytes.Buffer

		err := cmd.Execute(ctx, []string{tempDir}, nil, &outputBuffer, &errorBuffer)

		assert.NoError(t, err)
		assert.Empty(t, outputBuffer.String())
		assert.Empty(t, errorBuffer.String())
		mockRepo.AssertExpectations(t)
	})

	t.Run("success - change to valid directory (relative path)", func(t *testing.T) {
		mockRepo.On("GetSession").Return(shell.Session{WorkingDir: filepath.Dir(tempDir)}, nil).Once()
		mockRepo.On("SetSession", mock.MatchedBy(func(s shell.Session) bool {
			return s.WorkingDir == tempDir
		})).Return(nil).Once()

		var outputBuffer bytes.Buffer
		var errorBuffer bytes.Buffer

		err := cmd.Execute(ctx, []string{filepath.Base(tempDir)}, nil, &outputBuffer, &errorBuffer)

		assert.NoError(t, err)
		assert.Empty(t, outputBuffer.String())
		assert.Empty(t, errorBuffer.String())
		mockRepo.AssertExpectations(t)
	})

	t.Run("failure - directory does not exist", func(t *testing.T) {
		mockRepo.On("GetSession").Return(shell.Session{WorkingDir: "/tmp"}, nil).Once()

		var outputBuffer bytes.Buffer
		var errorBuffer bytes.Buffer

		err := cmd.Execute(ctx, []string{"/nonexistentdir"}, nil, &outputBuffer, &errorBuffer)

		assert.NoError(t, err)
		assert.Empty(t, outputBuffer.String())
		assert.Contains(t, errorBuffer.String(), "no such file or directory\n")
		mockRepo.AssertExpectations(t)
	})

	t.Run("failure - path is a file, not a directory", func(t *testing.T) {
		mockRepo.On("GetSession").Return(shell.Session{WorkingDir: "/tmp"}, nil).Once()

		var outputBuffer bytes.Buffer
		var errorBuffer bytes.Buffer

		err := cmd.Execute(ctx, []string{tempFile.Name()}, nil, &outputBuffer, &errorBuffer)

		assert.NoError(t, err)
		assert.Empty(t, outputBuffer.String())
		assert.Contains(t, errorBuffer.String(), "no such file or directory")
		mockRepo.AssertExpectations(t)
	})

	t.Run("failure - session retrieval error", func(t *testing.T) {
		mockRepo.On("GetSession").Return(shell.Session{}, errors.New("session error")).Once()

		var outputBuffer bytes.Buffer
		var errorBuffer bytes.Buffer

		err := cmd.Execute(ctx, []string{"somewhere"}, nil, &outputBuffer, &errorBuffer)

		assert.NoError(t, err)
		assert.Empty(t, outputBuffer.String())
		assert.Contains(t, errorBuffer.String(), "session error")
		mockRepo.AssertExpectations(t)
	})

	t.Run("failure - session update error", func(t *testing.T) {
		mockRepo.On("GetSession").Return(shell.Session{WorkingDir: "/"}, nil).Once()
		mockRepo.On("SetSession", mock.Anything).Return(errors.New("session update error")).Once()

		var outputBuffer bytes.Buffer
		var errorBuffer bytes.Buffer

		err := cmd.Execute(ctx, []string{tempDir}, nil, &outputBuffer, &errorBuffer)

		assert.NoError(t, err)
		assert.Empty(t, outputBuffer.String())
		assert.Contains(t, errorBuffer.String(), "session update error")
		mockRepo.AssertExpectations(t)
	})

	t.Run("failure - no arguments", func(t *testing.T) {
		var outputBuffer bytes.Buffer
		var errorBuffer bytes.Buffer

		err := cmd.Execute(ctx, []string{}, nil, &outputBuffer, &errorBuffer)

		assert.NoError(t, err)
		assert.Empty(t, outputBuffer.String())
		assert.Contains(t, errorBuffer.String(), "usage: cd <dir>")
	})
}
