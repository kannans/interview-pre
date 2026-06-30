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
| AI | OpenAI Responses API via Go backend only *(not wired yet)* |
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
│   └── internal/     # HTTP handlers, services, repositories, generated db code
├── db/               # Database source of truth (migrations, sqlc config, SQL queries)
├── docs/             # Product, strategy, and implementation docs
│   └── progress/     # Build status and next-step prompts
└── frontend/         # Next.js app (planned)
```

**Why `db/` is at the root:** migrations and SQL queries are shared infrastructure, not Go application code. `sqlc` reads from `db/` and generates typed Go into `backend/internal/db/`. This is a common pattern for Go + sqlc + golang-migrate projects.

## Current status

**In progress.** Backend foundation is running. Profile, baseline diagnosis, recovery plan, practice questions, and practice attempts APIs are implemented. Evaluations, progress, AI, and frontend are not started.

| Area | Status |
|------|--------|
| Product docs & schema | Done |
| Go backend foundation | Done |
| Profile + diagnosis APIs | Done |
| Recovery plan API | Done |
| Practice questions API | Done |
| Practice attempts API | Done |
| Evaluations, progress | Not started |
| AI integration | Not started |
| Frontend | Not started |
| Tests & deployment | Not started |

For the full checklist and next implementation step, see [`docs/progress/STATUS.md`](docs/progress/STATUS.md) and [`docs/progress/NEXT.md`](docs/progress/NEXT.md).

### Live API endpoints

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/health` | Health check |
| `POST` | `/api/v1/profile` | Create or update profile |
| `GET` | `/api/v1/profile` | Get profile |
| `POST` | `/api/v1/diagnosis` | Submit baseline diagnosis |
| `GET` | `/api/v1/diagnosis` | Get latest diagnosis |
| `POST` | `/api/v1/plan` | Generate or replace active recovery plan |
| `GET` | `/api/v1/plan` | Get active plan with tasks |
| `GET` | `/api/v1/plan/today` | Get today's recommended task |
| `PATCH` | `/api/v1/plan/tasks/{id}/complete` | Mark a plan task complete |
| `GET` | `/api/v1/questions` | List practice questions (optional filters) |
| `GET` | `/api/v1/questions/{id}` | Get a single question |
| `POST` | `/api/v1/attempts` | Submit a practice attempt |
| `GET` | `/api/v1/attempts` | List attempts (optional `question_id` filter) |
| `GET` | `/api/v1/attempts/{id}` | Get a single attempt |

## Local development

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

Run migrations:

```bash
migrate -path db/migrations -database "postgres://postgres:postgres@localhost:5432/interview_prep?sslmode=disable" up
```

### 2. Backend configuration

```bash
cp backend/.env.example backend/.env
```

Edit `backend/.env` with your `DATABASE_URL` and optionally `OPENAI_API_KEY`.

### 3. Generate sqlc code (after changing SQL)

```bash
cd db && sqlc generate
```

### 4. Run the API server

```bash
cd backend && go run ./cmd/api
```

The server starts on port `8080` by default. Verify with:

```bash
curl http://localhost:8080/health
```

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
