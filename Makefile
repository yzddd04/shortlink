# Link Shortener Backend Makefile

.PHONY: help build run test clean docker-build docker-run migrate

# Default target
help:
	@echo "Available commands:"
	@echo "  build        - Build the application"
	@echo "  run          - Run the application"
	@echo "  test         - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  clean        - Clean build artifacts"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run with Docker Compose"
	@echo "  migrate      - Run database migrations"
	@echo "  lint         - Run linter"
	@echo "  fmt          - Format code"

# Build the application
build:
	@echo "Building application..."
	go build -o bin/server cmd/server/main.go

# Run the application
run:
	@echo "Running application..."
	go run cmd/server/main.go

# Run tests
test:
	@echo "Running tests..."
	go test ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -cover ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	go clean

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t link-shortener .

# Run with Docker Compose
docker-run:
	@echo "Running with Docker Compose..."
	docker-compose up --build

# Run database migrations
migrate:
	@echo "Running database migrations..."
	psql -d link_shortener -f migrations/001_init.sql

# Run linter
lint:
	@echo "Running linter..."
	golangci-lint run

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod tidy
	go mod download

# Generate API documentation
docs:
	@echo "Generating API documentation..."
	swag init -g cmd/server/main.go

# Setup development environment
setup: deps
	@echo "Setting up development environment..."
	@if [ ! -f .env ]; then cp env.example .env; echo "Created .env file from env.example"; fi
	@echo "Please edit .env file with your configuration"

# Development mode with hot reload (requires air)
dev:
	@echo "Starting development mode with hot reload..."
	@if command -v air > /dev/null; then air; else echo "Air not found. Install with: go install github.com/cosmtrek/air@latest"; fi
