package commands_test

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Ali-Farhadnia/goshell/internal/service/shell/commands"
	"github.com/Ali-Farhadnia/goshell/internal/service/user"
	userRepository "github.com/Ali-Farhadnia/goshell/internal/service/user/repository"
	"github.com/stretchr/testify/assert"
)

func TestUsersCommand_Execute(t *testing.T) {
	ctx := context.Background()

	now := time.Now()
	formattedTime := now.Format(time.RFC822)

	cases := []struct {
		name           string
		setupRepo      func(repo *userRepository.UserRepositoryMock)
		expectedOutput string
		expectedError  string
	}{
		{
			name: "success - list users with last login",
			setupRepo: func(repo *userRepository.UserRepositoryMock) {
				repo.On("ListUsers").Return([]user.User{
					{Username: "user1", LastLogin: &now},
					{Username: "user2", LastLogin: nil},
				}, nil).Once()
			},
			expectedOutput: "Registered users:\n----------------\nuser1           Last login: " + formattedTime + "\nuser2           Last login: Never\n",
			expectedError:  "",
		},
		{
			name: "success - list users with no users",
			setupRepo: func(repo *userRepository.UserRepositoryMock) {
				repo.On("ListUsers").Return([]user.User{}, nil).Once()
			},
			expectedOutput: "Registered users:\n----------------\n",
			expectedError:  "",
		},
		{
			name: "failure - repo error",
			setupRepo: func(repo *userRepository.UserRepositoryMock) {
				repo.On("ListUsers").Return(nil, errors.New("repo error")).Once()
			},
			expectedOutput: "",
			expectedError:  "repo error\n",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(userRepository.UserRepositoryMock)
			tc.setupRepo(mockRepo)

			var outputBuffer bytes.Buffer
			var errorBuffer bytes.Buffer

			userSvc := user.New(mockRepo)
			cmd := commands.NewUsersCommand(userSvc)
			err := cmd.Execute(ctx, []string{}, nil, &outputBuffer, &errorBuffer)

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedOutput, outputBuffer.String())
			assert.Equal(t, tc.expectedError, errorBuffer.String())
			mockRepo.AssertExpectations(t)
		})
	}
}
