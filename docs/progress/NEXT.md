# Next Implementation Step

**Current focus:** Phase 6 — Practice questions module  
**Prerequisite:** Profile, baseline diagnosis, and recovery plan modules are complete.

Use the prompt below in a new agent session. Review and test the output before moving to practice attempts.

---

## Copy-paste prompt: Practice questions module

```text
Implement the practice questions backend module for Interview Preparation Coach.

Context:
- Personal-first interview recovery app for experienced engineers.
- Stack: Go + Chi + pgx + sqlc + PostgreSQL.
- Core loop: profile -> diagnosis -> plan -> practice -> evaluate -> progress -> adjust.
- Profile, baseline diagnosis, and recovery plan modules are already implemented.
- Recovery plan tasks use categories: programming_basics, dsa, concept.

Existing schema (already migrated):
- questions (id, type, category, level, title, content, constraints, starter_code, solution_hint, created_at)

Existing API patterns to follow:
- backend/internal/repository/ — database access via sqlc queries
- backend/internal/service/ — business logic and validation
- backend/internal/http/handler/ — thin HTTP handlers
- backend/internal/http/router.go — wire new routes under /api/v1

Task:
Implement the practice questions module only. Do not implement frontend, OpenAI, practice attempts, or evaluations yet.

Requirements:
1. Add sqlc queries in db/queries.sql for:
   - List questions (filter by type, category, level — all optional)
   - Get question by id

2. Regenerate sqlc output (document the command).

3. Implement repository, service, and handler layers matching existing profile/diagnosis/plan style.

4. API endpoints (suggest and implement):
   - GET /api/v1/questions — list with optional query params (type, category, level)
   - GET /api/v1/questions/{id} — get single question

5. Seed data:
   - Add a seed SQL file or migration with 10–15 starter questions across programming_basics, dsa, and concept
   - Cover level_1 through level_3 where appropriate
   - Align categories with recovery plan task categories

6. Validation and errors:
   - Return 404 for unknown question id
   - Validate query param values where applicable

Return:
1. Files created or changed
2. New API request/response JSON examples
3. sqlc regenerate command
4. How to test manually (curl examples)
5. What is left for the next module (practice attempts)

Rules:
- Implement one module only.
- Keep handlers thin; business logic in services.
- Match existing code conventions in backend/.
- Do not add OpenAI calls yet.
- Do not build frontend.
```

---

## After practice questions

Continue Phase 6 in this order:

1. **Practice attempts** — submit answer + time taken
2. **Evaluations** — store scores and feedback (placeholder AI until Phase 7)
3. **Progress snapshots** — aggregate stats and readiness index

Then proceed to Phase 7 (AI prompt layer) per `docs/prompts/interview-preparation-app-master-implementation-prompts.md`.
