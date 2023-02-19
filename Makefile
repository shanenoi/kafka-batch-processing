PRODUCER_CONTAINER_ID := $(shell docker ps -a | grep demo-kafka-producer | awk '{print $$1}')
CONSUMER_CONTAINER_ID := $(shell docker ps -a | grep demo-kafka-consumer | awk '{print $$1}')

build:
	go build -o bins/producer ./producer/main.go
	go build -o bins/consumer ./consumer/main.go

start-platform:
	docker compose up -d zookeeper
	docker compose up -d kafka1
	docker compose up -d kafka2
	docker compose up -d kafka3

producer-start: build start-platform
	docker compose build
	docker compose up -d producer

producer-upload: build
	docker cp ./bins/producer $(PRODUCER_CONTAINER_ID):/producer
	docker restart $(PRODUCER_CONTAINER_ID)

producer-logs:
	docker logs $(PRODUCER_CONTAINER_ID)

consumer-start: build start-platform
	docker compose build
	docker compose up -d consumer

consumer-upload: build
	docker cp ./bins/consumer $(CONSUMER_CONTAINER_ID):/consumer
	docker restart $(CONSUMER_CONTAINER_ID)

consumer-logs:
	docker logs $(CONSUMER_CONTAINER_ID)

down:
	docker compose down
