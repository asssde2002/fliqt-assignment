DOCKER_COMPOSE_FILE=docker-compose.yaml

.PHONY: all build stop deploy test

all: build deploy

build:
	@echo "Building..."
	@docker-compose -f $(DOCKER_COMPOSE_FILE) build > /dev/null 2>&1

deploy:
	@echo "Starting..."
	@docker-compose -f $(DOCKER_COMPOSE_FILE) up -d > /dev/null 2>&1

stop:
	@echo "Stopping..."
	@docker-compose -f $(DOCKER_COMPOSE_FILE) down -v > /dev/null 2>&1

test:
	@cd backend && go test ./internal/handlers -v
