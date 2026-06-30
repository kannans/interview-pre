# Next Implementation Step

**Current focus:** Phase 6 — Progress snapshots module  
**Prerequisite:** Profile, baseline diagnosis, recovery plan, practice questions, practice attempts, and evaluations modules are complete.

Use the prompt below in a new agent session. Review and test the output before moving to Phase 7 (AI prompt layer).

---

## Copy-paste prompt: Progress snapshots module

```text
Implement the progress snapshots backend module for Interview Preparation Coach.

Context:
- Personal-first interview recovery app for experienced engineers.
- Stack: Go + Chi + pgx + sqlc + PostgreSQL.
- Core loop: profile -> diagnosis -> plan -> practice -> evaluate -> progress -> adjust.
- Profile, baseline diagnosis, recovery plan, practice questions, practice attempts, and evaluations modules are already implemented.
- Evaluations API: POST /api/v1/attempts/{id}/evaluate, GET /api/v1/attempts/{id}/evaluation.

Existing schema (already migrated):
- progress_snapshots (id, profile_id, snapshot_date, stats JSONB, created_at)
- UNIQUE(profile_id, snapshot_date)

Existing API patterns to follow:
- backend/internal/repository/ — database access via sqlc queries
- backend/internal/service/ — business logic and validation
- backend/internal/http/handler/ — thin HTTP handlers
- backend/internal/http/router.go — wire new routes under /api/v1

Task:
Implement the progress snapshots module only. Do not implement frontend, real OpenAI calls, or weekly AI summaries yet.

Requirements:
1. Add sqlc queries in db/queries.sql for:
   - Create progress snapshot (linked to profile_id + snapshot_date)
   - Get latest progress snapshot for profile
   - Get progress snapshot by id
   - List progress snapshots for profile (optional date range or limit)

2. Regenerate sqlc output (document the command).

3. Implement repository, service, and handler layers matching existing module style.

4. API endpoints (suggest and implement):
   - POST /api/v1/progress/snapshot — compute and store today's snapshot from attempts + evaluations + plan task completion
   - GET /api/v1/progress/latest — get most recent snapshot for current profile
   - GET /api/v1/progress/snapshots — list snapshots (optional limit)

5. Placeholder stats logic (store in stats JSONB):
   - attempts_count, evaluations_count, average_scores (correctness, approach, edge_cases, communication)
   - plan_tasks_completed, plan_tasks_total
   - readiness_index (0–100 weighted placeholder from available data)
   - weak_patterns (placeholder array derived from lowest evaluation dimensions)
   - status: "pending_ai_summary" until Phase 7

6. Validation and errors:
   - Profile must exist
   - Return 404 when no snapshot exists yet (for GET latest)
   - Return 409 if snapshot already exists for today (or upsert — pick one and document)
   - All snapshots scoped to current profile

Return:
1. Files created or changed
2. New API request/response JSON examples
3. sqlc regenerate command
4. How to test manually (curl examples)
5. What is left for Phase 7 (AI prompt layer)

Rules:
- Implement one module only.
- Keep handlers thin; business logic in services.
- Match existing code conventions in backend/.
- Do not call OpenAI yet — use placeholder aggregation logic.
- Do not build frontend.
```

---

## After progress snapshots

Phase 6 backend modules are complete. Proceed to **Phase 7 — AI prompt layer + OpenAI integration** per `docs/prompts/interview-preparation-app-master-implementation-prompts.md`:

1. Versioned prompts for diagnosis, plan generation, attempt evaluation, and weekly summaries
2. `backend/internal/ai/` service with OpenAI client wrapper
3. Replace placeholder logic in diagnosis, plan, evaluations, and progress modules
