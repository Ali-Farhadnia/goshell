package commands_test

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/Ali-Farhadnia/goshell/internal/service/shell/commands"
	"github.com/Ali-Farhadnia/goshell/internal/service/user"
	userRepository "github.com/Ali-Farhadnia/goshell/internal/service/user/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAddUserCommand_Execute(t *testing.T) {
	ctx := context.Background()

	cases := []struct {
		name           string
		args           []string
		setupUserRepo  func(repo *userRepository.UserRepositoryMock)
		expectedOutput string
		expectedError  string
	}{
		{
			name: "success - add user without password",
			args: []string{"newuser"},
			setupUserRepo: func(repo *userRepository.UserRepositoryMock) {
				repo.On("FindUserByUsername", "newuser").Return(nil, user.ErrUserNotFound).Once() // User does not exist
				repo.On("CreateUser", mock.MatchedBy(func(u user.User) bool {
					return u.Username == "newuser" && u.PasswordHash == nil
				})).Return(nil).Once()
			},
			expectedOutput: "User created successfully\n",
			expectedError:  "",
		},
		{
			name: "success - add user with password",
			args: []string{"newuser", "securepassword"},
			setupUserRepo: func(repo *userRepository.UserRepositoryMock) {
				repo.On("FindUserByUsername", "newuser").Return(nil, user.ErrUserNotFound).Once() // User does not exist
				repo.On("CreateUser", mock.MatchedBy(func(u user.User) bool {
					return u.Username == "newuser" && u.PasswordHash != nil
				})).Return(nil).Once()
			},
			expectedOutput: "User created successfully\n",
			expectedError:  "",
		},
		{
			name:           "failure - missing username",
			args:           []string{},
			setupUserRepo:  func(repo *userRepository.UserRepositoryMock) {},
			expectedOutput: "",
			expectedError:  "usage: adduser <username> [password]\n",
		},
		{
			name: "failure - user already exists",
			args: []string{"existinguser"},
			setupUserRepo: func(repo *userRepository.UserRepositoryMock) {
				repo.On("FindUserByUsername", "existinguser").Return(user.User{Username: "existinguser"}, nil).Once()
			},
			expectedOutput: "",
			expectedError:  "error creating user: user already exists: existinguser\n",
		},
		{
			name: "failure - database error on user lookup",
			args: []string{"newuser"},
			setupUserRepo: func(repo *userRepository.UserRepositoryMock) {
				repo.On("FindUserByUsername", "newuser").Return(nil, errors.New("db error")).Once()
			},
			expectedOutput: "",
			expectedError:  "error creating user: db error\n",
		},
		{
			name: "failure - database error on user creation",
			args: []string{"newuser"},
			setupUserRepo: func(repo *userRepository.UserRepositoryMock) {
				repo.On("FindUserByUsername", "newuser").Return(nil, user.ErrUserNotFound).Once()
				repo.On("CreateUser", mock.Anything).Return(errors.New("db error")).Once()
			},
			expectedOutput: "",
			expectedError:  "error creating user: failed to create user in database: db error\n",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockUserRepo := new(userRepository.UserRepositoryMock)
			tc.setupUserRepo(mockUserRepo)

			var outputBuffer bytes.Buffer
			var errorBuffer bytes.Buffer

			userSvc := user.New(mockUserRepo)
			cmd := commands.NewAddUserCommand(userSvc)
			err := cmd.Execute(ctx, tc.args, nil, &outputBuffer, &errorBuffer)

			assert.NoError(t, err)

			assert.Equal(t, tc.expectedOutput, outputBuffer.String())
			assert.Equal(t, tc.expectedError, errorBuffer.String())

			mockUserRepo.AssertExpectations(t)
		})
	}
}
