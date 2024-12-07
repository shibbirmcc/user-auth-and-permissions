package mocks

import "errors"

type MockPasswordDeliveryService struct {
	ShouldFail bool
}

func (m *MockPasswordDeliveryService) SendPassword(email, firstName, middleName, lastName, password string) error {
	if m.ShouldFail {
		return errors.New("mock error: failed to send password")
	}
	return nil
}
