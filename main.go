package main

import (
	"log"
	"os"
	"path/filepath"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/shibbirmcc/user-auth-and-permissions/config"
	"github.com/shibbirmcc/user-auth-and-permissions/initializer"
)

func main() {
	absEnvPath, err := filepath.Abs(".env")
	if err != nil {
		log.Fatalf("Error finding absolute path of .env.test: %v", err)
	}
	config.LoadEnv(absEnvPath)
	db, err := config.GetDatabase()
	if err != nil {
		os.Exit(1)
	}
	initializer.ApplyMigrations(db, "migrations")

	passwordDeliveryService, err := initializer.InitializePasswordDeliveryService()
	if err != nil {
		log.Fatalf("Failed to initialize KafkaPasswordDeliveryService: %v", err)
	}
	userRegService, userLoginService := initializer.InitializeServices(db)
	userHandler := initializer.InitializeHandlers(userRegService, userLoginService, &passwordDeliveryService)

	router := initializer.SetupRouter(userHandler)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}
