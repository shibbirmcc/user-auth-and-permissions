package tests

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/segmentio/kafka-go"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func waitForDBConnection(db *sql.DB) error {
	maxAttempts := 5
	for i := 0; i < maxAttempts; i++ {
		if err := db.Ping(); err == nil {
			return nil
		}
		time.Sleep(2 * time.Second) // Wait before retrying
	}
	return fmt.Errorf("could not establish database connection after %d attempts", maxAttempts)
}

func GetGormDBFromSQLDB(sqlDB *sql.DB) (*gorm.DB, error) {
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return gormDB, nil
}

func DeleteTestData(sqlDB *sql.DB) {
	err := sqlDB.QueryRow("DELETE FROM user_details;")
	if err != nil {
		fmt.Printf("Error delete rows from user_details: %v\n", err)
	}
	err = sqlDB.QueryRow("DELETE FROM users;")
	if err != nil {
		fmt.Printf("Error delete rows from user_details: %v\n", err)
	}
}

func SetupPostgresContainer() (*gorm.DB, func()) {
	LoadEnvironmentVariables()

	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:14",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     os.Getenv("DB_USER"),
			"POSTGRES_PASSWORD": os.Getenv("DB_PASSWORD"),
			"POSTGRES_DB":       os.Getenv("DB_NAME"),
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").WithStartupTimeout(60 * time.Second),
	}

	var err error
	var postgresContainer testcontainers.Container
	var sqlDB *sql.DB

	postgresContainer, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatalf("Could not start postgres container: %s", err)
	}

	host, err := postgresContainer.Host(ctx)
	if err != nil {
		log.Fatalf("Could not get host: %s", err)
	}
	port, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		log.Fatalf("Could not get port: %s", err)
	}

	// Override environment variables for connecting to the testcontainer
	os.Setenv("DB_HOST", host)
	os.Setenv("DB_PORT", port.Port())

	// Connect to the database for verification in tests
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	sqlDB, err = sql.Open("postgres", dsn)
	if err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
		os.Exit(1)
	}

	if err := waitForDBConnection(sqlDB); err != nil {
		fmt.Printf("Testcontainer Database is not ready: %v\n", err)
		os.Exit(1)
	}

	// Return *gorm.DB instance based on *sql.DB
	gormDB, err := GetGormDBFromSQLDB(sqlDB)
	if err != nil {
		log.Fatalf("Could not get Gorm DB: %s", err)
	}

	TeardownPostgresContainer := func() {
		DeleteTestData(sqlDB)
		if sqlDB != nil {
			sqlDB.Close()
		}
		if postgresContainer != nil {
			ctx := context.Background()
			err := postgresContainer.Terminate(ctx)
			if err != nil {
				log.Fatalf("Could not terminate container: %s", err)
			}
		}
	}

	return gormDB, TeardownPostgresContainer
}

func SetupKafkaContainer() func() {
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

	zookeeperHost, err := zookeeperContainer.Host(ctx)
	if err != nil {
		log.Fatalf("Could not get Zookeeper host: %s", err)
	}
	zookeeperPort, err := zookeeperContainer.MappedPort(ctx, "2181")
	if err != nil {
		log.Fatalf("Could not get Zookeeper port: %s", err)
	}

	// Start Kafka container
	kafkaReq := testcontainers.ContainerRequest{
		Image:        "confluentinc/cp-kafka:latest",
		ExposedPorts: []string{"9092/tcp"},
		Env: map[string]string{
			"KAFKA_BROKER_ID":                                "1",
			"KAFKA_ZOOKEEPER_CONNECT":                        fmt.Sprintf("%s:%s", zookeeperHost, zookeeperPort.Port()),
			"KAFKA_ADVERTISED_LISTENERS":                     "PLAINTEXT://localhost:9092",
			"KAFKA_LISTENER_SECURITY_PROTOCOL_MAP":           "PLAINTEXT:PLAINTEXT",
			"KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR":         "1",
			"KAFKA_TRANSACTION_STATE_LOG_REPLICATION_FACTOR": "1",
			"KAFKA_TRANSACTION_STATE_LOG_MIN_ISR":            "1",
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
