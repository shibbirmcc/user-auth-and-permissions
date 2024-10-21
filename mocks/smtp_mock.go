package mocks

import (
	"net/smtp"

	"github.com/stretchr/testify/mock"
)

type SmtpMock struct {
	mock.Mock
}

func (m *SmtpMock) SendMail(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
	args := m.Called(addr, a, from, to, msg)
	return args.Error(0)
}
