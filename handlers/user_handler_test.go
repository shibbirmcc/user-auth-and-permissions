package handlers

import (
	"testing"

	"github.com/shibbirmcc/user-auth-and-permissions/mocks"
	"github.com/shibbirmcc/user-auth-and-permissions/services"
	"github.com/stretchr/testify/assert"
)

func TestNewUserHandler(t *testing.T) {
	mockDBService := new(mocks.MockDatabaseOperationService)
	mockRegService := services.NewUserRegistrationService(mockDBService)
	mockLoginService := services.NewUserLoginService(mockDBService)
	handler := NewUserHandler(*mockRegService, *mockLoginService)

	assert.NotNil(t, handler)

	// assert.Equal(t, mockRegService, handler.userRegistrationService)
	// assert.Equal(t, mockLoginService, handler.userLoginService)
}
