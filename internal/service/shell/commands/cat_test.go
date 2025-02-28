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
)

func TestCatCommand_Execute(t *testing.T) {
	ctx := context.Background()

	// Create a mock session repository
	mockRepo := new(repository.SessionRepositoryMock)
	cmd := commands.NewCatCommand(mockRepo)

	// Create a temporary file for testing
	tempFile, err := os.CreateTemp("", "testfile")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

	content := "Hello, world!"
	_, err = tempFile.WriteString(content)
	assert.NoError(t, err)
	tempFile.Close()

	t.Run("success - absolute path", func(t *testing.T) {
		mockRepo.On("GetSession").Return(shell.Session{WorkingDir: "/tmp"}, nil).Once()
		var outputBuffer bytes.Buffer
		var errorBuffer bytes.Buffer

		err := cmd.Execute(ctx, []string{tempFile.Name()}, nil, &outputBuffer, &errorBuffer)

		assert.NoError(t, err)
		assert.Equal(t, content, outputBuffer.String())
		assert.Empty(t, errorBuffer.String())
		mockRepo.AssertExpectations(t)
	})

	t.Run("success - relative path", func(t *testing.T) {
		mockRepo.On("GetSession").Return(shell.Session{WorkingDir: filepath.Dir(tempFile.Name())}, nil).Once()
		var outputBuffer bytes.Buffer
		var errorBuffer bytes.Buffer

		err := cmd.Execute(ctx, []string{filepath.Base(tempFile.Name())}, nil, &outputBuffer, &errorBuffer)

		assert.NoError(t, err)
		assert.Equal(t, content, outputBuffer.String())
		assert.Empty(t, errorBuffer.String())
		mockRepo.AssertExpectations(t)
	})

	t.Run("failure - file not found", func(t *testing.T) {
		mockRepo.On("GetSession").Return(shell.Session{WorkingDir: "/tmp"}, nil).Once()
		var outputBuffer bytes.Buffer
		var errorBuffer bytes.Buffer

		err := cmd.Execute(ctx, []string{"/nonexistent/file.txt"}, nil, &outputBuffer, &errorBuffer)

		assert.NoError(t, err)
		assert.Empty(t, outputBuffer.String())
		assert.Contains(t, errorBuffer.String(), "error reading file")
		mockRepo.AssertExpectations(t)
	})

	t.Run("failure - session error", func(t *testing.T) {
		mockRepo.On("GetSession").Return(shell.Session{}, errors.New("session error")).Once()
		var outputBuffer bytes.Buffer
		var errorBuffer bytes.Buffer

		err := cmd.Execute(ctx, []string{"somefile.txt"}, nil, &outputBuffer, &errorBuffer)

		assert.NoError(t, err)
		assert.Empty(t, outputBuffer.String())
		assert.Contains(t, errorBuffer.String(), "error getting session")
		mockRepo.AssertExpectations(t)
	})

	t.Run("failure - no arguments", func(t *testing.T) {
		var outputBuffer bytes.Buffer
		var errorBuffer bytes.Buffer

		err := cmd.Execute(ctx, []string{}, nil, &outputBuffer, &errorBuffer)

		assert.NoError(t, err)
		assert.Empty(t, outputBuffer.String())
		assert.Equal(t, "usage: cat <filename>\n", errorBuffer.String())
	})
}
