# Interview Preparation Coach — Build Status

**Last updated:** 2026-07-01  
**Core loop:** profile → diagnosis → plan → practice → evaluate → progress → adjust

## Summary

Backend (Phases 5–7), Next.js frontend (Phase 8), end-to-end workflow wiring (Phase 9), and structured loop review (Phase 10) are complete in code. The dashboard drives users through the full recovery loop with recommended practice, auto task completion, progress snapshot refresh, and workflow routing when plans are exhausted.

**Manual QA and deployment are the remaining work.** See [REVIEW.md](./REVIEW.md) for issues fixed, test checklist, and gaps.

## Phase checklist

| Phase | Name | Status |
|-------|------|--------|
| 1 | Freeze MVP scope | Done |
| 2 | Technical architecture | Done |
| 3 | Folder structure | Done — `backend/`, `db/`, `frontend/`, `deploy/` |
| 4 | PostgreSQL schema & migrations | Done |
| 5 | Go backend foundation | Done |
| 6 | Backend product modules | Done (7/7) |
| 7 | AI prompt layer + OpenAI | Done |
| 8 | Next.js frontend | Done |
| 9 | End-to-end recovery loop | Done |
| 10 | Testing & review | Done (automated); manual QA pending |
| 11 | EC2 / Docker deployment | Done |
| 12 | 30-day personal validation plan | Not started |
| 13 | Measurable recovery (failure patterns + velocity) | Not started |
| 14 | Pressure mode + proactive adjust + mock interview | Not started |
| 15 | Briefing + debrief + micro-drills | Not started |
| 16 | Company readiness lens + platform (conditional) | Not started |

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
- [x] Progress snapshots API: `POST` `/api/v1/progress/snapshot`, `GET` `/api/v1/progress/latest`, `GET` `/api/v1/progress/snapshots` (**upsert same-day snapshot** — 201 create / 200 refresh; degraded AI fallback)
- [x] AI layer: `backend/internal/ai/` — Responses API client, 8 versioned prompts, strict JSON validation, logging
- [x] Next.js frontend: all app pages, API client + React Query, shared UI, AI error handling
- [x] End-to-end workflow: diagnosis → plan generation, recommended practice deep-links, post-evaluate snapshot + task completion, workflow stage routing (practice / progress / adjust)
- [x] Loop review fixes: status codes, cache invalidation, workflow stage precedence, error handling (see [REVIEW.md](./REVIEW.md))
- [x] Docker deployment: `backend/Dockerfile`, `frontend/Dockerfile`, `docker-compose.yml`, `.env.example`, Nginx config, README deploy section

## Backend modules

| Module | Status | Notes |
|--------|--------|-------|
| Profile | Done | Single-user upsert/get |
| Baseline diagnosis | Done | OpenAI diagnostic_report JSON |
| Recovery plan | Done | AI-generated 5–10 tasks; replaces active plan on regenerate |
| Practice questions | Done | List/get with filters |
| Practice attempts | Done | Submit/list/get; profile-scoped |
| Evaluations | Done | AI scores/feedback; one per attempt; 409 on duplicate |
| Progress snapshots | Done | Daily upsert + AI summary; `status: complete` or `degraded_ai_summary` |
| AI service | Done | Mock interview prompts stubbed; API returns `ErrNotImplemented` |

## Frontend pages

| Page | Route | Status |
|------|-------|--------|
| Dashboard | `/` | Workflow stage cards + today's task + recommended practice |
| Onboarding / profile | `/onboarding` | Form wired to profile API |
| Baseline diagnosis | `/diagnosis` | Submit + plan generation; redirect to dashboard |
| Recovery plan | `/plan` | Today's tasks + generate plan |
| Daily practice | `/practice`, `/practice/[id]` | Recommended question + attempt/evaluate + auto task complete |
| Mock interview | `/mock-interview` | UI shell (backend stubs) |
| Progress | `/progress` | Snapshot + history + adjustment recommendation |
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
| GET | `/api/v1/plan/today` | Working (404 when exhausted) |
| PATCH | `/api/v1/plan/tasks/{id}/complete` | Working (409 if already done) |
| GET | `/api/v1/questions` | Working |
| GET | `/api/v1/questions/{id}` | Working |
| POST | `/api/v1/attempts` | Working |
| GET | `/api/v1/attempts` | Working |
| GET | `/api/v1/attempts/{id}` | Working |
| POST | `/api/v1/attempts/{id}/evaluate` | Working (OpenAI) |
| GET | `/api/v1/attempts/{id}/evaluation` | Working |
| POST | `/api/v1/progress/snapshot` | Working (201 create / 200 refresh) |
| GET | `/api/v1/progress/latest` | Working |
| GET | `/api/v1/progress/snapshots` | Working |

## Not started

- [ ] Manual end-to-end QA (checklist in [REVIEW.md](./REVIEW.md))
- [ ] Mock interview backend endpoints (prompts exist; service stubs only)
- [ ] Integration tests against real PostgreSQL (optional CI job)
- [ ] Phase 12 — 30-day personal validation playbook

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

- **Manual QA not run** — automated tests pass; live walkthrough still needed ([REVIEW.md](./REVIEW.md) checklist)

## Repo layout notes

- `db/` — migrations, sqlc config, and SQL queries (source of truth)
- `backend/internal/db/` — generated sqlc code + pgx connection pool
- `backend/internal/ai/prompts/` — versioned prompt JSON documents
- `frontend/src/` — Next.js app, API client, hooks, UI components
- `deploy/nginx/` — optional EC2 reverse proxy config
- `docs/progress/` — build status, review notes, and next-step prompts

## Recommended next step

**Phase 12 — 30-day personal validation.** Copy-paste prompt in [NEXT.md](./NEXT.md).

Post-MVP direction: [Innovation Plan](../product/interview-preparation-app-innovation-plan.md) · All prompts: [Master Implementation Prompts](../prompts/interview-preparation-app-master-implementation-prompts.md).
