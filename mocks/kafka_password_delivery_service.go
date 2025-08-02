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
