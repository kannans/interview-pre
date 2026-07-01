# Next Up

**Last updated:** 2026-07-01  
**Current focus:** Phase 12 — Execute 30-day personal validation (daily dogfood before any innovation build)

## Recently completed

- Phase 12 — validation playbook + working templates ([VALIDATION.md](./VALIDATION.md))
- Phase 11 — Docker / EC2 deployment (Dockerfiles, compose, env examples, Nginx config, README)
- Phase 10 loop review (automated tests + review doc)
- Practice workspace refactor (local run, resizable editor, draft persistence)
- Practice draft SSR hydration fix (`practice-drafts.ts`)
- [Innovation Plan](../product/interview-preparation-app-innovation-plan.md) + implementation prompts for Phases 11–16

## Your job now (Phase 12 execution)

**Do not write code.** Use the app daily for 30 days and decide GO / PIVOT / STOP.

| Step | Doc |
|------|-----|
| 1. Day 0 setup | [VALIDATION.md §1](./VALIDATION.md#1-pre-validation-setup-day-0) + [REVIEW.md](./REVIEW.md) QA |
| 2. Baseline | [validation/baseline-metrics.md](./validation/baseline-metrics.md) |
| 3. Daily practice | [VALIDATION.md §2](./VALIDATION.md#2-daily-usage-routine-monsun) → log in [session-logs.md](./validation/session-logs.md) |
| 4. Friction | [friction-log.md](./validation/friction-log.md) as issues arise |
| 5. Weekly review | [weekly-reviews.md](./validation/weekly-reviews.md) every Sunday |
| 6. Day 30 gate | [VALIDATION.md §7](./VALIDATION.md#7-decision-gates-at-day-30) |

**Targets:** ≥ 4 practice days/week · median open → editor < 2 min · readiness_index +5 Week 1 → Week 4.

Do **not** start Phase 13 until Day 30 gate = **GO**.

---

## Before Day 0

| Check | Command / doc |
|-------|----------------|
| Backend tests | `cd backend && go test ./...` |
| Frontend build | `cd frontend && npm run build` |
| Docker stack | `cp .env.example .env` → set `OPENAI_API_KEY` → `docker compose up --build -d` |
| Health | `curl http://localhost:8080/health` |
| Manual QA | [REVIEW.md](./REVIEW.md) checklist — required on Day 0 |
| Local smoke | `./runserver.sh` or Docker → profile → diagnosis → practice → evaluate |

---

## Phase 12 — Playbook (done)

Full playbook: **[VALIDATION.md](./VALIDATION.md)**

Working files:

```
docs/progress/validation/
├── baseline-metrics.md   # Day 0
├── session-logs.md       # per session
├── friction-log.md       # backlog
└── weekly-reviews.md     # Sundays + Day 30 decision
```

<details>
<summary>Phase 12 original prompt (archive)</summary>

See [Master Implementation Prompts § Phase 12](../prompts/interview-preparation-app-master-implementation-prompts.md#phase-12-30-day-personal-validation).

</details>

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

---

## Roadmap after Phase 12 gate

| Phase | Focus | When |
|-------|--------|------|
| **12** | 30-day validation — **you are here (execute)** | Now → Day 30 |
| **13** | Failure pattern memory + recovery velocity + question bank | After GO gate |
| **14** | Pressure mode + proactive plan adjust + mock interview | After Phase 13 |
| **15** | Pre-interview briefing + real debrief + micro-drills | After Phase 14 |
| **16** | Company readiness lens + platform (conditional) | After Phase 15 |

Prompts: [Master Implementation Prompts](../prompts/interview-preparation-app-master-implementation-prompts.md) · Product direction: [Innovation Plan](../product/interview-preparation-app-innovation-plan.md).

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
| [VALIDATION.md](./VALIDATION.md) | **30-day validation playbook** |
| [NEXT.md](./NEXT.md) | **This file** — what to do next |
| [Innovation Plan](../product/interview-preparation-app-innovation-plan.md) | Post-MVP product direction |
| [Implementation Prompts](../prompts/interview-preparation-app-master-implementation-prompts.md) | All phase copy-paste prompts |
