# SturdyTicket — CLAUDE.md

This file is read automatically by Claude Code at the start of every session.
Always follow the decisions and conventions here. If something is unclear, ask before implementing.

---

## Project Overview

**SturdyTicket** is a demo concert/event ticketing web app built for Sturdy.
It is intentionally simple in UX but must implement production-grade reliability patterns
to demonstrate how real ticketing systems handle high-load, concurrency, and abuse scenarios.

**This is a demo app — keep UI simple. Do not over-engineer the frontend.**
**Do not skip the reliability patterns — they are the point of this project.**

---

## Tech Stack

| Layer    | Technology                          |
|----------|-------------------------------------|
| Backend  | Go (latest stable)                  |
| Frontend | React + Vite + TypeScript           |
| Database | PostgreSQL (Cloud SQL on GCP)       |
| Auth     | Firebase Authentication             |
| Cache    | Redis (session management, waiting room)  |
| Infra    | GCP — Cloud Run, Cloud SQL              |
| IaC      | Terraform                           |
| CI/CD    | GitHub Actions                      |

---

## Monorepo Structure

```
/
├── CLAUDE.md
├── backend/
│   ├── cmd/server/         # main entrypoint
│   ├── internal/
│   │   ├── auth/               # Firebase token verification middleware
│   │   ├── event/              # Event bounded context (DDD)
│   │   │   ├── domain.go       #   Aggregate: Event, Entity: Ticket, Value Objects
│   │   │   ├── service.go      #   Domain service
│   │   │   ├── usecase.go      #   Application service / use cases
│   │   │   ├── handler.go      #   HTTP handler (delivery layer)
│   │   │   ├── repository.go   #   Repository interface (port)
│   │   │   └── postgres/
│   │   │       └── repository.go  # PostgreSQL adapter
│   │   ├── booking/            # Booking bounded context (DDD)
│   │   │   ├── domain.go       #   Aggregate: Booking, Value Objects
│   │   │   ├── usecase.go      #   Application service / use cases
│   │   │   ├── handler.go      #   HTTP handler (delivery layer)
│   │   │   ├── repository.go   #   Repository interface (port)
│   │   │   └── postgres/
│   │   │       └── repository.go  # PostgreSQL adapter
│   │   ├── session/            # Session bounded context (Redis-backed)
│   │   │   ├── domain.go       #   Session struct
│   │   │   ├── store.go        #   Store interface (port)
│   │   │   ├── service.go      #   Application service (create, heartbeat, validate)
│   │   │   ├── handler.go      #   HTTP handler
│   │   │   ├── router.go       #   Route registration
│   │   │   └── redis/
│   │   │       ├── store.go       # Redis adapter (Lua scripts for atomicity)
│   │   │       └── subscriber.go  # Keyspace notification listener (TTL cleanup)
│   │   ├── ratelimit/          # Rate limiting middleware
│   │   ├── circuitbreaker/     # Circuit breaker wrapper
│   │   └── middleware/         # Logging, recovery, CORS, bot detection
│   ├── pkg/                # shared utilities (response helpers, errors, config)
│   ├── migrations/         # SQL migration files (golang-migrate format)
│   ├── go.mod
│   └── go.sum
├── frontend/
│   ├── src/
│   │   ├── pages/          # Home, EventDetail, Checkout, MyTickets
│   │   ├── components/     # shared UI components
│   │   ├── hooks/          # custom React hooks
│   │   ├── api/            # typed API client (fetch-based)
│   │   └── auth/           # Firebase auth context and hooks
│   ├── index.html
│   ├── vite.config.ts
│   └── package.json
├── infra/
│   ├── terraform/
│   │   ├── main.tf
│   │   ├── variables.tf
│   │   ├── cloud_run.tf
│   │   ├── cloud_sql.tf
│   │   └── outputs.tf
│   └── cloudbuild.yaml
├── docs/
│   └── adr/                # Architecture Decision Records
└── .github/
    └── workflows/
        ├── backend.yml
        └── frontend.yml
```

---

## Backend Conventions (Go)

### Architecture: Domain-Driven Design (DDD)

The backend strictly follows DDD. Every feature must be modelled through this lens.

**Layering (per domain package, e.g. `internal/booking/`)**

```
Handler (Delivery layer)
  └─▶ UseCase / Application Service
        └─▶ Domain Service + Domain Model
              └─▶ Repository Interface (port)
                    └─▶ PostgreSQL Repository (adapter)
```

- **Domain Model** (`domain.go`) — pure Go structs and methods with no framework or DB dependencies. Holds business rules (e.g. `Booking.CanBeCancelled()`). No `pgx`, no `net/http` imports here.
- **Domain Service** (`service.go`) — orchestrates domain logic that spans multiple aggregates. Depends only on repository interfaces, never on concrete implementations.
- **Repository Interface** (`repository.go`) — Go interface defined in the domain package. The domain owns this interface; the DB adapter implements it.
- **PostgreSQL Adapter** (`postgres/repository.go`) — implements the repository interface. All SQL lives here.
- **Use Case / Application Service** (`usecase.go`) — coordinates domain services, repositories, and external ports (auth, session). One file per use case is acceptable for clarity.
- **Handler** (`handler.go`) — HTTP concerns only: parse request, call use case, write response. No business logic here.

**Aggregates and Bounded Contexts**

| Bounded Context | Aggregate Root | Key Entities         | Storage    |
|-----------------|---------------|----------------------|------------|
| Event           | `Event`       | `Ticket`             | PostgreSQL |
| Booking         | `Booking`     | —                    | PostgreSQL |
| Session         | `Session`     | —                    | Redis      |
| User (thin)     | Firebase UID  | no local aggregate   | Firebase   |

- Each bounded context maps to one package under `internal/`
- Aggregates are the only entry point for mutations — never modify child entities directly
- Cross-context communication goes through use cases or domain events, never by importing another domain's internals

**Value Objects**

- Use value objects for concepts like `SeatLabel`, `Price`, `BookingStatus` instead of raw primitives
- Value objects are immutable — no setters, constructed via constructor functions that validate input

**Rules**
- Business rules belong in the domain model or domain service — never in handlers or repositories
- If you find yourself writing an `if` that encodes business logic in a handler, move it to the domain
- Repository interfaces must not leak DB concepts (no `pgx.Tx`, no SQL error types) — wrap them in domain errors

### General
- Go version: use the latest stable (check `go.mod`)
- Package layout follows standard Go project layout with DDD layering above
- All business logic lives in `internal/` — nothing in `internal/` is imported by external packages
- Use `pkg/` only for truly reusable, side-effect-free utilities

### HTTP / REST
- Use `net/http` stdlib + `chi` router (github.com/go-chi/chi)
- All handlers follow this signature: `func (h *Handler) MethodResource(w http.ResponseWriter, r *http.Request)`
- JSON responses always use the shared response helper in `pkg/response`:
  ```go
  response.JSON(w, http.StatusOK, payload)
  response.Error(w, http.StatusBadRequest, "message")
  ```
- Route naming: `GET /events`, `GET /events/{id}`, `POST /bookings`, `DELETE /bookings/{id}`
- Always version the API: `/api/v1/...`

### Error Handling
- Never panic in handlers — always return errors up the stack
- Use a custom `AppError` type in `pkg/errors` that carries HTTP status + message + optional cause
- Middleware catches unhandled errors and returns a generic 500

### Database
- Use `pgx/v5` as the PostgreSQL driver (github.com/jackc/pgx)
- Use `golang-migrate` for schema migrations — migration files live in `backend/migrations/`
- Repository pattern: each domain has its own `repository.go` with an interface + postgres implementation
- Never write raw SQL in handlers or services — always go through the repository interface
- Always use parameterized queries — never string-concatenate SQL

### Configuration
- All config loaded from environment variables at startup via a `Config` struct in `pkg/config`
- No hardcoded secrets anywhere — fail fast if required env vars are missing
- Use `.env.example` to document all required variables

### Testing
- Unit test files sit next to the file they test (`foo_test.go`)
- Use `testify` for assertions (github.com/stretchr/testify)
- Integration tests use a real PostgreSQL (via Docker in CI) — tag them with `//go:build integration`

---

## Frontend Conventions (React)

- TypeScript strict mode — no `any` unless absolutely necessary with a comment explaining why
- Vite for bundling
- No CSS framework — use plain CSS modules (keep the demo UI simple)
- Firebase SDK for auth — wrap in `src/auth/AuthContext.tsx`
- All API calls go through `src/api/client.ts` — never call `fetch` directly in components
- Use React Query (`@tanstack/react-query`) for data fetching and cache management
- Keep components small — if a component exceeds ~150 lines, split it

### Pages (keep these minimal)
- `/` — event listing
- `/events/:id` — event detail + seat selector
- `/checkout` — booking confirmation flow
- `/my-tickets` — authenticated user's bookings

---

## Key Reliability Features

These are the core of this demo. Every one must be implemented properly.

### 1. Double Booking Prevention ✅
**Decision: Optimistic locking with PostgreSQL**

- `tickets` table has a `version` integer column
- On reserve: `UPDATE tickets SET status='reserved', version=version+1 WHERE id=$1 AND version=$2 AND status='available'`
- If `rows affected == 0`, another request won the race — return HTTP 409 Conflict
- Do NOT use application-level locks — the DB constraint is the source of truth
- 5-minute reservation hold with `reserved_by` (userID) and `reserved_until` (timestamp)
- Expired reservations treated as available at query time via SQL CASE

### 2. Session-Based Seat Map Access Control ✅
**Decision: Redis sessions with concurrency cap**

- When a user enters the seat map page, a session is created in Redis (30s TTL, 10s heartbeat)
- Per-event concurrency counter tracks how many users are on the seat map
- Max concurrency cap (default 100) — groundwork for waiting room queue
- Reserve endpoint requires valid session via `X-Session-ID` header
- Redis keyspace notifications + subscriber goroutine for accurate counter cleanup on TTL expiry
- Redis keys: `session:{eventID}:{sessionID}`, `event:{eventID}:active`, `user:{userID}:event:{eventID}:session`

### 3. Waiting Room Queue (TODO)
**Decision: FIFO queue when concurrency cap is hit**

- When `event:{eventID}:active` >= max concurrency, new users enter a waiting room
- Show queue position and estimated wait time
- As sessions end, admit users from the queue in FIFO order

### 4. Rate Limiting (TODO)
**Decision: Token bucket per user (Firebase UID) and per IP**

- Implement in `internal/ratelimit/` as Chi middleware
- Use an in-memory token bucket for single-instance demo; note in comments how this would be Redis in production
- Limits: 10 booking attempts per user per minute, 100 requests per IP per minute
- Return HTTP 429 with `Retry-After` header

### 5. Circuit Breaker (TODO)
**Decision: `sony/gobreaker` library**

- Wrap any external call (payment simulation) with a circuit breaker
- States: Closed → Open (after 5 failures in 10s) → Half-Open (after 30s)
- When open, return HTTP 503 with a meaningful error message — do not hang
- Log state transitions for observability

### 6. Bot Detection (TODO)
**Decision: Firebase App Check + request heuristics**

- Frontend integrates Firebase App Check with reCAPTCHA v3
- Backend middleware validates the `X-Firebase-AppCheck` token on booking endpoints
- Secondary heuristic check in middleware: flag requests with no `User-Agent`, suspiciously high rates, or sequential ticket IDs

### 7. Scaling
**Decision: Cloud Run with autoscaling**

- Backend deployed to Cloud Run — stateless, scales to zero
- Min instances: 1 (avoid cold start on demo), Max instances: 10
- All state in PostgreSQL or Redis — no in-process state that breaks horizontal scaling
- Note: in-memory rate limiter must be replaced with Redis if scaling beyond 1 instance

---

## Database Schema (key tables)

```sql
-- Events
CREATE TABLE events (
  id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name                TEXT NOT NULL,
  venue               TEXT NOT NULL,
  starts_at           TIMESTAMPTZ NOT NULL,
  ticketing_starts_at TIMESTAMPTZ NOT NULL,
  ticketing_ends_at   TIMESTAMPTZ NOT NULL,
  seat_layout         JSONB,              -- canvas-based seat map (positions, sections)
  created_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Tickets (one row per seat)
CREATE TABLE tickets (
  id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  event_id       UUID NOT NULL REFERENCES events(id),
  seat_label     TEXT NOT NULL,           -- e.g. "A-12"
  price_jpy      INTEGER NOT NULL,
  status         TEXT NOT NULL DEFAULT 'available', -- available | reserved | sold | cancelled
  version        INTEGER NOT NULL DEFAULT 0,         -- optimistic lock version
  reserved_by    TEXT,                    -- Firebase UID of user who reserved
  reserved_until TIMESTAMPTZ,            -- reservation expiry (5 min)
  created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE(event_id, seat_label)
);

-- Bookings
CREATE TABLE bookings (
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id     TEXT NOT NULL,           -- Firebase UID
  ticket_id   UUID NOT NULL REFERENCES tickets(id),
  status      TEXT NOT NULL DEFAULT 'pending', -- pending | confirmed | failed | cancelled
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Failed bookings (dead letter)
CREATE TABLE failed_bookings (
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  booking_id  UUID REFERENCES bookings(id),
  reason      TEXT,
  failed_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

---

## Auth

- Firebase Authentication handles all user identity
- Backend verifies Firebase ID tokens via the Firebase Admin SDK for Go (`firebase.google.com/go/v4`)
- Middleware in `internal/auth/middleware.go` extracts and verifies the token, injects `userID` into `context.Context`
- All booking endpoints require a valid token — return 401 if missing or invalid
- Public endpoints (event listing, event detail) do not require auth

---

## GCP / Infrastructure

- Project: set via `GCP_PROJECT_ID` env var — never hardcode
- Region: `asia-northeast1` (Tokyo)
- Cloud Run service name: `sturdyticket-backend`
- Cloud SQL instance: `sturdyticket-db` (PostgreSQL 15)
- All infra changes go through Terraform — never click in the console
- Terraform state in GCS bucket: `sturdyticket-tf-state`

---

## What NOT to do

- Do not add unnecessary abstractions — this is a demo, not a framework
- Do not use ORMs (no GORM) — use `pgx` directly with the repository pattern
- Do not store secrets in code or Terraform files — use Secret Manager
- Do not skip error handling — every error must be handled or explicitly propagated
- Do not use `context.Background()` in handlers — always propagate the request context
- Do not make the frontend pretty at the expense of the backend reliability features

---

## Session Startup Checklist

At the start of every Claude Code session:
1. Read this file fully
2. Run `git status` to understand current state
3. Ask the user what they want to work on if not specified
4. Before writing code for any reliability feature, confirm the approach matches the decisions above
