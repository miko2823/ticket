# SturdyTicket

A demo concert/event ticketing web app with production-grade reliability patterns.

## Prerequisites

- **Go** 1.25+ (`go version`)
- **Node.js** 20+ (`node --version`) — use `nvm use 24` if needed
- **PostgreSQL** 15+ — local instance or Docker
- **Docker** (optional, for running PostgreSQL locally)

## Project Structure

```
backend/     Go API server (chi + pgx, DDD architecture)
frontend/    React + Vite + TypeScript SPA
infra/       Terraform (GCP)
docs/        Architecture Decision Records
```

## Local Setup

### 1. PostgreSQL

Start a local PostgreSQL instance. The easiest way is Docker:

```bash
docker run --name sturdyticket-db \
  -e POSTGRES_USER=user \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=sturdyticket \
  -p 5432:5432 \
  -d postgres:15
```

### 2. Run Migrations

Install [golang-migrate](https://github.com/golang-migrate/migrate) if you don't have it:

```bash
brew install golang-migrate
```

Apply migrations:

```bash
migrate -path backend/migrations \
  -database "postgres://user:password@localhost:5432/sturdyticket?sslmode=disable" \
  up
```

### 3. Backend

```bash
cd backend
cp .env.example .env   # edit values as needed
```

Set the required environment variables (or edit `.env`):

| Variable         | Description              | Default           |
|------------------|--------------------------|--------------------|
| `PORT`           | Server listen port       | `8080`            |
| `DATABASE_URL`   | PostgreSQL connection URL| *(required)*       |
| `GCP_PROJECT_ID` | GCP project ID           | *(required)*       |
| `GCP_REGION`     | GCP region               | `asia-northeast1` |

Run the server:

```bash
# Load env vars (if using .env file)
export $(cat .env | xargs)

go run ./cmd/server/
```

For live reload during development, install [air](https://github.com/air-verse/air):

```bash
go install github.com/air-verse/air@latest
```

Then run `air` instead of `go run`:

```bash
export $(cat .env | xargs)
air
```

Verify it's running:

```bash
curl http://localhost:8080/health
# ok
```

### 4. Frontend

```bash
cd frontend
npm install
npm run dev
```

Opens at http://localhost:5173. API requests to `/api` are proxied to the backend at `localhost:8080`.

## Waiting Room Queue Architecture

When many users try to access the seat map simultaneously, a FIFO waiting room queue controls concurrency.

```
┌─────────────────────────────────────────────────────────────────────┐
│  User A (in seat map)                                               │
│  ┌──────────┐    heartbeat     ┌──────────┐     session:e1:abc      │
│  │ Browser  │ ──── PUT ──────▶ │ Go API   │ ──── EXPIRE ────▶ Redis│
│  └──────────┘    every 10s     └──────────┘     (TTL 30s)          │
└─────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────┐
│  User B (arrives while A is active)                                 │
│                                                                     │
│  ┌──────────┐  POST /session   ┌──────────┐  active >= max?        │
│  │ Browser  │ ───────────────▶ │ Go API   │ ──── GET ────────▶ Redis│
│  └──────────┘                  └──────────┘                        │
│       │                             │                               │
│       │  ◀── 202 Accepted ──────────┘  yes → ZADD queue:e1         │
│       │      { position: 1,                  (sorted set,          │
│       │        estimated_wait: 60s }          score=timestamp)      │
│       │                                                             │
│       │   poll every 3s                                             │
│       │  ──── GET /queue ─────▶ ZRANK → position                   │
│       │  ◀── { status: waiting, position: 1 }                      │
└─────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────┐
│  Session Expiry → Admission Flow                                    │
│                                                                     │
│  Redis               Subscriber              Service                │
│  ┌─────┐  expired   ┌──────────┐  callback  ┌──────────┐           │
│  │ TTL │ ─────────▶ │ PSubscribe│ ────────▶ │ AdmitNext│           │
│  │ hit │  keyspace  │ listener │            │          │           │
│  └─────┘  event     └──────────┘            └────┬─────┘           │
│                                                   │                 │
│                    Lua script (atomic):            │                 │
│                    ┌──────────────────────────┐    │                 │
│                    │ if active < max:         │◀───┘                 │
│                    │   ZRANGE queue:e1 0 0    │                      │
│                    │   ZREM  queue:e1 userB   │                      │
│                    │   SET   admitted token   │                      │
│                    │        (TTL 30s)         │                      │
│                    └──────────────────────────┘                      │
└─────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────┐
│  User B (admitted)                                                  │
│                                                                     │
│  ┌──────────┐  GET /queue      ┌──────────┐                        │
│  │ Browser  │ ───────────────▶ │ Go API   │ ── EXISTS admitted ──▶ │
│  └──────────┘                  └──────────┘    key? yes!           │
│       │  ◀── { status: admitted }                                   │
│       │                                                             │
│       │  POST /session         ┌──────────┐                        │
│       │ ────────────────────▶  │ Go API   │ ── clear admission     │
│       │                        │          │ ── create session       │
│       │  ◀── 201 Created       └──────────┘ ── start heartbeat     │
│       │      { session_id }                                         │
│       │                                                             │
│       ▼                                                             │
│  Seat map loads ✓                                                   │
└─────────────────────────────────────────────────────────────────────┘
```

**Redis keys involved:**

| Key | Type | TTL | Purpose |
|-----|------|-----|---------|
| `session:{eventID}:{sessionID}` | Hash | 30s | Active session (userID, createdAt) |
| `event:{eventID}:active` | String (counter) | — | Number of users on seat map |
| `queue:{eventID}` | Sorted Set | — | FIFO queue (member=userID, score=timestamp) |
| `queue:{eventID}:{userID}:admitted` | String | 30s | Admission token — "your turn" signal |
| `user:{userID}:event:{eventID}:session` | String | 30s | Prevents duplicate sessions per user |

## Development

### Backend

```bash
cd backend
go build ./...    # compile check
go test ./...     # run unit tests
```

### Frontend

```bash
cd frontend
npm run build     # production build
npm run preview   # preview production build
```

### comment
add comment
add comment