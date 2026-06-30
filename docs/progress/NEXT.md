# Next Implementation Step

**Current focus:** Phase 6 — Practice attempts module  
**Prerequisite:** Profile, baseline diagnosis, recovery plan, and practice questions modules are complete.

Use the prompt below in a new agent session. Review and test the output before moving to evaluations.

---

## Copy-paste prompt: Practice attempts module

```text
Implement the practice attempts backend module for Interview Preparation Coach.

Context:
- Personal-first interview recovery app for experienced engineers.
- Stack: Go + Chi + pgx + sqlc + PostgreSQL.
- Core loop: profile -> diagnosis -> plan -> practice -> evaluate -> progress -> adjust.
- Profile, baseline diagnosis, recovery plan, and practice questions modules are already implemented.
- Practice questions API: GET /api/v1/questions, GET /api/v1/questions/{id} (15 seeded questions).

Existing schema (already migrated):
- attempts (id, profile_id, question_id, submitted_content, time_taken_seconds, created_at)

Existing API patterns to follow:
- backend/internal/repository/ — database access via sqlc queries
- backend/internal/service/ — business logic and validation
- backend/internal/http/handler/ — thin HTTP handlers
- backend/internal/http/router.go — wire new routes under /api/v1

Task:
Implement the practice attempts module only. Do not implement frontend, OpenAI, evaluations, or progress snapshots yet.

Requirements:
1. Add sqlc queries in db/queries.sql for:
   - Create attempt
   - Get attempt by id
   - List attempts for profile (optional filter by question_id)

2. Regenerate sqlc output (document the command).

3. Implement repository, service, and handler layers matching existing module style.

4. API endpoints (suggest and implement):
   - POST /api/v1/attempts — submit answer + time_taken_seconds for a question
   - GET /api/v1/attempts — list attempts for the current profile (optional ?question_id=)
   - GET /api/v1/attempts/{id} — get single attempt

5. Validation and errors:
   - Profile must exist (412 Precondition Failed, matching diagnosis/plan pattern)
   - Question must exist (404 or 400 as appropriate)
   - submitted_content must be non-empty
   - time_taken_seconds must be > 0
   - Return 404 for unknown attempt id
   - Attempt must belong to the current profile when fetching by id

6. Single-user model: resolve profile the same way other modules do (GetProfile LIMIT 1).

Return:
1. Files created or changed
2. New API request/response JSON examples
3. sqlc regenerate command
4. How to test manually (curl examples)
5. What is left for the next module (evaluations)

Rules:
- Implement one module only.
- Keep handlers thin; business logic in services.
- Match existing code conventions in backend/.
- Do not add OpenAI calls yet.
- Do not build frontend.
- Do not implement evaluations yet — attempts are stored only.
```

---

## After practice attempts

Continue Phase 6 in this order:

1. **Evaluations** — store scores and feedback (placeholder AI until Phase 7)
2. **Progress snapshots** — aggregate stats and readiness index

Then proceed to Phase 7 (AI prompt layer) per `docs/prompts/interview-preparation-app-master-implementation-prompts.md`.
