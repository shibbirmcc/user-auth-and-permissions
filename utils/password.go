package utils

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	if password == "" {
		return "", errors.New("password cannot be empty")
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GetRandomPasswordAndHash() (string, string, error) {
	password, err := GenerateRandomPassword(12) // for example, a 12 character password
	if err != nil {
		return "", "", fmt.Errorf("failed to generate password: %w", err)
	}
	// Hash the generated password
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return "", "", fmt.Errorf("error hashing password: %w", err)
	}
	return password, hashedPassword, nil
}

func GenerateRandomPassword(length int) (string, error) {
	if length <= 0 {
		return "", errors.New("Password length must be greater than zero")
	}
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var password strings.Builder
	charsetSize := big.NewInt(int64(len(charset)))
	for i := 0; i < length; i++ {
		charIndex, err := rand.Int(rand.Reader, charsetSize)
		if err != nil {
			return "", err
		}
		password.WriteByte(charset[charIndex.Int64()])
	}
	return password.String(), nil
}
