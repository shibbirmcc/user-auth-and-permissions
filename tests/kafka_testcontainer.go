package tests

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/docker/docker/api/types/container"
	_ "github.com/lib/pq"
	"github.com/segmentio/kafka-go"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func SetupKafkaContainer() func() {
	LoadEnvironmentVariables()

	ctx := context.Background()

	kafkaTopic := os.Getenv("KAFKA_TOPIC")
	if kafkaTopic == "" {
		log.Fatal("Environment variable KAFKA_TOPIC is not set")
	}

	// Start Zookeeper container
	zookeeperReq := testcontainers.ContainerRequest{
		Image:        "confluentinc/cp-zookeeper:latest",
		ExposedPorts: []string{"2181/tcp"},
		Env: map[string]string{
			"ZOOKEEPER_CLIENT_PORT": "2181",
			"ZOOKEEPER_TICK_TIME":   "2000",
		},
		WaitingFor: wait.ForLog("binding to port").WithStartupTimeout(60 * time.Second),
	}

	zookeeperContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: zookeeperReq,
		Started:          true,
	})
	if err != nil {
		log.Fatalf("Could not start Zookeeper container: %s", err)
	}
	zookeeperPort, err := zookeeperContainer.MappedPort(ctx, "2181")
	if err != nil {
		log.Fatalf("Could not get Zookeeper port: %s", err)
	}

	kafkaReq := testcontainers.ContainerRequest{
		Image:        "confluentinc/cp-kafka:latest",
		ExposedPorts: []string{"9092/tcp"},
		Env: map[string]string{
			"KAFKA_BROKER_ID":                                "1",
			"KAFKA_ZOOKEEPER_CONNECT":                        fmt.Sprintf("host.docker.internal:%s", zookeeperPort.Port()),
			"KAFKA_ADVERTISED_LISTENERS":                     "PLAINTEXT://localhost:9092",
			"KAFKA_LISTENER_SECURITY_PROTOCOL_MAP":           "PLAINTEXT:PLAINTEXT",
			"KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR":         "1",
			"KAFKA_TRANSACTION_STATE_LOG_REPLICATION_FACTOR": "1",
			"KAFKA_TRANSACTION_STATE_LOG_MIN_ISR":            "1",
		},
		HostConfigModifier: func(config *container.HostConfig) {
			config.ExtraHosts = []string{"host.docker.internal:host-gateway"} // Add host mapping
		},
		WaitingFor: wait.ForLog("Kafka Server started").WithStartupTimeout(120 * time.Second),
	}

	kafkaContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: kafkaReq,
		Started:          true,
	})
	if err != nil {
		log.Fatalf("Could not start Kafka container: %s", err)
	}

	kafkaHost, err := kafkaContainer.Host(ctx)
	if err != nil {
		log.Fatalf("Could not get Kafka host: %s", err)
	}
	kafkaPort, err := kafkaContainer.MappedPort(ctx, "9092")
	if err != nil {
		log.Fatalf("Could not get Kafka port: %s", err)
	}

	// Set the Kafka brokers environment variable
	os.Setenv("KAFKA_BROKERS", fmt.Sprintf("%s:%s", kafkaHost, kafkaPort.Port()))

	// Create the Kafka topic
	err = createKafkaTopic(fmt.Sprintf("%s:%s", kafkaHost, kafkaPort.Port()), kafkaTopic)
	if err != nil {
		log.Fatalf("Could not create Kafka topic: %s", err)
	}

	// Teardown function
	teardown := func() {
		if kafkaContainer != nil {
			err := kafkaContainer.Terminate(ctx)
			if err != nil {
				log.Fatalf("Could not terminate Kafka container: %s", err)
			}
		}
		if zookeeperContainer != nil {
			err := zookeeperContainer.Terminate(ctx)
			if err != nil {
				log.Fatalf("Could not terminate Zookeeper container: %s", err)
			}
		}
	}

	return teardown
}

func createKafkaTopic(broker string, topic string) error {
	conn, err := kafka.Dial("tcp", broker)
	if err != nil {
		return fmt.Errorf("failed to connect to Kafka broker: %w", err)
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return fmt.Errorf("failed to get Kafka controller: %w", err)
	}

	controllerConn, err := kafka.Dial("tcp", fmt.Sprintf("%s:%d", controller.Host, controller.Port))
	if err != nil {
		return fmt.Errorf("failed to connect to Kafka controller: %w", err)
	}
	defer controllerConn.Close()

	err = controllerConn.CreateTopics(kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     1,
		ReplicationFactor: 1,
	})
	if err != nil {
		return fmt.Errorf("failed to create Kafka topic: %w", err)
	}

	log.Printf("Kafka topic %s created successfully", topic)
	return nil
}
