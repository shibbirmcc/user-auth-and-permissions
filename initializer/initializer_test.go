package initializer

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/shibbirmcc/user-auth-and-permissions/tests"
	"github.com/stretchr/testify/assert"
)

func TestInitializeServices(t *testing.T) {
	db, TeardownPostgresContainer := tests.SetupPostgresContainer()
	defer TeardownPostgresContainer()

	// Test InitializeServices function
	regService, loginService := InitializeServices(db)

	assert.NotNil(t, regService, "UserRegistrationService should not be nil")
	assert.NotNil(t, loginService, "UserLoginService should not be nil")
}

func TestInitializeHandlers(t *testing.T) {
	db, TeardownPostgresContainer := tests.SetupPostgresContainer()
	defer TeardownPostgresContainer()

	// Initialize services to pass into InitializeHandlers
	regService, loginService := InitializeServices(db)

	// Test InitializeHandlers function
	userHandler := InitializeHandlers(regService, loginService)
	assert.NotNil(t, userHandler, "UserHandler should not be nil")
}

func TestApplyMigrations(t *testing.T) {
	db, TeardownPostgresContainer := tests.SetupPostgresContainer()
	defer TeardownPostgresContainer()

	// Apply migrations and check for table existence
	ApplyMigrations(db, "../migrations")

	// Check for a specific table to verify migration success
	hasTable := db.Migrator().HasTable("users")
	assert.True(t, hasTable, "Expected users table to exist after migration")
}

func TestSetupRouter(t *testing.T) {
	db, TeardownPostgresContainer := tests.SetupPostgresContainer()
	defer TeardownPostgresContainer()
	regService, loginService := InitializeServices(db)
	userHandler := InitializeHandlers(regService, loginService)

	// Test SetupRouter function
	router := SetupRouter(userHandler)
	assert.NotNil(t, router, "Router should not be nil")

	// Test that the router has the expected routes
	w := performRequest(router, "POST", "/auth/login")
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// Helper function to perform requests in the router
func performRequest(r *gin.Engine, method, path string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
