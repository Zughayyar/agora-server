# 🍽️ Agora Restaurant Management API

A simple and scalable Go REST API server for restaurant management following Go best practices.

## 🚀 Quick Start

### Prerequisites

- Go 1.21 or higher
- Make (optional, for using Makefile commands)

### Installation

1. **Install dependencies**:

   ```bash
   go mod tidy
   ```

2. **Build and run the server**:

   ```bash
   make run
   ```

   This will build the binary to `bin/agora-server` and run it.

3. **Test the API**:

   ```bash
   curl http://localhost:8080/
   ```

## 📡 API Endpoints

### Health Check

- **GET** `/` - Root health check
- **GET** `/api/v1/health` - Versioned health check

Both endpoints return:

```json
{
    "message": "Hello from Agora Restaurant Management API! 🍽️",
    "service": "agora-server",
    "status": "healthy",
    "timestamp": "2025-06-22T22:59:39.569107+03:00"
}
```

## 🏗️ Project Structure

```text
server/
├── cmd/                    # Application entry points
│   └── agora-server/      # Main application
│       └── main.go        # Application entry point
├── internal/              # Private application code
│   ├── handlers/          # HTTP request handlers
│   │   └── health.go      # Health check handler
│   ├── middleware/        # HTTP middleware
│   │   └── middleware.go  # Logging and CORS middleware
│   └── router/            # Route configuration
│       └── router.go      # Route setup and organization
├── bin/                   # Built binaries (excluded from git)
├── go.mod                 # Go module definition
├── go.sum                 # Go module checksums
├── Makefile              # Build and run commands
├── .gitignore            # Git ignore rules
└── README.md             # This file
```

This structure follows the [Standard Go Project Layout](https://github.com/golang-standards/project-layout):

- `cmd/` - Main applications for this project
- `internal/` - Private application and library code
- `bin/` - Binary output directory (auto-created)

## 🛠️ Available Commands

Using the Makefile:

```bash
make deps      # Install dependencies
make run       # Build and run the server
make build     # Build binary to bin/agora-server
make start     # Build and run binary (alias for run)
make clean     # Clean build artifacts (removes bin/)
make fmt       # Format code
make vet       # Vet code for issues
make test      # Run tests (when added)
make dev-run   # Run directly without building binary
```

### Manual Commands

```bash
# Build manually
go build -o bin/agora-server ./cmd/agora-server

# Run development mode
go run ./cmd/agora-server

# Run built binary
./bin/agora-server
```

## 🔧 Configuration

The server uses environment variables for configuration:

- `PORT` - Server port (default: 8080)

Example:

```bash
PORT=3000 make run
```

## 🎯 Future Features

The current structure is designed to easily extend with:

- Database integration
- User authentication (JWT)
- Restaurant management endpoints
- Menu management
- Order processing
- Real-time updates

## 📝 Development

This server follows Go best practices:

- **Standard Go Project Layout** with `cmd/` and `internal/` directories
- **Clean architecture** with separated concerns
- **Middleware** for cross-cutting concerns
- **Structured logging** with request tracking
- **CORS support** for frontend integration
- **Graceful error handling**
- **Environment-based configuration**
- **Binary output isolation** in `bin/` directory

## 🔨 Build Information

- **Binary Output**: `bin/agora-server`
- **Main Package**: `./cmd/agora-server`
- **Module Name**: `agora-server`
- **Go Version**: 1.21+
