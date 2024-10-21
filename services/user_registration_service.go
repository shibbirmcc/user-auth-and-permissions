package services

import (
	"errors"
	"fmt"

	"github.com/shibbirmcc/user-auth-and-permissions/models"
	"github.com/shibbirmcc/user-auth-and-permissions/utils"
)

type UserRegistrationService struct {
	dbService IDatabaseOperationService
}

func NewUserRegistrationService(dbService IDatabaseOperationService) *UserRegistrationService {
	return &UserRegistrationService{
		dbService: dbService,
	}
}

func (s *UserRegistrationService) RegisterUser(input models.UserRegitrationRequest) error {
	generatedPassword, hashedPassword, err := utils.GetRandomPasswordAndHash()
	if err != nil {
		return errors.New("Error while generating temporary password and hash")
	}
	fmt.Println("Generated Password:", generatedPassword)

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
