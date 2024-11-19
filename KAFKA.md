
```sh
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
```