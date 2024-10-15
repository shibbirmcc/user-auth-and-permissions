package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shibbirmcc/user-auth-and-permissions/handlers"
)

func ConfigureRouteEndpoints(router *gin.Engine) {
	router.POST("/auth/register", handlers.RegisterUser)
	router.POST("/auth/login", handlers.LoginUser)
}
