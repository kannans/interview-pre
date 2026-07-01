# Next Up

**Last updated:** 2026-07-01  
**Current focus:** Phase 12 — 30-day personal validation (daily dogfood before any innovation build)

## Recently completed

- Phase 11 — Docker / EC2 deployment (Dockerfiles, compose, env examples, Nginx config, README)
- Phase 10 loop review (automated tests + review doc)
- Practice workspace refactor (local run, resizable editor, draft persistence)
- Practice draft SSR hydration fix (`practice-drafts.ts`)
- [Innovation Plan](../product/interview-preparation-app-innovation-plan.md) + implementation prompts for Phases 11–16

## Before you start Phase 12

| Check | Command / doc |
|-------|----------------|
| Backend tests | `cd backend && go test ./...` |
| Frontend build | `cd frontend && npm run build` |
| Docker stack | `cp .env.example .env` → set `OPENAI_API_KEY` → `docker compose up --build -d` |
| Health | `curl http://localhost:8080/health` |
| Manual QA | [REVIEW.md](./REVIEW.md) checklist (run or acknowledge gaps) |
| Local smoke | `./runserver.sh` or Docker → profile → diagnosis → practice → evaluate |

Do **not** start Phase 13 (failure patterns, velocity, pressure mode) until Phase 12 validation passes.

---

## Phase 12 — Copy-paste prompt

Use in a new agent session. Full version in [Master Implementation Prompts](../prompts/interview-preparation-app-master-implementation-prompts.md#phase-12-30-day-personal-validation).

```text
You are creating a 30-day personal validation plan for Interview Preparation Coach.
...
```

(See prompts doc for the full Phase 12 text.)

---

## Phase 11 — Done

Docker deployment artifacts added:

| Artifact | Purpose |
|----------|---------|
| `backend/Dockerfile` | Multi-stage Go API image |
| `frontend/Dockerfile` | Next.js standalone production image |
| `docker-compose.yml` | Postgres + API + frontend |
| `.env.example` | Compose secrets template (root) |
| `frontend/.env.example` | Frontend API URL for local dev |
| `deploy/nginx/interview-prep.conf` | Optional EC2 reverse proxy |
| `README.md` | Deploy steps, env vars, smoke test, RDS notes |

**Local Docker:**

```bash
cp .env.example .env   # set OPENAI_API_KEY
docker compose up --build -d
curl http://localhost:8080/health
open http://localhost:3000
```

<details>
<summary>Phase 11 original prompt (archive)</summary>

Context:
- Core loop is complete: profile -> diagnosis -> plan -> practice -> evaluate -> progress -> adjust
- Backend: Go + Chi at backend/cmd/api
- Frontend: Next.js at frontend/
- Database: PostgreSQL with golang-migrate migrations in db/migrations/
- Migrations should run on API startup or via a documented init step
- Single-user MVP — no auth layer
- OpenAI key only in backend environment

Task:
Add Docker and deployment artifacts. Do not add new product features.

Deliver:
1. Dockerfile for Go API (multi-stage build, production binary)
2. Dockerfile for Next.js frontend (production build)
3. docker-compose.yml — Postgres + API + frontend (local and single-server)
4. .env.example files for backend and frontend (no secrets committed)
5. README section: deploy steps, env vars, smoke test
6. Optional: Nginx reverse proxy config for EC2 single-host deploy
7. Optional: notes for RDS vs containerized Postgres

Implementation requirements:
- Frontend API base URL configurable at build/runtime
- CORS allows frontend origin in compose
- Health check: GET /health
- Document how to seed DB / run migrations
- Document ports: API :8080, frontend :3000, Postgres :5432

Return:
1. Files added or changed (list with purpose)
2. Commands to run locally with docker compose
3. EC2 setup steps (security groups, env, restart strategy)
4. Smoke-test checklist for deployed stack (profile -> diagnosis -> practice -> evaluate)
5. What differs when using RDS instead of local Postgres

Rules:
- Do not expose secrets in code or committed env files
- Keep it practical for a solo developer on macOS deploying to one EC2 instance
- Do not add auth, billing, or innovation features (Phases 13+)
- Match existing repo conventions (sqlc, migrate, Chi routes)
```

</details>

---

## Roadmap after Phase 11

| Phase | Focus | Prompt |
|-------|--------|--------|
| **12** | 30-day personal validation — prove the loop helps you daily | [Prompts § Phase 12](../prompts/interview-preparation-app-master-implementation-prompts.md#phase-12-30-day-personal-validation) |
| **13** | Failure pattern memory + recovery velocity + question bank | [Prompts § Phase 13](../prompts/interview-preparation-app-master-implementation-prompts.md#phase-13-measurable-recovery-failure-patterns--velocity) |
| **14** | Pressure mode + proactive plan adjust + mock interview | [Prompts § Phase 14](../prompts/interview-preparation-app-master-implementation-prompts.md#phase-14-pressure-proactive-adjust-mock-interview) |
| **15** | Pre-interview briefing + real debrief + micro-drills | [Prompts § Phase 15](../prompts/interview-preparation-app-master-implementation-prompts.md#phase-15-real-world-closure-briefing-debrief-micro-drills) |
| **16** | Company readiness lens + platform (conditional) | [Prompts § Phase 16](../prompts/interview-preparation-app-master-implementation-prompts.md#phase-16-strategic-expansion-company-lens--platform) |

Product direction for Phases 13–16: [Innovation Plan](../product/interview-preparation-app-innovation-plan.md).

**Motto:** Turn interview rust into measurable recovery before your next real round.

---

## Local dev

```bash
./runserver.sh
```

Or separately:

```bash
# Terminal 1 — API
cd backend && go run ./cmd/api

# Terminal 2 — frontend
cd frontend && cp .env.local.example .env.local   # if not done
cd frontend && npm run dev
```

- Frontend: http://localhost:3000  
- API: http://localhost:8080/api/v1  

See README for database setup and AI curl examples (`backend/internal/ai/README.md`).

## Progress docs

| File | Purpose |
|------|---------|
| [STATUS.md](./STATUS.md) | Phase checklist and module status |
| [REVIEW.md](./REVIEW.md) | Phase 10 loop review + manual QA checklist |
| [NEXT.md](./NEXT.md) | **This file** — what to build next |
| [Innovation Plan](../product/interview-preparation-app-innovation-plan.md) | Post-MVP product direction |
| [Implementation Prompts](../prompts/interview-preparation-app-master-implementation-prompts.md) | All phase copy-paste prompts |
