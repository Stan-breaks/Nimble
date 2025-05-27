# 🚀 NimbleStack

_A Modern Go + Templ + Tailwind CSS Starter Template with HTMX, Alpine.js & SQLC_

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)  
![Templ](https://img.shields.io/badge/Templ-0.2+-blue)  
![Tailwind CSS](https://img.shields.io/badge/Tailwind_CSS-3.3+-06B6D4?logo=tailwind-css)  
![HTMX](https://img.shields.io/badge/HTMX-1.9+-red)  
![Alpine.js](https://img.shields.io/badge/Alpine.js-3.13+-8BC0D0)  
![sqlc](https://img.shields.io/badge/sqlc-1.25+-brightgreen)  
![SQLite](https://img.shields.io/badge/SQLite-3+-003B57?logo=sqlite)

**NimbleStack** is a lightning-fast, full-stack starter template designed for developers who want to build modern web apps with minimal boilerplate. It features **SQLite + SQLC** for embedded database magic! ✨

---

## 🌟 Features

- **Go Backend**: Blazing-fast API and server logic with Go.
- **SQLite + SQLC**: Type-safe database access with a single-file embedded database.
- **Templ Templates**: Clean, type-safe HTML templating.
- **Tailwind CSS**: JIT-compiled, utility-first CSS.
- **HTMX + Alpine.js**: Dynamic UI without JavaScript fatigue.
- **Podman Containerization**: Run NimbleStack anywhere with a single binary, thanks to our multi-stage Containerfile.

---

## 🛠️ Why NimbleStack?

- **Zero Deployment Hassle**: Package your app as a single binary with an embedded SQLite database.
- **Full-Stack Type Safety**: Enjoy a seamless SQLC → Go → Templ workflow.
- **Local Development Bliss**: No need to install or configure separate database servers.
- **Portability with Podman or Docker**: Our provided Containerfile lets you build a container that runs consistently on any platform—whether on your local machine, in the cloud, or in CI/CD pipelines.
- **Modern UI/UX**: Use HTMX and Alpine.js to create responsive, reactive interfaces without heavy frameworks.

---

## 🚀 Getting Started

### Prerequisites

- **Go**: 1.21+
- **Node.js**: 18+ & pnpm
- **Tailwind CSS**: Although this template uses Tailwind for styling, please note that the Tailwind CLI is installed via the AUR. Users on other platforms will need to set up their own method for building the CSS.
- **SQLC**

### Installation

1. **Clone the repository:**

   ```bash
   git clone https://github.com/Stan-breaks/nimblestack.git
   cd nimblestack
   ```

2. **Install dependencies:**

   ```bash
   pnpm install
   go mod tidy
   ```

3. **Generate code:**

   ```bash
   templ generate ./views/
   sqlc generate
   ```

4. **Start the server:**

   ```bash
   just watch
   ```

5. **(Optional) Build and run with Docker/Podman:**

   The included Containerfile lets you containerize NimbleStack. For example, to build and run using Podman:

   ```bash
   podman build -t nimblestack .
   podman run -p 8080:8080 nimblestack
   ```

---

## 📂 Project Structure

```
nimblestack/
├── database/         # Generated Go models
├── sqlc/             # SQLC schema and queries
├── public/           # Static assets (CSS, images, etc.)
├── views/            # Templ components
├── router/           # HTTP handlers
├── src/              # Built tailwind css
├── Dockerfile        # Multi-stage Dockerfile for containerization
├── sqlc.yaml         # SQLC configuration
└── main.go           # Server entry point
```

---

## 🔌 Database Workflow (SQLite + SQLC)

### 1. Create Tables

`sqlc/schema.sql`:

```sql
CREATE TABLE users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  email TEXT NOT NULL UNIQUE,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### 2. Write Queries

`sqlc/queries.sql`:

```sql
-- name: CreateUser :one
INSERT INTO users (name, email)
VALUES (?, ?)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = ?;
```

### 3. Generate Code

```bash
sqlc generate
```

### 4. Use in apis

`apis/users.go`:

```go
func func (h *AuthApi) CreateUserHandler(w http.ResponseWriter, r *http.Request) {

    // Type-safe database operation
    user, err := h.queries.CreateUser(r.Context(), db.CreateUserParams{
        Name:  r.FormValue("name"),
        Email: r.FormValue("email"),
    })
    if err != nil {
        http.Error(w, "Database error", 500)
        return
    }
}
```

---

## 📈 Roadmap

- [ ] Add SQLite migration tool.
- [ ] HTMX CRUD example with optimistic UI.
- [ ] SQLite connection pool benchmarks.
- [ ] ARM64 build support.

---

## 📚 Learning Resources

- [SQLC SQLite Guide](https://docs.sqlc.dev/en/latest/howto/sqlite.html)
- [Modern SQLite Driver Docs](https://pkg.go.dev/modernc.org/sqlite)
- [HTMX Patterns](https://htmx.org/examples/)

---

## Docker & Portability

The provided **Containerfile** enables you to package NimbleStack into a container that runs anywhere—whether on local development machines, cloud servers, or within CI/CD pipelines. This offers several advantages:

- **Consistency**: The container ensures the environment (OS, dependencies, configuration) remains the same across different deployments.
- **Portability**: You can run your containerized app on any platform that supports Docker or Podman.
- **Ease of Deployment**: Single binary + container means minimal configuration and fewer moving parts.

---

## License

MIT © [Stan-breaks] | Made with ❤️ for fast web apps
