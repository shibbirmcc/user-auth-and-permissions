package services

import (
	"errors"
	"fmt"

	"github.com/shibbirmcc/user-auth-and-permissions/models"
	"github.com/shibbirmcc/user-auth-and-permissions/utils"
)

type UserLoginService struct {
	dbService IDatabaseOperationService
}

func NewUserLoginService(dbService IDatabaseOperationService) *UserLoginService {
	return &UserLoginService{
		dbService: dbService,
	}
}

func (s *UserLoginService) Login(input models.LoginRequest) (string, error) {
	user, err := s.dbService.FindUserByEmail(input.Email)
	if err != nil {
		return "", errors.New("Invalid credentials")
	}

	if !utils.CheckPasswordHash(input.Password, user.Password) {
		return "", errors.New("Invalid credentials")
	}

	userDetails, err := s.dbService.FindUserDetailsByUserID(user.ID)
	if err != nil {
		return "", errors.New("Invalid user Id")
	}

	token, err := utils.GenerateJWT(user.Email, *userDetails)
	if err != nil {
		fmt.Println("Error:", err)
		return "", errors.New("Could not generate token")
	}

	return token, nil
}
