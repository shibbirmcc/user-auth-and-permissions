# **Running Kafka during local execution and tests**

## Local Execution Kafka Commands

### Option 1: Modern KRaft Mode (Recommended)
Using the official Apache Kafka image with KRaft mode (no Zookeeper required):

```shell
# Clean up any existing containers
sudo docker container prune

# Run Kafka in KRaft mode
sudo docker run -d --name kafka \
  -p 9092:9092 \
  -e KAFKA_NODE_ID=1 \
  -e KAFKA_PROCESS_ROLES=broker,controller \
  -e KAFKA_CONTROLLER_QUORUM_VOTERS=1@localhost:9093 \
  -e KAFKA_LISTENERS=PLAINTEXT://0.0.0.0:9092,CONTROLLER://0.0.0.0:9093 \
  -e KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://localhost:9092 \
  -e KAFKA_LISTENER_SECURITY_PROTOCOL_MAP=PLAINTEXT:PLAINTEXT,CONTROLLER:PLAINTEXT \
  -e KAFKA_CONTROLLER_LISTENER_NAMES=CONTROLLER \
  -e KAFKA_INTER_BROKER_LISTENER_NAME=PLAINTEXT \
  -e KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR=1 \
  -e KAFKA_TRANSACTION_STATE_LOG_REPLICATION_FACTOR=1 \
  -e KAFKA_TRANSACTION_STATE_LOG_MIN_ISR=1 \
  -e KAFKA_LOG_DIRS=/tmp/kraft-combined-logs \
  -e KAFKA_AUTO_CREATE_TOPICS_ENABLE=true \
  apache/kafka:3.8.1

# Create topic
sudo docker exec -it kafka /opt/kafka/bin/kafka-topics.sh \
  --create --topic credentials \
  --bootstrap-server localhost:9092 \
  --partitions 1 --replication-factor 1

# Console consumer
sudo docker exec -it kafka /opt/kafka/bin/kafka-console-consumer.sh \
  --bootstrap-server localhost:9092 --topic credentials --from-beginning

# Console producer
sudo docker exec -it kafka /opt/kafka/bin/kafka-console-producer.sh \
  --bootstrap-server localhost:9092 --topic credentials
```

### Option 2: Legacy Zookeeper Mode
Using Confluent images with separate Zookeeper (for compatibility with older setups):
```shell
sudo docker container prune

sudo docker run -d --name zookeeper \
-p 2181:2181 \
-e ZOOKEEPER_CLIENT_PORT=2181 \
-e ZOOKEEPER_TICK_TIME=2000 \
confluentinc/cp-zookeeper:latest



sudo docker run -d \
--name kafka \
-p 9092:9092 \
-e KAFKA_BROKER_ID=1 \
-e KAFKA_ZOOKEEPER_CONNECT=host.docker.internal:2181 \
-e KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://localhost:9092 \
-e KAFKA_LISTENERS=PLAINTEXT://0.0.0.0:9092 \
-e KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR=1 \
confluentinc/cp-kafka:latest

sudo docker exec -it 6f45beaa852c kafka-topics --bootstrap-server localhost:9092 --create --replication-factor 1 --partitions 1 --topic credentials

sudo docker exec -it 6f45beaa852c kafka-console-consumer --bootstrap-server localhost:9092 --topic credentials --from-beginning

## Console producer
sudo docker exec -it 636092f818f5 kafka-console-producer.sh --bootstrap-server localhost:9092 --topic credentials
```

## TestContainer Execution

The testcontainer implementation now uses the official Apache Kafka image with KRaft mode, eliminating the need for Zookeeper. This provides a simpler, more modern setup.

### Key Features:
- **No Zookeeper Required**: Uses Kafka's KRaft mode (Kafka Raft metadata mode)
- **Official Apache Kafka Image**: Uses `apache/kafka:3.8.1` 
- **Automatic Topic Creation**: Creates the required topic programmatically
- **Simplified Configuration**: Fewer moving parts and dependencies

### TestContainer Configuration:
The testcontainer is configured with the following KRaft mode settings:

```go
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
```

### Benefits of KRaft Mode:
1. **Simplified Architecture**: No separate Zookeeper cluster to manage
2. **Better Performance**: Reduced latency and improved throughput
3. **Easier Scaling**: Simpler cluster management and scaling operations
4. **Official Support**: Uses the official Apache Kafka image maintained by the Kafka team
5. **Future-Proof**: KRaft is the future of Kafka (Zookeeper is being deprecated)

### Usage in Tests:
The testcontainer automatically:
- Starts a Kafka broker in KRaft mode
- Waits for Kafka to be fully ready
- Creates the required topic specified in `KAFKA_TOPIC` environment variable
- Sets the `KAFKA_BROKERS` environment variable for the application to use
- Provides a teardown function to clean up resources

### Environment Variables Required:
- `KAFKA_TOPIC`: The name of the Kafka topic to create for testing
- `KAFKA_BROKERS`: Automatically set by the testcontainer setup function

### Manual Testing with Official Image:
If you want to run Kafka manually for testing, you can use the same official image:

```shell
# Run Kafka in KRaft mode
docker run -d --name kafka \
  -p 9092:9092 \
  -e KAFKA_NODE_ID=1 \
  -e KAFKA_PROCESS_ROLES=broker,controller \
  -e KAFKA_CONTROLLER_QUORUM_VOTERS=1@localhost:9093 \
  -e KAFKA_LISTENERS=PLAINTEXT://0.0.0.0:9092,CONTROLLER://0.0.0.0:9093 \
  -e KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://localhost:9092 \
  -e KAFKA_LISTENER_SECURITY_PROTOCOL_MAP=PLAINTEXT:PLAINTEXT,CONTROLLER:PLAINTEXT \
  -e KAFKA_CONTROLLER_LISTENER_NAMES=CONTROLLER \
  -e KAFKA_INTER_BROKER_LISTENER_NAME=PLAINTEXT \
  -e KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR=1 \
  -e KAFKA_TRANSACTION_STATE_LOG_REPLICATION_FACTOR=1 \
  -e KAFKA_TRANSACTION_STATE_LOG_MIN_ISR=1 \
  -e KAFKA_LOG_DIRS=/tmp/kraft-combined-logs \
  -e KAFKA_AUTO_CREATE_TOPICS_ENABLE=true \
  apache/kafka:3.8.1

# Create topic manually if needed
docker exec -it kafka /opt/kafka/bin/kafka-topics.sh \
  --create --topic credentials \
  --bootstrap-server localhost:9092 \
  --partitions 1 --replication-factor 1

# Console consumer
docker exec -it kafka /opt/kafka/bin/kafka-console-consumer.sh \
  --bootstrap-server localhost:9092 --topic credentials --from-beginning

# Console producer
docker exec -it kafka /opt/kafka/bin/kafka-console-producer.sh \
  --bootstrap-server localhost:9092 --topic credentials
```
