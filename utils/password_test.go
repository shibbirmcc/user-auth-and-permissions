package utils_test

import (
	"testing"

	"github.com/shibbirmcc/user-auth-and-permissions/utils"
	"golang.org/x/crypto/bcrypt"
)

// TestHashPassword tests the HashPassword function
func TestHashPassword(t *testing.T) {
	password := "testpassword123"
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Check if the hashed password can be compared with the original password
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		t.Fatalf("expected password to match, but got error: %v", err)
	}
}

// TestHashPassword_Error tests error case in HashPassword
func TestHashPassword_Error(t *testing.T) {
	// Simulate an invalid password length to trigger bcrypt failure
	invalidPassword := ""
	_, err := utils.HashPassword(invalidPassword)
	if err == nil {
		t.Fatalf("expected error when hashing an empty password, got nil")
	}
}

// TestGetRandomPasswordAndHash tests the GetRandomPasswordAndHash function
func TestGetRandomPasswordAndHash(t *testing.T) {
	password, hashedPassword, err := utils.GetRandomPasswordAndHash()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(password) != 12 {
		t.Fatalf("expected password of length 12, but got %d", len(password))
	}

	// Check if the hashed password can be compared with the generated password
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		t.Fatalf("expected password to match, but got error: %v", err)
	}
}

// TestGenerateRandomPassword tests the GenerateRandomPassword function
func TestGenerateRandomPassword(t *testing.T) {
	length := 16
	password, err := utils.GenerateRandomPassword(length)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(password) != length {
		t.Fatalf("expected password of length %d, but got %d", length, len(password))
	}
}

// TestGenerateRandomPassword_Error tests the error case of GenerateRandomPassword
func TestGenerateRandomPassword_Error(t *testing.T) {
	// This test case would simulate failure of rand.Int function
	// however, since rand.Int cannot be easily mocked, this error scenario
	// is hard to simulate in standard testing.

	// Leaving this function as a placeholder in case mocking is introduced.
}
