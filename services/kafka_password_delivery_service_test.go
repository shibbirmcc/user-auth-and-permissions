package services

import (
	"encoding/json"
	"errors"
	"github.com/segmentio/kafka-go"
	"github.com/shibbirmcc/user-auth-and-permissions/mocks"
	"github.com/shibbirmcc/user-auth-and-permissions/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestKafkaPasswordDeliveryService_SendPassword_Success(t *testing.T) {
	mockProducer := new(mocks.MockProducer)
	service := &KafkaPasswordDeliveryService{
		Producer: mockProducer,
		Topic:    "test-topic",
	}

	mockProducer.On("WriteMessages", mock.Anything, mock.MatchedBy(func(msgs []kafka.Message) bool {
		if len(msgs) != 1 {
			return false
		}
		var credentials models.UserCredentials
		err := json.Unmarshal(msgs[0].Value, &credentials)
		return err == nil &&
			credentials.Email == "test@example.com" &&
			credentials.Password == "securePassword123"
	})).Return(nil)

	credentials := models.UserCredentials{
		Email:      "test@example.com",
		FirstName:  "John",
		MiddleName: "M",
		LastName:   "Doe",
		Password:   "securePassword123",
	}
	err := service.SendPassword(credentials)
	assert.NoError(t, err)
	mockProducer.AssertExpectations(t)
}

func TestKafkaPasswordDeliveryService_SendPassword_WriteError(t *testing.T) {
	mockProducer := new(mocks.MockProducer)
	service := &KafkaPasswordDeliveryService{
		Producer: mockProducer,
		Topic:    "test-topic",
	}

	mockProducer.On("WriteMessages", mock.Anything, mock.Anything).
		Return(errors.New("mock error: failed to write messages"))

	credentials := models.UserCredentials{
		Email:      "test@example.com",
		FirstName:  "John",
		MiddleName: "M",
		LastName:   "Doe",
		Password:   "securePassword123",
	}
	err := service.SendPassword(credentials)
	assert.Error(t, err)
	assert.EqualError(t, err, "mock error: failed to write messages")
	mockProducer.AssertExpectations(t)
}
