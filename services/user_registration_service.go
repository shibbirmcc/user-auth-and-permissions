package services

import (
	"errors"

	"github.com/shibbirmcc/user-auth-and-permissions/models"
	"github.com/shibbirmcc/user-auth-and-permissions/utils"
)

type IDatabaseOperationService interface {
	CreateUser(user *models.User, userDetail *models.UserDetail) error
}

type UserRegistrationService struct {
	dbService IDatabaseOperationService
}

func NewUserRegistrationService(dbService IDatabaseOperationService) *UserRegistrationService {
	return &UserRegistrationService{
		dbService: dbService,
	}
}

func (s *UserRegistrationService) RegisterUser(input models.UserRegitrationRequest) error {
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return errors.New("error while hashing password")
	}

	user := models.User{Email: input.Email, Password: hashedPassword}
	userDetail := models.UserDetail{
		FirstName:  input.FirstName,
		MiddleName: input.MiddleName,
		LastName:   input.LastName,
	}

	if err := s.dbService.CreateUser(&user, &userDetail); err != nil {
		return errors.New("error while registering user")
	}
	return nil
}
