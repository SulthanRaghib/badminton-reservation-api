.PHONY: help install run build test clean docker-build docker-run migrate seed

# Variables
APP_NAME=badminton-reservation-api
DOCKER_IMAGE=badminton-api:latest

help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

install: ## Install dependencies
	@echo "📦 Installing dependencies..."
	go mod download
	go mod tidy

run: ## Run the application
	@echo "🚀 Starting application..."
	go run main.go

build: ## Build the application
	@echo "🔨 Building application..."
	go build -o bin/$(APP_NAME) main.go
	@echo "✅ Build complete: bin/$(APP_NAME)"

test: ## Run tests
	@echo "🧪 Running tests..."
	go test -v ./...

clean: ## Clean build artifacts
	@echo "🧹 Cleaning..."
	rm -rf bin/
	rm -f $(APP_NAME)
	go clean

docker-build: ## Build Docker image
	@echo "🐳 Building Docker image..."
	docker build -t $(DOCKER_IMAGE) .

docker-run: ## Run Docker container
	@echo "🐳 Running Docker container..."
	docker-compose up -d

docker-stop: ## Stop Docker container
	@echo "🛑 Stopping Docker container..."
	docker-compose down

docker-logs: ## Show Docker logs
	@echo "📋 Showing logs..."
	docker-compose logs -f

dev: ## Run in development mode with hot reload (requires air)
	@echo "🔥 Starting development server with hot reload..."
	@if ! command -v air > /dev/null; then \
		echo "Installing air..."; \
		go install github.com/cosmtrek/air@latest; \
	fi
	air

lint: ## Run linter
	@echo "🔍 Running linter..."
	@if ! command -v golangci-lint > /dev/null; then \
		echo "Installing golangci-lint..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi
	golangci-lint run

fmt: ## Format code
	@echo "✨ Formatting code..."
	go fmt ./...

setup: install ## Initial setup
	@echo "⚙️  Setting up project..."
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo "✅ Created .env file. Please configure it!"; \
	else \
		echo "⚠️  .env already exists"; \
	fi
	@echo "✅ Setup complete!"

db-migrate: ## Show migration instructions
	@echo "📊 To run migrations, execute the SQL files in this order:"
	@echo "  1. database/migrations/001_create_courts.sql"
	@echo "  2. database/migrations/002_create_timeslots.sql"
	@echo "  3. database/migrations/003_create_reservations.sql"
	@echo "  4. database/migrations/004_create_payments.sql"
	@echo ""
	@echo "Run these in Neon Console SQL Editor or via psql"

gorm-migrate: ## Run GORM AutoMigrate using the small migrate CLI
	@echo "Running GORM AutoMigrate..."
	go run ./cmd/migrate

db-seed-run: ## Run SQL seed file using seed CLI
	@echo "Running seed data..."
	go run ./cmd/seed

db-seed: ## Show seed instructions
	@echo "🌱 To seed data, execute:"
	@echo "  database/seeds/seed_data.sql"
	@echo ""
	@echo "Run this in Neon Console SQL Editor or via psql after migrations"