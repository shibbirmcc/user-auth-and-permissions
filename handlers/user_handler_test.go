package handlers

import (
	"testing"

	"github.com/shibbirmcc/user-auth-and-permissions/models"
	"github.com/shibbirmcc/user-auth-and-permissions/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const testUserEmail string = "user1@testmail.com"
const testUserPassword string = "TestPasswordAdmin1234"
const testUserId uint = 1
const testUserFirstName string = "Test"
const testUserLastName string = "User"

type MockDatabaseOperationService struct {
	mock.Mock
}

func (m *MockDatabaseOperationService) CreateUser(user *models.User, userDetail *models.UserDetail) error {
	args := m.Called(user, userDetail)
	return args.Error(0)
}

func (m *MockDatabaseOperationService) FindUserByEmail(email string) (*models.User, error) {
	return &models.User{
		ID:       testUserId,
		Email:    testUserEmail,
		Password: "$2y$10$4SJKHqzU178iOYuqwj.wH.MdLzM7HuCX9eoJJgWFMxM6fwUHwSg32", // this hash is for the password: TestPasswordAdmin1234
	}, nil
}

func (m *MockDatabaseOperationService) FindUserDetailsByUserID(userID uint) (*models.UserDetail, error) {
	return &models.UserDetail{
		UserID:    testUserId,
		FirstName: testUserFirstName,
		LastName:  testUserLastName,
	}, nil
}

func TestNewUserHandler(t *testing.T) {
	mockDBService := new(MockDatabaseOperationService)
	mockRegService := services.NewUserRegistrationService(mockDBService)
	mockLoginService := services.NewUserLoginService(mockDBService)
	handler := NewUserHandler(*mockRegService, *mockLoginService)

	assert.NotNil(t, handler)

	// assert.Equal(t, mockRegService, handler.userRegistrationService)
	// assert.Equal(t, mockLoginService, handler.userLoginService)
}
