# Load environment variables from .env file
include .env
export

# Database migration commands
GOOSE := goose -dir internal/infrastructure/db/migrations

# SQLC command
SQLC := sqlc

.PHONY: migrate-create migrate-up migrate-down migrate-status db-reset sqlc-generate sqlc-verify

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

# Build commands
.PHONY: build run

build:
	go build -o bin/api cmd/api/main.go

run: build
	./bin/api

run-local:
	go run cmd/api/main.go

# Test commands
.PHONY: test

test:
	go test -v ./...

# Help command
.PHONY: help

help:
	@echo "Available commands:"
	@echo "  migrate-create  - Create a new migration file"
	@echo "  migrate-up      - Run all available migrations"
	@echo "  migrate-down    - Rollback the last migration"
	@echo "  migrate-status  - Show current migration status"
	@echo "  db-reset        - Reset the database (rollback all migrations)"
	@echo "  sqlc-generate   - Generate Go code from SQL"
	@echo "  sqlc-verify     - Verify SQL queries against the schema"
	@echo "  build           - Build the application"
	@echo "  run             - Build and run the application"
	@echo "  test            - Run tests"
