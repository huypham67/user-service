# user-service

Production-ready REST API service for user authentication and profile management. Part of the bookmark microservices architecture, extracted from `bookmark-service-monolithic` to enable independent scaling and deployment.

## Overview

**user-service** provides JWT-based authentication, user registration and login, profile management (get/update), and health checks. It uses PostgreSQL for persistence and integrates with `bookmark-common` library for shared middleware, utilities, and infrastructure.

## 🎯 Features

- **User Registration**: Create new user accounts with validation
- **User Authentication**: Login with JWT token generation (RSA-based)
- **User Profiles**: Get and update user information
- **JWT Middleware**: Token validation and user context extraction
- **Rate Limiting**: Redis-backed request throttling on all endpoints
- **Health Checks**: PostgreSQL connectivity monitoring
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
| Cache / Rate limit | Redis (go-redis) | v9.19.0 |
| Auth | JWT (RSA, issuer) | v5.3.1 |
| Password Hashing | bcrypt | (golang.org/x/crypto) |
| Logger | Zerolog | v1.35.1 |
| Shared Library | bookmark-common | v0.2.0 |
| API Docs | Swagger (swaggo/gin-swagger) | v1.6.1 |
| Testing | Testify | v1.11.1 |
| Migrations | golang-migrate | v4.19.1 |

## 🚀 Quick Start

### Prerequisites

- **Go 1.26** or higher
- **PostgreSQL** (required) and **Redis** (rate limiting)

### Run

```bash
cd user-service
go mod download
make gen-keys-local    # RSA key pair — user-service is the JWT signer/issuer
createdb user_db
make migrate-up        # apply migrations
make run               # build + run
make test              # local tests + coverage (80% gate)
```

API base: `http://localhost:8080/api/user_service/v1` · Swagger UI: `http://localhost:8080/swagger/index.html` (`make swagger` to regenerate docs)

### Environment (`.env`)

```env
APP_PORT=8080
SERVICE_NAME=user-service
APP_HOST_NAME=/api/user_service
DB_HOST=localhost
DB_PORT=5432
DB_USER=admin
DB_PASSWORD=admin
DB_NAME=user_db
REDIS_ADDR=localhost:6379
RATELIMIT_LIMIT=20
RATELIMIT_WINDOW=10s
JWT_PRIVATE_KEY_PATH=keys/private.pem
JWT_PUBLIC_KEY_PATH=keys/public.pem
JWT_ISSUER=user-service
JWT_AUDIENCE=bookmark-app
JWT_EXPIRATION_SECONDS=3600
```

## 🔌 API Endpoints

### Base URL
```
http://localhost:8080/api/user_service/v1
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
GET /api/user_service/health-check   # service root (not under /v1); pings PostgreSQL
```

## 🧪 Testing

```bash
make test              # Local tests + coverage
make docker-test       # Docker tests (as CI)
make test-coverage     # View HTML coverage report
```

Coverage threshold: **80%** on business logic

## 🛠️ Make Targets

Run `make help` for the full list. Grouped:

- **Dev**: `run` · `dev` (fmt→vet→test→swagger→run) · `fmt` · `vet` · `lint` · `tidy` · `vendor`
- **Database**: `migrate-up` · `migrate-down` · `migrate-force` · `migrate-version`
- **Testing**: `test` · `test-coverage`
- **Build**: `build` · `build-linux` · `build-macos` · `build-windows` · `build-prod` · `release`
- **Mocks**: `generate-mocks` · `clean-mocks`
- **Docker / CI**: `docker-test` · `docker-sonar` · `docker-build-push` · `docker-run` · `docker-stop` · `docker-logs` · `docker-shell` · `docker-clean`
- **Utilities**: `swagger` · `gen-keys-local` · `install-tools` · `info` · `clean` · `clean-all`

The **Makefile is the single source of truth** for coverage/quality-gate exclusions (`INFRA_DIRS` / `SYSTEM_DIRS`); `sonar-project.properties` carries identity/scope only.

## 🐳 Docker

```bash
make docker-build-push   # multi-stage build (base → build → test-exec → test → final)
make docker-run          # run the image with --env-file .env on port 8080
make docker-test         # run the test stage and extract coverage (as CI does)
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

Consumes `github.com/huypham67/bookmark-common` v0.2.0 for JWT middleware/provider, Redis-backed rate limiting, password hashing, structured logging, request/response helpers, and SQL/Redis clients. As the **JWT issuer**, user-service signs tokens with its private key; other services (e.g. `bookmark-service`) validate them with the matching public key.