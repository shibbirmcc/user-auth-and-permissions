package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/shibbirmcc/user-auth-and-permissions/models"
)

func GenerateJWT(email string, userDetails models.UserDetail) (string, error) {
	if email == "" {
		return "", errors.New("email cannot be empty")
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return "", errors.New("JWT_SECRET is missing")
	}
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &models.Claims{
		Email:      email,
		UserID:     userDetails.UserID,
		FirstName:  userDetails.FirstName,
		MiddleName: userDetails.MiddleName,
		LastName:   userDetails.LastName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}
