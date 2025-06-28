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
