# ⚡ Nimble

_A lightweight Go + SQLite + SQLC backend starter. Clone it, build on it._

![Go](https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go)
![sqlc](https://img.shields.io/badge/sqlc-2.0+-brightgreen)
![SQLite](https://img.shields.io/badge/SQLite-3+-003B57?logo=sqlite)

---

## Why Nimble?

- **Single binary**: Go compiles to one binary. SQLite is embedded. No external database servers.
- **Type-safe SQL**: SQLC generates Go code from your SQL — no ORMs, no magic, just SQL.
- **Concurrency-safe**: Pre-configured WAL mode + busy timeout so SQLite won't choke under load.
- **JWT auth included**: Register, login, and protected routes out of the box.
- **Container-ready**: Multi-stage Containerfile produces a minimal `scratch` image.

---

## Getting Started

### Prerequisites

- **Go** 1.24+
- **SQLC** ([install](https://docs.sqlc.dev/en/latest/overview/install.html))
- **Just** (optional, [install](https://github.com/casey/just))

### Quick Start

```bash
git clone https://github.com/Stan-breaks/nimble.git
cd nimble

# Start developing — auto-generates SQLC code and rebuilds on file changes
just watch
```

The server starts on `:8080` and auto-rebuilds when you edit `.go` or `.sql` files.

---

## Project Structure

```
nimble/
├── main.go           # Entry point, DB init, server start
├── database/         # Generated SQLC code (gitignored)
├── sqlc/
│   ├── schema.sql    # Table definitions
│   └── query.sql     # SQL queries → Go functions
├── router/
│   ├── router.go     # Route definitions
│   ├── apis/         # JSON API handlers
│   └── middleware/    # Auth middleware
├── example/          # Advanced patterns & reference code
├── sqlc.yaml         # SQLC config
├── justfile          # Task runner
└── Containerfile     # Multi-stage container build
```

---

## API Endpoints

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `POST` | `/api/register` | ✗ | Create a new user |
| `POST` | `/api/login` | ✗ | Login, get JWT cookie |
| `GET` | `/api/me` | ✓ | Get current user |

### Example

```bash
# Register
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{"email":"you@example.com", "password":"secret123"}'

# Login
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"you@example.com", "password":"secret123"}'
```

---

## Examples

The `example/` directory contains advanced patterns extracted from a real project built with Nimble:

| File | What It Shows |
|---|---|
| `schema.sql` | Role-based tables, foreign keys, status tracking with defaults |
| `query.sql` | FK assignment updates, filtered queries (`IS NULL`), bulk fetch |
| `api_patterns.go` | Role-based auth switch, safe response types, dashboard aggregation, nullable FKs, duplicate-check pattern |

Use these as reference when building your own features.

---

## Adding Your Own Features

### 1. Define tables in `sqlc/schema.sql`

```sql
CREATE TABLE IF NOT EXISTS posts (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL,
  title TEXT NOT NULL,
  body TEXT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users(id)
);
```

### 2. Write queries in `sqlc/query.sql`

```sql
-- name: CreatePost :one
INSERT INTO posts (user_id, title, body)
VALUES (?, ?, ?)
RETURNING *;

-- name: GetPostsByUser :many
SELECT * FROM posts WHERE user_id = ?;
```

### 3. Save & go

`just watch` picks up your `.sql` changes automatically and regenerates the Go code. The generated functions are ready to use in your API handlers.

---

## Container Build

```bash
# Build
podman build -t nimble .

# Run
podman run -p 8080:8080 nimble
```

Works with Docker too — just swap `podman` for `docker`.

---

## SQLite Concurrency

Nimble configures SQLite for safe concurrent access:

- **WAL mode**: Readers don't block writers (and vice versa)
- **Busy timeout**: 5s retry instead of immediate `SQLITE_BUSY` errors
- **Single connection pool**: Prevents Go's `database/sql` from opening multiple writers

This handles most web app workloads. For very high write throughput, consider PostgreSQL.

---

## License

MIT © [Stan-breaks]
