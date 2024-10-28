package tests

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

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

func GetGormDBFromSQLDB(sqlDB *sql.DB) (*gorm.DB, error) {
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return gormDB, nil
}

func DeleteTestData(sqlDB *sql.DB) {
	err := sqlDB.QueryRow("DELETE FROM user_details;")
	if err != nil {
		fmt.Printf("Error delete rows from user_details: %v\n", err)
	}
	err = sqlDB.QueryRow("DELETE FROM users;")
	if err != nil {
		fmt.Printf("Error delete rows from user_details: %v\n", err)
	}
}

func SetupPostgresContainer() (*gorm.DB, func()) {
	LoadEnvironmentVariables()

	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:14",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     os.Getenv("DB_USER"),
			"POSTGRES_PASSWORD": os.Getenv("DB_PASSWORD"),
			"POSTGRES_DB":       os.Getenv("DB_NAME"),
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").WithStartupTimeout(60 * time.Second),
	}

	var err error
	var postgresContainer testcontainers.Container
	var sqlDB *sql.DB

	postgresContainer, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatalf("Could not start postgres container: %s", err)
	}

	host, err := postgresContainer.Host(ctx)
	if err != nil {
		log.Fatalf("Could not get host: %s", err)
	}
	port, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		log.Fatalf("Could not get port: %s", err)
	}

	// Override environment variables for connecting to the testcontainer
	os.Setenv("DB_HOST", host)
	os.Setenv("DB_PORT", port.Port())

	// Connect to the database for verification in tests
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	sqlDB, err = sql.Open("postgres", dsn)
	if err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
		os.Exit(1)
	}

	if err := waitForDBConnection(sqlDB); err != nil {
		fmt.Printf("Testcontainer Database is not ready: %v\n", err)
		os.Exit(1)
	}

	// Return *gorm.DB instance based on *sql.DB
	gormDB, err := GetGormDBFromSQLDB(sqlDB)
	if err != nil {
		log.Fatalf("Could not get Gorm DB: %s", err)
	}

	TeardownPostgresContainer := func() {
		DeleteTestData(sqlDB)
		if sqlDB != nil {
			sqlDB.Close()
		}
		if postgresContainer != nil {
			ctx := context.Background()
			err := postgresContainer.Terminate(ctx)
			if err != nil {
				log.Fatalf("Could not terminate container: %s", err)
			}
		}
	}

	return gormDB, TeardownPostgresContainer
}
