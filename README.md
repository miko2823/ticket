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
