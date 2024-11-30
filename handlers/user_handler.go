package handlers

import "github.com/shibbirmcc/user-auth-and-permissions/services"

type UserHandler struct {
	userRegistrationService services.UserRegistrationService
	userLoginService        services.UserLoginService
	passwordDeliveryService services.PasswordDeliveryService
}

func NewUserHandler(regService services.UserRegistrationService,
	loginService services.UserLoginService,
	passwordDeliveryService services.PasswordDeliveryService) *UserHandler {
	return &UserHandler{
		userRegistrationService: regService,
		userLoginService:        loginService,
		passwordDeliveryService: passwordDeliveryService,
	}
}
