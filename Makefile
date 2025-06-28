# Agora Server Makefile

.PHONY: run build clean test deps dev dev-run

# Default target
all: deps build

# Install dependencies
deps:
	@echo "📦 Installing dependencies..."
	go mod tidy
	go mod download
	@echo "✅ Dependencies installed successfully!"

# Setup environment file
setup-env:
	@echo "⚙️ Setting up environment file..."
	@if [ ! -f .env ]; then \
		cp env.example .env; \
		echo "✅ .env file created from env.example"; \
		echo "🔧 Please edit .env file with your specific configuration"; \
	else \
		echo "⚠️ .env file already exists"; \
	fi

# Build and run the server
run: build
	mkdir -p bin
	@echo "🚀 Starting Agora server..."
	./bin/agora-server

# Build the server binary
build:
	@echo "🔨 Building Agora server binary..."
	mkdir -p bin
	go build -o bin/agora-server ./cmd/agora
	@echo "✅ Binary built successfully at bin/agora-server"

# Run the built binary
start: build
	@echo "🚀 Starting Agora server..."
	./bin/agora-server

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
		echo "📍 Using Air from Go bin directory..."; \
		APP_ENV=development $(shell go env GOPATH)/bin/air; \
	fi

# Run in development mode without building binary
dev-run:
	@echo "🚀 Running in development mode..."
	APP_ENV=development go run ./cmd/agora

# Docker commands
docker-build:
	@echo "🐳 Building Docker image..."
	docker build -t agora-server .
	@echo "✅ Docker image built successfully!"

docker-up:
	@echo "🐳 Starting services with Docker Compose..."
	@if [ ! -f .env ]; then \
		echo "⚠️ No .env file found. Creating one from env.example..."; \
		cp env.example .env; \
		echo "🔧 Please edit .env file with your configuration before running again"; \
		exit 1; \
	fi
	docker-compose up -d
	@echo "✅ Services started successfully!"

docker-down:
	@echo "🐳 Stopping Docker services..."
	docker-compose down
	@echo "✅ Services stopped successfully!"

docker-logs:
	@echo "📋 Showing Docker logs..."
	docker-compose logs -f

docker-restart:
	@echo "🔄 Restarting Docker services..."
	docker-compose restart
	@echo "✅ Services restarted successfully!"

docker-clean:
	@echo "🧹 Cleaning Docker resources..."
	docker-compose down -v
	docker system prune -f
	@echo "✅ Docker cleanup completed!" 