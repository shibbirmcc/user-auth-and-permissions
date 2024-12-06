# **Running Kakfa during local execution and tests**

## Local Execution Kakfa Commands
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
### Create a docker image with embedded zookeeper using below Dockerfile:
```dockerfile
# Use a lightweight base image
FROM ubuntu:20.04

# Install dependencies
RUN apt-get update && \
    apt-get install -y openjdk-11-jre wget tar netcat && \
    apt-get clean

# Set Kafka and Zookeeper versions
ENV KAFKA_VERSION=3.9.0
ENV SCALA_VERSION=2.13
ENV KAFKA_HOME=/opt/kafka

# Download and extract Kafka
RUN wget https://downloads.apache.org/kafka/${KAFKA_VERSION}/kafka_${SCALA_VERSION}-${KAFKA_VERSION}.tgz -O /tmp/kafka.tgz && \
    mkdir -p ${KAFKA_HOME} && \
    tar -xvzf /tmp/kafka.tgz --strip-components=1 -C ${KAFKA_HOME} && \
    rm /tmp/kafka.tgz

# Add Kafka binaries to PATH
ENV PATH="${KAFKA_HOME}/bin:${PATH}"

# Copy start script to container
COPY start-kafka-zookeeper.sh /usr/bin/start-kafka-zookeeper.sh
RUN chmod +x /usr/bin/start-kafka-zookeeper.sh

# Expose Zookeeper and Kafka ports
EXPOSE 2181 9092

# Set environment variables
ENV ZOOKEEPER_CLIENT_PORT=2181
ENV KAFKA_BROKER_ID=1
ENV KAFKA_ZOOKEEPER_CONNECT=localhost:2181
ENV KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://localhost:9092
ENV KAFKA_LISTENERS=PLAINTEXT://0.0.0.0:9092
ENV KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR=1

# Start both services
CMD ["start-kafka-zookeeper.sh"]
```

### Create the startup script called ``start-kafka-zookeeper.sh``
```shell
#!/bin/bash

# Start Zookeeper
echo "Starting Zookeeper..."
${KAFKA_HOME}/bin/zookeeper-server-start.sh ${KAFKA_HOME}/config/zookeeper.properties > /var/log/zookeeper.log 2>&1 &

# Wait for Zookeeper to start
echo "Waiting for Zookeeper to start..."
while ! nc -z localhost 2181; do   
  sleep 1
done
echo "Zookeeper started successfully."

# Start Kafka
echo "Starting Kafka..."
${KAFKA_HOME}/bin/kafka-server-start.sh ${KAFKA_HOME}/config/server.properties > /var/log/kafka.log 2>&1 &

# Wait for Kafka to start
echo "Waiting for Kafka to start..."
while ! grep -q "started (kafka.server.KafkaServer)" /var/log/kafka.log; do
  sleep 1
done
echo "Kafka started successfully."

# Keep the container running
tail -f /dev/null
```
### Create and push the image to an public repository
N.B: You have to have an account to any public repository to push any custom image.
```shell
docker build -t <repository_username>/kafka-with-zookeeper:latest .
docker push <repository_username>/kafka-with-zookeeper:latest
```

The embedded zookeeper anf kafka image is already available publicly, just pull this image:
```shell
docker pull shibbirmcc/kafka-with-zookeeper:latest
```