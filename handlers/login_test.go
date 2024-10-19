package handlers

import (
	"bytes"
	"encoding/json"
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
}
