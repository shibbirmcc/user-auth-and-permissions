package mocks

import (
	"errors"
	"github.com/shibbirmcc/user-auth-and-permissions/models"
)

type MockPasswordDeliveryService struct {
	ShouldFail bool
}

func (m *MockPasswordDeliveryService) SendPassword(credentials models.UserCredentials) error {
	if m.ShouldFail {
		return errors.New("mock error: failed to send password")
	}
	return nil
}
