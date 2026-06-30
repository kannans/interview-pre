# Phase 10 — End-to-end loop review

**Date:** 2026-06-30  
**Scope:** profile → diagnosis → plan → practice → evaluate → progress → adjust  
**Automated checks:** `go test ./...` pass · `npm run build` pass

---

## Summary

Phases 9 and 10 are complete in code. The recovery loop is wired front-to-back: diagnosis auto-generates a plan, practice deep-links to today's task, evaluation refreshes progress snapshots, and the dashboard routes users through `practice` → `progress` → `adjust` when tasks are exhausted.

Phase 10 review fixes (backend status codes, snapshot upsert, workflow routing, React Query cache invalidation, error handling) are committed. **Manual QA against a live DB + OpenAI is still pending** before deployment.

---

## Verified working

- Precondition chain: profile → diagnosis → plan (**412** when missing)
- AI failures on diagnosis/plan/evaluate → **502**; progress degrades (`status: degraded_ai_summary`)
- Plan exhaustion: all tasks done → `GET /plan/today` **404** (`ErrNoTasksRemaining`)
- Task already complete → **PATCH** **409**
- Evaluation idempotency → **409** on re-evaluate
- Post-evaluate snapshot refresh (backend upsert + frontend `createSnapshot()`)
- Recommended question matching (`task.category` ↔ `question.type`)
- Diagnostic report parsing (JSON + base64 JSONB)
- Snapshot POST: **201** create / **200** refresh; concurrent create falls back to update
- Evaluate 409 returns existing evaluation in UI; progress page surfaces non-404 errors

---

## Manual test checklist

Run with API on `:8080`, frontend on `:3000`, Postgres seeded.

| Step | Action | Expected | Status |
|------|--------|----------|--------|
| Profile | Submit onboarding form | `POST /profile` → 201; dashboard advances | ⬜ |
| Profile missing | `GET /diagnosis` before profile | **412** precondition failed | ⬜ |
| Diagnosis | Submit all three answers | Diagnosis + plan saved; redirect to `/` | ⬜ |
| Diagnosis partial fail | Plan AI down after diagnosis succeeds | Diagnosis visible; error names plan step | ⬜ |
| Plan today | Load dashboard | Today's task + recommended question link | ⬜ |
| Practice | Submit answer on recommended question | Attempt saved; evaluation runs | ⬜ |
| Evaluate retry | Re-evaluate same attempt | **409**; existing feedback shown | ⬜ |
| Task complete | After evaluation | Task marked complete; next task on refresh | ⬜ |
| Progress refresh | Complete practice | Latest progress + history update | ⬜ |
| Snapshot upsert | Evaluate twice same day | Second snapshot POST → **200** | ⬜ |
| Plan exhausted | Complete all tasks | `GET /plan/today` 404; stage → progress/adjust | ⬜ |
| Adjust with active task | Old snapshot has recommendation but tasks remain | Stage stays `practice` | ⬜ |
| AI degraded | Break OpenAI key; create snapshot | Stats saved; degraded banner shown | ⬜ |
| Regenerate plan | From `/plan` after exhaustion | New tasks; progress stats refresh | ⬜ |

---

## Remaining gaps before deployment

| Gap | Severity | Notes |
|-----|----------|-------|
| Manual QA not run | **High** | Automated tests pass; live DB + OpenAI walkthrough still needed |
| No integration tests against real Postgres | Medium | Unique-constraint/race paths untested end-to-end |
| Diagnosis resubmit creates orphan rows | Low | `GET` returns latest; always 201 — acceptable for MVP |
| Plan status never flips to `completed` | Low | Exhaustion inferred from `GET /plan/today` 404 only |
| Tasks completable out of order | Low | `GetToday` returns lowest incomplete sequence |
| Diagnosis/plan/evaluate have no degraded fallback | Low | Only progress degrades; hard 502 elsewhere is intentional |
| Mock interview backend stubbed | Out of scope | UI shell only |
| 412 not surfaced with guided UI on all pages | Low | `isPreconditionFailed` exists; generic errors on some pages |
| No deploy artifacts | **Blocker** | Docker/compose/EC2 not started (Phase 11) |

---

## Loop flow

```mermaid
flowchart LR
  A[POST /profile] --> B[POST /diagnosis]
  B --> C[POST /plan]
  C --> D[GET /plan/today]
  D --> E[POST /attempts]
  E --> F[POST /evaluate]
  F --> G[POST /progress/snapshot]
  D -->|404 exhausted| H[progress / adjust]
  H -->|regenerate| C
```

**Adjust** = user reads `adjustment_recommendation` on the progress page and manually regenerates via `POST /plan`.

---

## Recommended follow-ups

1. Run the manual test checklist above with live DB and OpenAI key.
2. Optional: map `isPreconditionFailed` on diagnosis/plan/progress to guided next-step links.
3. Optional: one Postgres integration test for snapshot upsert + plan exhaustion.
4. Proceed to **Phase 11 — deployment** (Docker compose, env docs, smoke test).

---

## Related docs

- [STATUS.md](./STATUS.md) — phase checklist and module status
- [NEXT.md](./NEXT.md) — next implementation prompt (Phase 11)
