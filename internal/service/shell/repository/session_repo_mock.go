package repository

import (
	"github.com/Ali-Farhadnia/goshell/internal/service/shell"
	"github.com/stretchr/testify/mock"
)

// SessionRepositoryMock is a mock implementation of SessionRepository.
type SessionRepositoryMock struct {
	mock.Mock
}

// GetSession mocks the GetSession method.
func (m *SessionRepositoryMock) GetSession() (shell.Session, error) {
	args := m.Called()
	return args.Get(0).(shell.Session), args.Error(1)
}

// SetSession mocks the SetSession method.
func (m *SessionRepositoryMock) SetSession(s shell.Session) error {
	args := m.Called(s)
	return args.Error(0)
}
