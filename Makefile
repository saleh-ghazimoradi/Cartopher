.PHONY: help build run dev lint docker-up docker-down migrate-up migrate-down

help:
	@echo "Available commands:"
	@echo "  make build        - Build the application"
	@echo "  make run          - Run the application"
	@echo "  make dev          - Run in development mode (same as run)"
	@echo "  make lint         - Lint AND auto-fix formatting/issues"
	@echo "  make docker-up    - Start Docker containers"
	@echo "  make docker-down  - Stop Docker containers"
	@echo "  make migrate-up   - Apply database migrations"
	@echo "  make migrate-down - Rollback database migrations"

build:
	go build -o bin/app

run:
	go run . run

dev:
	go run . run

lint:
	golangci-lint run --fix ./...

docker-up:
	docker compose up -d

docker-down:
	docker compose down

migrate-up:
	go run . migrateUp

migrate-down:
	go run . migrateDown