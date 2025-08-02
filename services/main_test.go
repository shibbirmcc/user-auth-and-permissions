package services

import (
	"os"
	"testing"

	"github.com/shibbirmcc/user-auth-and-permissions/migrations"
	"github.com/shibbirmcc/user-auth-and-permissions/tests"
)

var (
	DBOperationService DatabaseOperationService
)

/*
This TestMain method will be executed before starting to execute tests of this package
*/
func TestMain(m *testing.M) {
	db, TeardownPostgresContainer := tests.SetupPostgresContainer()
	migrations.RunMigrations(db, "../migrations")
	DBOperationService = *NewDatabaseOperationService(db)

	code := m.Run()
	TeardownPostgresContainer()
	os.Exit(code)
}
