package services

import (
	"os"
	"testing"

	"github.com/shibbirmcc/user-auth-and-permissions/migrations"
	"github.com/shibbirmcc/user-auth-and-permissions/tests"
	"gorm.io/gorm"
)

var (
	DBOperationService DatabaseOperationService
)

/*
This TestMain method will be executed before starting to execute tests of this package
*/
func TestMain(m *testing.M) {
	tests.SetupPostgresContainer()
	var db *gorm.DB
	db, _ = tests.GetGormDBFromSQLDB(tests.DB)
	migrations.RunMigrations(db, "../migrations")
	DBOperationService = *NewDatabaseOperationService(db)

	code := m.Run()
	tests.TeardownPostgresContainer()
	os.Exit(code)
}
