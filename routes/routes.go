package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shibbirmcc/user-auth-and-permissions/handlers"
)

func InitRoutes() *gin.Engine {
	router := gin.Default()

	// Auth routes
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/register", handlers.RegisterUser)
		authGroup.POST("/login", handlers.LoginUser)
	}

	// Other routes like permissions, roles, etc. can go here
	return router
}
