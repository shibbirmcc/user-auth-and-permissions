package migrations

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/shibbirmcc/user-auth-and-permissions/tests"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// func TestMain(m *testing.M) {
// 	tests.SetupPostgresContainer()
// 	code := m.Run()
// 	tests.TeardownPostgresContainer()
// 	os.Exit(code)
// }

func TestRunMigrations(t *testing.T) {
	gormDB, TeardownPostgresContainer := tests.SetupPostgresContainer()
	defer TeardownPostgresContainer()
	RunMigrations(gormDB, "./")

	sqlDB, err := gormDB.DB()
	if err != nil {
		log.Printf("Failed to connect to database for migrations: %v", err)
	}

	var exists bool
	err = sqlDB.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'users');").Scan(&exists)
	assert.NoError(t, err)
	assert.True(t, exists, "Expected table 'users' to exist after migration")

	err = sqlDB.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'user_details');").Scan(&exists)
	assert.NoError(t, err)
	assert.True(t, exists, "Expected table 'user_details' to exist after migration")

	err = sqlDB.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'roles');").Scan(&exists)
	assert.NoError(t, err)
	assert.True(t, exists, "Expected table 'roles' to exist after migration")

	err = sqlDB.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'permissions');").Scan(&exists)
	assert.NoError(t, err)
	assert.True(t, exists, "Expected table 'permissions' to exist after migration")

	err = sqlDB.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'role_permissions');").Scan(&exists)
	assert.NoError(t, err)
	assert.True(t, exists, "Expected table 'role_permissions' to exist after migration")
}

func TestRunMigrations_DBConnectionFailure(t *testing.T) {
	// Use an unsupported driver to create an invalid *gorm.DB instance
	invalidDB, err := gorm.Open(nil, &gorm.Config{NamingStrategy: schema.NamingStrategy{}})
	assert.NoError(t, err)

	// Capture log output
	var logBuffer bytes.Buffer
	log.SetOutput(&logBuffer)
	defer func() {
		log.SetOutput(os.Stderr) // Restore log output
	}()

	// Call RunMigrations with the invalid *gorm.DB instance
	err = RunMigrations(invalidDB, "./")

	// Assert that an error is returned and check the error message
	assert.Error(t, err)
	assert.EqualError(t, err, "Failed to connect to database for migrations")

	// Check log output to confirm that the correct message was logged
	assert.Contains(t, logBuffer.String(), "Failed to connect to database for migrations")
}

func TestRunMigrations_PingFailure(t *testing.T) {
	tests.LoadEnvironmentVariables()

	// Set up in-memory SQLite and close the connection to simulate Ping failure
	gormDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// Close the underlying database to simulate a Ping failure
	db, err := gormDB.DB()
	assert.NoError(t, err)
	db.Close() // Close to force Ping to fail

	err = RunMigrations(gormDB, "./")
	assert.Error(t, err)
	assert.EqualError(t, err, "Failed to ping database")
}

func TestRunMigrations_DriverCreationFailure(t *testing.T) {
	// Set up an in-memory SQLite database, which is incompatible with the PostgreSQL driver
	gormDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// Call RunMigrations with the incompatible SQLite database
	err = RunMigrations(gormDB, "./")

	// Assert that it returns the expected error message due to driver creation failure
	assert.Error(t, err)
	assert.EqualError(t, err, "Failed to create migration driver")
}

func TestRunMigrations_MigrationInitializationFailure(t *testing.T) {
	gormDB, TeardownPostgresContainer := tests.SetupPostgresContainer()
	defer TeardownPostgresContainer()

	// Use an invalid migrations directory to simulate initialization failure
	var err error
	err = RunMigrations(gormDB, "migrations")
	assert.Error(t, err)
	assert.EqualError(t, err, "Failed to initialize migrations")
}

func TestRunMigrations_MigrationApplicationFailure(t *testing.T) {
	gormDB, TeardownPostgresContainer := tests.SetupPostgresContainer()
	defer TeardownPostgresContainer()

	// Create a temporary directory for migration files
	tmpDir, err := ioutil.TempDir("", "migrations")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir) // Clean up after the test

	// Create an invalid migration file
	migrationFile := tmpDir + "/0001_invalid_migration.up.sql"
	err = ioutil.WriteFile(migrationFile, []byte("INVALID SQL SYNTAX;"), 0644)
	assert.NoError(t, err)

	// Run migrations using the directory with the invalid migration file
	err = RunMigrations(gormDB, tmpDir)

	// Assert that an error occurs due to the invalid migration
	assert.Error(t, err)
	assert.EqualError(t, err, "Failed to apply migrations")
}
