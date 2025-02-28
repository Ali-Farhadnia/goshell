package commands_test

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/Ali-Farhadnia/goshell/internal/service/shell"
	"github.com/Ali-Farhadnia/goshell/internal/service/shell/commands"
	shellRepository "github.com/Ali-Farhadnia/goshell/internal/service/shell/repository"
	"github.com/Ali-Farhadnia/goshell/internal/service/user"
	userRepository "github.com/Ali-Farhadnia/goshell/internal/service/user/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestLoginCommand_Execute(t *testing.T) {
	ctx := context.Background()

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
	hashedPasswordStr := string(hashedPassword)

	cases := []struct {
		name           string
		args           []string
		setupUserRepo  func(repo *userRepository.UserRepositoryMock)
		setupSession   func(repo *shellRepository.SessionRepositoryMock)
		expectedOutput string
		expectedError  string
	}{
		{
			name: "success - login user without password",
			args: []string{"testuser"},
			setupUserRepo: func(repo *userRepository.UserRepositoryMock) {
				repo.On("FindUserByUsername", "testuser").Return(user.User{Username: "testuser", ID: 1}, nil).Once()
				repo.On("UpdateLastLogin", int64(1)).Return(nil).Once()
			},
			setupSession: func(repo *shellRepository.SessionRepositoryMock) {
				repo.On("GetSession").Return(shell.Session{}, nil).Once()
				repo.On("SetSession", mock.MatchedBy(func(s shell.Session) bool {
					return s.User != nil && s.User.Username == "testuser"
				})).Return(nil).Once()
			},
			expectedOutput: "Logged in as: testuser\n",
			expectedError:  "",
		},
		{
			name: "success - login user with correct password",
			args: []string{"testuser", "correctpassword"},
			setupUserRepo: func(repo *userRepository.UserRepositoryMock) {
				repo.On("FindUserByUsername", "testuser").Return(user.User{
					Username:     "testuser",
					ID:           1,
					PasswordHash: &hashedPasswordStr,
				}, nil).Once()
				repo.On("UpdateLastLogin", int64(1)).Return(nil).Once()
			},
			setupSession: func(repo *shellRepository.SessionRepositoryMock) {
				repo.On("GetSession").Return(shell.Session{}, nil).Once()
				repo.On("SetSession", mock.Anything).Return(nil).Once()
			},
			expectedOutput: "Logged in as: testuser\n",
			expectedError:  "",
		},
		{
			name: "failure - incorrect password",
			args: []string{"testuser", "wrongpassword"},
			setupUserRepo: func(repo *userRepository.UserRepositoryMock) {
				repo.On("FindUserByUsername", "testuser").Return(user.User{
					Username:     "testuser",
					ID:           1,
					PasswordHash: &hashedPasswordStr,
				}, nil).Once()
			},
			setupSession:   func(repo *shellRepository.SessionRepositoryMock) {},
			expectedOutput: "",
			expectedError:  "login failed: invalid password\n",
		},
		{
			name:           "failure - missing username",
			args:           []string{},
			setupUserRepo:  func(repo *userRepository.UserRepositoryMock) {},
			setupSession:   func(repo *shellRepository.SessionRepositoryMock) {},
			expectedOutput: "",
			expectedError:  "usage: login <username> [password]\n",
		},
		{
			name: "failure - user not found",
			args: []string{"testuser"},
			setupUserRepo: func(repo *userRepository.UserRepositoryMock) {
				repo.On("FindUserByUsername", "testuser").Return(nil, user.ErrUserNotFound).Once()
			},
			setupSession:   func(repo *shellRepository.SessionRepositoryMock) {},
			expectedOutput: "",
			expectedError:  "login failed: error on finding user: user not found\n",
		},
		{
			name: "failure - session get error",
			args: []string{"testuser"},
			setupUserRepo: func(repo *userRepository.UserRepositoryMock) {
				repo.On("FindUserByUsername", "testuser").Return(user.User{Username: "testuser", ID: 1}, nil).Once()
				repo.On("UpdateLastLogin", int64(1)).Return(nil).Once()
			},
			setupSession: func(repo *shellRepository.SessionRepositoryMock) {
				repo.On("GetSession").Return(shell.Session{}, errors.New("session get error")).Once()
			},
			expectedOutput: "",
			expectedError:  "session error: session get error\n",
		},
		{
			name: "failure - session save error",
			args: []string{"testuser"},
			setupUserRepo: func(repo *userRepository.UserRepositoryMock) {
				repo.On("FindUserByUsername", "testuser").Return(user.User{Username: "testuser", ID: 1}, nil).Once()
				repo.On("UpdateLastLogin", int64(1)).Return(nil).Once()
			},
			setupSession: func(repo *shellRepository.SessionRepositoryMock) {
				repo.On("GetSession").Return(shell.Session{}, nil).Once()
				repo.On("SetSession", mock.Anything).Return(errors.New("session save error")).Once()
			},
			expectedOutput: "",
			expectedError:  "session save error: session save error\n",
		},
		{
			name: "failure - update last login error",
			args: []string{"testuser"},
			setupUserRepo: func(repo *userRepository.UserRepositoryMock) {
				repo.On("FindUserByUsername", "testuser").Return(user.User{Username: "testuser", ID: 1}, nil).Once()
				repo.On("UpdateLastLogin", int64(1)).Return(errors.New("last login error")).Once()
			},
			setupSession:   func(repo *shellRepository.SessionRepositoryMock) {},
			expectedOutput: "",
			expectedError:  "login failed: failed to update last login time: last login error\n",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockUserRepo := new(userRepository.UserRepositoryMock)
			mockSessionRepo := new(shellRepository.SessionRepositoryMock)
			tc.setupUserRepo(mockUserRepo)
			tc.setupSession(mockSessionRepo)

			var outputBuffer bytes.Buffer
			var errorBuffer bytes.Buffer

			userSvc := user.New(mockUserRepo)
			cmd := commands.NewLoginCommand(userSvc, mockSessionRepo)
			err := cmd.Execute(ctx, tc.args, nil, &outputBuffer, &errorBuffer)

			assert.NoError(t, err)

			assert.Equal(t, tc.expectedOutput, outputBuffer.String())
			assert.Equal(t, tc.expectedError, errorBuffer.String())

			mockUserRepo.AssertExpectations(t)
			mockSessionRepo.AssertExpectations(t)
		})
	}
}
