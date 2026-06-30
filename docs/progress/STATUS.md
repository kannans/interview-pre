# Interview Preparation Coach — Build Status

**Last updated:** 2026-06-30  
**Core loop:** profile → diagnosis → plan → practice → evaluate → progress → adjust

## Summary

Planning, database schema, Go backend (Phases 5–7), and **Phase 8 Next.js frontend foundation** are complete. The frontend includes all seven app pages, API client layer, React Query hooks, shared UI components, and AI error handling. End-to-end workflow integration, comprehensive testing, and deployment are not started.

## Phase checklist

| Phase | Name | Status |
|-------|------|--------|
| 1 | Freeze MVP scope | Done |
| 2 | Technical architecture | Done |
| 3 | Folder structure | Partial — `backend/`, `db/`, `frontend/`; no `deploy/` |
| 4 | PostgreSQL schema & migrations | Done |
| 5 | Go backend foundation | Done |
| 6 | Backend product modules | Done (7/7) |
| 7 | AI prompt layer + OpenAI | Done |
| 8 | Next.js frontend | Done |
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
- [x] Next.js frontend: onboarding, diagnosis, plan, practice, mock interview (UI shell), progress, settings; API client + React Query; app shell layout

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

## Frontend pages

| Page | Route | Status |
|------|-------|--------|
| Dashboard | `/` | Shell + workflow stage cards |
| Onboarding / profile | `/onboarding` | Form wired to profile API |
| Baseline diagnosis | `/diagnosis` | Submit + view diagnosis |
| Recovery plan | `/plan` | Today's tasks + generate plan |
| Daily practice | `/practice`, `/practice/[id]` | Question list + attempt submit |
| Mock interview | `/mock-interview` | UI shell (backend stubs) |
| Progress | `/progress` | Snapshot + latest stats |
| Settings | `/settings` | API base URL config |

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

- [ ] End-to-end workflow integration (verify full loop in UI)
- [ ] Mock interview backend endpoints (prompts exist; service stubs only)
- [ ] Integration tests against real PostgreSQL (optional CI job)
- [ ] Deployment (Dockerfiles, docker-compose, EC2)

## Testing

Backend unit + handler tests (no DB/OpenAI required):

```bash
cd backend && go test ./...
```

Frontend production build:

```bash
cd frontend && npm run build
```

## Blockers

- **End-to-end loop not verified** — pages exist but full profile → progress flow needs integration QA

## Repo layout notes

- `db/` — migrations, sqlc config, and SQL queries (source of truth)
- `backend/internal/db/` — generated sqlc code + pgx connection pool
- `backend/internal/ai/prompts/` — versioned prompt JSON documents
- `frontend/src/` — Next.js app, API client, hooks, UI components
- `docs/progress/` — build status and next-step prompts

## Recommended next step

Wire the **end-to-end recovery workflow** (Phase 9). See [NEXT.md](./NEXT.md) for the copy-paste prompt.
