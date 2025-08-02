package services

import (
	"errors"
	"testing"

	"github.com/shibbirmcc/user-auth-and-permissions/mocks"
	"github.com/shibbirmcc/user-auth-and-permissions/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegisterUser_Success(t *testing.T) {
	// Create the mock database service
	mockDB := new(mocks.MockDatabaseOperationService)

	// Define the input for RegisterUser
	input := models.UserRegitrationRequest{
		Email:      "test@example.com",
		FirstName:  "John",
		MiddleName: "A",
		LastName:   "Doe",
	}

	// Set up expected calls and return values on the mock database
	mockDB.On("CreateUser", mock.AnythingOfType("*models.User"), mock.AnythingOfType("*models.UserDetail")).Return(nil)

	// Create the service with the mock database
	mockPasswordDeliveryService := &mocks.MockPasswordDeliveryService{ShouldFail: false}
	service := NewUserRegistrationService(mockPasswordDeliveryService, mockDB)

	// Call the RegisterUser method
	err := service.RegisterUser(input)

	// Assertions
	assert.NoError(t, err, "RegisterUser should succeed without errors")
	mockDB.AssertExpectations(t) // Ensure that all expectations were met
}

func TestRegisterUser_FailOnCreateUser(t *testing.T) {
	// Create the mock database service
	mockDB := new(mocks.MockDatabaseOperationService)

	// Define the input for RegisterUser
	input := models.UserRegitrationRequest{
		Email:      "fail@example.com",
		FirstName:  "Jane",
		MiddleName: "B",
		LastName:   "Doe",
	}

	// Set up the mock to return an error on CreateUser
	mockDB.On("CreateUser", mock.AnythingOfType("*models.User"), mock.AnythingOfType("*models.UserDetail")).Return(errors.New("database error"))

	// Create the service with the mock database
	mockPasswordDeliveryService := &mocks.MockPasswordDeliveryService{ShouldFail: false}
	service := NewUserRegistrationService(mockPasswordDeliveryService, mockDB)

	// Call the RegisterUser method
	err := service.RegisterUser(input)

	// Assertions
	assert.Error(t, err, "RegisterUser should return an error when CreateUser fails")
	assert.EqualError(t, err, "error while registering user", "Expected error message for database failure")
	mockDB.AssertExpectations(t) // Ensure that all expectations were met
}
