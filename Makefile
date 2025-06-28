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
	./bin/agora-server

# Build the server binary
build: clean
	@echo "🔨 Building Agora server binary..."
	mkdir -p bin
	go build -o bin/agora-server ./cmd/agora
	@echo "✅ Binary built successfully at bin/agora-server"

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
	go build -o bin/migrate ./cmd/migrate
	@echo "✅ Migration tool built successfully at bin/migrate"

# Database Migration Commands
migrate: build-migrate
	@echo "🗃️ Running database migrations..."
	./bin/migrate -action=migrate

migrate-rollback: build-migrate
	@echo "↩️ Rolling back database migrations..."
	./bin/migrate -action=rollback

migrate-status: build-migrate
	@echo "📊 Checking migration status..."
	./bin/migrate -action=status

# Verify database schema
verify-db:
	@echo "🔍 Verifying database schema..."
	go build -o bin/verify cmd/verify/main.go
	./bin/verify

# Docker database commands
db-up:
	@echo "🐘 Starting PostgreSQL database..."
	docker-compose up -d postgres

db-down:
	@echo "🛑 Stopping PostgreSQL database..."
	docker-compose stop postgres

db-logs:
	@echo "📜 Showing database logs..."
	docker-compose logs -f postgres

# Full Docker commands
docker-up:
	@echo "🐳 Starting all services with Docker Compose..."
	docker-compose up -d

docker-down:
	@echo "🛑 Stopping all Docker services..."
	docker-compose down

docker-logs:
	@echo "📜 Showing all Docker logs..."
	docker-compose logs -f

# Environment setup
env-setup:
	@if [ ! -f .env ]; then \
		echo "📝 Creating .env file from template..."; \
		cp env.example .env; \
		echo "✅ Created .env file. Please update it with your configuration."; \
	else \
		echo "⚠️ .env file already exists"; \
	fi

# Full setup for new developers
setup: env-setup db-up
	@echo "⏳ Waiting for database to be ready..."
	@sleep 5
	@$(MAKE) migrate
	@echo "🎉 Setup complete! You can now run 'make dev' to start development"
