# Simple Music Streaming API with Go (Fiber) + PostgreSQL

This project is a monolithic REST API implementation built with **Go** using the **Fiber** framework and **PostgreSQL** as the primary database. It provides a simple music streaming system without complex features, designed as a starting point for learning and exploring **DevOps** practices such as Microservices, Kubernetes, CI/CD pipelines, and more.

**Key Features**

- Tech-driven code structure to maintain clean and scalable architecture.
- Internal JWT authentication with GitHub OAuth integration.
- Automatically generated API documentation using Swagger version 2.0.
- Live-reloading support during development using Air.
- Database migrations are managed through an integrated manual migration system.
- Dockerized for easy deployment and testing.

**Purpose**

This application is built from scratch and is ideal for developers looking to learn modern backend development with Go. It is a foundational project before transitioning to Microservices Architecture and Fully Embracing DevOps Practices.

## Getting Started

### Project Structure

```bash
.
├── docker/
│   ├── Dockerfile.dev          # Development Dockerfile
│   └── Dockerfile.prod         # Production Dockerfile
├── migrations/                 # SQL migration folder
├── .env.development            # Environment file for development
├── .env.production             # Environment file for production
├── docker-compose-dev.yml      # Docker Compose Development configuration
├── docker-compose-prod.yml     # Docker Compose Production configuration
└── ...
```

---

### Development Environment

#### 1. Build and Run
Ensure .env.development exists and contains the appropriate environment variables.

```bash
docker compose -f docker-compose.dev.yml up -d --build
```

This command will:

- Build the image from `Dockerfile.dev`
- Start the following services:
  - `app`: Go application
  - `db`: PostgreSQL

#### 2. View Logs

```bash
docker compose -f docker-compose.dev.yml logs
```

---

### Production Deployment

#### 1. Build and Run

Ensure `.env.production` exists and contains the appropriate environment variables.

```bash
docker compose -f docker-compose.prod.yml up -d --build
```

This will use:

- `Dockerfile.prod` for optimized production builds
- Your production-ready Go binary and static settings
- Detached mode (`-d`) for background running

---

### Run Migrations
Place your migration SQL files inside the `migrations/` folder using the format:

  ```
  001_create_users_table.up.sql
  001_create_users_table.down.sql
  ```

**Up (Apply Migrations)**

```bash
docker run --rm -v $(pwd)/migrations:/migrations \
  --network mulo-go-api_mulo-api-dev \
  migrate/migrate -path=/migrations \
  -database "postgres://tungtungsahur:tralalelotralalala@mulo-db-dev:5432/mulo_bombardino?sslmode=disable" up
```

---

### Run Swagger Docs

```bash
swag init -g cmd/main.go --parseDependency --parseInternal
```

---

### Running Tests

This project uses Go's standard testing framework along with the [Testify](https://github.com/stretchr/testify) package for writing expressive unit tests.

#### ✅ Run All Tests in `app/services`

To run all test files inside the `app/services` directory and its sub-packages, use:

```bash
go test -v ./app/services/...
```
- **-v** Verbose output
- **./app/services/...** Run all tests recursively under app/services