DOCKER_COMPOSE_FILE=docker-compose.yaml

.PHONY: all build stop deploy

all: build deploy

build:
	docker-compose -f $(DOCKER_COMPOSE_FILE) build

deploy:
	docker-compose -f $(DOCKER_COMPOSE_FILE) up -d

stop:
	docker-compose -f $(DOCKER_COMPOSE_FILE) down -v
