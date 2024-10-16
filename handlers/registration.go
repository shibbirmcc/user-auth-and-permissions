package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shibbirmcc/user-auth-and-permissions/models"
)

func (h *UserHandler) RegisterUser(c *gin.Context) {
	var input models.UserRegitrationRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.userRegistrationService.RegisterUser(input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Registration successful, please confirm your email"})
}
