DOCKER_COMPOSE=docker compose

all: build run

build:
	@echo "Building Docker images..."
	$(DOCKER_COMPOSE) build

run:
	@echo "Starting application..."
	$(DOCKER_COMPOSE) up -d

stop:
	@echo "Stopping application..."
	$(DOCKER_COMPOSE) down -v

clear:
	@echo "Destroy cached data..."
	rm -rf ./debug/pg_data/
	rm -rf ./debug/redis_data/