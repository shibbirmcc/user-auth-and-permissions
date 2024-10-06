package utils_test

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestLoadEnvVariables(t *testing.T) {
	// Load the environment variables from the .env file
	err := godotenv.Load("../../.env") // Adjust the relative path to .env file if necessary
	if err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	// Define the expected environment variables
	expectedEnvVars := map[string]string{
		"DB_HOST":     "127.0.0.1",
		"DB_USER":     "serviceuser",
		"DB_PASSWORD": "servicepassword",
		"DB_NAME":     "servicedatabase",
		"DB_PORT":     "5432",
	}

	// Loop through the expected environment variables and check their values
	for key, expectedValue := range expectedEnvVars {
		actualValue := os.Getenv(key)
		if actualValue != expectedValue {
			t.Errorf("Expected %s to be %s, but got %s", key, expectedValue, actualValue)
		}
	}
}

func TestMissingEnvVariable(t *testing.T) {
	// Load the environment variables from the .env file
	err := godotenv.Load("../../.env") // Adjust the relative path to .env file if necessary
	if err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	// Test for an environment variable that doesn't exist
	nonExistentEnv := os.Getenv("NON_EXISTENT_ENV")
	if nonExistentEnv != "" {
		t.Errorf("Expected NON_EXISTENT_ENV to be empty, but got %s", nonExistentEnv)
	}
}
