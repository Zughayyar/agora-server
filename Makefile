# Agora Server Makefile

prettier: fmt vet lint

# Default target
all: deps build

# Install dependencies
deps:
	@echo "📦 Installing dependencies..."
	go mod tidy
	go mod download
	@echo "✅ Dependencies installed successfully!"


# Build and run the server
run: build
	@echo "🚀 Starting Agora server..."
	./bin/server

# Build the server binary
build: clean
	@echo "🔨 Building Agora server binary..."
	mkdir -p bin
	go build -o bin/server ./cmd/server
	@echo "✅ Binary built successfully at bin/server"

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	rm -rf bin/ tmp/
	@echo "✅ Clean completed!"

# Run tests (when we add them later)
test:
	@echo "🧪 Running tests..."
	go test ./...
	@echo "✅ Tests completed!"

# Format code
fmt:
	@echo "🎨 Formatting code..."
	go fmt ./...
	@echo "✅ Code formatted!"

# Vet code for potential issues
vet:
	@echo "🔍 Vetting code for issues..."
	go vet ./...
	@echo "✅ Code vetting completed!"

# Run linter (requires golangci-lint to be installed)
lint:
	@echo "🔍 Running linter..."
	golangci-lint run
	@echo "✅ Linting completed!"

# Development with live reload using Air
dev:
	@echo "🔥 Starting development mode with live reload..."
	@mkdir -p tmp
	@if command -v air >/dev/null 2>&1; then \
		APP_ENV=development air; \
	else \
		echo "❌ Air not found. Installing..."; \
		go install github.com/air-verse/air@latest; \
		APP_ENV=development air; \
	fi

# Build migration tool
build-migrate: clean
	@echo "🔨 Building migration tool..."
	mkdir -p bin
	go build -o bin/migration ./cmd/migration
	@echo "✅ Migration tool built successfully at bin/migrate"

# Database Migration Commands
migrate: build-migrate
	@echo "🗃️ Running database migrations..."
	./bin/migration -action=migrate

migrate-rollback: build-migrate
	@echo "↩️ Rolling back database migrations..."
	./bin/migration -action=rollback

migrate-status: build-migrate
	@echo "📊 Checking migration status..."
	./bin/migration -action=status

# Docker Commands
docker-build:
	@echo "🐳 Building Docker image..."
	docker build -t agora-server .

docker-run: docker-build
	@echo "🚀 Running Docker container..."
	docker run -p 3000:3000 --env-file .env agora-server

docker-compose-up:
	@echo "🐳 Starting services with Docker Compose..."
	docker-compose up -d

docker-compose-down:
	@echo "🛑 Stopping Docker Compose services..."
	docker-compose down

docker-compose-logs:
	@echo "📋 Showing Docker Compose logs..."
	docker-compose logs -f

docker-compose-migrate:
	@echo "🗃️ Running migrations via Docker Compose..."
	docker-compose run --rm server ./bin/migration -action=migrate

# Deployment Commands
deploy-check:
	@echo "🔍 Checking deployment readiness..."
	@if [ ! -f .env ]; then echo "❌ .env file missing"; exit 1; fi
	@echo "✅ .env file exists"
	@if [ ! -f Dockerfile ]; then echo "❌ Dockerfile missing"; exit 1; fi
	@echo "✅ Dockerfile exists"
	@if [ ! -f docker-compose.yml ]; then echo "❌ docker-compose.yml missing"; exit 1; fi
	@echo "✅ docker-compose.yml exists"
	@echo "🎉 Deployment ready!"

.PHONY: all deps run build clean test fmt vet lint prettier dev build-migrate migrate migrate-rollback migrate-status docker-build docker-run docker-compose-up docker-compose-down docker-compose-logs docker-compose-migrate deploy-check
