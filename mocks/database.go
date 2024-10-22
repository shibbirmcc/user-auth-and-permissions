package mocks

import (
	"github.com/shibbirmcc/user-auth-and-permissions/models"
	"github.com/stretchr/testify/mock"
)

type MockDatabaseOperationService struct {
	mock.Mock
}

func (m *MockDatabaseOperationService) CreateUser(user *models.User, userDetail *models.UserDetail) error {
	args := m.Called(user, userDetail)
	return args.Error(0)
}

func (m *MockDatabaseOperationService) FindUserByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) != nil {
		return args.Get(0).(*models.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockDatabaseOperationService) FindUserDetailsByUserID(userID uint) (*models.UserDetail, error) {
	args := m.Called(userID)
	if args.Get(0) != nil {
		return args.Get(0).(*models.UserDetail), args.Error(1)
	}
	return nil, args.Error(1)
}
