APP_PACKAGE := ./cmd/app
COMPOSE := docker compose
APP_SERVICE := exchange-rate-service

.DEFAULT_GOAL := help

.PHONY: help run build test test-client test-service test-handler test-core test-race test-cover clean \
	docker-build docker-up docker-up-build docker-down docker-stop docker-start \
	docker-restart docker-logs docker-app-logs docker-kafka-logs docker-ui-logs \
	docker-ps docker-health docker-kafka-up docker-kafka-down docker-tools-up docker-tools-down docker-config docker-pull \
	docker-ps docker-health docker-tools-up docker-tools-down docker-config docker-pull \
	docker-create docker-remove docker-reset

help:
	@echo "Available commands:"
	@echo "  make help               Show all available commands"
	@echo "  make run                Run the application locally"
	@echo "  make build              Compile the application"
	@echo "  make test               Run all project tests"
	@echo "  make test-core          Run client, service, and handler tests"
	@echo "  make test-client        Run Monobank client tests"
	@echo "  make test-service       Run exchange-rate service tests"
	@echo "  make test-handler       Run exchange-rate handler tests"
	@echo "  make test-race          Run all tests with the race detector"
	@echo "  make test-cover         Run all tests and print coverage"
	@echo "  make clean              Clear Go build and test caches"
	@echo "  make docker-build       Build Docker Compose images"
	@echo "  make docker-up          Start all containers in the background"
	@echo "  make docker-up-build    Build images and start all containers"
	@echo "  make docker-down        Stop and remove containers and networks"
	@echo "  make docker-stop        Stop all running containers"
	@echo "  make docker-start       Start existing containers"
	@echo "  make docker-restart     Restart all containers"
	@echo "  make docker-logs        Follow logs from all containers"
	@echo "  make docker-app-logs    Follow application container logs"
	@echo "  make docker-kafka-logs  Follow Kafka container logs"
	@echo "  make docker-ui-logs     Follow Kafka UI container logs"
	@echo "  make docker-ps          Show Compose container status"
	@echo "  make docker-health      Show application container health"
	@echo "  make docker-kafka-up    Start the optional Kafka broker"
	@echo "  make docker-kafka-down  Stop the optional Kafka broker"
	@echo "  make docker-tools-up    Start optional tools, including Kafka UI"
	@echo "  make docker-tools-down  Stop optional tools"
	@echo "  make docker-config      Validate and render Compose configuration"
	@echo "  make docker-pull        Pull service images"
	@echo "  make docker-create      Create containers without starting them"
	@echo "  make docker-remove      Remove stopped containers"
	@echo "  make docker-reset       Remove containers, networks, and volumes"

run:
	go run $(APP_PACKAGE)

build:
	go build $(APP_PACKAGE)

test:
	go test ./...

test-core:
	go test ./internal/external/monobank ./internal/services/exchange_rate ./internal/delivery/grpc/handler/exchange_rate

test-client:
	go test ./internal/external/monobank

test-service:
	go test ./internal/services/exchange_rate

test-handler:
	go test ./internal/delivery/grpc/handler/exchange_rate

test-race:
	go test -race ./...

test-cover:
	go test -cover ./...

clean:
	go clean -cache -testcache

docker-build:
	$(COMPOSE) build

docker-up:
	$(COMPOSE) up -d

docker-up-build:
	$(COMPOSE) up -d --build

docker-down:
	$(COMPOSE) down --remove-orphans

docker-stop:
	$(COMPOSE) stop

docker-start:
	$(COMPOSE) start

docker-restart:
	$(COMPOSE) restart

docker-logs:
	$(COMPOSE) logs -f

docker-app-logs:
	$(COMPOSE) logs -f $(APP_SERVICE)

docker-kafka-logs:
	$(COMPOSE) logs -f kafka

docker-ui-logs:
	$(COMPOSE) logs -f kafka-ui

docker-ps:
	$(COMPOSE) ps

docker-health:
	$(COMPOSE) ps $(APP_SERVICE)

docker-kafka-up:
	$(COMPOSE) --profile kafka up -d kafka

docker-kafka-down:
	$(COMPOSE) --profile kafka stop kafka

docker-tools-up:
	$(COMPOSE) --profile tools up -d kafka-ui

docker-tools-down:
	$(COMPOSE) --profile tools stop kafka-ui

docker-config:
	$(COMPOSE) config

docker-pull:
	$(COMPOSE) pull

docker-create:
	$(COMPOSE) create

docker-remove:
	$(COMPOSE) rm -f

docker-reset:
	$(COMPOSE) down --volumes --remove-orphans
