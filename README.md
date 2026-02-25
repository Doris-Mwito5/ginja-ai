# Ginja AI — Health Claims Intelligence Service

A backend service for real-time health claims validation. Built with Go, Gin, PostgreSQL, and JWT authentication.

---

## Architecture Decisions

### Layered (Hexagonal) Architecture

The codebase is organized into three distinct layers that separate concerns and make the system easy to test and maintain:

**Domain Layer** (`internal/domain`)
Responsible for all database interactions. Each entity (claims, members, providers, procedures, users) has its own domain file containing raw SQL queries and row scanning logic. This layer knows nothing about business rules — it only reads and writes data.

**Service Layer** (`internal/services`)
Encapsulates all business logic and validations. The claims service runs a four-stage validation pipeline on every submission: member eligibility → procedure verification → fraud detection → benefit limit check. Services depend on the domain layer via the `domain.Store`, which is a single struct that holds all domain interfaces — making it easy to pass around and mock in tests.

**API Layer** (`web/handlers`)
Handles HTTP concerns only — binding request bodies, calling the appropriate service, and returning responses. Handlers are kept thin; no business logic lives here.

### Claims Validation Pipeline

Every claim submission goes through these stages in order:

```
1. Member Eligibility   → REJECTED if member not found or inactive
2. Procedure Check      → REJECTED if procedure code not in system
3. Fraud Detection      → fraud_flag = true if amount > 2× average procedure cost
4. Benefit Limit Check  → PARTIAL if amount exceeds remaining benefit (benefit_limit - used_amount)
                       → APPROVED if within limit and no fraud
                       → PARTIAL if within limit but fraud flagged (pending manual review)
```

All DB writes in a single claim submission (inserting the claim + updating the member's `used_amount`) are wrapped in a single database transaction — if either fails, both are rolled back.

### Authentication

JWT Bearer tokens using `golang-jwt`. Tokens are issued on register and login, expire after 24 hours, and are validated on every protected route via a Gin middleware. Passwords are hashed with `bcrypt` before storage.

### Database

PostgreSQL with raw SQL queries (no ORM). This keeps queries explicit, predictable, and easy to optimize. The `procedures` table drives the fraud detection threshold via `average_cost`, meaning fraud rules can be updated with a data change rather than a code deployment.

---

## How to Run Locally

### Prerequisites
- Docker and Docker Compose
- Go 1.24+ (only needed if running without Docker)
- `goose` for migrations: `go install github.com/pressly/goose/v3/cmd/goose@latest`

### With Docker 

```bash
# 1. Clone the repository
git clone https://github.com/Doris-Mwito5/ginja-ai.git
cd ginja-ai

# 2. Build the dockerfile
docker-compose up --build -d

# 3. Run migrations
make migrate

# 4. Start the API
go run cmd/api/main.go
```

The API will be available at `http://localhost:8080`

### Without Docker

```bash
# 1. Start a local PostgreSQL instance and create the database
createdb ginja_claims

# 2. Copy and configure environment variables
cp .env.example .env
# Edit .env with your database credentials

# 3. Run migrations
make migrate

# 4. Start the server
go run cmd/api/*.go
```

---

## API Endpoints

### Auth (public)
```
POST /v1/register    — create account
POST /v1/login       — get JWT token
```

### Claims (requires Bearer token)
```
POST /v1/claims                      — submit a claim
GET  /v1/claims/:id                  — get claim by ID
GET  /v1/claims/member/:memberID     — list claims for a member
```

### Admin (requires Bearer token)
```
POST /v1/members     — create member
POST /v1/providers   — create provider
POST /v1/procedures  — create procedure
```

### Sample Requests

**Register**
```json
POST /v1/register
{
  "username": "alice",
  "email": "alice@example.com",
  "password": "secret1234"
}
```

**Submit Claim**
```json
POST /v1/claims
Authorization: Bearer <token>

{
  "member_id": 1,
  "provider_id": 1,
  "procedure_code": "P001",
  "diagnosis_code": "D001",
  "requested_amount": 30000
}
```

**Sample Response**
```json
{
  "claim_id": 1,
  "status": "APPROVED",
  "fraud_flag": false,
  "approved_amount": 30000
}
```

---

## What I Would Improve for Production

**Security**
- Store JWT secrets in a secrets manager (AWS Secrets Manager or HashiCorp Vault) and rotate them on a schedule
- Add rate limiting per IP and per user to prevent brute force and submission floods
- Add role-based access control — separate admin roles (who can create members/providers) from insurer roles (who can submit claims)
- Enforce HTTPS and mutual TLS for hospital and insurer integrations

**Reliability**
- Add Redis caching for member eligibility lookups — member status rarely changes and this is the hottest query path
- Move claim processing to an async queue — accept the claim synchronously, process validation asynchronously, return result via webhook or polling. This decouples submission volume from processing throughput
- Add database connection pooling via PgBouncer for high-concurrency scenarios

**Observability**
- Add Prometheus metrics endpoint — track `claims_submitted_total`, `claims_fraud_flagged_total`, `claims_processing_duration`
- Add distributed tracing with OpenTelemetry to trace a claim across services
- Set up alerts on fraud flag rate spikes and DB connection pool exhaustion

**Testing**
- Add unit tests for the full claims validation pipeline using `sqlmock`
- Add integration tests against a real test database
- Add load tests to validate throughput under peak submission periods (end of month)

**Operations**
- Add Kubernetes manifests with `HorizontalPodAutoscaler` for auto-scaling
- The graceful shutdown is already implemented — extend it with a readiness probe so Kubernetes knows when the pod is ready to serve traffic

