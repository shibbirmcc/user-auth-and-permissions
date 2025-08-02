package handlers

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

/*
This TestMain method will be executed before starting to execute tests of this package
*/
func TestMain(m *testing.M) {
	err := godotenv.Load("../.env.test")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Run the tests
	code := m.Run()

	// Exit with the appropriate code
	os.Exit(code)
}
