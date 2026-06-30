# Interview Preparation Coach — Build Status

**Last updated:** 2026-06-30  
**Core loop:** profile → diagnosis → plan → practice → evaluate → progress → adjust

## Summary

Planning and database schema are complete. The Go backend foundation is running with all Phase 6 product modules and **Phase 7 AI integration** wired through `backend/internal/ai`. Diagnosis, plan generation, evaluations, and progress snapshots call OpenAI Responses API with versioned prompts and strict JSON validation. Frontend, end-to-end integration, tests, and deployment are not started.

## Phase checklist

| Phase | Name | Status |
|-------|------|--------|
| 1 | Freeze MVP scope | Done |
| 2 | Technical architecture | Done |
| 3 | Folder structure | Partial — backend + `db/` + `backend/internal/ai/prompts/`; no `frontend/` or `deploy/` |
| 4 | PostgreSQL schema & migrations | Done |
| 5 | Go backend foundation | Done |
| 6 | Backend product modules | Done (7/7) |
| 7 | AI prompt layer + OpenAI | Done |
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
- [x] Diagnosis API: `POST` / `GET` `/api/v1/diagnosis` (**OpenAI baseline diagnosis**)
- [x] Recovery plan API: `POST` / `GET` `/api/v1/plan`, `GET` `/api/v1/plan/today`, `PATCH` `/api/v1/plan/tasks/{id}/complete` (**AI-driven tasks from diagnostic_report**)
- [x] Practice questions API: `GET` `/api/v1/questions`, `GET` `/api/v1/questions/{id}` (15 seeded questions)
- [x] Practice attempts API: `POST` / `GET` `/api/v1/attempts`, `GET` `/api/v1/attempts/{id}`
- [x] Evaluations API: `POST` `/api/v1/attempts/{id}/evaluate`, `GET` `/api/v1/attempts/{id}/evaluation` (**OpenAI scoring by question type**)
- [x] Progress snapshots API: `POST` `/api/v1/progress/snapshot`, `GET` `/api/v1/progress/latest`, `GET` `/api/v1/progress/snapshots` (**AI weekly summary**; `status: complete` or `degraded_ai_summary`)
- [x] AI layer: `backend/internal/ai/` — Responses API client, 8 versioned prompts, strict JSON validation, logging

## Backend modules

| Module | Status | Notes |
|--------|--------|-------|
| Profile | Done | Single-user upsert/get |
| Baseline diagnosis | Done | OpenAI diagnostic_report JSON |
| Recovery plan | Done | AI-generated 5–10 tasks; replaces active plan on regenerate |
| Practice questions | Done | List/get with filters |
| Practice attempts | Done | Submit/list/get; profile-scoped |
| Evaluations | Done | AI scores/feedback; one per attempt |
| Progress snapshots | Done | Daily snapshot + AI summary; 409 if already exists for today |
| AI service | Done | Mock interview prompts stubbed; API returns `ErrNotImplemented` |

## API surface today

| Method | Path | Status |
|--------|------|--------|
| GET | `/health` | Working |
| POST | `/api/v1/profile` | Working |
| GET | `/api/v1/profile` | Working |
| POST | `/api/v1/diagnosis` | Working (OpenAI) |
| GET | `/api/v1/diagnosis` | Working |
| POST | `/api/v1/plan` | Working (OpenAI) |
| GET | `/api/v1/plan` | Working |
| GET | `/api/v1/plan/today` | Working |
| PATCH | `/api/v1/plan/tasks/{id}/complete` | Working |
| GET | `/api/v1/questions` | Working |
| GET | `/api/v1/questions/{id}` | Working |
| POST | `/api/v1/attempts` | Working |
| GET | `/api/v1/attempts` | Working |
| GET | `/api/v1/attempts/{id}` | Working |
| POST | `/api/v1/attempts/{id}/evaluate` | Working (OpenAI) |
| GET | `/api/v1/attempts/{id}/evaluation` | Working |
| POST | `/api/v1/progress/snapshot` | Working (OpenAI summary) |
| GET | `/api/v1/progress/latest` | Working |
| GET | `/api/v1/progress/snapshots` | Working |

## Not started

- [ ] Next.js frontend (onboarding, diagnosis, plan, practice, mock interview, progress, settings)
- [ ] Mock interview backend endpoints (prompts exist; service stubs only)
- [ ] End-to-end workflow integration
- [ ] Integration tests with mocked OpenAI
- [ ] Deployment (Dockerfiles, docker-compose, EC2)

## Blockers

- **No frontend** — blocks personal daily use

## Repo layout notes

- `db/` — migrations, sqlc config, and SQL queries (source of truth)
- `backend/internal/db/` — generated sqlc code + pgx connection pool
- `backend/internal/ai/prompts/` — versioned prompt JSON documents
- `docs/progress/` — build status and next-step prompts

## Recommended next step

Implement the **Next.js frontend foundation**. See [NEXT.md](./NEXT.md) for the copy-paste prompt.
