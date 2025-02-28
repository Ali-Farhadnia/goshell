package history_test

import (
	"errors"
	"testing"
	"time"

	"github.com/Ali-Farhadnia/goshell/internal/service/history"
	"github.com/Ali-Farhadnia/goshell/internal/service/history/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_SaveCommandHistory(t *testing.T) {
	guestID := int64(123)
	mockRepo := new(repository.HistoryRepositoryMock)
	mockGuestRepo := new(repository.HistoryRepositoryMock)
	service := history.New(mockRepo, mockGuestRepo, guestID)

	t.Run("guest user", func(t *testing.T) {
		command := "ls -l"
		expectedHistory := &history.CommandHistory{
			UserID:    guestID,
			Command:   command,
			CreatedAt: time.Now(),
		}

		mockGuestRepo.On("SaveCommand", mock.MatchedBy(func(history *history.CommandHistory) bool {
			return history.UserID == expectedHistory.UserID && history.Command == expectedHistory.Command
		})).Return(nil).Once()

		err := service.SaveCommandHistory(nil, command)
		assert.NoError(t, err)
		mockGuestRepo.AssertExpectations(t)
	})

	t.Run("regular user", func(t *testing.T) {
		userID := int64(456)
		command := "git commit -m 'test'"
		expectedHistory := &history.CommandHistory{
			UserID:    userID,
			Command:   command,
			CreatedAt: time.Now(),
		}

		mockRepo.On("SaveCommand", mock.MatchedBy(func(history *history.CommandHistory) bool {
			return history.UserID == expectedHistory.UserID && history.Command == expectedHistory.Command
		})).Return(nil).Once()

		err := service.SaveCommandHistory(&userID, command)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("save command error", func(t *testing.T) {
		userID := int64(456)
		command := "git commit -m 'test'"
		expectedError := errors.New("save error")

		mockRepo.On("SaveCommand", mock.Anything).Return(expectedError).Once()

		err := service.SaveCommandHistory(&userID, command)
		assert.ErrorIs(t, err, expectedError)
		mockRepo.AssertExpectations(t)
	})

	t.Run("guest save command error", func(t *testing.T) {
		command := "ls -l"
		expectedError := errors.New("guest save error")

		mockGuestRepo.On("SaveCommand", mock.Anything).Return(expectedError).Once()

		err := service.SaveCommandHistory(nil, command)
		assert.ErrorIs(t, err, expectedError)
		mockGuestRepo.AssertExpectations(t)
	})
}

func TestService_GetCommandHistoryStats(t *testing.T) {
	guestID := int64(123)
	mockRepo := new(repository.HistoryRepositoryMock)
	mockGuestRepo := new(repository.HistoryRepositoryMock)
	service := history.New(mockRepo, mockGuestRepo, guestID)

	t.Run("guest user", func(t *testing.T) {
		expectedStats := []history.CommandStats{{Command: "ls", Count: 1}}
		mockGuestRepo.On("GetUserCommandStats", guestID).Return(expectedStats, nil).Once()

		stats, err := service.GetCommandHistoryStats(nil, 10)
		assert.NoError(t, err)
		assert.Equal(t, expectedStats, stats)
		mockGuestRepo.AssertExpectations(t)
	})

	t.Run("regular user", func(t *testing.T) {
		userID := int64(456)
		expectedStats := []history.CommandStats{{Command: "git", Count: 2}}
		mockRepo.On("GetUserCommandStats", userID).Return(expectedStats, nil).Once()

		stats, err := service.GetCommandHistoryStats(&userID, 10)
		assert.NoError(t, err)
		assert.Equal(t, expectedStats, stats)
		mockRepo.AssertExpectations(t)
	})

	t.Run("get stats error", func(t *testing.T) {
		userID := int64(456)
		expectedError := errors.New("stats error")
		mockRepo.On("GetUserCommandStats", userID).Return(nil, expectedError).Once()

		_, err := service.GetCommandHistoryStats(&userID, 10)
		assert.ErrorIs(t, err, expectedError)
		mockRepo.AssertExpectations(t)
	})

	t.Run("guest get stats error", func(t *testing.T) {
		expectedError := errors.New("guest stats error")
		mockGuestRepo.On("GetUserCommandStats", guestID).Return(nil, expectedError).Once()

		_, err := service.GetCommandHistoryStats(nil, 10)
		assert.ErrorIs(t, err, expectedError)
		mockGuestRepo.AssertExpectations(t)
	})
}

func TestService_ClearCommandHistory(t *testing.T) {
	guestID := int64(123)
	mockRepo := new(repository.HistoryRepositoryMock)
	mockGuestRepo := new(repository.HistoryRepositoryMock)
	service := history.New(mockRepo, mockGuestRepo, guestID)

	t.Run("guest user", func(t *testing.T) {
		mockGuestRepo.On("ClearUserHistory", guestID).Return(nil).Once()

		err := service.ClearCommandHistory(nil)
		assert.NoError(t, err)
		mockGuestRepo.AssertExpectations(t)
	})

	t.Run("regular user", func(t *testing.T) {
		userID := int64(456)
		mockRepo.On("ClearUserHistory", userID).Return(nil).Once()

		err := service.ClearCommandHistory(&userID)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("clear error", func(t *testing.T) {
		userID := int64(456)
		expectedError := errors.New("clear error")
		mockRepo.On("ClearUserHistory", userID).Return(expectedError).Once()

		err := service.ClearCommandHistory(&userID)
		assert.ErrorIs(t, err, expectedError)
		mockRepo.AssertExpectations(t)
	})

	t.Run("guest clear error", func(t *testing.T) {
		expectedError := errors.New("guest clear error")
		mockGuestRepo.On("ClearUserHistory", guestID).Return(expectedError).Once()

		err := service.ClearCommandHistory(nil)
		assert.ErrorIs(t, err, expectedError)
		mockGuestRepo.AssertExpectations(t)
	})
}
