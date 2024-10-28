package tests

import (
	"log"

	"github.com/joho/godotenv"
)

func loadEnvironmentVariables() {
	err := godotenv.Load("../.env.test")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}
