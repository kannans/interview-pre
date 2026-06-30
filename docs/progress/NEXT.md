# Next Implementation Step

**Current focus:** Phase 10 — Testing and review  
**Prerequisite:** Phase 9 end-to-end recovery loop wired; backend tests pass (`cd backend && go test ./...`); frontend builds (`cd frontend && npm run build`).

Use the prompt below in a new agent session. Review and test the output before moving to Phase 11 (deployment).

---

## Copy-paste prompt: Testing and review

```text
Review the end-to-end recovery loop for correctness, gaps, and error handling.

The loop is:
profile -> diagnosis -> plan -> practice -> evaluate -> progress -> adjust

Verify and fix as needed:
1. backend API correctness (status codes, validation, edge cases)
2. AI response parsing and degraded-mode fallbacks
3. frontend integration (hooks, cache invalidation, workflow stage)
4. progress snapshot refresh after evaluation
5. workflow transitions when today's task is complete or plan is exhausted

Return:
1. issues found and fixes applied
2. manual test checklist with pass/fail notes
3. remaining gaps before deployment
4. recommended follow-ups

Rules:
- Do not add new feature scope.
- Focus on reliability, error handling, and test coverage for the core loop.
```

---

## After testing and review

Phase 10 validates the full loop. Proceed to **Phase 11 — EC2 / Docker deployment** per `docs/prompts/interview-preparation-app-master-implementation-prompts.md`:

1. Add Dockerfiles and docker-compose for API + Postgres
2. Document environment variables and deployment steps
3. Smoke-test the deployed stack

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
