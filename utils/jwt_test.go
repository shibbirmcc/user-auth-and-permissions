package utils

import (
	"os"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/shibbirmcc/user-auth-and-permissions/models"
	"github.com/stretchr/testify/assert"
)

// Mock user details for the tests
var mockUserDetails = models.UserDetail{
	UserID:     12,
	FirstName:  "John",
	MiddleName: "F.",
	LastName:   "Doe",
}

// TestGenerateJWT_Success tests the successful generation of a JWT.
func TestGenerateJWT_Success(t *testing.T) {
	// Set the JWT secret for testing
	os.Setenv("JWT_SECRET", "mysecretkey")

	email := "johndoe@example.com"

	// Generate JWT token
	tokenString, err := GenerateJWT(email, mockUserDetails)

	// Assertions
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	// Parse the token to verify claims
	token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	// Ensure the token is valid and claims are correct
	assert.NoError(t, err)
	assert.True(t, token.Valid)

	if claims, ok := token.Claims.(*models.Claims); ok && token.Valid {
		assert.Equal(t, email, claims.Email)
		assert.Equal(t, mockUserDetails.UserID, claims.UserID)
		assert.Equal(t, mockUserDetails.FirstName, claims.FirstName)
		assert.Equal(t, mockUserDetails.MiddleName, claims.MiddleName)
		assert.Equal(t, mockUserDetails.LastName, claims.LastName)
	} else {
		t.Fail()
	}
}

// TestGenerateJWT_EmptyEmail tests that an error is returned when email is empty.
func TestGenerateJWT_EmptyEmail(t *testing.T) {
	// Set the JWT secret for testing
	os.Setenv("JWT_SECRET", "mysecretkey")

	email := "" // Empty email

	// Generate JWT token
	tokenString, err := GenerateJWT(email, mockUserDetails)

	// Assertions
	assert.EqualError(t, err, "email cannot be empty")
	assert.Empty(t, tokenString)
}

func TestGenerateJWT_MissingSecret(t *testing.T) {
	// Unset the JWT secret
	os.Unsetenv("JWT_SECRET")

	email := "johndoe@example.com"

	// Generate JWT token
	tokenString, err := GenerateJWT(email, mockUserDetails)

	// Assertions
	assert.EqualError(t, err, "JWT_SECRET is missing")
	assert.Empty(t, tokenString)
}
