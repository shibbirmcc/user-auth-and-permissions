package initializer

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shibbirmcc/user-auth-and-permissions/handlers"
	"github.com/shibbirmcc/user-auth-and-permissions/middlewares"
	"github.com/shibbirmcc/user-auth-and-permissions/migrations"
	"github.com/shibbirmcc/user-auth-and-permissions/routes"
	"github.com/shibbirmcc/user-auth-and-permissions/services"
	"gorm.io/gorm"
	"os"
)

func InitializeServices(db *gorm.DB) (*services.UserRegistrationService, *services.UserLoginService) {
	databaseOperationService := services.NewDatabaseOperationService(db)
	passwordDeliveryService, _ := InitializePasswordDeliveryService()
	userRegistrationService := services.NewUserRegistrationService(passwordDeliveryService, databaseOperationService)
	userLoginService := services.NewUserLoginService(databaseOperationService)
	return userRegistrationService, userLoginService
}

func InitializeHandlers(regService *services.UserRegistrationService, loginService *services.UserLoginService) *handlers.UserHandler {
	return handlers.NewUserHandler(*regService, *loginService)
}

func ApplyMigrations(db *gorm.DB, migrationDirectory string) {
	migrations.RunMigrations(db, migrationDirectory)
}

func SetupRouter(userHandler *handlers.UserHandler) *gin.Engine {
	router := gin.Default()
	router.Use(middlewares.CORSMiddleware())            // Add CORS middleware
	routes.ConfigureRouteEndpoints(router, userHandler) // Set up route handlers
	return router
}

func InitializePasswordDeliveryService() (services.PasswordDeliveryService, error) {
	deliveryType := os.Getenv("PASSWORD_DELIVERY_TYPE")
	switch services.PasswordDeliveryType(deliveryType) {
	case services.KAFKA_TOPIC:
		return services.NewKafkaPasswordDeliveryService()
	default:
		return nil, fmt.Errorf("unsupported password delivery type: %s", deliveryType)
	}
}
