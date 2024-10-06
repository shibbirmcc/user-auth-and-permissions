package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shibbirmcc/user-auth-and-permissions/models"
	"github.com/shibbirmcc/user-auth-and-permissions/utils"
	"gorm.io/gorm"
)

func RegisterUser(c *gin.Context) {
	db, _ := c.MustGet("db").(*gorm.DB)
	var input models.UserRegitrationRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while hashing password"})
		return
	}

	user := models.User{Email: input.Email, Password: hashedPassword}
	userDetail := models.UserDetail{
		FirstName:  input.FirstName,
		MiddleName: input.MiddleName,
		LastName:   input.LastName,
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}
		userDetail.UserID = user.ID
		return tx.Create(&userDetail).Error
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while registering user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Registration successful, please confirm your email"})
}
