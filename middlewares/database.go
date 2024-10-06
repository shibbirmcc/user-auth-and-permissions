package middlewares

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// InjectDBMiddleware injects the provided database instance into the Gin context.
// This middleware allows the database instance to be available in subsequent handlers
// by setting it in the context under the "db" key. Handlers can retrieve the database
// instance using `c.MustGet("db")` to perform database operations.
//
// Usage Example:
//
//	router := gin.Default()
//	db := setupDatabase() // Assume setupDatabase returns a *gorm.DB instance
//	router.Use(InjectDBMiddleware(db))
//
// Parameters:
// - db (*gorm.DB): The GORM database instance to be injected into the context.
//
// Context Key:
// - "db" (*gorm.DB): The key used to retrieve the database instance from the Gin context.
//
// This middleware does not modify the response or terminate the request; it simply injects
// the database instance and passes control to the next handler in the chain.
func InjectDBMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set the database instance in the context
		c.Set("db", db)

		// Process the request
		c.Next()
	}
}
