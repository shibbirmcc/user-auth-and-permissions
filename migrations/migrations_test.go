package migrations

import (
	"os"
	"testing"

	"github.com/shibbirmcc/user-auth-and-permissions/tests"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestMain(m *testing.M) {
	tests.SetupPostgresContainer()
	code := m.Run()
	tests.TeardownPostgresContainer()
	os.Exit(code)
}

func TestRunMigrations(t *testing.T) {
	var db *gorm.DB
	db, _ = tests.GetGormDBFromSQLDB(tests.DB)
	RunMigrations(db, "./")

	var exists bool
	err := tests.DB.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'users');").Scan(&exists)
	assert.NoError(t, err)
	assert.True(t, exists, "Expected table 'users' to exist after migration")

	err = tests.DB.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'user_details');").Scan(&exists)
	assert.NoError(t, err)
	assert.True(t, exists, "Expected table 'user_details' to exist after migration")

	err = tests.DB.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'roles');").Scan(&exists)
	assert.NoError(t, err)
	assert.True(t, exists, "Expected table 'roles' to exist after migration")

	err = tests.DB.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'permissions');").Scan(&exists)
	assert.NoError(t, err)
	assert.True(t, exists, "Expected table 'permissions' to exist after migration")

	err = tests.DB.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'role_permissions');").Scan(&exists)
	assert.NoError(t, err)
	assert.True(t, exists, "Expected table 'role_permissions' to exist after migration")
}
