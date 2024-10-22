package migrations

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var db *sql.DB
var postgresC testcontainers.Container

func waitForDBConnection(db *sql.DB) error {
	maxAttempts := 5
	for i := 0; i < maxAttempts; i++ {
		if err := db.Ping(); err == nil {
			return nil
		}
		time.Sleep(2 * time.Second) // Wait before retrying
	}
	return fmt.Errorf("could not establish database connection after %d attempts", maxAttempts)
}

/*
This TestMain method will be executed before starting to execute tests of this package
*/
func TestMain(m *testing.M) {
	ctx := context.Background()

	env_err := godotenv.Load("../.env.test")
	if env_err != nil {
		log.Fatalf("Error loading .env file: %v", env_err)
	}

	// Create a PostgreSQL container
	req := testcontainers.ContainerRequest{
		Image:        "postgres:13",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     os.Getenv("DB_USER"),
			"POSTGRES_PASSWORD": os.Getenv("DB_PASSWORD"),
			"POSTGRES_DB":       os.Getenv("DB_NAME"),
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").WithStartupTimeout(60 * time.Second),
	}

	var err error
	postgresC, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		fmt.Printf("Failed to start container: %v\n", err)
		os.Exit(1)
	}

	host, _ := postgresC.Host(ctx)
	port, _ := postgresC.MappedPort(ctx, "5432")

	fmt.Printf("Postgres TestContainer started at %s:%s\n", host, port.Port())

	// Set environment variables for the migration
	// os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME")
	os.Setenv("DB_HOST", host)
	os.Setenv("DB_PORT", port.Port())

	cwd, _ := os.Getwd()
	fmt.Printf("Current working directory: %s\n", cwd)
	migrationPath := filepath.Join(cwd, "../")
	if err := os.Chdir(migrationPath); err != nil {
		log.Fatalf("Failed to change working directory: %v", err)
	}
	fmt.Printf("Changed working directory to: %s\n", migrationPath)

	// Connect to the database for verification in tests
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
		os.Exit(1)
	}

	if err := waitForDBConnection(db); err != nil {
		fmt.Printf("Database not ready: %v\n", err)
		os.Exit(1)
	}

	// Run the tests
	code := m.Run()

	// Cleanup resources
	db.Close()
	postgresC.Terminate(ctx)

	// Exit with the result code from `m.Run()`
	os.Exit(code)
}

func TestRunMigrations(t *testing.T) {
	// Run the migrations
	RunMigrations()

	// Verify if the migrations ran successfully (e.g., check if a table exists)
	var exists bool
	err := db.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'users');").Scan(&exists)
	assert.NoError(t, err)
	assert.True(t, exists, "Expected table 'users' to exist after migration")

	err = db.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'user_details');").Scan(&exists)
	assert.NoError(t, err)
	assert.True(t, exists, "Expected table 'user_details' to exist after migration")

	err = db.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'roles');").Scan(&exists)
	assert.NoError(t, err)
	assert.True(t, exists, "Expected table 'roles' to exist after migration")

	err = db.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'permissions');").Scan(&exists)
	assert.NoError(t, err)
	assert.True(t, exists, "Expected table 'permissions' to exist after migration")

	err = db.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'role_permissions');").Scan(&exists)
	assert.NoError(t, err)
	assert.True(t, exists, "Expected table 'role_permissions' to exist after migration")
}
