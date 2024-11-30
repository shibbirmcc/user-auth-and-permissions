package routes

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/shibbirmcc/user-auth-and-permissions/handlers"
	"github.com/shibbirmcc/user-auth-and-permissions/mocks"
	"github.com/shibbirmcc/user-auth-and-permissions/models"
	"github.com/shibbirmcc/user-auth-and-permissions/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestConfigureRouteEndPoints(t *testing.T) {

	mockDBService := new(mocks.MockDatabaseOperationService)
	mockPasswordDeliveryService := new(mocks.MockPasswordDeliveryService)
	mockRegService := services.NewUserRegistrationService(mockDBService)
	mockLoginService := services.NewUserLoginService(mockDBService)
	userHandler := handlers.NewUserHandler(*mockRegService, *mockLoginService, mockPasswordDeliveryService)

	router := gin.Default()
	ConfigureRouteEndpoints(router, userHandler)

	t.Run("RegisterUser endpoint", func(t *testing.T) {
		input := models.UserRegitrationRequest{
			Email:     "test@example.com",
			FirstName: "FirstName",
			LastName:  "LastName",
		}
		body, _ := json.Marshal(input)

		// Mock the CreateUser method of the mockDBService to return no error.
		mockDBService.On("CreateUser", mock.Anything, mock.Anything).Return(nil)

		// Create HTTP request and record response
		req := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		// Call router
		router.ServeHTTP(resp, req)

		// Assert expectations
		assert.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("LoginUser endpoint", func(t *testing.T) {
		loginRequest := models.LoginRequest{
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

		mockDBService.On("FindUserByEmail", loginRequest.Email).Return(user, nil)
		mockDBService.On("FindUserDetailsByUserID", user.ID).Return(userDetails, nil)

		// Create a request
		requestBody, _ := json.Marshal(loginRequest)

		// Create HTTP request and record response
		req := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		// Call router
		router.ServeHTTP(resp, req)

		// Assert expectations
		assert.Equal(t, http.StatusOK, resp.Code)
		// You may add more assertions if necessary
	})

}
