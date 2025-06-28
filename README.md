# 🍽️ Agora Restaurant Management API

A simple and scalable Go REST API server for restaurant management following Go best practices.

## 🚀 Quick Start

### Prerequisites

- Go 1.21 or higher
- Make (optional, for using Makefile commands)

### Installation

1. **Install dependencies**:

   ```bash
   make deps
   ```

2. **Start development with live reload**:

   ```bash
   make dev
   ```

3. **Test the API**:
   ```bash
   curl http://localhost:3000/health
   ```

## 📡 API Endpoints

### Health Check

- **GET** `/health` - Root health check
- **GET** `/api/v1/health` - Versioned health check

Both endpoints return:

```json
{
  "service": "agora-server",
  "status": "healthy",
  "timestamp": "2025-06-28T18:44:41.864+03:00"
}
```

## 🛠️ Available Commands

### Development

```bash
make dev       # 🔥 Start development with live reload (recommended)
```

### Build & Run

```bash
make run       # 🚀 Build and run production binary
make build     # 🔨 Build binary to bin/server
make clean     # 🧹 Clean build artifacts
```

### Database Migrations

```bash
make migrate           # 🗃️ Run database migrations
make migrate-rollback  # ↩️ Rollback last database migration
make migrate-status    # 📊 Check migration status
make build-migrate     # 🔨 Build migration tool only
```

### Code Quality

```bash
make fmt       # 🎨 Format code
make vet       # 🔍 Vet code for issues
make lint      # 🔍 Run linter (requires golangci-lint)
make prettier  # ✨ Run fmt + vet + lint
make test      # 🧪 Run tests
```

### Setup

```bash
make deps      # 📦 Install dependencies
make all       # 🚀 Install deps + build (default)
```

## 🔧 Configuration

Create a `.env` file from the template:

```bash
cp env.example .env
```

Edit the `.env` file with your configuration:

```bash
# Application Configuration
APP_ENV=development
APP_PORT=3000
APP_VERSION=1.0.0

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_NAME=agora_db
DB_USER=agora_user
DB_PASSWORD=agora_password
DB_SSL_MODE=disable
```

## 🐳 Docker Support

Start with PostgreSQL database:

```bash
docker-compose up -d
```

## 🔥 Development Workflow

1. **Start live reload development**:

   ```bash
   make dev
   ```

2. **Edit your code** - Air automatically rebuilds and restarts the server

3. **Test your changes** - Server restarts instantly on file changes

4. **Quality checks**:
   ```bash
   make prettier  # Format, vet, and lint
   ```
