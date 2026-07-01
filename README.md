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
| Frontend | Next.js, TypeScript, Tailwind CSS |
| Backend | Go, Chi router |
| Database | PostgreSQL |
| DB access | pgx, sqlc |
| Migrations | golang-migrate |
| AI | OpenAI Responses API via Go backend only |
| Deployment | Docker + EC2 |

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
└── frontend/         # Next.js app
```

**Why `db/` is at the root:** migrations and SQL queries are shared infrastructure, not Go application code. `sqlc` reads from `db/` and generates typed Go into `backend/internal/db/`. This is a common pattern for Go + sqlc + golang-migrate projects.

Deployment artifacts: `docker-compose.yml`, `backend/Dockerfile`, `frontend/Dockerfile`, `deploy/nginx/` (optional EC2 reverse proxy).

## Current status

**Full stack complete through Phase 11.** Backend, frontend, tests, and Docker deployment artifacts are in place.

| Area | Status |
|------|--------|
| Product docs & schema | Done |
| Go backend foundation | Done |
| Profile, diagnosis, plan, questions, attempts | Done |
| Evaluations & progress snapshots | Done |
| AI integration (`backend/internal/ai/`) | Done |
| Backend unit + handler tests | Done |
| Next.js frontend | Done |
| Docker + compose + EC2 docs | Done |

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

## Docker deployment

Reproducible stack: **Postgres :5432**, **API :8080**, **frontend :3000**.

### Ports

| Service | Container port | Default host port |
|---------|----------------|-------------------|
| PostgreSQL | 5432 | 5432 |
| API | 8080 | 8080 |
| Frontend | 3000 | 3000 |

### Quick start (local Docker)

From the repo root:

```bash
cp .env.example .env
# Edit .env — set OPENAI_API_KEY (required for AI endpoints)

docker compose up --build -d
docker compose ps
```

Open http://localhost:3000. API health:

```bash
curl http://localhost:8080/health
```

Expected when Postgres is connected:

```json
{"status":"ok","database":"connected"}
```

Stop and remove containers (keeps the Postgres volume):

```bash
docker compose down
```

Remove containers and database volume:

```bash
docker compose down -v
```

### Migrations and seed data

The API runs **golang-migrate** on startup (`backend/cmd/api/main.go`) using SQL files baked into the image at `/app/db/migrations` (`MIGRATIONS_PATH`).

Migration files in `db/migrations/`:

| File | Purpose |
|------|---------|
| `000001_init_schema.up.sql` | Tables |
| `000002_seed_questions.up.sql` | 15 practice questions |
| `000003_attempt_language.up.sql` | Attempt language column |

No separate seed step is required for a fresh stack — `docker compose up` applies migrations automatically.

Manual migration (optional, e.g. against RDS from your laptop):

```bash
migrate -path db/migrations -database "$DATABASE_URL" up
```

If a migration is **dirty**, fix the schema and force the version per the error message from the API logs.

### Environment variables

**Root `.env`** (for `docker compose` — copy from `.env.example`):

| Variable | Required | Description |
|----------|----------|-------------|
| `OPENAI_API_KEY` | Yes (AI features) | OpenAI key — backend only, never in frontend image |
| `OPENAI_MODEL` | No | Default `gpt-4o-mini` |
| `NEXT_PUBLIC_API_BASE_URL` | No | Frontend → API URL; **baked at image build**. Default `http://localhost:8080/api/v1` |
| `CORS_ALLOWED_ORIGINS` | No | Comma-separated origins for API CORS. Default local Next.js URLs |
| `POSTGRES_USER` / `POSTGRES_PASSWORD` / `POSTGRES_DB` | No | Containerized Postgres credentials |

**Backend** (`backend/.env` for local `go run`, or API service env in compose):

| Variable | Required | Description |
|----------|----------|-------------|
| `DATABASE_URL` | Yes | Postgres connection string |
| `PORT` | No | Default `8080` |
| `ENV` | No | `development` or `production` |
| `MIGRATIONS_PATH` | No | Set automatically in Docker image |
| `CORS_ALLOWED_ORIGINS` | No | Comma-separated browser origins |
| `OPENAI_API_KEY` | Yes (AI) | Backend only |
| `OPENAI_MODEL` | No | Default `gpt-4o-mini` |

**Frontend** (`frontend/.env.local` for `npm run dev`, or build arg in Docker):

| Variable | Description |
|----------|-------------|
| `NEXT_PUBLIC_API_BASE_URL` | API base including `/api/v1` suffix |

For **EC2 behind Nginx** on one host, rebuild the frontend with `NEXT_PUBLIC_API_BASE_URL=/api/v1` so the browser uses same-origin requests through the proxy (see below).

### EC2 single-instance setup

Practical path for a solo developer: one **Amazon Linux 2023** or **Ubuntu** EC2 instance running Docker Compose.

#### 1. Launch EC2

- Instance type: `t3.small` or larger (AI calls are external; app itself is light)
- Storage: 20 GB+ gp3
- Key pair for SSH

#### 2. Security groups

| Direction | Port | Source | Purpose |
|-----------|------|--------|---------|
| Inbound | 22 | Your IP | SSH |
| Inbound | 80 | `0.0.0.0/0` (or your IP) | HTTP via Nginx |
| Inbound | 443 | `0.0.0.0/0` | HTTPS (after TLS) |
| Inbound | 3000, 8080 | — | **Do not expose** publicly when using Nginx |
| Inbound | 5432 | — | **Do not expose** Postgres to the internet |

Outbound: allow HTTPS (OpenAI API, package installs).

#### 3. Install Docker on the instance

```bash
# Example: Amazon Linux 2023
sudo dnf update -y
sudo dnf install -y docker git
sudo systemctl enable --now docker
sudo usermod -aG docker ec2-user
# Log out and back in so docker group applies
```

Install [Docker Compose plugin](https://docs.docker.com/compose/install/linux/) per Docker docs.

#### 4. Deploy the app

```bash
git clone <your-repo-url> interview-prep
cd interview-prep
cp .env.example .env
```

Edit `.env` on the server:

- Set `OPENAI_API_KEY`
- For Nginx same-origin routing: `NEXT_PUBLIC_API_BASE_URL=/api/v1`
- Set `CORS_ALLOWED_ORIGINS` to your public origin, e.g. `http://YOUR_EC2_DNS` or `https://your.domain`

```bash
docker compose up --build -d
docker compose logs -f api   # confirm migrations applied
curl http://127.0.0.1:8080/health
```

#### 5. Nginx reverse proxy (recommended)

Copy `deploy/nginx/interview-prep.conf` to the host:

```bash
sudo cp deploy/nginx/interview-prep.conf /etc/nginx/conf.d/interview-prep.conf
sudo nginx -t && sudo systemctl reload nginx
```

Nginx listens on **:80** and proxies:

- `/` → frontend `:3000`
- `/api/` → API `:8080`
- `/health` → API health check

Add TLS later with Certbot (`server_name` + `listen 443 ssl`).

#### 6. Restart strategy

`docker-compose.yml` sets `restart: unless-stopped` on all services. After reboot:

```bash
sudo systemctl enable docker
cd ~/interview-prep && docker compose up -d
```

To deploy a new version:

```bash
git pull
docker compose up --build -d
```

### RDS vs containerized Postgres

| Topic | Container Postgres (compose) | Amazon RDS |
|-------|------------------------------|------------|
| Service | `postgres` service in compose | Managed RDS instance |
| `DATABASE_URL` | Built from compose env (`postgres:5432` hostname) | RDS endpoint, e.g. `postgres://user:pass@mydb.xxx.rds.amazonaws.com:5432/interview_prep?sslmode=require` |
| SSL | `sslmode=disable` inside compose network | Use `sslmode=require` (or stricter) |
| Volumes | `postgres_data` Docker volume | RDS automated backups |
| Security group | Not exposed publicly | Allow **only** EC2 security group on 5432 |
| Compose change | Default `docker-compose.yml` | Remove or disable `postgres` service; set `DATABASE_URL` on `api` |
| Migrations | API on startup (same) | Same — API runs migrate on boot |
| Seed questions | `000002_seed_questions.up.sql` on first migrate | Same migration file |

For RDS, point the API at the instance and keep Postgres off the public internet. The frontend and Nginx setup are unchanged.

### Deployed smoke-test checklist

Run after `docker compose up` locally or on EC2 (via browser or Nginx URL).

1. **Health** — `curl http://localhost:8080/health` → `{"status":"ok","database":"connected"}`
2. **Profile** — Open app → complete onboarding / save profile → reload → profile persists
3. **Diagnosis** — Submit baseline diagnosis → scores and weak areas appear (requires `OPENAI_API_KEY`)
4. **Plan** — Generate recovery plan → tasks listed on dashboard / plan page
5. **Practice** — Open a question → write solution → submit attempt
6. **Evaluate** — Evaluate attempt → feedback and score render
7. **Progress** — Dashboard / progress page loads without errors

If AI steps fail, check API logs: `docker compose logs api` and confirm `OPENAI_API_KEY` is set in `.env`.

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
