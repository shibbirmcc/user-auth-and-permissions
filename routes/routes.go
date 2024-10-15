package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shibbirmcc/user-auth-and-permissions/handlers"
)

func ConfigureRouteEndpoints(router *gin.Engine, userHandler *handlers.UserHandler) {
	router.POST("/auth/register", userHandler.RegisterUser)
}
