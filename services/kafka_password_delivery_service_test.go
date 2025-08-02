package services

import (
	"encoding/json"
	"errors"
	"github.com/segmentio/kafka-go"
	"github.com/shibbirmcc/user-auth-and-permissions/mocks"
	"github.com/shibbirmcc/user-auth-and-permissions/models"
	"github.com/shibbirmcc/user-auth-and-permissions/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"os"
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

func TestNewKafkaPasswordDeliveryService_Success_WithTestcontainers(t *testing.T) {
	tearDownKafkaContainer := tests.SetupKafkaContainer()
	defer tearDownKafkaContainer()
	service, err := NewKafkaPasswordDeliveryService()

	assert.NoError(t, err)
	assert.NotNil(t, service)
	assert.NotNil(t, service.Producer)
}

func TestNewKafkaPasswordDeliveryService_MissingBrokers(t *testing.T) {
	// Set only the topic environment variable
	os.Unsetenv("KAFKA_BROKERS")
	os.Setenv("KAFKA_TOPIC", "test-topic")
	defer os.Unsetenv("KAFKA_TOPIC")
	service, err := NewKafkaPasswordDeliveryService()

	// Assertions
	assert.Nil(t, service)
	assert.EqualError(t, err, "no kafka brokers found in environment variable")
}

func TestNewKafkaPasswordDeliveryService_MissingTopic(t *testing.T) {
	os.Setenv("KAFKA_BROKERS", "localhost:9092")
	os.Unsetenv("KAFKA_TOPIC")
	defer os.Unsetenv("KAFKA_BROKERS")
	service, err := NewKafkaPasswordDeliveryService()

	assert.Nil(t, service)
	assert.EqualError(t, err, "no kafka topic found in environment variable")
}

func TestNewKafkaPasswordDeliveryService_BrokerConnectionFailure(t *testing.T) {
	os.Setenv("KAFKA_BROKERS", "invalid:9092")
	os.Setenv("KAFKA_TOPIC", "test-topic")
	defer os.Unsetenv("KAFKA_BROKERS")
	defer os.Unsetenv("KAFKA_TOPIC")
	service, err := NewKafkaPasswordDeliveryService()

	assert.Nil(t, service)
	assert.Contains(t, err.Error(), "failed to connect to kafka broker")
}
