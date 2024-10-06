package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shibbirmcc/user-auth-and-permissions/models"
	"github.com/shibbirmcc/user-auth-and-permissions/utils"
	"gorm.io/gorm"
)

func LoginUser(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var input models.LoginRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := db.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if !utils.CheckPasswordHash(input.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	var userDetails models.UserDetail
	if err := db.Where("user_id = ?", user.ID).First(&userDetails).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user Id"})
		return
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.Email, userDetails)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
