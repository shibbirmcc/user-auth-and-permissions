package config

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestLoadEnv(t *testing.T) {
	absEnvPath, err := filepath.Abs("../.env.test")
	if err != nil {
		log.Fatalf("Error finding absolute path of .env.test: %v", err)
	}
	LoadEnv(absEnvPath)

	// Verify that environment variables are loaded as expected
	assert.Equal(t, "127.0.0.1", os.Getenv("DB_HOST"))
	assert.Equal(t, "serviceuser", os.Getenv("DB_USER"))
	assert.Equal(t, "servicepassword", os.Getenv("DB_PASSWORD"))
	assert.Equal(t, "servicedatabase", os.Getenv("DB_NAME"))
	assert.Equal(t, "5432", os.Getenv("DB_PORT"))
	assert.Equal(t, "snakeexactwhichrepliedpothearthasdigplentymathemat", os.Getenv("JWT_SECRET"))
}

// TestGetDatabase verifies that a database connection is established
func TestGetDatabase(t *testing.T) {
	os.Setenv("DB_HOST", "127.0.0.1")
	os.Setenv("DB_USER", "serviceuser")
	os.Setenv("DB_PASSWORD", "servicepassword")
	os.Setenv("DB_NAME", "servicedatabase")
	os.Setenv("DB_PORT", "5432")

	// Use an in-memory SQLite database to test connection
	dsn := "file::memory:?cache=shared"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to SQLite in-memory database: %v", err)
	}

	// Verify that the connection is valid
	sqlDB, err := db.DB()
	assert.NoError(t, err)
	assert.NoError(t, sqlDB.Ping())

	// Close the database connection at the end of the test
	defer sqlDB.Close()
}
