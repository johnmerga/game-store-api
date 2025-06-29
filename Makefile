# Makefile for RealgamingMarketplace Backend

.PHONY: help build run test clean migrate-up migrate-down docker-up docker-down

# Default environment variables
DB_HOST ?= localhost
DB_PORT ?= 5432
DB_USER ?= postgres
DB_PASSWORD ?= password
DB_NAME ?= marketplace
DB_SSL_MODE ?= disable

# Help command
help:
	@echo "Available commands:"
	@echo "  build         - Build the application"
	@echo "  run           - Run the application"
	@echo "  test          - Run tests"
	@echo "  clean         - Clean build artifacts"
	@echo "  migrate-up    - Run database migrations up"
	@echo "  migrate-down  - Run database migrations down"
	@echo "  docker-up     - Start PostgreSQL with Docker"
	@echo "  docker-down   - Stop PostgreSQL Docker container"
	@echo "  deps          - Install dependencies"

# Install dependencies
deps:
	go mod tidy
	go mod download

# Build the application
build:
	go build -o bin/api cmd/api/main.go

# Run the application
run:
	@export DB_HOST=$(DB_HOST) && \
	export DB_PORT=$(DB_PORT) && \
	export DB_USER=$(DB_USER) && \
	export DB_PASSWORD=$(DB_PASSWORD) && \
	export DB_NAME=$(DB_NAME) && \
	export DB_SSL_MODE=$(DB_SSL_MODE) && \
	export APP_ENV=development && \
	go run cmd/api/main.go

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Database migrations (using golang-migrate)
migrate-up:
	migrate -path db/migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)" up

migrate-down:
	migrate -path db/migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)" down

# Docker commands for local development
docker-up:
	docker run --name postgres-marketplace \
		-e POSTGRES_USER=$(DB_USER) \
		-e POSTGRES_PASSWORD=$(DB_PASSWORD) \
		-e POSTGRES_DB=$(DB_NAME) \
		-p $(DB_PORT):5432 \
		-d postgres:15-alpine

docker-down:
	docker stop postgres-marketplace
	docker rm postgres-marketplace

# Create database (if it doesn't exist)
create-db:
	docker exec -it postgres-marketplace createdb -U $(DB_USER) $(DB_NAME) || true
