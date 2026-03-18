# VaultPad Server

[Русская версия (README.ru.md)](README.ru.md) | [Client repository](https://github.com/Terciuss/VaultPad)

Backend server for [VaultPad](https://github.com/Terciuss/VaultPad) -- encrypted data vault with zero-knowledge architecture. The server stores and synchronizes encrypted blobs between clients. It never has access to plaintext data; all encryption and decryption happens on the client side.

## Features

- **Zero-knowledge storage** -- server only stores encrypted names and content, no plaintext ever touches the backend
- **JWT authentication** -- HS256 tokens with 7-day expiration
- **User management** -- admin panel for creating, editing, and deleting users
- **File access sharing** -- per-file access grants for individual users
- **Auto-seeded admin** -- first administrator is created automatically on first launch (no registration endpoint)
- **Rate limiting** -- 100 requests per IP per minute
- **CORS support** -- configurable cross-origin requests
- **Docker ready** -- multi-stage Dockerfile and docker-compose with MySQL

## Tech Stack

| Component | Technology |
|---|---|
| Language | Go 1.25 |
| Database | MySQL 8 |
| Authentication | JWT (`golang-jwt/jwt/v5`) |
| Password hashing | bcrypt (`golang.org/x/crypto`) |
| Configuration | Environment variables via `.env` (`godotenv`) |
| Containerization | Docker, Docker Compose |

## Quick Start

### Docker Compose (recommended)

```bash
cp .env.example .env
# Edit .env: set JWT_SECRET and MYSQL_ROOT_PASSWORD
docker compose up
```

The server starts on port `8080` with MySQL. First admin is auto-seeded with email from `SEED_ADMIN_EMAIL`.

### Local Development

```bash
cp .env.example .env
# Edit .env: replace "db" with "127.0.0.1" in DATABASE_DSN
# Ensure MySQL is running locally

go run ./cmd/server
```

Requires: Go 1.25+, MySQL 8+

## Environment Variables

| Variable | Default | Description |
|---|---|---|
| `DATABASE_DSN` | — | MySQL connection string (e.g. `root:password@tcp(127.0.0.1:3306)/vault_pad?parseTime=true`) |
| `JWT_SECRET` | — | Secret key for JWT signing (min 32 characters recommended) |
| `LISTEN_ADDR` | `:8080` | Server listen address and port |
| `SEED_ADMIN_EMAIL` | `admin@local.local` | Email for the auto-seeded admin user |
| `MYSQL_ROOT_PASSWORD` | `root` | MySQL root password (Docker only) |
| `MYSQL_DATABASE` | `vault_pad` | MySQL database name (Docker only) |

## API Endpoints

### Public

| Method | Path | Description |
|---|---|---|
| `POST` | `/api/auth/login` | Authenticate and receive JWT token |
| `GET` | `/api/health` | Health check |

### Authenticated (JWT required)

| Method | Path | Description |
|---|---|---|
| `GET` | `/api/auth/me` | Current user info |
| `GET` | `/api/projects/meta` | Project metadata (IDs, timestamps) for sync |
| `GET` | `/api/projects` | List all accessible projects |
| `GET` | `/api/projects/{id}` | Get project by ID |
| `POST` | `/api/projects` | Create project |
| `PUT` | `/api/projects/{id}` | Update project |
| `DELETE` | `/api/projects/{id}` | Delete project (owner/admin only) |

### Admin (JWT + admin role required)

| Method | Path | Description |
|---|---|---|
| `GET` | `/api/admin/users` | List all users |
| `POST` | `/api/admin/users` | Create user |
| `PUT` | `/api/admin/users/{id}` | Update user |
| `DELETE` | `/api/admin/users/{id}` | Delete user (cannot delete last admin) |
| `GET` | `/api/admin/users/{id}/shares` | List project shares for user |
| `POST` | `/api/projects/{id}/share` | Share project with user |
| `DELETE` | `/api/projects/{id}/share/{userId}` | Revoke project access |

## Project Structure

```
server/
├── cmd/
│   └── server/
│       └── main.go               # Entry point, routing, auto-seeder
├── internal/
│   ├── handler/
│   │   ├── auth.go               # Login, current user
│   │   ├── project.go            # Project CRUD
│   │   ├── admin.go              # User management, sharing
│   │   └── helpers.go            # JSON response helpers
│   ├── middleware/
│   │   ├── auth.go               # JWT verification
│   │   └── ratelimit.go          # Per-IP rate limiting
│   ├── model/
│   │   ├── user.go               # User model
│   │   └── project.go            # Project, ProjectShare models
│   ├── repository/
│   │   ├── db.go                 # Database connection and migrations
│   │   ├── user.go               # User queries
│   │   ├── project.go            # Project queries
│   │   └── share.go              # Share queries
│   └── service/
│       ├── auth.go               # Authentication logic
│       ├── project.go            # Project business logic
│       └── admin.go              # Admin business logic
├── .env.example
├── docker-compose.yml
├── Dockerfile
├── go.mod
└── go.sum
```

## Database Schema

### users

| Column | Type | Description |
|---|---|---|
| `id` | BIGINT PK | Auto-increment ID |
| `email` | VARCHAR(255) UNIQUE | User email |
| `password_hash` | VARCHAR(255) | bcrypt hash |
| `is_admin` | BOOLEAN | Admin flag |
| `created_at` | TIMESTAMP | Creation date |

### projects

| Column | Type | Description |
|---|---|---|
| `id` | BIGINT PK | Auto-increment ID |
| `user_id` | BIGINT FK | Owner |
| `encrypted_name` | BLOB | Encrypted file name |
| `encrypted_content` | MEDIUMBLOB | Encrypted file content |
| `sort_order` | INT | Display order |
| `created_at` | TIMESTAMP | Creation date |
| `updated_at` | TIMESTAMP | Last update |

### project_shares

| Column | Type | Description |
|---|---|---|
| `id` | BIGINT PK | Auto-increment ID |
| `project_id` | BIGINT FK | Shared project |
| `user_id` | BIGINT FK | Recipient user |
| `shared_by` | BIGINT | User who shared |
| `created_at` | TIMESTAMP | Share date |

## Access Control

- **Owner** -- full access to own projects
- **Admin** -- full access to all projects
- **Shared user** -- read/update access to shared projects (no delete)
- Last admin cannot be deleted or demoted

## Security

- All data stored on server is encrypted client-side (AES-256-GCM)
- Passwords hashed with bcrypt
- JWT tokens expire after 7 days
- Rate limiting: 100 req/min per IP
- No registration endpoint -- users are created by admins or via auto-seeder

## License

PolyForm Noncommercial License 1.0.0. See [LICENSE](LICENSE) for details.

Copyright (c) 2026 Pavel <mr.terks@yandex.ru>
