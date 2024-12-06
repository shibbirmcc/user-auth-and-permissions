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

	kafkaReq := testcontainers.ContainerRequest{
		Image:        "shibbirmcc/kafka-with-zookeeper:latest", // Kafka with Zookeeper included
		ExposedPorts: []string{"9092/tcp", "2181/tcp"},         // Kafka and Zookeeper ports
		WaitingFor:   wait.ForLog("Kafka started successfully.").WithStartupTimeout(120 * time.Second),
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
