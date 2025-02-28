package commands_test

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/Ali-Farhadnia/goshell/internal/service/shell"
	"github.com/Ali-Farhadnia/goshell/internal/service/shell/commands"
	"github.com/Ali-Farhadnia/goshell/internal/service/shell/repository"
	"github.com/Ali-Farhadnia/goshell/internal/service/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLogoutCommand_Execute(t *testing.T) {
	ctx := context.Background()

	cases := []struct {
		name           string
		setupSession   func(repo *repository.SessionRepositoryMock)
		expectedOutput string
		expectedError  string
	}{
		{
			name: "success - logout user",
			setupSession: func(repo *repository.SessionRepositoryMock) {
				repo.On("GetSession").Return(shell.Session{User: &user.User{Username: "testuser"}}, nil).Once()
				repo.On("SetSession", mock.MatchedBy(func(s shell.Session) bool {
					return s.User == nil
				})).Return(nil).Once()
			},
			expectedOutput: "Logged out.\n",
			expectedError:  "",
		},
		{
			name: "failure - session get error",
			setupSession: func(repo *repository.SessionRepositoryMock) {
				repo.On("GetSession").Return(shell.Session{}, errors.New("session get error")).Once()
			},
			expectedOutput: "",
			expectedError:  "session error: session get error\n",
		},
		{
			name: "failure - session save error",
			setupSession: func(repo *repository.SessionRepositoryMock) {
				repo.On("GetSession").Return(shell.Session{User: &user.User{Username: "testuser"}}, nil).Once()
				repo.On("SetSession", mock.Anything).Return(errors.New("session save error")).Once()
			},
			expectedOutput: "",
			expectedError:  "session save error: session save error\n",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockSessionRepo := new(repository.SessionRepositoryMock)
			tc.setupSession(mockSessionRepo)

			var outputBuffer bytes.Buffer
			var errorBuffer bytes.Buffer

			cmd := commands.NewLogoutCommand(mockSessionRepo)
			err := cmd.Execute(ctx, []string{}, nil, &outputBuffer, &errorBuffer)

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedOutput, outputBuffer.String())
			assert.Equal(t, tc.expectedError, errorBuffer.String())
			mockSessionRepo.AssertExpectations(t)
		})
	}
}
