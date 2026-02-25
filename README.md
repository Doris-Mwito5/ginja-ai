   

- **`cmd/api`** — Single entry point: loads config, initializes DB and domain store, builds the router, and runs the HTTP server with graceful shutdown.
- **`web/`** — HTTP surface: `routes` wires the app and mounts handlers; `handlers/<resource>` define URL groups and view functions. Handlers only parse input, call services, and write responses.
- **`internal/services`** — Application/use-case layer. Services depend on the domain store and `db.DB`; they orchestrate workflows and translate between DTOs and domain models.
- **`internal/domain`** — Domain logic and persistence. A single **`domain.Store`** holds all domain interfaces (User, Member, Provider, Procedure, Claim). Domain types contain SQL and speak in terms of `db.SQLOperations`, so the same code works with a transaction or the global DB.
- **`internal/db`** — Database abstraction: `DB` interface (with `InTransaction`), connection pooling (pgx), and goose-based migrations in `internal/db/migrations`.

This gives you:

- Testability: services and domains can be tested with a fake or in-memory DB by implementing `db.SQLOperations` / `db.DB`.
- Consistent transactions: `dB.InTransaction(ctx, fn)` is used wherever a unit of work spans multiple domain calls.
- A single place to wire dependencies: the router builds one store and passes it into all services.

### API and auth

- **Versioned API**: All routes live under `/v1`. Public and protected groups are split; protected routes use `AuthMiddleware(jwtMaker)` so only valid JWT bearers can access them.
- **JWT**: Tokens are created at login (e.g. 24h expiry) and validated by the middleware; the maker is built from `JWT_SECRET` in config.
- **CORS**: Handled in middleware so the same app can be used by a separate frontend.

### Configuration and logging

- **Config**: Viper loads a `.env` file (or falls back to environment variables). Required values: `DATABASE_URL`, `JWT_SECRET`, `HTTP_PORT`, `ENVIRONMENT`. This keeps environment-specific and secret values out of the binary.
- **Logging**: Centralized logger (e.g. zerolog) is initialized at startup and used across the app for consistent, structured logs.

### Database

- **PostgreSQL** with **goose** for migrations (up/down in `internal/db/migrations`), so schema changes are versioned and repeatable.
- **Connection pooling**: pgx pool plus `database/sql` with bounded connections and idle timeouts to avoid exhausting the DB.

---

## How to run the application locally

### Prerequisites

- Go 1.24+
- PostgreSQL (or use Docker for Postgres only)
- [goose](https://github.com/pressly/goose) for migrations (optional if using Docker Compose, which runs migrations via init scripts)

### 1. Clone and install dependencies

```bash
cd /path/to/ginja-ai
go mod download
```

### 2. Environment variables

Create a `.env` file in the project root (or export the variables):

```env
DATABASE_URL=postgres://ginja:ginja_secret@localhost:5433/ginja_claims?sslmode=disable
JWT_SECRET=your-secret-key-at-least-32-chars
HTTP_PORT=8080
ENVIRONMENT=development
DEBUG=false
```

For a **local Postgres** instance, adjust `DATABASE_URL` (host, port, user, password, dbname) as needed. The Docker Compose setup uses port `5433` for Postgres to avoid clashing with a local PostgreSQL on 5432.

### 3. Database and migrations

**Option A — Docker Compose (Postgres + optional API)**

Start Postgres (and optionally the API) and run migrations via the init script:

```bash
docker-compose up -d postgres
# Wait for Postgres to be healthy, then run migrations if not using init script:
make migrate
```

Then run the API on your host:

```bash
make api
# or: go run cmd/api/*.go
```

**Option B — Local Postgres only**

Ensure PostgreSQL is running and the database exists, then:

```bash
make migrate
make api
```

### 4. Run the API

From the project root:

```bash
make api
```

The server listens on the port set in `HTTP_PORT` (default `8080`). Health/readiness are implied by the server responding; you can hit a known route (e.g. `POST /v1/users/login` or a protected endpoint with a valid token) to confirm.

### 5. Run everything with Docker

To build and run both the API and Postgres in containers:

```bash
make build
# or: docker-compose up --build -d
```

The API is exposed on port `8080`; Postgres is on `5433` on the host (mapped from 5432 in the container). Migrations in `internal/db/migrations` are applied by the Postgres container’s `docker-entrypoint-initdb.d` (only on first run).

### Make targets

| Target           | Description                          |
|------------------|--------------------------------------|
| `make api`       | Run the API locally (`go run`)       |
| `make migrate`   | Run goose migrations up              |
| `make rollback`  | Run goose migrations down            |
| `make status`    | Show migration status                |
| `make build`     | Build and start services with Docker |
| `make migration name=my_migration` | Create a new migration file   |

---

## What we would improve for production

- **Secrets and config**
  - Avoid committing `.env` or any real secrets. Use a secret manager (e.g. AWS Secrets Manager, HashiCorp Vault) or platform env/injectors and rotate `JWT_SECRET` and DB credentials.
  - Prefer explicit env vars (e.g. `POSTGRES_HOST`, `POSTGRES_PASSWORD`) or a single injected `DATABASE_URL` from the platform.

- **Observability**
  - Add **metrics** (e.g. Prometheus) for request counts, latency, and error rates per route and status.
  - Add **tracing** (e.g. OpenTelemetry) with a single trace ID per request and propagate it in logs and to downstream calls.
  - Ensure **structured logging** (e.g. zerolog JSON) with request ID, user ID where relevant, and error codes so log aggregation (e.g. ELK, Datadog) is effective.

- **Resilience and performance**
  - **Rate limiting** (per IP and/or per user) on login and expensive endpoints to reduce abuse and load.
  - **Timeouts**: enforce timeouts on all HTTP and DB calls (context with deadline) and size limits on request bodies.
  - **Health endpoints**: explicit `/health` and `/ready` (e.g. DB ping, dependency checks) for load balancers and orchestrators.

- **Security**
  - **HTTPS only** in production; use TLS termination at the load balancer or in the app.
  - **CORS**: restrict origins and methods in production instead of allowing broad CORS.
  - **Auth**: consider short-lived access tokens plus refresh tokens and secure, HTTP-only cookie options for web clients.
  - **Input validation**: validate and sanitize all inputs (e.g. stricter validation on DTOs, limit string lengths and numeric ranges).

- **Testing and quality**
  - **Unit tests** for domain and services with a fake `db.SQLOperations` or test DB.
  - **Integration tests** against a real Postgres (e.g. testcontainers or CI Postgres) to validate migrations and critical flows.
  - **CI**: run tests, linters, and security checks on every PR.

- **Database and migrations**
  - Run migrations as a separate step in deployment (e.g. init container or release phase), not only via container init, so you can control order and rollback.
  - Consider connection pool sizing (e.g. from env) and read replicas for read-heavy endpoints.

- **Deployment**
  - Run the app as a non-root user in the container; the Dockerfile can be extended to add a user and switch to it.
  - Use the multi-stage Docker build (as in the existing Dockerfile) to keep the image small and avoid shipping source code.
  - In Kubernetes or similar, set resource limits/requests and use the health/ready endpoints for probes.

These improvements would make the service more secure, observable, and easier to operate at scale in production.
# ginja-ai
