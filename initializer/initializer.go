package initializer

import (
	"github.com/gin-gonic/gin"
	"github.com/shibbirmcc/user-auth-and-permissions/handlers"
	"github.com/shibbirmcc/user-auth-and-permissions/middlewares"
	"github.com/shibbirmcc/user-auth-and-permissions/migrations"
	"github.com/shibbirmcc/user-auth-and-permissions/routes"
	"github.com/shibbirmcc/user-auth-and-permissions/services"
	"gorm.io/gorm"
)

func InitializeServices(db *gorm.DB) (*services.UserRegistrationService, *services.UserLoginService) {
	databaseOperationService := services.NewDatabaseOperationService(db)
	userRegistrationService := services.NewUserRegistrationService(databaseOperationService)
	userLoginService := services.NewUserLoginService(databaseOperationService)
	return userRegistrationService, userLoginService
}

func InitializeHandlers(regService *services.UserRegistrationService, loginService *services.UserLoginService) *handlers.UserHandler {
	return handlers.NewUserHandler(*regService, *loginService)
}

func ApplyMigrations(db *gorm.DB) {
	migrations.RunMigrations(db, "migrations")
}

func SetupRouter(userHandler *handlers.UserHandler) *gin.Engine {
	router := gin.Default()
	router.Use(middlewares.CORSMiddleware())            // Add CORS middleware
	routes.ConfigureRouteEndpoints(router, userHandler) // Set up route handlers
	return router
}
