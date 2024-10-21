package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/shibbirmcc/user-auth-and-permissions/models"
	"github.com/shibbirmcc/user-auth-and-permissions/services"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	err := godotenv.Load("../.env.test")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Run the tests
	code := m.Run()

	// Exit with the appropriate code
	os.Exit(code)
}

func TestLoginUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("successful login", func(t *testing.T) {
		mockDBService := new(MockDatabaseOperationService)
		mockRegService := services.NewUserRegistrationService(mockDBService)
		mockLoginService := services.NewUserLoginService(mockDBService)

		handler := &UserHandler{userRegistrationService: *mockRegService, userLoginService: *mockLoginService}

		loginRequest := models.LoginRequest{
			Email:    testUserEmail,
			Password: testUserPassword,
		}

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
		mockDBService := new(MockDatabaseOperationService)
		mockRegService := services.NewUserRegistrationService(mockDBService)
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
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "invalid character")
	})

	t.Run("Invalid Email format Error Test", func(t *testing.T) {
		mockDBService := new(MockDatabaseOperationService)
		mockRegService := services.NewUserRegistrationService(mockDBService)
		mockLoginService := services.NewUserLoginService(mockDBService)
		handler := &UserHandler{userRegistrationService: *mockRegService, userLoginService: *mockLoginService}

		// Sample request data
		loginRequest := models.LoginRequest{
			Email:    "testuser",
			Password: "wrongpass",
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
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "Error:Field validation for 'Email' failed on the 'email' tag")
	})

	t.Run("Unauthorized Login Test", func(t *testing.T) {
		mockDBService := new(MockDatabaseOperationService)
		mockRegService := services.NewUserRegistrationService(mockDBService)
		mockLoginService := services.NewUserLoginService(mockDBService)
		handler := &UserHandler{userRegistrationService: *mockRegService, userLoginService: *mockLoginService}

		// Sample request data
		loginRequest := models.LoginRequest{
			Email:    testUserEmail,
			Password: "wrongpass",
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

		fmt.Println("Response Body:", w.Body.String())

		// Assert the results
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "Invalid credentials")
	})
}
