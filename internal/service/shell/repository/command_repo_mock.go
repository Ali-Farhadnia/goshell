package repository

import (
	"github.com/Ali-Farhadnia/goshell/internal/service/shell"
	"github.com/stretchr/testify/mock"
)

// CommandRepositoryMock is a mock implementation of CommandRepository.
type CommandRepositoryMock struct {
	mock.Mock
}

// Register mocks the Register method.
func (m *CommandRepositoryMock) Register(cmd shell.Command) error {
	args := m.Called(cmd)
	return args.Error(0)
}

// Get mocks the Get method.
func (m *CommandRepositoryMock) Get(name string) (shell.Command, error) {
	args := m.Called(name)
	return args.Get(0).(shell.Command), args.Error(1)
}

// List mocks the List method.
func (m *CommandRepositoryMock) List() ([]shell.Command, error) {
	args := m.Called()
	return args.Get(0).([]shell.Command), args.Error(1)
}
