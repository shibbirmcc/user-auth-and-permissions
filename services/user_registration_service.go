package services

import (
	"errors"

	"github.com/shibbirmcc/user-auth-and-permissions/models"
	"github.com/shibbirmcc/user-auth-and-permissions/utils"
)

type UserRegistrationService struct {
	dbService    IDatabaseOperationService
	emailService IEmailService
}

func NewUserRegistrationService(dbService IDatabaseOperationService, emailService IEmailService) *UserRegistrationService {
	return &UserRegistrationService{
		dbService:    dbService,
		emailService: emailService,
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
