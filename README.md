# Interview Preparation Coach

A personal-first interview recovery app for experienced engineers who are strong in real work but want to rebuild coding confidence, practice consistently, and perform better in interviews.

This is not a job portal. The first goal is personal improvement — not monetization.

## Core product loop

```
profile → diagnosis → plan → practice → evaluate → progress → adjust
```

1. **Profile** — capture experience, target role, weak areas, and timeline
2. **Diagnosis** — baseline assessment across coding, concepts, and behavioral answers
3. **Plan** — turn diagnosis into a structured recovery plan with daily tasks
4. **Practice** — work through programming and DSA questions
5. **Evaluate** — AI scores attempts and gives direct feedback
6. **Progress** — track improvement, weak patterns, and readiness over time
7. **Adjust** — update the plan based on what is working

## Stack

| Layer | Technology |
|-------|------------|
| Frontend | Next.js, TypeScript, Tailwind CSS *(not started)* |
| Backend | Go, Chi router |
| Database | PostgreSQL |
| DB access | pgx, sqlc |
| Migrations | golang-migrate |
| AI | OpenAI Responses API via Go backend only |
| Deployment | Docker + EC2 *(planned)* |

## What we are building in v1

- Preparation profile and baseline diagnosis
- AI-assisted recovery plan
- Daily practice with attempt logging
- AI evaluation and feedback
- Progress dashboard and weak-area tracking

## What is out of scope for v1

- Job portal or vacancy search
- CV/JD analysis
- Public community features
- Voice interviews
- Billing or multi-user auth

See `docs/product/` for the full PRD, MVP scope, and technical design.

## Repository layout

```
interview-prep/
├── backend/          # Go API server
│   ├── cmd/api/      # Application entrypoint
│   └── internal/     # HTTP handlers, services, repositories, AI, generated db code
├── db/               # Database source of truth (migrations, sqlc config, SQL queries)
├── docs/             # Product, strategy, and implementation docs
│   └── progress/     # Build status and next-step prompts
└── frontend/         # Next.js app (planned)
```

**Why `db/` is at the root:** migrations and SQL queries are shared infrastructure, not Go application code. `sqlc` reads from `db/` and generates typed Go into `backend/internal/db/`. This is a common pattern for Go + sqlc + golang-migrate projects.

## Current status

**Backend complete through Phase 7.** All API modules, OpenAI integration, and unit/handler tests are in place. Frontend is not started.

| Area | Status |
|------|--------|
| Product docs & schema | Done |
| Go backend foundation | Done |
| Profile, diagnosis, plan, questions, attempts | Done |
| Evaluations & progress snapshots | Done |
| AI integration (`backend/internal/ai/`) | Done |
| Backend unit + handler tests | Done |
| Frontend | Not started |
| Deployment | Not started |

For the full checklist and next implementation step, see [`docs/progress/STATUS.md`](docs/progress/STATUS.md) and [`docs/progress/NEXT.md`](docs/progress/NEXT.md).

### Live API endpoints

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/health` | Health check |
| `POST` | `/api/v1/profile` | Create or update profile |
| `GET` | `/api/v1/profile` | Get profile |
| `POST` | `/api/v1/diagnosis` | Submit baseline diagnosis (AI) |
| `GET` | `/api/v1/diagnosis` | Get latest diagnosis |
| `POST` | `/api/v1/plan` | Generate or replace active recovery plan (AI) |
| `GET` | `/api/v1/plan` | Get active plan with tasks |
| `GET` | `/api/v1/plan/today` | Get today's recommended task |
| `PATCH` | `/api/v1/plan/tasks/{id}/complete` | Mark a plan task complete |
| `GET` | `/api/v1/questions` | List practice questions (optional filters) |
| `GET` | `/api/v1/questions/{id}` | Get a single question |
| `POST` | `/api/v1/attempts` | Submit a practice attempt |
| `GET` | `/api/v1/attempts` | List attempts (optional `question_id` filter) |
| `GET` | `/api/v1/attempts/{id}` | Get a single attempt |
| `POST` | `/api/v1/attempts/{id}/evaluate` | Evaluate an attempt (AI) |
| `GET` | `/api/v1/attempts/{id}/evaluation` | Get evaluation for an attempt |
| `POST` | `/api/v1/progress/snapshot` | Create daily progress snapshot (AI summary) |
| `GET` | `/api/v1/progress/latest` | Get latest progress snapshot |
| `GET` | `/api/v1/progress/snapshots` | List snapshots (optional `limit`, `start_date`, `end_date`) |

## Local development

### Quick start (backend + frontend)

From the repo root:

```bash
./runserver.sh
```

This starts:
- **Frontend** at http://localhost:3000 (`next dev` — hot reload built in)
- **Backend** at http://localhost:8080 (`air` if installed, otherwise `go run` with a file watcher)

Install [air](https://github.com/air-verse/air) for faster Go reloads:

```bash
go install github.com/air-verse/air@latest
```

Press **Ctrl+C** to stop both servers.

### Prerequisites

- Go 1.23+
- PostgreSQL
- [golang-migrate](https://github.com/golang-migrate/migrate) CLI
- [sqlc](https://sqlc.dev/) CLI

### 1. Database setup

Create a local database:

```bash
createdb interview_prep
```

Run migrations (optional — the API also runs pending migrations on startup):

```bash
migrate -path db/migrations -database "postgres://postgres:postgres@localhost:5432/interview_prep?sslmode=disable" up
```

Or apply a single file with `psql` if you do not install the `migrate` CLI.

### 2. Backend configuration

```bash
cp backend/.env.example backend/.env
```

Edit `backend/.env` with your `DATABASE_URL` and `OPENAI_API_KEY` (required for diagnosis, plan, evaluate, and progress snapshot endpoints).

### 3. Generate sqlc code (after changing SQL)

```bash
cd db && sqlc generate
```

### 4. Start the API server

From the repo root:

```bash
cd backend && go run ./cmd/api
```

On startup the API applies any pending migrations from `db/migrations` before accepting requests.

The server listens on port `8080` by default (override with `PORT` in `.env`).

Verify it is running:

```bash
curl http://localhost:8080/health
```

Expected response when the database is connected:

```json
{"status":"ok","database":"connected"}
```

### 5. Stop the API server

If the server is running in the foreground (e.g. in a terminal tab), press **`Ctrl+C`**. The server performs a graceful shutdown (drains in-flight requests, closes the DB pool).

If it is running in the background, find and kill the process:

```bash
# Find the process on port 8080 (macOS)
lsof -ti :8080 | xargs kill
```

### 6. Run tests

Tests use mocked services — no PostgreSQL or OpenAI key required.

```bash
cd backend && go test ./...
```

Run with coverage:

```bash
cd backend && go test ./... -cover
```

Run a single package:

```bash
cd backend && go test ./internal/service/... -v
```

Tests cover all 19 API routes (happy path + key error statuses), service validation/business logic, and AI schema parsing.

## Build order

Implementation follows the phased plan in `docs/prompts/interview-preparation-app-master-implementation-prompts.md`:

1. Freeze MVP scope
2. Design architecture
3. Create folder structure
4. Design database
5. Build Go backend foundation
6. Build backend modules (profile → diagnosis → plan → practice → evaluate → progress)
7. Add AI prompt layer
8. Build Next.js frontend
9. Connect end-to-end workflow
10. Test and fix
11. Deploy to EC2

## Documentation

| Doc | Purpose |
|-----|---------|
| [`docs/product/interview-preparation-app-personal-first-mvp.md`](docs/product/interview-preparation-app-personal-first-mvp.md) | Frozen MVP scope |
| [`docs/product/interview-preparation-app-technical-design.md`](docs/product/interview-preparation-app-technical-design.md) | Architecture and modules |
| [`docs/product/interview-preparation-app-master-build-plan.md`](docs/product/interview-preparation-app-master-build-plan.md) | Step-by-step build plan |
| [`docs/prompts/interview-preparation-app-master-implementation-prompts.md`](docs/prompts/interview-preparation-app-master-implementation-prompts.md) | Copy-paste prompts per phase |
| [`docs/progress/STATUS.md`](docs/progress/STATUS.md) | Current build status |
