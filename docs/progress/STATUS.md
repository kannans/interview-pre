# Interview Preparation Coach — Build Status

**Last updated:** 2026-06-30  
**Core loop:** profile → diagnosis → plan → practice → evaluate → progress → adjust

## Summary

Planning and database schema are complete. The Go backend foundation is running with Profile, Baseline Diagnosis, Recovery Plan, Practice Questions, Practice Attempts, and Evaluations APIs. Diagnosis, plan generation, and evaluations use placeholder/static logic (no OpenAI yet). Progress snapshots, frontend, AI, tests, and deployment are not started.

## Phase checklist

| Phase | Name | Status |
|-------|------|--------|
| 1 | Freeze MVP scope | Done |
| 2 | Technical architecture | Done |
| 3 | Folder structure | Partial — backend + `db/` exist; no `frontend/`, `deploy/`, or AI prompts folder |
| 4 | PostgreSQL schema & migrations | Done |
| 5 | Go backend foundation | Done |
| 6 | Backend product modules | Partial (6/7) |
| 7 | AI prompt layer + OpenAI | Not started |
| 8 | Next.js frontend | Not started |
| 9 | End-to-end recovery loop | Not started |
| 10 | Testing & review | Not started |
| 11 | EC2 / Docker deployment | Not started |
| 12 | 30-day personal validation plan | Not started |

## Done

- [x] MVP scope, PRD, technical design, build plan (`docs/`)
- [x] PostgreSQL initial migration — 8 tables (`db/migrations/000001_init_schema.up.sql`)
- [x] Go backend: config, Chi router, pgx pool, slog, health, CORS, graceful shutdown
- [x] Profile API: `POST` / `GET` `/api/v1/profile`
- [x] Diagnosis API: `POST` / `GET` `/api/v1/diagnosis` (placeholder AI report)
- [x] Recovery plan API: `POST` / `GET` `/api/v1/plan`, `GET` `/api/v1/plan/today`, `PATCH` `/api/v1/plan/tasks/{id}/complete` (static plan from `weak_areas_self`)
- [x] Practice questions API: `GET` `/api/v1/questions`, `GET` `/api/v1/questions/{id}` (15 seeded questions in `000002_seed_questions`)
- [x] Practice attempts API: `POST` / `GET` `/api/v1/attempts`, `GET` `/api/v1/attempts/{id}`
- [x] Evaluations API: `POST` `/api/v1/attempts/{id}/evaluate`, `GET` `/api/v1/attempts/{id}/evaluation` (placeholder scoring)

## Backend modules (Phase 6)

| Module | Status | Notes |
|--------|--------|-------|
| Profile | Done | Single-user upsert/get |
| Baseline diagnosis | Done | Saves answers; placeholder `diagnostic_report` JSON |
| Recovery plan | Done | Static 7–10 tasks; replaces existing active plan on regenerate |
| Practice questions | Done | List/get with filters; seed migration `000002_seed_questions` |
| Practice attempts | Done | Submit/list/get; profile-scoped |
| Evaluations | Done | Placeholder scores/feedback; one per attempt; profile-scoped |
| Progress snapshots | Not started | Table exists; no API |

## API surface today

| Method | Path | Status |
|--------|------|--------|
| GET | `/health` | Working |
| POST | `/api/v1/profile` | Working |
| GET | `/api/v1/profile` | Working |
| POST | `/api/v1/diagnosis` | Working (no real AI) |
| GET | `/api/v1/diagnosis` | Working |
| POST | `/api/v1/plan` | Working (static generation) |
| GET | `/api/v1/plan` | Working |
| GET | `/api/v1/plan/today` | Working |
| PATCH | `/api/v1/plan/tasks/{id}/complete` | Working |
| GET | `/api/v1/questions` | Working (optional `type`, `category`, `level` filters) |
| GET | `/api/v1/questions/{id}` | Working |
| POST | `/api/v1/attempts` | Working |
| GET | `/api/v1/attempts` | Working (optional `question_id` filter) |
| GET | `/api/v1/attempts/{id}` | Working |
| POST | `/api/v1/attempts/{id}/evaluate` | Working (placeholder scoring) |
| GET | `/api/v1/attempts/{id}/evaluation` | Working |

## Not started

- [ ] Progress snapshots
- [ ] AI prompt layer + OpenAI integration (`backend/internal/ai/`)
- [ ] Next.js frontend (onboarding, diagnosis, plan, practice, mock interview, progress, settings)
- [ ] End-to-end workflow integration
- [ ] Tests (`*_test.go`, migration smoke tests, AI parsing tests)
- [ ] Deployment (Dockerfiles, docker-compose, EC2)

## Blockers

- **OpenAI not wired** — blocks real diagnosis, plan generation, and attempt evaluations
- **No frontend** — blocks personal daily use
- **No progress snapshots API** — blocks readiness tracking and weekly trends

## Repo layout notes

- `db/` — migrations, sqlc config, and SQL queries (source of truth)
- `backend/internal/db/` — generated sqlc code + pgx connection pool
- `docs/progress/` — build status and next-step prompts

## Recommended next step

Implement the **progress snapshots** backend module. See [NEXT.md](./NEXT.md) for the copy-paste prompt.
