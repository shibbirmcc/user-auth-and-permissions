package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"

	"github.com/segmentio/kafka-go"
)

type MockProducer struct {
	mock.Mock
}

func (m *MockProducer) WriteMessages(ctx context.Context, msgs ...kafka.Message) error {
	args := m.Called(ctx, msgs)
	return args.Error(0)
}

func (m *MockProducer) Close() error {
	return nil
}

//type MockKafkaPasswordDeliveryService struct {
//	mock.Mock
//	Producer *MockKafkaWriter
//	Topic    string
//}
//
//func (m *MockKafkaPasswordDeliveryService) SendPassword(email, firstName, middleName, lastName, password string) error {
//	args := m.Called(email, firstName, middleName, lastName, password)
//	return args.Error(0)
//}
