# user-service

REST API service for user authentication and profile management. Part of the bookmark microservices architecture.

## Overview

**user-service** provides JWT-based authentication, user registration, login, and profile management. It uses PostgreSQL for persistence and Redis for rate limiting. As the **JWT issuer**, it signs tokens with an RSA private key; all other services (e.g. `bookmark-service`) validate tokens using only the matching public key.

Migrations run **automatically on startup** via `sqldb.RunMigration` in `NewContainer`. The manual CLI (`cmd/migrate`) is available for rollbacks and recovery.

## Tech Stack

| Component | Technology | Version |
|---|---|---|
| Language | Go | 1.26 |
| Web framework | Gin | v1.12.0 |
| Database | PostgreSQL (GORM) | v1.31.1 / v1.6.0 |
| Cache / Rate limit | Redis (go-redis) | v9.19.0 |
| Auth | JWT RS256 (issuer) | v5.3.1 |
| Password hashing | bcrypt | golang.org/x/crypto v0.52.0 |
| Logger | Zerolog | v1.35.1 |
| Shared library | bookmark-common | v0.3.0 |
| API docs | Swagger (swaggo/gin-swagger) | v1.6.1 |
| Testing | Testify + SQLite | v1.11.1 |

## Quick Start

```bash
cd user-service
cp .env.example .env     # edit as needed
make gen-keys-local      # generate RSA keypair → keys/private.pem + keys/public.pem
make run                 # starts on :8080; migrations run automatically
```

Share `keys/public.pem` with `bookmark-service` (copy to `bookmark-service/keys/public.pem`). Never share `keys/private.pem`.

Swagger UI: `http://localhost:8080/swagger/index.html`  
API base: `http://localhost:8080/api/user_service/v1`

### Environment variables

| Variable | Default | Description |
|---|---|---|
| `APP_PORT` | `8080` | HTTP listen port |
| `SERVICE_NAME` | _(required)_ | Service identifier |
| `APP_HOST_NAME` | `/api/user_service` | API base path |
| `APP_ENV` | `development` | Environment (`development` / `production`) |
| `DB_HOST` | `localhost` | PostgreSQL host |
| `DB_PORT` | `5432` | PostgreSQL port |
| `DB_USER` | — | PostgreSQL user |
| `DB_PASSWORD` | — | PostgreSQL password |
| `DB_NAME` | `user_db` | PostgreSQL database |
| `DB_SSLMODE` | `disable` | PostgreSQL SSL mode |
| `DB_TIMEZONE` | `UTC` | PostgreSQL timezone |
| `REDIS_ADDR` | `localhost:6379` | Redis address |
| `REDIS_PASSWORD` | — | Redis password |
| `REDIS_DATABASE` | `0` | Redis logical DB index |
| `RATELIMIT_LIMIT` | `20` | Requests per window |
| `RATELIMIT_WINDOW` | `10s` | Rate-limit window |
| `JWT_PRIVATE_KEY_PATH` | `keys/private.pem` | RSA private key (sign tokens) |
| `JWT_PUBLIC_KEY_PATH` | `keys/public.pem` | RSA public key (verify tokens) |
| `JWT_ISSUER` | `user-service` | JWT `iss` claim |
| `JWT_AUDIENCE` | `bookmark-app` | JWT `aud` claim |
| `JWT_EXPIRATION_SECONDS` | `3600` | Token lifetime in seconds |

## API Endpoints

All routes under `/api/user_service`.

### Health
| Method | Path | Auth | Description |
|---|---|---|---|
| GET | `/health-check` | — | Pings PostgreSQL |

### Auth (under `/v1`)
| Method | Path | Auth | Description |
|---|---|---|---|
| POST | `/v1/users/register` | — | Register new user |
| POST | `/v1/users/login` | — | Login, returns JWT |

### Profile (under `/v1`, require `Authorization: Bearer <token>`)
| Method | Path | Description |
|---|---|---|
| GET | `/v1/self/info` | Get own profile |
| PUT | `/v1/self/info` | Update own profile |

## Database Migrations

Migrations run automatically on startup. The `cmd/migrate` binary provides manual control:

```bash
make migrate-up              # apply all pending
make migrate-up STEPS=1      # apply 1 step
make migrate-down STEPS=1    # roll back 1 step
make migrate-version         # show current version + dirty flag
make migrate-force           # clear dirty flag (interactive prompt)
```

### Migration history

| Version | File | Description |
|---|---|---|
| 0001 | `create_users_table` | `users` table with email/username unique constraints |

## Testing

```bash
make test              # unit + integration tests, 80% coverage gate
make test-coverage     # open HTML coverage report
make docker-test       # test inside Docker (CI parity)
```

Integration tests use SQLite — no external services required.

**Coverage threshold: 80%** on business logic. Infrastructure packages (`cmd`, `bootstrap`, `api`, `dto`, `model`, `repository/ping`) are excluded from the threshold but still scanned by SonarCloud.

## Make Targets

```
Development:
  make run             Run locally (auto-migrates on start)
  make dev             fmt → vet → test → swagger → run
  make fmt / vet / lint / tidy / vendor

Database:
  make migrate-up [STEPS=n]
  make migrate-down [STEPS=n]
  make migrate-version
  make migrate-force

Testing:
  make test
  make test-coverage

Build:
  make build / build-linux / build-macos / build-windows / build-prod / release

Mocks:
  make generate-mocks
  make clean-mocks

Docker / CI:
  make docker-test / docker-sonar / docker-build-push
  make docker-run / docker-stop / docker-logs / docker-shell / docker-clean

Keys:
  make gen-keys-local  Generate RSA keypair locally

Utilities:
  make swagger / install-tools / info / clean / clean-docs / clean-all
```

## CI/CD

| Trigger | CI | CD |
|---|---|---|
| PR to `main` | test + SonarCloud (no push) | — |
| Push to `main` | test + SonarCloud + build + push `main`/`<sha7>` tags | deploy via self-hosted runner |
| Git tag `v*.*.*` | test + SonarCloud + build + push `<tag>` + `latest` | deploy via self-hosted runner |

CD runner working directory: `/opt/bookmark-system`. Updates `USER_SERVICE_TAG` in `.env` and runs `docker compose up -d --force-recreate user-service`.

## Project Structure

```
user-service/
├── cmd/
│   ├── api/main.go           # HTTP server entrypoint
│   └── migrate/main.go       # Migration CLI (up/down/version/force)
├── internal/
│   ├── api/                  # router.go, swagger.go
│   ├── bootstrap/            # app.go, container.go (DI), config.go, routes.go
│   ├── dto/                  # auth/, profile/, health/ request+response structs
│   ├── handler/              # auth/, profile/, health/ HTTP handlers
│   ├── model/                # base.go (UUID PK), user.go
│   ├── repository/
│   │   ├── user/             # PostgreSQL read/write + mocks
│   │   └── ping/             # sqldb.go pinger + mocks
│   ├── service/
│   │   ├── auth/             # register, login + mocks
│   │   ├── profile/          # get, update + mocks
│   │   └── health/           # check + mocks
│   └── test/
│       ├── integration/      # end-to-end handler tests
│       └── fixtures/         # SQLite testdb helpers
├── migrations/               # 000001 up/down SQL files
├── keys/                     # private.pem + public.pem (git-ignored)
├── docs/                     # swagger generated output
├── Makefile
├── Dockerfile
├── sonar-project.properties
└── .github/workflows/ci.yaml, cd.yaml
```
