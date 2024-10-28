package migrations

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	migratePostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/gorm"
)

// RunMigrations applies all database migrations.
func RunMigrations(gormDB *gorm.DB, migrationSirectory string) {
	db, err := gormDB.DB()
	if err != nil {
		log.Fatalf("Failed to connect to database for migrations: %v", err)
	}

	// Ensure the connection is working
	if err = db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Create a Postgres driver instance for migrations
	driver, err := migratePostgres.WithInstance(db, &migratePostgres.Config{})
	if err != nil {
		log.Fatalf("Failed to create migration driver: %v", err)
	}

	// Create a new migration instance
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationSirectory),
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatalf("Failed to initialize migrations: %v", err)
	}

	// Apply all migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to apply migrations: %v", err)
	}

	log.Println("Database migrations applied successfully")
}
