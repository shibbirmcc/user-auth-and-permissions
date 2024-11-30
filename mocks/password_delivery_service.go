package mocks

import "github.com/stretchr/testify/mock"

type MockPasswordDeliveryService struct {
	mock.Mock
}

func (m *MockPasswordDeliveryService) SendPassword(email, firstName, middleName, lastName, password string) error {
	args := m.Called(email, firstName, middleName, lastName, password)
	return args.Error(0)
}
