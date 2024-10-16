package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shibbirmcc/user-auth-and-permissions/models"
)

func (h *UserHandler) LoginUser(c *gin.Context) {
	var input models.LoginRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.userLoginService.Login(input)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
