package migrations

import (
	"errors"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	migratePostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/gorm"
)

// RunMigrations applies all database migrations.
func RunMigrations(gormDB *gorm.DB, migrationSirectory string) error {
	db, err := gormDB.DB()
	if err != nil {
		log.Printf("Failed to connect to database for migrations: %v", err)
		return errors.New("Failed to connect to database for migrations")
	}

	// Ensure the connection is working
	if err = db.Ping(); err != nil {
		log.Printf("logging Failed to ping database: %v", err)
		return errors.New("Failed to ping database")
	}

	// Create a Postgres driver instance for migrations
	driver, err := migratePostgres.WithInstance(db, &migratePostgres.Config{})
	if err != nil {
		log.Printf("Failed to create migration driver: %v", err)
		return errors.New("Failed to create migration driver")
	}

	// Create a new migration instance
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationSirectory),
		"postgres",
		driver,
	)
	if err != nil {
		log.Printf("Failed to initialize migrations: %v", err)
		return errors.New("Failed to initialize migrations")
	}

	// Apply all migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Printf("Failed to apply migrations: %v", err)
		return errors.New("Failed to apply migrations")
	}

	log.Println("Database migrations applied successfully")
	return nil
}
