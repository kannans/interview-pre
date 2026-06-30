# Next Implementation Step

**Current focus:** Phase 9 — End-to-end recovery workflow  
**Prerequisite:** Phase 8 frontend foundation complete; backend tests pass (`cd backend && go test ./...`); frontend builds (`cd frontend && npm run build`).

Use the prompt below in a new agent session. Review and test the output before moving to Phase 10 (testing and review).

---

## Copy-paste prompt: End-to-end recovery workflow

```text
Connect the frontend and backend into the full personal recovery loop.

The loop is:
profile -> diagnosis -> plan -> practice -> evaluate -> progress -> adjust

Implement or verify:
1. profile creation saves correctly
2. baseline diagnosis can be submitted
3. diagnosis result creates or updates recovery plan
4. user can start today's practice
5. user can submit answer or code
6. backend evaluates attempt
7. evaluation is saved
8. progress dashboard updates
9. next recommended task is shown

Return:
1. files changed
2. integration flow summary
3. API calls used by each page
4. data persisted at each step
5. remaining gaps
6. manual test checklist

Rules:
- Do not add new feature scope.
- Focus on making the core loop reliable.
```

---

## After end-to-end integration

Phase 9 wires the full loop. Proceed to **Phase 10 — Testing and review** per `docs/prompts/interview-preparation-app-master-implementation-prompts.md`:

1. Review backend API correctness, AI parsing, and frontend integration
2. Fix workflow gaps and error handling
3. Document findings before deployment (Phase 11)

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
