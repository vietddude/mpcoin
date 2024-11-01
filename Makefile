# Load environment variables from .env file
-include .env
export

# Go commands
GO := go
GOFLAGS := -v
API_CMD := cmd/api/main.go
WORKER_CMD := cmd/worker/main.go
MAIL_CMD := cmd/mail/main.go

# Database migration commands
GOOSE := goose -dir internal/infrastructure/db/migrations

# SQLC command
SQLC := sqlc

# Docker Compose files
DOCKER_COMPOSE_FILE := docker-compose.yml

# Build output directories
BIN_DIR := bin
$(shell mkdir -p $(BIN_DIR))

.PHONY: all build clean test help docker-* run-* migrate-* sqlc-*

# Build commands
build-api:
	CGO_ENABLED=0 GOOS=linux $(GO) build $(GOFLAGS) -o $(BIN_DIR)/api $(API_CMD)

build-worker:
	CGO_ENABLED=0 GOOS=linux $(GO) build $(GOFLAGS) -o $(BIN_DIR)/worker $(WORKER_CMD)

build-mailworker:
	CGO_ENABLED=0 GOOS=linux $(GO) build $(GOFLAGS) -o $(BIN_DIR)/mailworker $(MAIL_CMD)

build-all: build-api build-worker build-mailworker

# Run commands
run-api:
	$(GO) run $(API_CMD)

run-worker:
	$(GO) run $(WORKER_CMD)

run-mailworker:
	$(GO) run $(MAIL_CMD)

# Docker commands
docker-build-api:
	docker build -f Dockerfile.api -t mpc-api .

docker-build-worker:
	docker build -f Dockerfile.worker -t mpc-worker .

docker-build-mailworker:
	docker build -f Dockerfile.mail -t mpc-mailworker .

docker-build: docker-build-api docker-build-worker docker-build-mailworker

up:
	docker-compose up -d

down:
	docker-compose down

restart:
	docker-compose restart

logs:
	docker-compose logs -f

ps:
	docker-compose ps

# Database migration commands
migrate-create:
	@read -p "Enter migration name: " name; \
	$(GOOSE) -s create $$name sql

migrate-up:
	$(GOOSE) up

migrate-down:
	$(GOOSE) down

migrate-status:
	$(GOOSE) status

db-reset:
	$(GOOSE) reset

# SQLC commands
sqlc-generate:
	$(SQLC) generate

sqlc-verify:
	$(SQLC) verify

# Test commands
test:
	$(GO) test ./... -v

test-coverage:
	$(GO) test ./... -v -coverprofile=coverage.out && \
	$(GO) tool cover -html=coverage.out

# Clean command
clean:
	rm -rf $(BIN_DIR)
	docker-compose down -v

# Development commands
dev-api: up
	$(MAKE) run-api

dev-worker: up
	$(MAKE) run-worker

dev-mailworker: up
	$(MAKE) run-mailworker

# Help command
help:
	@echo "Available commands:"
	@echo "Build commands:"
	@echo "  build-api         - Build the API server"
	@echo "  build-worker      - Build the worker"
	@echo "  build-mailworker  - Build the mail worker"
	@echo "  build-all         - Build all services"
	@echo
	@echo "Run commands:"
	@echo "  run-api           - Run the API server locally"
	@echo "  run-worker        - Run the worker locally"
	@echo "  run-mailworker    - Run the mail worker locally"
	@echo
	@echo "Docker commands:"
	@echo "  docker-build      - Build all Docker images"
	@echo "  up                - Start all Docker services"
	@echo "  down              - Stop all Docker services"
	@echo "  restart           - Restart all Docker services"
	@echo "  logs              - View Docker logs"
	@echo "  ps                - List running Docker services"
	@echo
	@echo "Development commands:"
	@echo "  dev-api           - Run API in development mode"
	@echo "  dev-worker        - Run worker in development mode"
	@echo "  dev-mailworker    - Run mail worker in development mode"
	@echo
	@echo "Database commands:"
	@echo "  migrate-create    - Create a new migration"
	@echo "  migrate-up        - Run all migrations"
	@echo "  migrate-down      - Rollback last migration"
	@echo "  migrate-status    - Show migration status"
	@echo "  db-reset         - Reset database"
	@echo
	@echo "SQLC commands:"
	@echo "  sqlc-generate     - Generate Go code from SQL"
	@echo "  sqlc-verify       - Verify SQL queries"
	@echo
	@echo "Test commands:"
	@echo "  test             - Run tests"
	@echo "  test-coverage    - Run tests with coverage"
	@echo
	@echo "Other commands:"
	@echo "  clean            - Clean build artifacts and Docker volumes"