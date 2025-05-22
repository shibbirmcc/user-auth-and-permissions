package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/shibbirmcc/user-auth-and-permissions/mocks"
	"github.com/shibbirmcc/user-auth-and-permissions/models"
	"github.com/shibbirmcc/user-auth-and-permissions/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestRegisterUser_Success tests the successful registration flow.
func TestRegisterUser_Success(t *testing.T) {
	mockDBService := new(mocks.MockDatabaseOperationService)
	mockPasswordDeliveryService := &mocks.MockPasswordDeliveryService{ShouldFail: false}
	mockRegService := services.NewUserRegistrationService(mockPasswordDeliveryService, mockDBService)
	mockLoginService := services.NewUserLoginService(mockDBService)
	handler := &UserHandler{userRegistrationService: *mockRegService, userLoginService: *mockLoginService}

	// Define the request body.
	input := models.UserRegitrationRequest{
		Email:     "test@example.com",
		FirstName: "FirstName",
		LastName:  "LastName",
	}
	body, _ := json.Marshal(input)

	// Mock the CreateUser method of the mockDBService to return no error.
	mockDBService.On("CreateUser", mock.Anything, mock.Anything).Return(nil)

	// Create a test request.
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	// Call the RegisterUser handler.
	handler.RegisterUser(c)

	// Check the response.
	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"message": "Registration successful, please confirm your email"}`, w.Body.String())

	// Assert that the expected methods were called.
	mockDBService.AssertExpectations(t)
}

// TestRegisterUser_InternalServerError tests the case where the CreateUser method fails.
func TestRegisterUser_InternalServerError(t *testing.T) {
	mockDBService := new(mocks.MockDatabaseOperationService)
	mockPasswordDeliveryService := &mocks.MockPasswordDeliveryService{ShouldFail: false}
	mockRegService := services.NewUserRegistrationService(mockPasswordDeliveryService, mockDBService)
	mockLoginService := services.NewUserLoginService(mockDBService)
	handler := &UserHandler{userRegistrationService: *mockRegService, userLoginService: *mockLoginService}

	// Define the request body.
	input := models.UserRegitrationRequest{
		Email:     "test@example.com",
		FirstName: "FirstName",
		LastName:  "LastName",
	}
	body, _ := json.Marshal(input)

	// Mock the CreateUser method to return an error.
	mockDBService.On("CreateUser", mock.Anything, mock.Anything).Return(errors.New("failed to create user"))

	// Create a test request.
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	// Call the RegisterUser handler.
	handler.RegisterUser(c)

	// Check the response for an internal server error.
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "error while registering user")

	// Assert that the expected methods were called.
	mockDBService.AssertExpectations(t)
}

func TestRegisterUser_InvalidJson(t *testing.T) {
	mockDBService := new(mocks.MockDatabaseOperationService)
	mockPasswordDeliveryService := &mocks.MockPasswordDeliveryService{ShouldFail: false}
	mockRegService := services.NewUserRegistrationService(mockPasswordDeliveryService, mockDBService)
	mockLoginService := services.NewUserLoginService(mockDBService)
	handler := &UserHandler{userRegistrationService: *mockRegService, userLoginService: *mockLoginService}

	body := `{"email": "test@example.com", "firstName": 123}`

	// Create a test request.
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer([]byte(body)))
	c.Request.Header.Set("Content-Type", "application/json")

	// Call the RegisterUser handler.
	handler.RegisterUser(c)

	// Check the response.
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}
