package commands_test

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/Ali-Farhadnia/goshell/internal/service/shell"
	"github.com/Ali-Farhadnia/goshell/internal/service/shell/commands"
	"github.com/Ali-Farhadnia/goshell/internal/service/shell/repository"
	"github.com/stretchr/testify/assert"
)

func TestPWDCommand_Execute(t *testing.T) {
	ctx := context.Background()

	cases := []struct {
		name           string
		args           []string
		setupSession   func(repo *repository.SessionRepositoryMock)
		expectedOutput string
		expectedError  string
	}{
		{
			name: "success - print working directory",
			setupSession: func(repo *repository.SessionRepositoryMock) {
				repo.On("GetSession").Return(shell.Session{WorkingDir: "/test/dir"}, nil).Once()
			},
			expectedOutput: "/test/dir\n",
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
			name: "failure - too many arguments",
			args: []string{"extra"},
			setupSession: func(repo *repository.SessionRepositoryMock) {
				repo.On("GetSession").Return(shell.Session{WorkingDir: "/test/dir"}, nil).Once()
			},
			expectedOutput: "",
			expectedError:  "pwd: too many arguments\n",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockSessionRepo := new(repository.SessionRepositoryMock)
			tc.setupSession(mockSessionRepo)

			var outputBuffer bytes.Buffer
			var errorBuffer bytes.Buffer

			cmd := commands.NewPWDCommand(mockSessionRepo)
			err := cmd.Execute(ctx, tc.args, nil, &outputBuffer, &errorBuffer)

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedOutput, outputBuffer.String())
			assert.Equal(t, tc.expectedError, errorBuffer.String())
			mockSessionRepo.AssertExpectations(t)
		})
	}
}
