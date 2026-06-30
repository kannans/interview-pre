# Next Implementation Step

**Current focus:** Phase 8 — Next.js frontend foundation  
**Prerequisite:** Phase 7 AI integration complete; backend endpoint + service tests pass (`cd backend && go test ./...`).

Use the prompt below in a new agent session. Review and test the output before moving to Phase 9 (end-to-end recovery loop).

---

## Copy-paste prompt: Next.js frontend foundation

```text
Implement the Next.js frontend foundation for Interview Preparation Coach.

Stack:
- Next.js
- TypeScript
- Tailwind CSS
- Backend API: Go REST API at http://localhost:8080/api/v1

Build pages for:
1. onboarding/profile setup
2. baseline diagnosis
3. recovery plan
4. daily practice
5. mock interview practice (UI shell OK — mock interview AI stubs remain backend-only)
6. progress dashboard
7. settings/basic config

Return:
1. frontend file structure
2. page components
3. shared UI components
4. API client layer
5. state management approach
6. loading and error states
7. responsive layout notes

Design requirements:
- App UI, not a marketing landing page
- User quickly sees today's task, progress, weak areas, and next action
- Handle 502 from AI endpoints gracefully (retry messaging)

Rules:
- Connect to backend API contracts (no direct OpenAI calls from frontend)
- Do not build payment or job portal UI
```

---

## After frontend foundation

Phase 8 completes the user-facing shell. Proceed to **Phase 9 — End-to-end recovery workflow** per `docs/prompts/interview-preparation-app-master-implementation-prompts.md`:

1. Wire full loop: profile → diagnosis → plan → practice → evaluate → progress
2. Mock interview flow once backend mock endpoints exist
3. Manual QA checklist across all pages

## Manual AI testing (Phase 7)

With `OPENAI_API_KEY` set and backend running on `:8080`:

```bash
# 1. Profile
curl -s -X POST http://localhost:8080/api/v1/profile \
  -H 'Content-Type: application/json' \
  -d '{"experience_years":8,"target_role":"Senior Backend Engineer","target_companies":"Google","weak_areas_self":"DSA, system design","timeline":"60 days"}'

# 2. Diagnosis (AI)
curl -s -X POST http://localhost:8080/api/v1/diagnosis \
  -H 'Content-Type: application/json' \
  -d '{"coding_answer":"Two-sum with hash map O(n)...","conceptual_answer":"Load balancer + cache...","behavioral_answer":"STAR story on production incident..."}'

# 3. Plan (AI)
curl -s -X POST http://localhost:8080/api/v1/plan

# 4. Practice + evaluate (AI)
curl -s -X POST http://localhost:8080/api/v1/attempts \
  -H 'Content-Type: application/json' \
  -d '{"question_id":"a1000001-0000-4000-8000-000000000004","submitted_content":"function findMax(nums){...}","time_taken_seconds":300}'
# Use returned attempt id:
curl -s -X POST http://localhost:8080/api/v1/attempts/{ATTEMPT_ID}/evaluate

# 5. Progress snapshot (AI summary; degrades if AI fails)
curl -s -X POST http://localhost:8080/api/v1/progress/snapshot
curl -s http://localhost:8080/api/v1/progress/latest
```

See `backend/internal/ai/README.md` for prompt versioning, logging, and fallback behavior.
