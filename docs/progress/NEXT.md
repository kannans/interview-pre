# Next Implementation Step

**Current focus:** Phase 11 — EC2 / Docker deployment  
**Prerequisite:** Phase 10 review complete; `go test ./...` and `npm run build` pass; manual QA checklist in [REVIEW.md](./REVIEW.md) run (or acknowledged).

Use the prompt below in a new agent session after committing any pending review fixes.

---

## Copy-paste prompt: Deployment

```text
Add Docker and deployment artifacts for the interview prep coach MVP.

Deliver:
1. Dockerfile for Go API (backend/cmd/api)
2. Dockerfile for Next.js frontend (production build)
3. docker-compose.yml — Postgres + API + frontend
4. Environment variable documentation (.env.example files)
5. Update README with deploy and smoke-test steps

Constraints:
- Single-user MVP; no auth layer
- Postgres migrations run on API startup or via documented init step
- Frontend talks to API via configurable base URL
- Do not add new product features

Return:
1. files added/changed
2. how to run locally with docker compose
3. smoke-test checklist for deployed stack
```

---

## After deployment

Phase 11 ships the stack. Proceed to **Phase 12 — 30-day personal validation plan** per `docs/prompts/interview-preparation-app-master-implementation-prompts.md`:

1. Use the app daily for real interview prep
2. Log friction, AI quality issues, and missing features
3. Decide what to cut or improve post-MVP

## Local dev (frontend + backend)

```bash
# Terminal 1 — API (from repo root)
cd backend && go run ./cmd/api

# Terminal 2 — frontend
cd frontend && cp .env.local.example .env.local   # if not done
cd frontend && npm run dev
```

Frontend: http://localhost:3000 — API: http://localhost:8080/api/v1

See README for database setup and manual AI curl examples (`backend/internal/ai/README.md`).

## Progress docs

| File | Purpose |
|------|---------|
| [STATUS.md](./STATUS.md) | Phase checklist and module status |
| [REVIEW.md](./REVIEW.md) | Phase 10 loop review, fixes, manual QA checklist |
| [NEXT.md](./NEXT.md) | Current implementation prompt |
