# Variables for the environment
DEV_COMPOSE_FILE = docker-compose.dev.yml
PROD_COMPOSE_FILE = docker-compose.prod.yml
GO_APP_DIR = cmd/main.go  # Path to your Go application

# Default target
.PHONY: help
help:
	@echo "Usage:"
	@echo "  make up-dev       Start development environment"
	@echo "  make up-prod      Start production environment"
	@echo "  make down         Stop all running containers"
	@echo "  make build-dev    Build development containers"
	@echo "  make build-prod   Build production containers"
	@echo "  make logs         View logs of all running containers"
	@echo "  make clean        Stop and remove all containers, networks, and volumes"
	@echo "  make build-go     Build Go application"
	@echo "  make run-go       Run Go application"

# Start development environment
.PHONY: up-dev
up-dev:
	@docker-compose -f $(DEV_COMPOSE_FILE) up --build

# Start production environment
.PHONY: up-prod
up-prod:
	@docker-compose -f $(PROD_COMPOSE_FILE) up --build -d

# Stop all running containers
.PHONY: down
down:
	@docker-compose -f $(DEV_COMPOSE_FILE) down
	@docker-compose -f $(PROD_COMPOSE_FILE) down

# Build development containers
.PHONY: build-dev
build-dev:
	@docker-compose -f $(DEV_COMPOSE_FILE) build

# Build production containers
.PHONY: build-prod
build-prod:
	@docker-compose -f $(PROD_COMPOSE_FILE) build

# View logs of all running containers
.PHONY: logs
logs:
	@docker-compose logs -f

# Clean up containers, networks, and volumes
.PHONY: clean
clean:
	@docker-compose -f $(DEV_COMPOSE_FILE) down -v
	@docker-compose -f $(PROD_COMPOSE_FILE) down -v
	@docker system prune --volumes -f

# Build Go application (optional)
.PHONY: build-go
build-go:
	@go build -o app $(GO_APP_DIR)

# Run Go application
.PHONY: run-go
run-go:
	@go run $(GO_APP_DIR)
