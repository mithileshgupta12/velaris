.PHONY: build run

build:
	@go build -o ./target/api ./cmd/api/main.go

run: build
	@./target/api
