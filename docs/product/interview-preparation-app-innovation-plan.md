# Interview Preparation Coach — Innovation Plan

## Document Purpose

This document defines the **next product direction** after the core recovery loop MVP (Phases 1–10).

It answers one question:

> How do we become **the system that turns interview rust into measurable recovery before your next real round** — and not another prep app with AI feedback?

This is not a feature wish list. It is a **quality and differentiation plan** grounded in the personal-first mission: help one experienced engineer recover interview performance with visible, trustworthy progress.

## Product Motto

**Turn interview rust into measurable recovery before your next real round.**

Every feature in this plan must support at least one of:

- **Rust** — the user struggles under interview-style pressure despite real-world experience
- **Recovery** — structured, consistent improvement in weak areas
- **Measurable** — the user can see proof they are getting better, not just feel busy
- **Before the real round** — preparation is timed, targeted, and confidence-building for an actual upcoming interview

## What We Are Not Building

To protect product quality and focus, these are explicitly out of scope for this innovation phase:

- job portal or vacancy search
- generic LeetCode-style problem grinding without coaching context
- AI chat for its own sake (no “ask anything” copilot)
- broad multi-persona expansion before personal validation succeeds
- billing, auth, or community features

If a feature does not strengthen the recovery loop or make progress more measurable, it does not ship.

## Current Foundation (What We Already Have)

The MVP loop is complete in code:

```
profile → diagnosis → plan → practice → evaluate → progress → adjust
```

| Capability | Status |
|------------|--------|
| Baseline diagnosis (coding, concepts, behavioral) | Done |
| AI recovery plan with daily tasks | Done |
| Practice + code run + attempt logging | Done |
| Multi-axis AI evaluation | Done |
| Progress snapshots + readiness index | Done |
| Workflow routing (dashboard next action) | Done |
| Mock interview | UI shell only |
| Question bank | ~15 seeded questions |
| Deployment + 30-day personal validation | Not started |

**Gap:** The architecture is a coach skeleton. The **experience** does not yet feel unmistakably different from “AI + practice problems.” This plan closes that gap.

## Strategic Position

### The problem we own

Experienced engineers fail interviews not because they lack job skills, but because **interview performance is a different skill** — coding under time pressure, structured explanation, edge-case discipline, and round-specific expectations.

### Our wedge

| Generic prep apps | Interview Preparation Coach |
|-------------------|----------------------------|
| Problem banks | Prescribed recovery path |
| Score on one attempt | Trend across weak patterns |
| Practice when motivated | Daily task + next action |
| “Here’s feedback” | “Here’s what changed in your plan and why” |
| Readiness as a feeling | Readiness as a measured trajectory |

### The moat we are building

Longitudinal **failure pattern memory** tied to a closed recovery loop. Prompts are copyable; a user’s attempt history, weak-pattern trends, and plan adjustments over weeks are not.

---

## Innovation Pillars

Four pillars turn the motto into product reality:

### 1. Pressure Realism

Interview rust shows up under pressure, not in calm study sessions.

### 2. Pattern Memory

Recovery requires knowing *what* keeps going wrong, not just *that* something went wrong once.

### 3. Proactive Coaching

The system adjusts the plan; the user does not hunt for what to do next.

### 4. Outcome Closure

The loop extends into real interviews (debrief) and pre-interview rituals (briefing).

---

## Feature Initiatives

### Initiative A — Failure Pattern Memory

**Problem:** Evaluations give per-attempt feedback, but the user cannot see recurring mistakes across sessions.

**Solution:** Aggregate evaluation signals into persistent weak patterns and surface them on the dashboard, progress page, and plan regeneration.

**Examples of patterns:**

- skips edge cases on array problems
- correct but brute-force; rarely optimizes
- strong code, weak verbal explanation
- rushes to implementation without clarifying constraints

**User-visible outputs:**

- “Your top 3 interview anti-patterns (last 20 attempts)”
- pattern trend: improving / stable / regressing
- pattern-linked practice recommendations

**Data we already collect:**

- `score_correctness`, `score_approach`, `score_edge_cases`, `score_communication`
- per-attempt feedback text by dimension
- question category and difficulty

**New work (high level):**

- backend aggregation job or snapshot enrichment
- `weak_patterns` expansion beyond current progress stats
- dashboard panel: “Recurring mistakes”
- feed patterns into recovery plan prompt on regenerate

**Success criteria:**

- user can name their top recurring weakness without reading raw feedback
- plan regeneration explicitly references at least one detected pattern

**Priority:** P0 — highest leverage, uses existing data

---

### Initiative B — Recovery Velocity

**Problem:** `readiness_index` is a point-in-time number. The user’s real question is: **“Am I getting better?”**

**Solution:** Add velocity and regression metrics alongside readiness.

**Metrics:**

| Metric | Meaning |
|--------|---------|
| Recovery velocity | Score improvement rate per category over 7 / 14 / 30 days |
| Regression risk | Days since last practice in a weak category; alert when stale |
| Calibration gap | Self-rated confidence before submit vs AI score (over/under-confidence) |
| Consistency streak | Practice days per week vs plan target |

**User-visible outputs:**

- progress page: “Velocity” card next to readiness
- dashboard: “At risk of rust” nudge when a weak area is neglected
- weekly snapshot narrative includes velocity, not just averages

**Success criteria:**

- user can answer “am I improving?” with a trend, not a single percentage
- regression nudge triggers after configurable idle days on a focus area

**Priority:** P0 — core to “measurable recovery”

---

### Initiative C — Interview Pressure Mode

**Problem:** Timer alone does not simulate interview conditions. Rust appears with interruption, explanation demands, and time anxiety.

**Solution:** A dedicated practice mode that simulates a real coding round.

**Session rules:**

- fixed duration (e.g. 25–35 minutes)
- phases: clarify → approach → code → review
- optional “interviewer” prompts: “Can you explain before coding?” / “What about edge cases?”
- scoring weights communication and process, not only final code
- no hints; degraded score if user skips approach phase

**UX:**

- entry from dashboard: “Pressure round” alongside normal practice
- post-session debrief aligned with evaluation dimensions
- link session to today’s plan task when applicable

**Success criteria:**

- session feels meaningfully different from untimed practice
- user reports higher confidence transfer to real interviews (tracked in Phase 12 validation)

**Priority:** P1 — strongest experiential differentiator

---

### Initiative D — Proactive Plan Adjust

**Problem:** Today, adjust = read `adjustment_recommendation` on progress and manually regenerate plan. The coach feels passive.

**Solution:** System proposes plan changes with a clear diff and one-click accept.

**Flow:**

1. After N evaluations or weekly snapshot, detect shift in weak patterns
2. Generate proposed plan diff: tasks added, removed, reprioritized
3. Show rationale: “Last 5 evaluations: behavioral strong, coding flat → shift 70% to coding drills”
4. User accepts or dismisses; accepted plan replaces active plan

**Success criteria:**

- user spends less time deciding what to study next
- every plan change has an explicit “why” tied to data

**Priority:** P1 — completes the adjust step in the loop

---

### Initiative E — Mock Interview (Complete the Stub)

**Problem:** Mock interview is in the roadmap and pitch, but backend is stubbed. The product story is incomplete.

**Solution:** Ship text-based mock interview rounds using existing prompts (`mock_interview_questions`, `mock_interview_grading`).

**MVP scope:**

- start session: round type (coding / DSA / behavioral)
- 3–5 questions with follow-ups
- timed or untimed mode
- grading rubric aligned with evaluation dimensions
- session summary feeds progress snapshot and pattern memory

**Out of scope for first ship:**

- voice / video
- live human interviewer

**Success criteria:**

- full mock session completes without manual workaround
- session results appear in progress history

**Priority:** P1 — already promised in product vision

---

### Initiative F — Pre-Interview Briefing

**Problem:** Peak value moment is 1–3 days before a real interview. Generic daily practice does not optimize for that.

**Solution:** User enters upcoming interview details; system generates a focused briefing.

**Inputs:**

- company name (optional)
- round type (phone screen, coding, system design, behavioral)
- interview date

**Outputs:**

- 3 warm-up questions (easy → medium) matched to weak patterns
- recap of top 2 recurring mistakes to watch for
- 10-minute behavioral story refresh from past strong answers
- short confidence script based on recent wins

**Success criteria:**

- briefing is generated in under one minute
- user completes briefing in under 30 minutes

**Priority:** P2 — high emotional value, builds trust

---

### Initiative G — Real Interview Debrief Loop

**Problem:** The app only sees practice. Real interview outcomes are the ultimate signal.

**Solution:** Post-interview log that feeds back into diagnosis and plan.

**Inputs:**

- company / round type
- what was asked (free text)
- self-rated performance per dimension
- outcome (advanced / rejected / waiting)
- optional notes

**Outputs:**

- debrief stored in history
- weak areas updated when real round exposes new gaps
- optional: trigger plan refresh or re-diagnosis after rejection

**Success criteria:**

- user logs at least one real debrief during 30-day validation
- next plan cycle references debrief data when present

**Priority:** P2 — closes loop outside the app

---

### Initiative H — Micro-Drills for Busy Seniors

**Problem:** 15+ YOE candidates will not grind 2-hour LeetCode sessions daily. Inconsistent practice breaks recovery.

**Solution:** 5–15 minute drill types that still count toward plan progress.

**Drill types:**

- fix this buggy function
- what is wrong with this approach?
- explain this solution in 60 seconds (text)
- spot the missing edge case

**Success criteria:**

- drill completable on mobile or in a short break
- drills tagged to weak patterns and plan tasks

**Priority:** P2 — supports consistency

---

### Initiative I — Company Readiness Lens (Later)

**Problem:** Users care about target companies, but readiness is generic today.

**Solution:** Readiness view filtered by company or company type (startup vs big tech vs product mid-size).

**Depends on:** larger question bank, pattern memory, enough attempt history.

**Priority:** P3 — after personal validation

---

## Delivery Phases

Phases align with [build status](../progress/STATUS.md). Innovation work starts after deployment and daily dogfooding are unblocked.

### Phase 11 — Deploy and Stabilize (prerequisite)

- Docker / compose / EC2
- manual QA checklist complete
- daily personal use begins

**Gate:** app runs reliably for daily practice without friction.

### Phase 12 — 30-Day Personal Validation (prerequisite)

- use the app as primary prep tool for 30 days
- log attempts, snapshots, and at least one real interview or mock outcome
- document what helped vs what felt like noise

**Gate:** honest answer to “did measurable recovery happen?”

### Phase 13 — Measurable Recovery (P0 initiatives)

| Deliverable | Initiative |
|-------------|------------|
| Failure pattern aggregation + dashboard | A |
| Recovery velocity + regression alerts | B |
| Expanded question bank (topic coverage for patterns) | content |

**Theme:** Make progress impossible to fake.

### Phase 14 — Pressure and Coaching (P1 initiatives)

| Deliverable | Initiative |
|-------------|------------|
| Interview Pressure Mode | C |
| Proactive plan adjust with diff | D |
| Mock interview backend + UI | E |

**Theme:** Feel like a coach, not a problem site.

### Phase 15 — Real-World Closure (P2 initiatives)

| Deliverable | Initiative |
|-------------|------------|
| Pre-interview briefing | F |
| Real interview debrief | G |
| Micro-drills | H |

**Theme:** Bridge practice and real rounds.

### Phase 16 — Strategic Expansion (P3)

| Deliverable | Initiative |
|-------------|------------|
| Company readiness lens | I |
| Multi-user / auth (only if personal validation succeeds) | platform |

---

## Quality Bar

A feature ships only when it passes these checks:

### 1. Recovery test

Does it help the user improve a specific weak area, not just consume content?

### 2. Measurability test

Can the user point to a metric or trend that changed because of this feature?

### 3. Coach test

Does the system tell the user what to do next, or does it add another page to explore?

### 4. Anti-bloat test

Can we ship a smaller version in one week of focused work? If not, split the initiative.

### 5. Personal validation test

Would the primary user (15 YOE, interview rust, job switch) use this weekly without reminders?

---

## Success Metrics (Personal-First)

Before scaling to other users, track these during Phase 12–15:

| Metric | Target direction |
|--------|------------------|
| Practice days per week | ≥ 4 |
| Average evaluation scores by dimension | upward trend over 30 days |
| Top weak pattern severity | downward trend |
| Recovery velocity | positive in at least 2 categories |
| Time to “what should I do today?” | < 30 seconds from dashboard |
| Real interview self-rated performance | improved vs baseline diagnosis |
| Feature abandonment | pressure mode and drills used more than once |

---

## What Stays the Same

These principles from the MVP remain non-negotiable:

- personal-first until validation succeeds
- closed loop: profile → diagnosis → plan → practice → evaluate → progress → adjust
- AI behind Go backend with versioned prompts and schema validation
- store every attempt, score, and mistake signal
- no job portal complexity

---

## Related Docs

- [Personal-First MVP](./interview-preparation-app-personal-first-mvp.md) — original scope and user definition
- [PRD](./interview-preparation-app-prd.md) — full product requirements
- [Roadmap](./interview-preparation-app-roadmap.md) — long-term feature phases
- [Master Build Plan](./interview-preparation-app-master-build-plan.md) — implementation phases 1–11
- [Build Status](../progress/STATUS.md) — current completion state
- [Loop Review](../progress/REVIEW.md) — end-to-end QA checklist
- [Implementation Prompts](../prompts/interview-preparation-app-master-implementation-prompts.md) — copy-paste prompts for Phases 11–16

---

## Summary

The MVP built the **recovery loop**. The innovation plan builds the **recovery system**:

1. **Remember** what keeps going wrong (Failure Pattern Memory)
2. **Prove** improvement is happening (Recovery Velocity)
3. **Simulate** real interview pressure (Pressure Mode)
4. **Act** on data, not passive text (Proactive Adjust)
5. **Close the loop** with real rounds (Briefing + Debrief)

That is how we earn the motto: **turn interview rust into measurable recovery before your next real round.**
