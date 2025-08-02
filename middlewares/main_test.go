package middlewares

import (
	"os"
	"testing"

	"github.com/shibbirmcc/user-auth-and-permissions/tests"
)

/*
This TestMain method will be executed before starting to execute tests of this package
*/
func TestMain(m *testing.M) {
	tests.LoadEnvironmentVariables()
	code := m.Run()
	os.Exit(code)
}
