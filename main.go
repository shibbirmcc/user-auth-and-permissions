package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4" // Aliasing to avoid conflict
	migratePostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/shibbirmcc/user-auth-and-permissions/middlewares"
	"github.com/shibbirmcc/user-auth-and-permissions/routes"

	"github.com/joho/godotenv"
	gormPostgres "gorm.io/driver/postgres" // Aliasing to avoid conflict
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func runMigrations() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	// Open a connection to the database
	db, err := sql.Open("postgres", dsn) // Use the standard "database/sql" package to open a connection
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
		"file://migrations", // Path to your migrations folder
		"postgres",          // The name of the database driver
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

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var db *gorm.DB
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
	db, err = gorm.Open(gormPostgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to database")
	}

	runMigrations()

	router := routes.InitRoutes()

	router.Use(middlewares.InjectDBMiddleware(db))
	router.Use(middlewares.CORSMiddleware())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}
