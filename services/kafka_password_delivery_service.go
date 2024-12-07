package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/segmentio/kafka-go"
	"github.com/shibbirmcc/user-auth-and-permissions/models"
)

type MessageProducer interface {
	WriteMessages(ctx context.Context, msgs ...kafka.Message) error
	Close() error
}

type KafkaPasswordDeliveryService struct {
	Producer MessageProducer
	Topic    string
}

func NewKafkaPasswordDeliveryService() (*KafkaPasswordDeliveryService, error) {
	brokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	if len(brokers) == 0 || brokers[0] == "" {
		return nil, errors.New("no kafka brokers found in environment variable")
	}
	fmt.Printf("brokers: %v\n", brokers)
	topic := os.Getenv("KAFKA_TOPIC")
	if topic == "" {
		return nil, errors.New("no kafka topic found in environment variable")
	}
	fmt.Printf("topic: %s\n", topic)

	conn, err := kafka.Dial("tcp", brokers[0])
	if err != nil {
		return nil, fmt.Errorf("failed to connect to kafka broker: %v", err)
	}
	conn.Close()

	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers[0]),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	return &KafkaPasswordDeliveryService{
		Producer: writer,
		Topic:    topic,
	}, nil
}

func (s *KafkaPasswordDeliveryService) SendPassword(credentials models.UserCredentials) error {
	message, err := json.Marshal(credentials)
	if err != nil {
		return err
	}

	err = s.Producer.WriteMessages(context.Background(), kafka.Message{
		Value: message,
	})
	if err != nil {
		log.Printf("Failed to send message to Kafka: %v", err)
		return err
	}

	log.Printf("Password sent to Kafka topic %s for user %s", s.Topic, credentials.Email)
	return nil
}
