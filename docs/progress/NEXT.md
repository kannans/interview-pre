# Next Implementation Step

**Current focus:** Phase 7 — AI prompt layer + OpenAI integration  
**Prerequisite:** All Phase 6 backend modules are complete (profile through progress snapshots).

Use the prompt below in a new agent session. Review and test the output before moving to Phase 8 (Next.js frontend).

---

## Copy-paste prompt: AI prompt layer + OpenAI integration

```text
Build the AI prompt layer for Interview Preparation Coach.

Context:
- Personal-first interview recovery app for experienced engineers.
- Stack: Go + Chi + pgx + sqlc + PostgreSQL.
- Core loop: profile -> diagnosis -> plan -> practice -> evaluate -> progress -> adjust.
- All Phase 6 backend modules are implemented with placeholder logic (no real OpenAI calls yet).
- Progress snapshots API: POST /api/v1/progress/snapshot, GET /api/v1/progress/latest, GET /api/v1/progress/snapshots.

Existing placeholder modules to replace (Phase 7 wiring):
- Diagnosis: placeholder diagnostic_report JSON in backend/internal/service/diagnosis.go
- Plan: static task derivation in backend/internal/service/plan.go
- Evaluations: placeholder scoring in backend/internal/service/evaluation.go
- Progress: placeholder stats aggregation in backend/internal/service/progress.go (status: "pending_ai_summary")

Task:
Implement the AI prompt layer and OpenAI integration only. Do not build frontend.

Requirements:
1. Create backend/internal/ai/ with:
   - OpenAI Responses API client wrapper
   - AI service interface used by domain services
   - Config from OPENAI_API_KEY (already in backend/.env.example)

2. Create versioned prompt files (e.g. backend/internal/ai/prompts/) for:
   - baseline diagnosis
   - recovery plan generation
   - basic programming evaluation
   - DSA evaluation
   - technical answer evaluation
   - mock interview question generation (stub OK if no API yet)
   - mock interview grading (stub OK if no API yet)
   - weekly progress summary

3. For each prompt document: file name, system instruction, user input shape, required JSON output schema, example input/output.

4. Wire AI service into existing modules:
   - Diagnosis Submit -> real diagnostic_report
   - Plan Generate -> AI-driven tasks from diagnostic_report
   - Evaluation Evaluate -> real scores and feedback
   - Progress CreateSnapshot -> enrich stats with weekly AI summary; set status to "complete"

5. Operational concerns:
   - Prompt versioning strategy (filename or embedded version field)
   - Logging: prompt version, model, latency, token usage where available
   - Fallback behavior if AI call fails (graceful error or degraded placeholder — document choice)
   - Strict JSON parsing with validation

Rules:
- AI must never be called directly from the frontend.
- JSON output must be strict and parseable.
- Feedback must be practical, direct, and non-generic.
- Match existing backend conventions (repository/service/handler layers).
- Do not build Next.js frontend in this phase.

Return:
1. Files created or changed
2. AI service interface and client design
3. Prompt versioning and logging strategy
4. How to test manually (curl examples for diagnosis, plan, evaluate, progress)
5. What is left for Phase 8 (frontend)
```

---

## After AI integration

Phase 7 completes the backend intelligence layer. Proceed to **Phase 8 — Next.js frontend foundation** per `docs/prompts/interview-preparation-app-master-implementation-prompts.md`:

1. Onboarding/profile, diagnosis, plan, daily practice, mock interview, progress dashboard, settings
2. API client layer against existing `/api/v1` contracts
3. Loading and error states; responsive layout focused on today's task and weak areas
