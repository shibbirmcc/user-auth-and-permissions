package main

import (
	"log"
	"os"
	"path/filepath"

	// Aliasing to avoid conflict
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/shibbirmcc/user-auth-and-permissions/config"
	"github.com/shibbirmcc/user-auth-and-permissions/initializer"
	// Aliasing to avoid conflict
)

func main() {
	absEnvPath, err := filepath.Abs(".env")
	if err != nil {
		log.Fatalf("Error finding absolute path of .env.test: %v", err)
	}
	config.LoadEnv(absEnvPath)
	db := config.GetDatabase()
	initializer.ApplyMigrations(db, "migrations")

	userRegService, userLoginService := initializer.InitializeServices(db)
	userHandler := initializer.InitializeHandlers(userRegService, userLoginService)
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
