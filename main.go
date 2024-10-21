package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin" // Aliasing to avoid conflict
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"github.com/shibbirmcc/user-auth-and-permissions/handlers"
	"github.com/shibbirmcc/user-auth-and-permissions/middlewares"
	"github.com/shibbirmcc/user-auth-and-permissions/migrations"
	"github.com/shibbirmcc/user-auth-and-permissions/routes"
	"github.com/shibbirmcc/user-auth-and-permissions/services"
	gormPostgres "gorm.io/driver/postgres" // Aliasing to avoid conflict
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Set up the database connection using Gorm
	var db *gorm.DB
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
	db, err = gorm.Open(gormPostgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to database")
	}

	// Run database migrations
	migrations.RunMigrations()

	// Set up services
	databaseOperationService := services.NewDatabaseOperationService(db)
	emailService := services.NewEmailService()
	userRegistrationService := services.NewUserRegistrationService(databaseOperationService, emailService)
	userLoginService := services.NewUserLoginService(databaseOperationService)

	// Set up handlers
	userHandler := handlers.NewUserHandler(*userRegistrationService, *userLoginService)

	// Set up the Gin router
	router := gin.Default()
	router.Use(middlewares.CORSMiddleware()) // Add middleware for CORS

	// Set up route handlers (using the userHandler)
	routes.ConfigureRouteEndpoints(router, userHandler)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}
