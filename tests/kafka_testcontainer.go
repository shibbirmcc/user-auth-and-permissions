package tests

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
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

	// Use official Apache Kafka image with Kraft mode (no Zookeeper)
	kafkaReq := testcontainers.ContainerRequest{
		Image:        "apache/kafka:3.8.1", // Latest stable Kafka image with Kraft mode
		ExposedPorts: []string{"9092/tcp"},  // Only Kafka port needed (no Zookeeper)
		Env: map[string]string{
			// Kraft mode configuration
			"KAFKA_NODE_ID":                         "1",
			"KAFKA_PROCESS_ROLES":                   "broker,controller",
			"KAFKA_CONTROLLER_QUORUM_VOTERS":        "1@localhost:9093",
			"KAFKA_LISTENERS":                       "PLAINTEXT://0.0.0.0:9092,CONTROLLER://0.0.0.0:9093",
			"KAFKA_ADVERTISED_LISTENERS":            "PLAINTEXT://localhost:9092",
			"KAFKA_LISTENER_SECURITY_PROTOCOL_MAP":  "PLAINTEXT:PLAINTEXT,CONTROLLER:PLAINTEXT",
			"KAFKA_CONTROLLER_LISTENER_NAMES":       "CONTROLLER",
			"KAFKA_INTER_BROKER_LISTENER_NAME":      "PLAINTEXT",
			"KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR": "1",
			"KAFKA_TRANSACTION_STATE_LOG_REPLICATION_FACTOR": "1",
			"KAFKA_TRANSACTION_STATE_LOG_MIN_ISR":   "1",
			"KAFKA_LOG_DIRS":                        "/tmp/kraft-combined-logs",
			"KAFKA_AUTO_CREATE_TOPICS_ENABLE":       "true",
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

	// Wait a bit for Kafka to be fully ready
	time.Sleep(5 * time.Second)

	// Create the Kafka topic
	err = createKafkaTopic(kafkaContainer, kafkaTopic)
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
	}
	return teardown
}

func createKafkaTopic(container testcontainers.Container, topic string) error {
	command := []string{
		"/opt/kafka/bin/kafka-topics.sh",
		"--create",
		"--topic", topic,
		"--bootstrap-server", "localhost:9092",
		"--partitions", "1",
		"--replication-factor", "1",
	}
	exitCode, stdout, stderr := container.Exec(context.Background(), command)
	if exitCode != 0 {
		return fmt.Errorf("failed to create Kafka topic: %s (stderr: %s)", stdout, stderr)
	}
	log.Printf("Kafka topic %s created successfully", topic)
	return nil
}
