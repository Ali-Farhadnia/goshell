package commands_test

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"text/tabwriter"

	"github.com/Ali-Farhadnia/goshell/internal/service/history"
	historyRepository "github.com/Ali-Farhadnia/goshell/internal/service/history/repository"
	"github.com/Ali-Farhadnia/goshell/internal/service/shell"
	"github.com/Ali-Farhadnia/goshell/internal/service/shell/commands"
	shellRepository "github.com/Ali-Farhadnia/goshell/internal/service/shell/repository"
	"github.com/Ali-Farhadnia/goshell/internal/service/user"
	"github.com/stretchr/testify/assert"
)

func TestHistoryCommand_Execute(t *testing.T) {
	ctx := context.Background()
	guestID := int64(999)

	cases := []struct {
		name           string
		args           []string
		setupSession   func(repo *shellRepository.SessionRepositoryMock)
		setupHistory   func(repo, guestRepo *historyRepository.HistoryRepositoryMock)
		expectedOutput string
		expectedError  string
	}{
		{
			name: "success - show history",
			setupSession: func(repo *shellRepository.SessionRepositoryMock) {
				repo.On("GetSession").Return(shell.Session{User: &user.User{ID: 123}}, nil).Once()
			},
			setupHistory: func(repo, guestRepo *historyRepository.HistoryRepositoryMock) {
				repo.On("GetUserCommandStats", int64(123)).Return([]history.CommandStats{
					{Command: "ls", Count: 5},
					{Command: "cd", Count: 3},
				}, nil).Once()
			},
			expectedOutput: "Command   Count\n---       ---\nls        5\ncd        3\n",
			expectedError:  "",
		},
		{
			name: "success - show history with limit",
			args: []string{"-n", "1"},
			setupSession: func(repo *shellRepository.SessionRepositoryMock) {
				repo.On("GetSession").Return(shell.Session{User: &user.User{ID: 123}}, nil).Once()
			},
			setupHistory: func(repo, guestRepo *historyRepository.HistoryRepositoryMock) {
				repo.On("GetUserCommandStats", int64(123)).Return([]history.CommandStats{
					{Command: "ls", Count: 5},
				}, nil).Once()
			},
			expectedOutput: "Command   Count\n---       ---\nls        5\n",
			expectedError:  "",
		},
		{
			name: "success - clear history",
			args: []string{"clear"},
			setupSession: func(repo *shellRepository.SessionRepositoryMock) {
				repo.On("GetSession").Return(shell.Session{User: &user.User{ID: 123}}, nil).Once()
			},
			setupHistory: func(repo, guestRepo *historyRepository.HistoryRepositoryMock) {
				repo.On("ClearUserHistory", int64(123)).Return(nil).Once()
			},
			expectedOutput: "History cleared.\n",
			expectedError:  "",
		},
		{
			name: "failure - invalid limit",
			args: []string{"-n", "abc"},
			setupSession: func(repo *shellRepository.SessionRepositoryMock) {
				repo.On("GetSession").Return(shell.Session{User: &user.User{ID: 123}}, nil).Once()
			},
			setupHistory:   func(repo, guestRepo *historyRepository.HistoryRepositoryMock) {},
			expectedOutput: "",
			expectedError:  "invalid limit: abc\n",
		},
		{
			name: "failure - missing limit",
			args: []string{"-n"},
			setupSession: func(repo *shellRepository.SessionRepositoryMock) {
				repo.On("GetSession").Return(shell.Session{User: &user.User{ID: 123}}, nil).Once()
			},
			setupHistory:   func(repo, guestRepo *historyRepository.HistoryRepositoryMock) {},
			expectedOutput: "",
			expectedError:  "usage: history -n <limit>\n",
		},
		{
			name: "failure - session error",
			setupSession: func(repo *shellRepository.SessionRepositoryMock) {
				repo.On("GetSession").Return(shell.Session{}, errors.New("session error")).Once()
			},
			setupHistory:   func(repo, guestRepo *historyRepository.HistoryRepositoryMock) {},
			expectedOutput: "",
			expectedError:  "error getting session: session error\n",
		},
		{
			name: "failure - history service error",
			setupSession: func(repo *shellRepository.SessionRepositoryMock) {
				repo.On("GetSession").Return(shell.Session{User: &user.User{ID: 123}}, nil).Once()
			},
			setupHistory: func(repo, guestRepo *historyRepository.HistoryRepositoryMock) {
				repo.On("GetUserCommandStats", int64(123)).Return(nil, errors.New("history error")).Once()
			},
			expectedOutput: "",
			expectedError:  "error retrieving history: history error\n",
		},
		{
			name: "failure - clear history service error",
			args: []string{"clear"},
			setupSession: func(repo *shellRepository.SessionRepositoryMock) {
				repo.On("GetSession").Return(shell.Session{User: &user.User{ID: 123}}, nil).Once()
			},
			setupHistory: func(repo, guestRepo *historyRepository.HistoryRepositoryMock) {
				repo.On("ClearUserHistory", int64(123)).Return(errors.New("clear error")).Once()
			},
			expectedOutput: "",
			expectedError:  "error clearing history: clear error\n",
		},
		{
			name: "success - guest history",
			setupSession: func(repo *shellRepository.SessionRepositoryMock) {
				repo.On("GetSession").Return(shell.Session{}, nil).Once()
			},
			setupHistory: func(repo, guestRepo *historyRepository.HistoryRepositoryMock) {
				guestRepo.On("GetUserCommandStats", guestID).Return([]history.CommandStats{
					{Command: "ls", Count: 5},
					{Command: "cd", Count: 3},
				}, nil).Once()
			},
			expectedOutput: "Command   Count\n---       ---\nls        5\ncd        3\n",
			expectedError:  "",
		},
		{
			name: "success - clear guest history",
			args: []string{"clear"},
			setupSession: func(repo *shellRepository.SessionRepositoryMock) {
				repo.On("GetSession").Return(shell.Session{}, nil).Once()
			},
			setupHistory: func(repo, guestRepo *historyRepository.HistoryRepositoryMock) {
				guestRepo.On("ClearUserHistory", guestID).Return(nil).Once()
			},
			expectedOutput: "History cleared.\n",
			expectedError:  "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockSessionRepo := new(shellRepository.SessionRepositoryMock)
			mockHistoryRepo := new(historyRepository.HistoryRepositoryMock)
			mockGuestHistoryRepo := new(historyRepository.HistoryRepositoryMock)

			tc.setupSession(mockSessionRepo)
			tc.setupHistory(mockHistoryRepo, mockGuestHistoryRepo)

			var outputBuffer bytes.Buffer
			var errorBuffer bytes.Buffer

			historySvc := history.New(mockHistoryRepo, mockGuestHistoryRepo, guestID)
			cmd := commands.NewHistoryCommand(historySvc, mockSessionRepo)
			err := cmd.Execute(ctx, tc.args, nil, &outputBuffer, &errorBuffer)

			assert.NoError(t, err)

			// Normalize tabwriter output for consistent testing
			var normalizedOutput bytes.Buffer
			w := tabwriter.NewWriter(&normalizedOutput, 0, 0, 3, ' ', 0)
			w.Write(outputBuffer.Bytes())
			w.Flush()

			assert.Equal(t, tc.expectedOutput, normalizedOutput.String())
			assert.Equal(t, tc.expectedError, errorBuffer.String())
			mockSessionRepo.AssertExpectations(t)
			mockHistoryRepo.AssertExpectations(t)
			mockGuestHistoryRepo.AssertExpectations(t)
		})
	}
}
