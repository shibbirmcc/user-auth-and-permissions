package models

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	Email      string `json:"email"`
	UserID     uint   `json:"user_id"`
	FirstName  string `json:"first_name"`
	MiddleName string `json:"middle_name"`
	LastName   string `json:"last_name"`
	jwt.RegisteredClaims
}
