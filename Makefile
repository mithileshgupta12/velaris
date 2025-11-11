ifneq (,$(wildcard ./.env))
    include .env
    export
endif

MIGRATIONS_DIR=internal/db/migrations

DB_DSN=postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

.PHONY: build run migrate-create migrate-up migrate-down test lint

build:
	@go build -o ./target/main ./main.go

run: build
	@./target/main \
		-app-port=$(APP_PORT) \
		-frontend-url=$(FRONTEND_URL) \
		-db-host=$(DB_HOST) \
		-db-port=$(DB_PORT) \
		-db-name=$(DB_NAME) \
		-db-username=$(DB_USERNAME) \
		-db-password=$(DB_PASSWORD) \
		-db-sslmode=$(DB_SSLMODE)

migrate-create:
	@if [ -z "$(NAME)" ]; then \
		echo "Error: NAME is required. Usage: make migrate-create NAME=your_migration_name"; \
		exit 1; \
	fi
	@echo "Creating migration: $(NAME)"
	@goose -dir $(MIGRATIONS_DIR) create $(NAME) sql

migrate-up:
	@echo "Running migrations..."
	@goose -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" up

migrate-down:
	@echo "Running migrations..."
	@goose -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" down

test:
	@go test ./...

lint:
	@golangci-lint run