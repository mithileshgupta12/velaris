.PHONY: build run create-migration

MIGRATION_DIR = ./database/migrations

build:
	@mkdir -p build
	@go build -o build/main .

run: build
	@./build/main

create-migration:
ifndef NAME
	$(error NAME is required. Usage: make create-migration NAME=migration_name)		
endif
	@goose -dir $(MIGRATION_DIR) create $(NAME) sql
