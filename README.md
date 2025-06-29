# 🍽️ Agora Restaurant Management API

A modern, scalable Go REST API server for restaurant menu management built with best practices, featuring comprehensive CRUD operations, soft delete support, and advanced filtering capabilities.

## 🎯 Overview

The Agora Restaurant Management API is a production-ready Go application that provides a complete solution for managing restaurant menu items. It features a clean architecture with proper separation of concerns, comprehensive error handling, and modern development practices.

### ✨ Key Features

- **🍽️ Menu Item Management**: Full CRUD operations with validation
- **📂 Category Organization**: Items categorized as appetizer, main, dessert, drink, side, or fast food
- **🔄 Soft Delete Support**: Delete and restore items without data loss
- **🔍 Advanced Filtering**: Search by name, filter by category, availability, and more
- **💰 Price Management**: Decimal-based pricing with validation
- **📊 Swagger Documentation**: Interactive API documentation
- **🐳 Docker Support**: Containerized deployment with Docker Compose
- **🚀 CI/CD Pipeline**: Automated deployment to AWS via GitHub Actions
- **📈 Health Monitoring**: Multiple health check endpoints
- **🔒 Production Ready**: Graceful shutdown, connection pooling, and error handling

## 🚀 Quick Start

### Prerequisites

- Go 1.23 or higher
- PostgreSQL 16.9 or higher
- Make (optional, for using Makefile commands)
- Docker & Docker Compose (for containerized setup)

### Installation

1. **Clone and setup**:

   ```bash
   git clone <repository-url>
   cd agora-server
   cp env.example .env
   # Edit .env with your database configuration
   ```

2. **Install dependencies**:

   ```bash
   make deps
   ```

3. **Start with Docker Compose** (recommended):

   ```bash
   docker-compose up -d
   ```

4. **Or start development with live reload**:

   ```bash
   make dev
   ```

5. **Test the API**:
   ```bash
   curl http://localhost:3000/health
   ```

## 📡 API Endpoints

### Health Check

- **GET** `/health` - Root health check
- **GET** `/api/v1/health` - Versioned health check with database status

Both endpoints return:

```json
{
  "service": "agora-server",
  "status": "healthy",
  "timestamp": "2025-06-28T18:44:41.864+03:00"
}
```

### Menu Items Management

#### Core CRUD Operations

- **GET** `/api/v1/items` - List all menu items
- **POST** `/api/v1/items` - Create a new menu item
- **GET** `/api/v1/items/{id}` - Get specific menu item
- **PUT** `/api/v1/items/{id}` - Update menu item
- **DELETE** `/api/v1/items/{id}` - Soft delete menu item

#### Advanced Operations

- **GET** `/api/v1/items/category/{category}` - Filter by category
- **GET** `/api/v1/items/deleted` - List soft-deleted items
- **POST** `/api/v1/items/{id}/restore` - Restore deleted item

#### Query Parameters

- `?category=main` - Filter by category
- `?available=true` - Show only available items
- `?include_deleted=true` - Include soft-deleted items
- `?search=pizza` - Search items by name

### API Documentation

- **GET** `/swagger/` - Interactive Swagger UI documentation

## 🍽️ Menu Item Structure

```json
{
  "id": 1,
  "name": "Margherita Pizza",
  "description": "Classic tomato and mozzarella pizza",
  "price": "15.99",
  "category": "main",
  "is_available": true,
  "created_at": "2025-06-28T18:44:41.864+03:00",
  "updated_at": "2025-06-28T18:44:41.864+03:00"
}
```

### Supported Categories

- `appetizer` - Starters and appetizers
- `main` - Main course dishes
- `dessert` - Sweet treats and desserts
- `drink` - Beverages and drinks
- `side` - Side dishes and accompaniments
- `fast food` - Quick service items

## 🛠️ Available Commands

### Development

```bash
make dev       # 🔥 Start development with live reload (recommended)
make prettier  # ✨ Format, vet, and lint code
make test      # 🧪 Run tests
```

### Build & Run

```bash
make run       # 🚀 Build and run production binary
make build     # 🔨 Build binary to bin/server
make clean     # 🧹 Clean build artifacts
```

### Database Operations

```bash
make migrate           # 🗃️ Run database migrations
make migrate-rollback  # ↩️ Rollback last migration
make migrate-status    # 📊 Check migration status
```

### Docker Operations

```bash
docker-compose up -d   # 🐳 Start all services
docker-compose down    # 🛑 Stop all services
docker-compose logs    # 📋 View service logs
```

### Code Quality

```bash
make fmt       # 🎨 Format code
make vet       # 🔍 Vet code for issues
make lint      # 🔍 Run linter (requires golangci-lint)
```

## 🔧 Configuration

### Local Development

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

# Database Connection Pool
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5
DB_CONN_MAX_LIFETIME_MINUTES=15
DB_CONN_MAX_IDLE_TIME_MINUTES=5
```

### Production Deployment

For production deployment, environment variables are managed through GitHub repository secrets. See the deployment section for the complete list of required secrets.

## 🏗️ Technical Architecture

### Project Structure

```
agora-server/
├── cmd/                    # Application entry points
│   ├── server/            # Main server application
│   └── migration/         # Database migration tool
├── internal/              # Private application code
│   ├── database/          # Database models and migrations
│   ├── handlers/          # HTTP request handlers
│   ├── middlewares/       # HTTP middlewares
│   ├── routers/           # Route definitions
│   └── services/          # Business logic layer
├── docs/                  # Swagger documentation
├── docker-compose.yml     # Docker services configuration
└── Makefile              # Development commands
```

### Technology Stack

- **Language**: Go 1.23+
- **Framework**: Standard `net/http` with custom routing
- **Database**: PostgreSQL 16.9 with Bun ORM
- **Documentation**: Swagger/OpenAPI 3.0
- **Containerization**: Docker & Docker Compose
- **CI/CD**: GitHub Actions
- **Deployment**: AWS EC2

### Key Libraries

- **Bun ORM**: Type-safe database operations
- **Swaggo**: Swagger documentation generation
- **Godotenv**: Environment variable management
- **Decimal**: Precise decimal arithmetic for pricing

## 🐳 Docker Support

### Quick Start with Docker

```bash
# Start all services (PostgreSQL + API)
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

### Docker Services

- **postgres**: PostgreSQL 16.9 database
- **migrate**: Database migration runner
- **server**: Agora API server

## 🚀 Deployment

The application uses GitHub Actions for automated deployment to AWS. The deployment process:

1. **Builds** the Docker image on every push to main
2. **Pushes** the image to GitHub Container Registry
3. **Deploys** to AWS server using the pre-built image
4. **Creates** the `.env` file from GitHub secrets on the server
5. **Starts** the application using Docker Compose

### Deployment Benefits

- ✅ **No git pull required** on the AWS server
- ✅ **Environment variables** managed in GitHub secrets
- ✅ **Immutable deployments** using pre-built Docker images
- ✅ **Rollback capability** by changing image tags
- ✅ **Secure** - no sensitive data in version control

### Required GitHub Repository Secrets

Navigate to your repository settings → Secrets and variables → Actions, and add:

**Application Configuration:**

- `APP_ENV`: `production`
- `APP_PORT`: `3000`
- `APP_VERSION`: `1.0.0`

**Database Configuration:**

- `DB_NAME`: Your database name
- `DB_USER`: Your database user
- `DB_PASSWORD`: Your secure database password
- `DB_PORT`: `5432`
- `DB_SSL_MODE`: `disable` (or `require` for production)
- `DB_LOG_QUERIES`: `false`

**Database Connection Pool:**

- `DB_MAX_OPEN_CONNS`: `25`
- `DB_MAX_IDLE_CONNS`: `5`
- `DB_CONN_MAX_LIFETIME_MINUTES`: `15`
- `DB_CONN_MAX_IDLE_TIME_MINUTES`: `5`

**AWS Deployment:**

- `AWS_IP`: Your AWS server IP
- `AWS_USER`: `ubuntu`
- `AWS_SSH_KEY`: Your SSH private key

## 🔥 Development Workflow

1. **Start development**:

   ```bash
   make dev
   ```

2. **Edit your code** - Air automatically rebuilds and restarts the server

3. **Test your changes** - Server restarts instantly on file changes

4. **Quality checks**:

   ```bash
   make prettier  # Format, vet, and lint
   ```

5. **Deploy** - Push to main branch triggers automatic deployment

## 📚 API Examples

### Create a Menu Item

```bash
curl -X POST http://localhost:3000/api/v1/items \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Margherita Pizza",
    "description": "Classic tomato and mozzarella pizza",
    "price": "15.99",
    "category": "main",
    "is_available": true
  }'
```

### List All Items

```bash
curl http://localhost:3000/api/v1/items
```

### Filter by Category

```bash
curl "http://localhost:3000/api/v1/items?category=main"
```

### Search Items

```bash
curl "http://localhost:3000/api/v1/items?search=pizza"
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run quality checks: `make prettier`
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🆘 Support

For support and questions:

- 📧 Email: support@agora-restaurant.com
- 🌐 Website: https://agora-restaurant.com/support
- 📖 Documentation: http://localhost:3000/swagger/ (when running locally)
