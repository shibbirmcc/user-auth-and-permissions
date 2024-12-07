package config

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/shibbirmcc/user-auth-and-permissions/tests"
	"github.com/stretchr/testify/assert"
)

func TestLoadEnv(t *testing.T) {
	absEnvPath, err := filepath.Abs("../.env.test")
	if err != nil {
		log.Fatalf("Error finding absolute path of .env.test: %v", err)
	}
	err = LoadEnv(absEnvPath)
	assert.Nil(t, err)

	// Verify that environment variables are loaded as expected
	assert.Equal(t, "127.0.0.1", os.Getenv("DB_HOST"))
	assert.Equal(t, "serviceuser", os.Getenv("DB_USER"))
	assert.Equal(t, "servicepassword", os.Getenv("DB_PASSWORD"))
	assert.Equal(t, "servicedatabase", os.Getenv("DB_NAME"))
	assert.Equal(t, "5432", os.Getenv("DB_PORT"))
	assert.Equal(t, "snakeexactwhichrepliedpothearthasdigplentymathemat", os.Getenv("JWT_SECRET"))
}

func TestLoadEnv_Error_Loading_EnvFile(t *testing.T) {
	err := LoadEnv("unrecognized_path/.env")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Error loading .env file at")
}

func TestGetDatabase_Success(t *testing.T) {
	_, TeardownPostgresContainer := tests.SetupPostgresContainer()
	defer TeardownPostgresContainer()

	// Use an in-memory SQLite database for testing connection
	db, err := GetDatabase()
	defer func() {
		if db != nil {
			sqlDB, _ := db.DB()
			sqlDB.Close()
		}
	}()
	assert.NoError(t, err)
	assert.NotNil(t, db, "Database should be successfully connected")
}

func TestGetDatabase_Failure(t *testing.T) {
	// Clear environment variables to simulate missing connection information
	os.Setenv("DB_HOST", "")
	os.Setenv("DB_USER", "")
	os.Setenv("DB_PASSWORD", "")
	os.Setenv("DB_NAME", "")
	os.Setenv("DB_PORT", "")

	// Attempt to get a database connection
	db, err := GetDatabase()

	// Assertions
	assert.Error(t, err, "Expected an error due to missing environment variables")
	assert.Nil(t, db, "Database should not be connected on failure")
}
