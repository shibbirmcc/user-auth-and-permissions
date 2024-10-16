package handlers

import "github.com/shibbirmcc/user-auth-and-permissions/services"

type UserHandler struct {
	userRegistrationService services.UserRegistrationService
	userLoginService        services.UserLoginService
}

func NewUserHandler(regService services.UserRegistrationService, loginService services.UserLoginService) *UserHandler {
	return &UserHandler{
		userRegistrationService: regService,
		userLoginService:        loginService,
	}
}
