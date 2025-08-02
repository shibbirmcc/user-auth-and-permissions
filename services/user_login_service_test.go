package services

import (
	"errors"
	"os"
	"testing"

	"github.com/shibbirmcc/user-auth-and-permissions/mocks"
	"github.com/shibbirmcc/user-auth-and-permissions/models"
	"github.com/stretchr/testify/assert"
)

func TestLogin_Success(t *testing.T) {
	mockDBService := new(mocks.MockDatabaseOperationService)
	loginService := NewUserLoginService(mockDBService)

	input := models.LoginRequest{
		Email:    mocks.TestUserEmail,
		Password: mocks.TestUserPassword,
	}
	user := &models.User{
		Email:    mocks.TestUserEmail,
		Password: mocks.TestUserPasswordHash,
	}
	userDetails := &models.UserDetail{
		FirstName: mocks.TestUserFirstName,
		LastName:  mocks.TestUserLastName,
	}

	mockDBService.On("FindUserByEmail", input.Email).Return(user, nil)
	mockDBService.On("FindUserDetailsByUserID", user.ID).Return(userDetails, nil)

	result, err := loginService.Login(input)

	assert.NoError(t, err)
	assert.NotNil(t, result)

	mockDBService.AssertExpectations(t)
}

func TestLogin_InvalidCredentials(t *testing.T) {
	mockDBService := new(mocks.MockDatabaseOperationService)
	loginService := NewUserLoginService(mockDBService)

	input := models.LoginRequest{
		Email:    mocks.TestUserEmail,
		Password: "wrongpassword",
	}

	user := &models.User{
		Email:    mocks.TestUserEmail,
		Password: mocks.TestUserPasswordHash,
	}
	mockDBService.On("FindUserByEmail", user.Email).Return(user, nil)

	result, err := loginService.Login(input)

	assert.EqualError(t, err, "Invalid credentials")
	assert.Empty(t, result)

	mockDBService.AssertExpectations(t)
}

func TestLogin_FailToFindUserDetails(t *testing.T) {
	mockDBService := new(mocks.MockDatabaseOperationService)
	loginService := NewUserLoginService(mockDBService)

	input := models.LoginRequest{
		Email:    mocks.TestUserEmail,
		Password: mocks.TestUserPassword,
	}
	user := &models.User{
		Email:    mocks.TestUserEmail,
		Password: mocks.TestUserPasswordHash,
	}
	mockDBService.On("FindUserByEmail", user.Email).Return(user, nil)
	mockDBService.On("FindUserDetailsByUserID", user.ID).Return(nil, errors.New("user details not found"))

	result, err := loginService.Login(input)

	assert.EqualError(t, err, "Invalid user Id")
	assert.Empty(t, result)

	mockDBService.AssertExpectations(t)
}

func TestLogin_FailToGenerateToken(t *testing.T) {
	originalSecret := os.Getenv("JWT_SECRET")
	defer os.Setenv("JWT_SECRET", originalSecret) // Reset after the test
	// Unset the JWT_SECRET environment variable
	os.Unsetenv("JWT_SECRET")

	mockDBService := new(mocks.MockDatabaseOperationService)
	loginService := NewUserLoginService(mockDBService)

	input := models.LoginRequest{
		Email:    mocks.TestUserEmail,
		Password: mocks.TestUserPassword,
	}
	user := &models.User{
		Email:    mocks.TestUserEmail,
		Password: mocks.TestUserPasswordHash,
	}
	userDetails := &models.UserDetail{
		FirstName: mocks.TestUserFirstName,
		LastName:  mocks.TestUserLastName,
	}

	mockDBService.On("FindUserByEmail", user.Email).Return(user, nil)
	mockDBService.On("FindUserDetailsByUserID", user.ID).Return(userDetails, nil)

	// Call the Login method
	result, err := loginService.Login(input)

	// Assert that an error is returned due to token generation failure
	assert.EqualError(t, err, "Could not generate token")
	assert.Empty(t, result)

	mockDBService.AssertExpectations(t)
}
