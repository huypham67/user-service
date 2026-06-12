# user-service

Production-ready REST API service for user authentication and profile management. Part of the bookmark microservices architecture, extracted from `bookmark-service-monolithic` to enable independent scaling and deployment.

## Overview

**user-service** provides JWT-based authentication, user registration and login, profile management (get/update), and health checks. It uses PostgreSQL for persistence and integrates with `bookmark-common` library for shared middleware, utilities, and infrastructure.

## 🎯 Features

- **User Registration**: Create new user accounts with validation
- **User Authentication**: Login with JWT token generation (RSA-based)
- **User Profiles**: Get and update user information
- **JWT Middleware**: Token validation and user context extraction
- **Health Checks**: Database connectivity monitoring
- **PostgreSQL Backend**: GORM ORM with migrations
- **Structured Logging**: Zerolog integration
- **Comprehensive Testing**: Unit and integration tests with 80% coverage gate
- **Docker Ready**: Multi-stage Dockerfile for testing and production
- **CI/CD Pipeline**: GitHub Actions with Docker layer caching and SonarCloud scanning
- **Environment Configuration**: Flexible setup via environment variables

## 📋 Tech Stack

| Component | Technology | Version |
|-----------|-----------|---------|
| Language | Go | 1.26 |
| Web Framework | Gin | v1.12.0 |
| Database | PostgreSQL | (via GORM) |
| ORM | GORM | v1.31.1 (postgres driver v1.6.0) |
| Auth | JWT (RSA) | v5.3.1 |
| Password Hashing | bcrypt | (golang.org/x/crypto) |
| Logger | Zerolog | v1.35.1 |
| Shared Library | bookmark-common | local |
| Testing | Testify | v1.11.1 |
| Migrations | golang-migrate | v4.19.1 |

## 🚀 Quick Start

### Prerequisites

- **Go 1.26** or higher
- **PostgreSQL 12+** (required)
- **Make** or **PowerShell** (Windows)

### Installation

```bash
cd user-service
go mod download
go mod tidy
```

### Environment Setup

Create `.env` file:

```env
APP_PORT=8081
SERVICE_NAME=user-service
JWT_PRIVATE_KEY_PATH=keys/private.pem
JWT_PUBLIC_KEY_PATH=keys/public.pem
DB_HOST=localhost
DB_PORT=5432
DB_USER=admin
DB_PASSWORD=admin
DB_NAME=user_db
```

### Database Setup

```bash
createdb user_db
make run-migrations
```

### Running

```bash
make run           # Build and run
make test          # Run tests
```

API available at: `http://localhost:8081/api/user_service/v1`

## 🔌 API Endpoints

### Base URL
```
http://localhost:8081/api/user_service/v1
```

#### Register User
```http
POST /users/register
Content-Type: application/json

{
  "display_name": "John Doe",
  "username": "john_doe",
  "email": "john@example.com",
  "password": "SecurePassword123"
}
```

#### Login
```http
POST /users/login
Content-Type: application/json

{
  "username": "john_doe",
  "password": "SecurePassword123"
}
```

Returns JWT token in response.

#### Get Profile (requires JWT)
```http
GET /self/info
Authorization: Bearer <token>
```

#### Update Profile (requires JWT)
```http
PUT /self/info
Authorization: Bearer <token>
Content-Type: application/json

{
  "display_name": "Updated Name",
  "email": "newemail@example.com"
}
```

#### Health Check
```http
GET /health-check
```

## 🧪 Testing

```bash
make test              # Local tests + coverage
make docker-test       # Docker tests (as CI)
make test-coverage     # View HTML coverage report
```

Coverage threshold: **80%** on business logic

## 🛠️ Available Make Targets

```bash
make help              # Show all targets
make test              # Run tests with coverage
make docker-test       # Test in Docker
make docker-sonar      # SonarCloud scan
make build             # Build binary
make run               # Run service
make fmt               # Format code
make vet               # Run vet
make lint              # Run linter
make clean             # Remove artifacts
```

## 🐳 Docker

### Build
```bash
docker build -t user-service:latest .
```

### Run
```bash
docker run -d \
  -e APP_PORT=8081 \
  -e DB_HOST=host.docker.internal \
  -e JWT_PRIVATE_KEY_PATH=/keys/private.pem \
  -e JWT_PUBLIC_KEY_PATH=/keys/public.pem \
  -v /path/to/keys:/keys \
  -p 8081:8081 \
  user-service:latest
```

## 🔄 CI/CD

GitHub Actions workflow:
1. Docker build → test (extract coverage)
2. Upload coverage to Codecov
3. SonarCloud quality scan

Secrets needed: `SONAR_TOKEN`

## 📁 Structure

```
user-service/
├── cmd/
│   ├── api/main.go           # API server
│   └── migrate/main.go        # DB migrations
├── internal/
│   ├── api/                   # HTTP routing
│   ├── bootstrap/             # DI configuration
│   ├── dto/                   # Request/response objects
│   ├── handler/               # HTTP handlers (auth, profile, health)
│   ├── model/                 # Domain models (User)
│   ├── repository/            # Database layer
│   ├── service/               # Business logic
│   └── test/                  # Integration tests
├── migrations/                # SQL migrations
├── keys/                      # JWT keys (git-ignored)
├── Makefile
├── Dockerfile
├── sonar-project.properties
└── .env (local, git-ignored)
```

## 🏗️ Architecture

Clean architecture with 4 layers:

```
Handler (HTTP I/O)
   ↓
Service (Business Logic)
   ↓
Repository (Data Access)
   ↓
Database (PostgreSQL)
```

## 🔐 Security

- **JWT**: RSA-based asymmetric encryption
- **Passwords**: bcrypt hashing
- **Input Validation**: Struct-tag based
- **SQL Injection**: Prevented via GORM ORM
- **Secrets**: Loaded from environment variables
- **CI Actions**: Pinned to commit SHAs

## 🔗 Integration

Uses shared library `bookmark-common`:
- JWT middleware for authentication
- Zerolog for structured logging
- Common error handling
- Database utilities

## 📄 License

Part of bookmark microservices architecture.

## 🔗 Useful Links

- [Go](https://golang.org/)
- [Gin Framework](https://gin-gonic.com/)
- [GORM](https://gorm.io/)
- [PostgreSQL](https://www.postgresql.org/)
- [JWT](https://jwt.io/)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)