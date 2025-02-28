package repository

import (
	"github.com/Ali-Farhadnia/goshell/internal/service/user"
	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	mock.Mock
}

func (m *UserRepositoryMock) FindUserByUsername(username string) (user.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return user.User{}, args.Error(1)
	}
	return args.Get(0).(user.User), args.Error(1)
}

func (m *UserRepositoryMock) CreateUser(user user.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *UserRepositoryMock) UpdateUser(user user.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *UserRepositoryMock) ListUsers() ([]user.User, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]user.User), args.Error(1)
}

func (m *UserRepositoryMock) UpdateLastLogin(userID int64) error {
	args := m.Called(userID)
	return args.Error(0)
}
