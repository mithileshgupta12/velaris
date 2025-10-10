.PHONY: build run

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

build:
	@go build -o ./target/api ./cmd/api/main.go

run: build
	@./target/api \
		-db-host=$(DB_HOST) \
		-db-port=$(DB_PORT) \
		-db-name=$(DB_NAME) \
		-db-username=$(DB_USERNAME) \
		-db-password=$(DB_PASSWORD) \
		-db-sslmode=$(DB_SSLMODE)
