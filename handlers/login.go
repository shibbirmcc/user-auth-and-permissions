package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/shibbirmcc/user-auth-and-permissions/models"
)

// LoginResponse represents the structure of a successful login response
type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
}

// ErrorResponse represents the structure of an error response
type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

func (h *UserHandler) LoginUser(c *gin.Context) {
	var input models.LoginRequest

	// Bind and validate JSON input
	if err := c.ShouldBindJSON(&input); err != nil {
		// Log the validation error for debugging (without sensitive data)
		log.Printf("Login validation error for email %s: %v", input.Email, err)
		
		// Return user-friendly validation error
		errorMsg := "Invalid input data"
		if strings.Contains(err.Error(), "email") {
			errorMsg = "Please provide a valid email address"
		} else if strings.Contains(err.Error(), "password") {
			errorMsg = "Password is required"
		}
		
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Message: errorMsg,
			Error:   "validation_failed",
		})
		return
	}

	// Additional input validation
	if strings.TrimSpace(input.Email) == "" {
		log.Printf("Login attempt with empty email from IP: %s", c.ClientIP())
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Message: "Email is required",
			Error:   "missing_email",
		})
		return
	}

	if strings.TrimSpace(input.Password) == "" {
		log.Printf("Login attempt with empty password for email: %s from IP: %s", input.Email, c.ClientIP())
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Message: "Password is required",
			Error:   "missing_password",
		})
		return
	}

	// Attempt login
	token, err := h.userLoginService.Login(input)
	if err != nil {
		// Log failed login attempt for security monitoring
		log.Printf("Failed login attempt for email: %s from IP: %s - Error: %v", input.Email, c.ClientIP(), err)
		
		// Return generic error message to prevent information disclosure
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Success: false,
			Message: "Invalid email or password",
			Error:   "authentication_failed",
		})
		return
	}

	// Log successful login for security monitoring
	log.Printf("Successful login for email: %s from IP: %s", input.Email, c.ClientIP())

	// Set security headers
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("X-Frame-Options", "DENY")
	c.Header("X-XSS-Protection", "1; mode=block")

	// Return successful response
	c.JSON(http.StatusOK, LoginResponse{
		Success: true,
		Message: "Login successful",
		Token:   token,
	})
}
