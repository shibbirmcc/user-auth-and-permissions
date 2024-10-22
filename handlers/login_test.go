package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/shibbirmcc/user-auth-and-permissions/mocks"
	"github.com/shibbirmcc/user-auth-and-permissions/models"
	"github.com/shibbirmcc/user-auth-and-permissions/services"
	"github.com/stretchr/testify/assert"
)

func TestLoginUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("successful login", func(t *testing.T) {
		mockDBService := new(mocks.MockDatabaseOperationService)
		mockRegService := services.NewUserRegistrationService(mockDBService)
		mockLoginService := services.NewUserLoginService(mockDBService)

		handler := &UserHandler{userRegistrationService: *mockRegService, userLoginService: *mockLoginService}

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
		req, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		// Create a ResponseRecorder to capture the response
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Call the handler
		handler.LoginUser(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("bad request with invalid JSON", func(t *testing.T) {
		mockDBService := new(mocks.MockDatabaseOperationService)
		mockRegService := services.NewUserRegistrationService(mockDBService)
		mockLoginService := services.NewUserLoginService(mockDBService)
		handler := &UserHandler{userRegistrationService: *mockRegService, userLoginService: *mockLoginService}

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

		// Create a request with invalid JSON
		req, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer([]byte("{invalid-json")))
		req.Header.Set("Content-Type", "application/json")

		// Create a ResponseRecorder to capture the response
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Call the handler
		handler.LoginUser(c)

		// Assert the results
		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "invalid character")
	})

	t.Run("Invalid Email format Error Test", func(t *testing.T) {
		mockDBService := new(mocks.MockDatabaseOperationService)
		mockRegService := services.NewUserRegistrationService(mockDBService)
		mockLoginService := services.NewUserLoginService(mockDBService)
		handler := &UserHandler{userRegistrationService: *mockRegService, userLoginService: *mockLoginService}

		// Sample request data
		loginRequest := models.LoginRequest{
			Email:    "testuser",
			Password: "wrongpass",
		}

		mockDBService.On("FindUserByEmail", loginRequest.Email).Return(nil, nil)

		// Create a request
		requestBody, _ := json.Marshal(loginRequest)
		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		// Create a ResponseRecorder to capture the response
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		handler.LoginUser(c)

		// Assert the results
		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "Error:Field validation for 'Email' failed on the 'email' tag")
	})

	t.Run("Unauthorized Login Test", func(t *testing.T) {
		mockDBService := new(mocks.MockDatabaseOperationService)
		mockRegService := services.NewUserRegistrationService(mockDBService)
		mockLoginService := services.NewUserLoginService(mockDBService)
		handler := &UserHandler{userRegistrationService: *mockRegService, userLoginService: *mockLoginService}

		// Sample request data
		loginRequest := models.LoginRequest{
			Email:    mocks.TestUserEmail,
			Password: "wrongpass",
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

		// Create a request
		requestBody, _ := json.Marshal(loginRequest)
		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		// Create a ResponseRecorder to capture the response
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		handler.LoginUser(c)

		// Assert the results
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "Invalid credentials")
	})
}
