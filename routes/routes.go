package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shibbirmcc/user-auth-and-permissions/handlers"
	"github.com/shibbirmcc/user-auth-and-permissions/middlewares"
	"gorm.io/gorm"
)

func InitRoutes(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	// router.Use(middlewares.InjectDBMiddleware(db))
	router.Use(middlewares.CORSMiddleware())

	router.POST("/auth/register", handlers.RegisterUser)
	router.POST("/auth/login", handlers.LoginUser)

	return router
}
