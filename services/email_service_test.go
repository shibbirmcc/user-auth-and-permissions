package services

import (
	"errors"
	"os"
	"testing"

	"github.com/shibbirmcc/user-auth-and-permissions/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestEmailService_SendMail_Success(t *testing.T) {
	os.Setenv("SMTP_USER", "test@example.com")
	os.Setenv("SMTP_HOST", "smtp.example.com")
	os.Setenv("SMTP_PORT", "587")
	os.Setenv("SMTP_PASSWORD", "password")

	emailService := NewEmailService()

	mockSmtp := new(mocks.SmtpMock)
	smtpSendMail = mockSmtp.SendMail

	// Define expectations for the mock.
	mockSmtp.On("SendMail", "smtp.example.com:587", mock.Anything, "test@example.com", []string{"receiver@example.com"}, []byte("Hello, World!")).
		Return(nil)

	// Call the method.
	err := emailService.SendMail("receiver@example.com", "Receiver Name", "Hello, World!")

	// Validate the result.
	assert.NoError(t, err)
	mockSmtp.AssertExpectations(t)
}

func TestEmailService_SendMail_Error(t *testing.T) {
	os.Setenv("SMTP_USER", "test@example.com")
	os.Setenv("SMTP_HOST", "smtp.example.com")
	os.Setenv("SMTP_PORT", "587")
	os.Setenv("SMTP_PASSWORD", "password")

	emailService := NewEmailService()

	mockSmtp := new(mocks.SmtpMock)
	smtpSendMail = mockSmtp.SendMail

	// Define expectations for the mock.
	mockSmtp.On("SendMail", "smtp.example.com:587", mock.Anything, "test@example.com", []string{"receiver@example.com"}, []byte("Hello, World!")).
		Return(errors.New("failed to send email"))

	// Call the method.
	err := emailService.SendMail("receiver@example.com", "Receiver Name", "Hello, World!")

	// Validate the result.
	assert.Error(t, err)
	assert.Equal(t, "failed to send email", err.Error())
	mockSmtp.AssertExpectations(t)
}
