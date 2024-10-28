package mocks

import (
	"github.com/stretchr/testify/mock"
)

type MockUserHandler struct {
	userRegistrationService mock.Mock
	userLoginService        mock.Mock
}
