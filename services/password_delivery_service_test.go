package services

import (
	"github.com/shibbirmcc/user-auth-and-permissions/mocks"
	"github.com/shibbirmcc/user-auth-and-permissions/models"
	"testing"
)

func TestPasswordDeliveryType_String(t *testing.T) {
	tests := []struct {
		name     string
		input    PasswordDeliveryType
		expected string
	}{
		{"PostgreSQL type", POSTGRESQL, "POSTGRESQL"},
		{"Redis type", REDIS, "REDIS"},
		{"Kafka topic type", KAFKA_TOPIC, "KAFKA_TOPIC"},
		{"Unknown type", PasswordDeliveryType("UNKNOWN"), ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.String()
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestPasswordDeliveryService_SendPassword(t *testing.T) {
	tests := []struct {
		name          string
		service       PasswordDeliveryService
		email         string
		firstName     string
		middleName    string
		lastName      string
		password      string
		expectedError bool
	}{
		{
			name:          "Successful password delivery",
			service:       &mocks.MockPasswordDeliveryService{ShouldFail: false},
			email:         "test@example.com",
			firstName:     "John",
			middleName:    "M",
			lastName:      "Doe",
			password:      "securePassword123",
			expectedError: false,
		},
		{
			name:          "Failed password delivery",
			service:       &mocks.MockPasswordDeliveryService{ShouldFail: true},
			email:         "test@example.com",
			firstName:     "John",
			middleName:    "M",
			lastName:      "Doe",
			password:      "securePassword123",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			credentials := models.UserCredentials{
				Email:      tt.email,
				FirstName:  tt.firstName,
				MiddleName: tt.middleName,
				LastName:   tt.lastName,
				Password:   tt.password,
			}
			err := tt.service.SendPassword(credentials)
			if (err != nil) != tt.expectedError {
				t.Errorf("Expected error: %v, got: %v", tt.expectedError, err != nil)
			}
		})
	}
}
