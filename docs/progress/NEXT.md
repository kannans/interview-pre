# Next Implementation Step

**Current focus:** Phase 6 — Evaluations module  
**Prerequisite:** Profile, baseline diagnosis, recovery plan, practice questions, and practice attempts modules are complete.

Use the prompt below in a new agent session. Review and test the output before moving to progress snapshots.

---

## Copy-paste prompt: Evaluations module

```text
Implement the evaluations backend module for Interview Preparation Coach.

Context:
- Personal-first interview recovery app for experienced engineers.
- Stack: Go + Chi + pgx + sqlc + PostgreSQL.
- Core loop: profile -> diagnosis -> plan -> practice -> evaluate -> progress -> adjust.
- Profile, baseline diagnosis, recovery plan, practice questions, and practice attempts modules are already implemented.
- Practice attempts API: POST/GET /api/v1/attempts, GET /api/v1/attempts/{id}.

Existing schema (already migrated):
- evaluations (id, attempt_id, score_correctness, score_approach, score_edge_cases, score_communication, feedback_*, raw_ai_response, created_at)

Existing API patterns to follow:
- backend/internal/repository/ — database access via sqlc queries
- backend/internal/service/ — business logic and validation
- backend/internal/http/handler/ — thin HTTP handlers
- backend/internal/http/router.go — wire new routes under /api/v1

Task:
Implement the evaluations module only. Do not implement frontend, real OpenAI calls, or progress snapshots yet.

Requirements:
1. Add sqlc queries in db/queries.sql for:
   - Create evaluation (linked to attempt_id)
   - Get evaluation by attempt id
   - Get evaluation by id

2. Regenerate sqlc output (document the command).

3. Implement repository, service, and handler layers matching existing module style.

4. API endpoints (suggest and implement):
   - POST /api/v1/attempts/{id}/evaluate — create evaluation for an attempt (placeholder scoring until Phase 7)
   - GET /api/v1/attempts/{id}/evaluation — get evaluation for an attempt

5. Placeholder evaluation logic:
   - Generate static/placeholder scores (0–10) and feedback text based on attempt metadata
   - Store placeholder raw_ai_response JSON indicating pending real AI
   - One evaluation per attempt (unique constraint on attempt_id)

6. Validation and errors:
   - Profile must exist
   - Attempt must exist and belong to current profile
   - Return 404 for unknown attempt or evaluation
   - Return 409 if evaluation already exists for attempt
   - Validate score ranges if accepting input (prefer server-generated scores for v1)

Return:
1. Files created or changed
2. New API request/response JSON examples
3. sqlc regenerate command
4. How to test manually (curl examples)
5. What is left for the next module (progress snapshots)

Rules:
- Implement one module only.
- Keep handlers thin; business logic in services.
- Match existing code conventions in backend/.
- Do not call OpenAI yet — use placeholder evaluation logic.
- Do not build frontend.
```

---

## After evaluations

Continue Phase 6:

1. **Progress snapshots** — aggregate stats and readiness index

Then proceed to Phase 7 (AI prompt layer) per `docs/prompts/interview-preparation-app-master-implementation-prompts.md`.
