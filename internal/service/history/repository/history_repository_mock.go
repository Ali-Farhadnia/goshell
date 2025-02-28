package repository

import (
	"github.com/Ali-Farhadnia/goshell/internal/service/history"
	"github.com/stretchr/testify/mock"
)

type HistoryRepositoryMock struct {
	mock.Mock
}

func (m *HistoryRepositoryMock) SaveCommand(command *history.CommandHistory) error {
	args := m.Called(command)
	return args.Error(0)
}

func (m *HistoryRepositoryMock) GetUserHistory(userID int64, limit int) ([]history.CommandHistory, error) {
	args := m.Called(userID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]history.CommandHistory), args.Error(1)
}

func (m *HistoryRepositoryMock) GetUserCommandStats(userID int64) ([]history.CommandStats, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]history.CommandStats), args.Error(1)
}

func (m *HistoryRepositoryMock) ClearUserHistory(userID int64) error {
	args := m.Called(userID)
	return args.Error(0)
}
