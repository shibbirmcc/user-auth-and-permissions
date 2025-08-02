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
		mockPasswordDeliveryService := &mocks.MockPasswordDeliveryService{ShouldFail: false}
		mockRegService := services.NewUserRegistrationService(mockPasswordDeliveryService, mockDBService)
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
		
		// Verify response structure
		var response LoginResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
		assert.Equal(t, "Login successful", response.Message)
		assert.NotEmpty(t, response.Token)
		
		// Verify security headers
		assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
		assert.Equal(t, "DENY", w.Header().Get("X-Frame-Options"))
		assert.Equal(t, "1; mode=block", w.Header().Get("X-XSS-Protection"))
	})

	t.Run("bad request with invalid JSON", func(t *testing.T) {
		mockDBService := new(mocks.MockDatabaseOperationService)
		mockPasswordDeliveryService := &mocks.MockPasswordDeliveryService{ShouldFail: false}
		mockRegService := services.NewUserRegistrationService(mockPasswordDeliveryService, mockDBService)
		mockLoginService := services.NewUserLoginService(mockDBService)
		handler := &UserHandler{userRegistrationService: *mockRegService, userLoginService: *mockLoginService}

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
		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response.Success)
		assert.Equal(t, "Invalid input data", response.Message)
		assert.Equal(t, "validation_failed", response.Error)
	})

	t.Run("Invalid Email format Error Test", func(t *testing.T) {
		mockDBService := new(mocks.MockDatabaseOperationService)
		mockPasswordDeliveryService := &mocks.MockPasswordDeliveryService{ShouldFail: false}
		mockRegService := services.NewUserRegistrationService(mockPasswordDeliveryService, mockDBService)
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
		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response.Success)
		assert.Equal(t, "Please provide a valid email address", response.Message)
		assert.Equal(t, "validation_failed", response.Error)
	})

	t.Run("Unauthorized Login Test", func(t *testing.T) {
		mockDBService := new(mocks.MockDatabaseOperationService)
		mockPasswordDeliveryService := &mocks.MockPasswordDeliveryService{ShouldFail: false}
		mockRegService := services.NewUserRegistrationService(mockPasswordDeliveryService, mockDBService)
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
		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response.Success)
		assert.Equal(t, "Invalid email or password", response.Message)
		assert.Equal(t, "authentication_failed", response.Error)
	})

	t.Run("Empty email validation", func(t *testing.T) {
		mockDBService := new(mocks.MockDatabaseOperationService)
		mockPasswordDeliveryService := &mocks.MockPasswordDeliveryService{ShouldFail: false}
		mockRegService := services.NewUserRegistrationService(mockPasswordDeliveryService, mockDBService)
		mockLoginService := services.NewUserLoginService(mockDBService)
		handler := &UserHandler{userRegistrationService: *mockRegService, userLoginService: *mockLoginService}

		// Sample request data with empty email (spaces are treated as invalid email format by binding validation)
		loginRequest := models.LoginRequest{
			Email:    "   ",
			Password: "password123",
		}

		// Create a request
		requestBody, _ := json.Marshal(loginRequest)
		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		// Create a ResponseRecorder to capture the response
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		handler.LoginUser(c)

		// Assert the results - binding validation catches this first as invalid email format
		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response.Success)
		assert.Equal(t, "Please provide a valid email address", response.Message)
		assert.Equal(t, "validation_failed", response.Error)
	})

	t.Run("Missing email field validation", func(t *testing.T) {
		mockDBService := new(mocks.MockDatabaseOperationService)
		mockPasswordDeliveryService := &mocks.MockPasswordDeliveryService{ShouldFail: false}
		mockRegService := services.NewUserRegistrationService(mockPasswordDeliveryService, mockDBService)
		mockLoginService := services.NewUserLoginService(mockDBService)
		handler := &UserHandler{userRegistrationService: *mockRegService, userLoginService: *mockLoginService}

		// Create request with missing email field entirely
		requestBody := []byte(`{"password": "password123"}`)
		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		// Create a ResponseRecorder to capture the response
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		handler.LoginUser(c)

		// Assert the results - missing required field triggers general validation error
		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response.Success)
		assert.Equal(t, "Invalid input data", response.Message)
		assert.Equal(t, "validation_failed", response.Error)
	})

	t.Run("Empty password validation", func(t *testing.T) {
		mockDBService := new(mocks.MockDatabaseOperationService)
		mockPasswordDeliveryService := &mocks.MockPasswordDeliveryService{ShouldFail: false}
		mockRegService := services.NewUserRegistrationService(mockPasswordDeliveryService, mockDBService)
		mockLoginService := services.NewUserLoginService(mockDBService)
		handler := &UserHandler{userRegistrationService: *mockRegService, userLoginService: *mockLoginService}

		// Sample request data with empty password
		loginRequest := models.LoginRequest{
			Email:    mocks.TestUserEmail,
			Password: "   ",
		}

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
		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response.Success)
		assert.Equal(t, "Password is required", response.Message)
		assert.Equal(t, "missing_password", response.Error)
	})
}
