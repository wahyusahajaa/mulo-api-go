# Mulo API

## Project Structure

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

## Getting Started

### Development

#### 1. Build and Run with Development Settings

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

### Production

#### 1. Build and Run with Production Settings

Ensure `.env.production` exists and contains the appropriate environment variables.

```bash
docker compose -f docker-compose.prod.yml up -d --build
```

This will use:

- `Dockerfile.prod` for optimized production builds
- Your production-ready Go binary and static settings
- Detached mode (`-d`) for background running

---

### Running Migrations Manually

#### Up (Apply Migrations)

```bash
docker run --rm -v $(pwd)/migrations:/migrations \
  --network mulo-go-api_mulo-api-dev \
  migrate/migrate -path=/migrations \
  -database "postgres://tungtungsahur:tralalelotralalala@mulo-db-dev:5432/mulo_bombardino?sslmode=disable" up
```

---

### Important Notes

- Make sure your environment files (e.g. `.env.development` or `.env.production`) are correctly configured.

- Place your migration SQL files inside the `migrations/` folder using the format:

  ```
  001_create_users_table.up.sql
  001_create_users_table.down.sql
  ```
