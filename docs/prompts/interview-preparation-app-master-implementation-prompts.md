# Interview Preparation Coach Master Implementation Prompts

## Purpose

This file contains detailed copy-paste prompts for building the Interview Preparation Coach app from planning to deployment.

Use these prompts one by one. Each prompt is intentionally detailed so the model has enough context to produce useful implementation work.

## Phase Index

| Phase | Name | Status (as of 2026-06-30) |
|-------|------|---------------------------|
| 1–10 | MVP foundation through loop review | Done |
| 11 | Docker / EC2 deployment | Not started |
| 12 | 30-day personal validation | Not started |
| 13 | Measurable recovery (patterns + velocity) | Not started |
| 14 | Pressure mode + proactive adjust + mock interview | Not started |
| 15 | Briefing + debrief + micro-drills | Not started |
| 16 | Company readiness lens + platform (conditional) | Not started |

Phases 13–16 implement [Innovation Plan](../product/interview-preparation-app-innovation-plan.md). Do not skip Phase 12 validation before Phase 13.

## Core Product Context

The product is a personal-first interview recovery app for experienced engineers.

The first user is an experienced engineer with around 15 years of experience who is strong in real work but struggles with basic coding, DSA, and interview-style answering.

The first goal is not revenue. The first goal is to help the user prepare better, practice consistently, track progress, rebuild coding confidence, and perform better in real interviews.

The core product loop is:

`profile -> diagnosis -> plan -> practice -> evaluate -> progress -> adjust`

## Stack Decisions

- Frontend: `Next.js`
- Language: `TypeScript`
- Styling: `Tailwind CSS`
- Backend: `Go`
- Router: `Chi`
- Database: `PostgreSQL`
- DB driver: `pgx`
- Query generation: `sqlc`
- Migrations: `golang-migrate`
- AI: `OpenAI Responses API` called only from the Go backend
- Deployment target: `EC2`

## Important Rules

- Build this as a personal-first MVP.
- Do not build job portal features in v1.
- Do not build vacancy search in v1.
- Do not build public community features in v1.
- Do not build voice interview features in v1.
- Do not add CV/JD analysis until the core recovery loop works.
- Keep every implementation step small and reviewable.
- Store attempts, scores, feedback, and progress history from the beginning.

## Phase 1: Freeze The MVP Scope

```text
You are helping me build an app called Interview Preparation Coach.

Context:
- This is a personal-first interview recovery app for experienced engineers.
- The first user has around 15 years of experience.
- The user is strong in real work but struggles with basic coding, DSA, and interview-style answers.
- The first goal is personal improvement, not monetization.
- The app must help the user prepare better, practice daily, track progress, and improve interview readiness.

Stack:
- Frontend: Next.js + TypeScript + Tailwind CSS
- Backend: Go + Chi
- Database: PostgreSQL
- AI: OpenAI Responses API through the backend

Task:
Freeze the MVP scope before any coding starts.

Return:
1. MVP goal in one paragraph
2. Primary user profile
3. Must-have features for v1
4. Explicitly excluded features for v1
5. First user journey
6. Success metrics for 30 days
7. First build slice
8. Assumptions

Rules:
- Do not write code.
- Do not include job portal features.
- Do not include CV/JD analysis in v1.
- Keep the scope personal-first and small enough to build.
```

## Phase 2: Technical Architecture

```text
You are the senior software architect for Interview Preparation Coach.

Use the frozen MVP scope from Phase 1.

Design the technical architecture for a personal-first interview recovery app.

Stack:
- Next.js + TypeScript + Tailwind CSS frontend
- Go + Chi backend
- PostgreSQL database
- pgx for database access
- sqlc for typed SQL
- golang-migrate for migrations
- OpenAI Responses API through the Go backend only

The core app loop is:
profile -> diagnosis -> plan -> practice -> evaluate -> progress -> adjust

Return:
1. High-level architecture
2. Mermaid architecture diagram
3. Main frontend pages
4. Main backend modules
5. Main database entities
6. AI service boundaries
7. API groups
8. Data flow for the core loop
9. Security and privacy notes
10. What should be deferred until later

Rules:
- Do not write implementation code.
- Keep the architecture practical for a solo developer on macOS.
- Assume deployment later to EC2.
```

## Phase 3: Project Folder Structure

```text
Create the project folder structure for Interview Preparation Coach.

Context:
- Frontend: Next.js + TypeScript + Tailwind CSS
- Backend: Go + Chi
- Database: PostgreSQL
- Migrations: golang-migrate
- Queries: sqlc
- AI prompts should be versioned and stored cleanly
- Deployment target: EC2 using Docker later

Return:
1. Full file tree
2. Short explanation of each top-level folder
3. Backend folder responsibilities
4. Frontend folder responsibilities
5. Where prompts should live
6. Where database migrations should live
7. Where deployment files should live

Rules:
- Do not write full code yet.
- Keep the folder structure clean but not overengineered.
- Make it suitable for a real build from scratch.
```

## Phase 4: PostgreSQL Schema And Migrations

```text
Design the PostgreSQL database schema for the personal-first MVP.

The app must support:
- preparation profile
- baseline diagnosis
- recovery plan
- practice questions
- practice attempts
- AI evaluations
- progress snapshots
- mistake tracking

Return:
1. Table list with purpose
2. SQL migration for initial schema
3. Primary keys and foreign keys
4. Useful indexes
5. Notes for sqlc query generation
6. Example seed data for a few starter questions

Rules:
- Keep the schema practical for v1.
- Store every attempt and evaluation so progress can be tracked over time.
- Do not design job portal tables.
- Do not design billing tables yet.
```

## Phase 5: Go Backend Foundation

```text
Implement the Go backend foundation for Interview Preparation Coach.

Use:
- Go
- Chi router
- pgx connection pool
- environment-based config
- slog or simple structured logging
- health endpoint

Required structure:
- cmd/api main entrypoint
- internal/config
- internal/http
- internal/db
- internal/service
- internal/repository
- internal/ai

Build only the foundation in this phase.

Return:
1. Files to create
2. Code for each file
3. Environment variables
4. Health endpoint behavior
5. How to run locally
6. How to verify it works

Rules:
- Do not implement product modules yet.
- Do not call OpenAI yet.
- Do not add authentication yet unless necessary.
- Keep the code small and clear.
```

## Phase 6: Backend Product Modules

```text
Implement backend product modules for the personal-first MVP.

Existing backend foundation:
- Go
- Chi
- pgx
- PostgreSQL
- sqlc-ready schema

Implement modules in this order:
1. profile
2. baseline diagnosis
3. recovery plan
4. practice questions
5. practice attempts
6. evaluations
7. progress snapshots

For each module return:
1. API endpoints
2. request and response JSON
3. service logic
4. repository queries
5. validation rules
6. error cases
7. tests or test plan

Rules:
- Implement one module at a time.
- Keep handlers thin.
- Put business logic in services.
- Put database access in repositories or sqlc-generated query wrappers.
- Do not implement frontend in this phase.
```

## Phase 7: AI Prompt Layer And OpenAI Integration

```text
Build the AI prompt layer for Interview Preparation Coach.

The backend must call OpenAI Responses API through a dedicated AI service.

Create versioned prompts for:
1. baseline diagnosis
2. recovery plan generation
3. basic programming evaluation
4. DSA evaluation
5. technical answer evaluation
6. mock interview question generation
7. mock interview grading
8. weekly progress summary

For each prompt return:
1. prompt file name
2. system instruction
3. user input shape
4. required JSON output schema
5. example input
6. example output

Also return:
1. Go AI service interface
2. OpenAI client wrapper design
3. prompt versioning strategy
4. logging strategy for prompt version, model, and evaluation result
5. fallback behavior if AI call fails

Rules:
- AI must never be called directly from the frontend.
- JSON output should be strict and parseable.
- Feedback must be practical, direct, and non-generic.
```

## Phase 8: Next.js Frontend Foundation

```text
Implement the Next.js frontend foundation for Interview Preparation Coach.

Stack:
- Next.js
- TypeScript
- Tailwind CSS
- Backend API: Go REST API

Build pages for:
1. onboarding/profile setup
2. baseline diagnosis
3. recovery plan
4. daily practice
5. mock interview practice
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
- This is an app UI, not a marketing landing page.
- Keep it clean, focused, and easy to scan.
- The user should quickly see today's task, progress, weak areas, and next action.
- Avoid decorative clutter.

Rules:
- Connect to backend API contracts.
- Use mock API data only if backend is not ready.
- Do not build payment or job portal UI.
```

## Phase 9: End-To-End Recovery Workflow

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

## Phase 10: Testing And Review

```text
Review and test the current implementation of Interview Preparation Coach.

Focus on:
- backend API correctness
- database persistence
- AI output parsing
- frontend API integration
- progress calculations
- user-facing workflow gaps
- data integrity
- error handling

Return findings first.

For each finding include:
1. severity
2. affected file or module
3. issue
4. why it matters
5. recommended fix

Then return:
1. missing tests
2. manual test checklist
3. smoke test checklist
4. deployment readiness status

Rules:
- Be strict.
- Do not refactor unrelated parts.
- Prioritize bugs and workflow blockers.
```

## Phase 11: EC2 / Docker Deployment

**Prerequisite:** Phase 10 complete — `go test ./...` and `npm run build` pass. Run or acknowledge the manual QA checklist in `docs/progress/REVIEW.md`.

**Reference:** `runserver.sh` for local dev; `backend/.env` for API config; `docs/progress/STATUS.md` for module inventory.

```text
You are deploying Interview Preparation Coach (personal-first MVP) to a reproducible Docker stack and documenting EC2 deployment.

Context:
- Core loop is complete: profile -> diagnosis -> plan -> practice -> evaluate -> progress -> adjust
- Backend: Go + Chi at backend/cmd/api
- Frontend: Next.js at frontend/
- Database: PostgreSQL with golang-migrate migrations in db/migrations/
- Migrations should run on API startup or via a documented init step
- Single-user MVP — no auth layer
- OpenAI key only in backend environment

Task:
Add Docker and deployment artifacts. Do not add new product features.

Deliver:
1. Dockerfile for Go API (multi-stage build, production binary)
2. Dockerfile for Next.js frontend (production build)
3. docker-compose.yml — Postgres + API + frontend (local and single-server)
4. .env.example files for backend and frontend (no secrets committed)
5. README section: deploy steps, env vars, smoke test
6. Optional: Nginx reverse proxy config for EC2 single-host deploy
7. Optional: notes for RDS vs containerized Postgres

Implementation requirements:
- Frontend API base URL configurable at build/runtime
- CORS allows frontend origin in compose
- Health check: GET /health
- Document how to seed DB / run migrations
- Document ports: API :8080, frontend :3000, Postgres :5432

Return:
1. Files added or changed (list with purpose)
2. Commands to run locally with docker compose
3. EC2 setup steps (security groups, env, restart strategy)
4. Smoke-test checklist for deployed stack (profile -> diagnosis -> practice -> evaluate)
5. What differs when using RDS instead of local Postgres

Rules:
- Do not expose secrets in code or committed env files
- Keep it practical for a solo developer on macOS deploying to one EC2 instance
- Do not add auth, billing, or innovation features (Phases 13+)
- Match existing repo conventions (sqlc, migrate, Chi routes)
```

## Phase 12: 30-Day Personal Validation

**Prerequisite:** Phase 11 complete — app runs reliably via Docker or `./runserver.sh` for daily use.

**Reference:** Product motto in `docs/product/interview-preparation-app-innovation-plan.md` — turn interview rust into measurable recovery before your next real round.

```text
You are creating a 30-day personal validation plan for Interview Preparation Coach.

Context:
- Primary user: ~15 YOE engineer, strong at work, weak under interview pressure
- First success = personal improvement, not revenue
- MVP loop exists but innovation features (Phases 13–16) are not built yet
- Goal: decide honestly whether the app helps before building more features

Task:
Produce a practical 30-day validation playbook the user can follow daily. Include templates and decision gates.

Return:

1. Pre-validation setup (Day 0)
   - Run full manual QA from docs/progress/REVIEW.md
   - Baseline metrics to record: diagnosis scores, readiness_index, self-rated confidence (1–10)
   - Friction log file or doc template location

2. Daily usage routine (Mon–Sun)
   - Morning: dashboard -> today's task -> practice (target: 30–60 min)
   - What to do when plan/today returns 404
   - Minimum viable day when busy (15 min)

3. Weekly review checklist (4 weeks)
   - Week 1: loop reliability and friction
   - Week 2: AI feedback quality and plan relevance
   - Week 3: progress measurability — can you see improvement?
   - Week 4: real interview readiness — mock or actual round

4. Session log template (copy-paste per practice session)
   - Date, question, time taken, scores, one thing that went wrong, one win
   - Self-rated confidence before submit (1–10)
   - Would you use this feedback in a real interview?

5. Success metrics (quantitative)
   - Practice days per week (target >= 4)
   - Average evaluation scores by dimension (trend)
   - readiness_index delta week 1 vs week 4
   - Plan tasks completed
   - Time from open app to start practice (< 2 min)

6. Success metrics (qualitative)
   - "Am I improving?" — honest yes/no with evidence
   - "Does the coach tell me what to do next?" — yes/no
   - "Would I recommend this to another senior engineer?" — yes/no/maybe

7. Decision gates at Day 30
   - GO: measurable recovery + daily habit -> proceed to Phase 13
   - PIVOT: loop works but feedback weak -> fix AI/prompts before Phase 13
   - STOP: no habit or no perceived value -> document why before more build

8. Friction and improvement backlog template
   - Columns: date, area, severity, keep/fix/cut, notes

9. Optional: docs/progress/VALIDATION.md file structure if creating a repo doc

Rules:
- Focus on personal interview recovery, not monetization
- Do not design Phase 13+ features yet — only log them as ideas if they emerge from friction
- Be strict: validation failure is valuable data
```

## Phase 13: Measurable Recovery (Failure Patterns + Velocity)

**Prerequisite:** Phase 12 gate passed — honest "yes" to measurable improvement OR clear list of blockers fixed.

**Reference:** Innovation Plan Initiative A (Failure Pattern Memory) and Initiative B (Recovery Velocity) in `docs/product/interview-preparation-app-innovation-plan.md`.

```text
You are implementing Phase 13 of Interview Preparation Coach: Measurable Recovery.

Product motto:
Turn interview rust into measurable recovery before your next real round.

Context — existing codebase:
- Evaluations store per-attempt scores: correctness, approach, edge_cases, communication (evaluations table)
- Progress snapshots store stats JSONB including readiness_index and weak_patterns (progress_snapshots)
- deriveWeakPatterns() in backend/internal/service/progress.go only ranks lowest average dimensions — not true recurring failure memory
- Dashboard shows weak areas from diagnosis or progress (frontend/src/app/(app)/page.tsx)
- ~15 seeded questions — expand for pattern coverage

Task:
Implement Initiative A (Failure Pattern Memory) and Initiative B (Recovery Velocity) as the smallest shippable version.

Part 1 — Backend: Failure Pattern Memory
1. Design pattern model — either:
   - enrich SnapshotStats JSONB with structured failure_patterns[], OR
   - new table failure_patterns / attempt_pattern_tags with migration
2. Aggregate from evaluations + feedback text + question category:
   - recurring themes (e.g. skips edge cases, brute-force, weak explanation)
   - trend per pattern: improving | stable | regressing
   - attempt count and last_seen_at per pattern
3. Computation: on snapshot create/refresh OR dedicated GET /api/v1/progress/patterns
4. Optional: versioned AI prompt to extract/classify patterns from recent feedback (fallback to rule-based if AI fails)
5. Feed top patterns into recovery plan regeneration prompt input

Part 2 — Backend: Recovery Velocity
1. Extend SnapshotStats or add GET /api/v1/progress/velocity:
   - score delta per dimension over 7 / 14 / 30 days (from snapshots or raw evaluations)
   - recovery_velocity: overall trend label (positive | flat | negative)
   - regression_risks: categories not practiced in N days (configurable, default 7)
   - optional: calibration_gap if self-rated confidence is stored on attempts
2. Add self_rated_confidence (1–10) to attempt submit if not present — migration + API + frontend field

Part 3 — Content
1. Add seed migration or SQL file with 20–30 more questions across arrays, strings, hash maps, loops — level_1 and level_2
2. Tag questions by category for regression detection

Part 4 — Frontend
1. Dashboard: "Recurring mistakes" panel (top 3 patterns with trend badges)
2. Progress page: Recovery Velocity card next to readiness index
3. Regression alert banner when a focus category is stale
4. Practice submit: optional confidence slider (1–10) before evaluate

Part 5 — Tests
1. Unit tests for pattern aggregation and velocity calculation
2. Handler tests for new endpoints
3. Manual test checklist

Return:
1. Schema/migration changes
2. API endpoints (method, path, request/response JSON)
3. Files changed (backend + frontend)
4. How patterns differ from current weak_patterns
5. Manual test checklist
6. Remaining gaps for a fuller Phase 13

Rules:
- Personal-first MVP; no auth
- AI behind Go backend only; degraded fallback if pattern AI fails
- Do not build Pressure Mode, mock interview, or debrief yet (Phase 14–15)
- Pass innovation quality bar: Recovery test, Measurability test, Coach test
- Keep diff reviewable — prefer enriching snapshots before many new tables unless necessary
```

## Phase 14: Pressure, Proactive Adjust, Mock Interview

**Prerequisite:** Phase 13 complete — failure patterns and velocity visible on dashboard/progress.

**Reference:** Innovation Plan Initiatives C, D, E in `docs/product/interview-preparation-app-innovation-plan.md`.

```text
You are implementing Phase 14 of Interview Preparation Coach: Pressure and Coaching.

Theme: Feel like a coach, not a problem site.

Context — existing codebase:
- Practice page: PracticeWorkspace, timer, run code, submit, evaluate (frontend/src/app/(app)/practice/[id]/page.tsx)
- Plan: POST /api/v1/plan regenerates; adjustment_recommendation on progress is read-only today
- Mock interview: UI shell at /mock-interview; AI stubs return ErrNotImplemented (backend/internal/ai/service.go)
- Prompts exist: mock_interview_questions.v1.json, mock_interview_grading.v1.json

Task:
Implement three initiatives in order: C (Pressure Mode), D (Proactive Adjust), E (Mock Interview).

Part 1 — Interview Pressure Mode (Initiative C)
1. New session type: pressure_session (table or attempts metadata flag)
2. Session flow:
   - User picks duration (25 or 35 min) and optional link to today's plan task
   - Phases: clarify (notes) -> approach (required text before code unlocks) -> code -> review
   - Timer enforces hard stop; submit auto-triggers on expiry optional
3. Scoring: weight communication and approach higher than normal practice evaluate
4. Optional: static "interviewer prompts" at phase boundaries (no live AI mid-session for v1)
5. Frontend: entry from dashboard "Pressure round" + distinct post-session debrief UI
6. API: POST start session, PATCH phase, POST complete -> evaluate with pressure rubric flag

Part 2 — Proactive Plan Adjust (Initiative D)
1. Detect when to propose adjust: after N evaluations (e.g. 5) or on snapshot refresh when patterns shift
2. POST /api/v1/plan/propose — returns proposed plan diff:
   - tasks_added, tasks_removed, tasks_reprioritized
   - rationale tied to failure_patterns and velocity data
3. POST /api/v1/plan/accept — applies proposed plan (replaces active plan)
4. Versioned AI prompt: plan_adjustment.v1.json using patterns + velocity + current plan
5. Frontend on /progress and dashboard:
   - "Suggested plan changes" card with diff view
   - Accept / Dismiss actions
   - Explain why each task changed

Part 3 — Mock Interview (Initiative E)
1. Wire GenerateMockInterviewQuestions and GradeMockInterview in ai/service.go (remove stubs)
2. Tables: mock_interview_sessions, mock_interview_turns (or JSONB transcript on session)
3. API:
   - POST /api/v1/mock-interview/sessions — start (round_type: coding | dsa | behavioral)
   - GET /api/v1/mock-interview/sessions/{id}
   - POST /api/v1/mock-interview/sessions/{id}/answer — append turn
   - POST /api/v1/mock-interview/sessions/{id}/complete — grade session
4. Session: 3–5 questions with follow-ups from AI; store transcript
5. Grading aligned with evaluation dimensions; results feed progress snapshot + pattern memory
6. Frontend: replace mock-interview shell with working flow (start -> Q&A -> summary)

Part 4 — Tests and docs
1. Backend tests for session state machine, plan propose/accept, mock grading parse
2. Update docs/progress/STATUS.md module table
3. Manual test checklist for all three features

Return:
1. Database migrations
2. API surface (all new endpoints)
3. New prompt files if any
4. Frontend routes and components
5. How pressure mode differs from normal practice in scoring
6. Manual test checklist

Rules:
- No voice/video mock interview in v1
- AI only from Go backend; JSON schema validation on all new prompts
- Do not build briefing, debrief, or micro-drills yet (Phase 15)
- Proactive adjust must show rationale — never silent plan replacement
- Single-user MVP; no auth
```

## Phase 15: Real-World Closure (Briefing, Debrief, Micro-Drills)

**Prerequisite:** Phase 14 complete — pressure mode, plan adjust, and mock interview working.

**Reference:** Innovation Plan Initiatives F, G, H in `docs/product/interview-preparation-app-innovation-plan.md`.

```text
You are implementing Phase 15 of Interview Preparation Coach: Real-World Closure.

Theme: Bridge practice and real rounds.

Context:
- Failure patterns and velocity exist (Phase 13)
- Pressure mode and mock interview exist (Phase 14)
- User motto: measurable recovery before the next real interview

Task:
Implement Initiatives F (Pre-Interview Briefing), G (Real Interview Debrief), H (Micro-Drills).

Part 1 — Pre-Interview Briefing (Initiative F)
1. User input: company (optional), round_type, interview_date
2. Store upcoming_interviews or briefing_requests table
3. Versioned AI prompt: pre_interview_briefing.v1.json
   - Inputs: profile, failure_patterns, velocity, recent wins, round_type
   - Outputs: warm_up_questions[3], top_mistakes_to_watch[2], behavioral_refresh, confidence_script
4. API: POST /api/v1/briefings/generate, GET /api/v1/briefings/latest
5. Frontend: /briefing or section on dashboard when interview_date within 7 days
   - Completable in < 30 minutes
   - Links to recommended practice questions for warm-ups

Part 2 — Real Interview Debrief (Initiative G)
1. Table: interview_debriefs
   - company, round_type, questions_asked (text), self_scores per dimension, outcome (advanced|rejected|waiting), notes
2. API: POST /api/v1/debriefs, GET /api/v1/debriefs, GET /api/v1/debriefs/{id}
3. On create: optionally trigger plan/propose or flag new weak areas for pattern memory
4. Feed debrief summary into next plan regeneration and weekly progress prompt
5. Frontend: /debrief form + history list on progress page

Part 3 — Micro-Drills (Initiative H)
1. New question type or drill_flag on questions: drill_subtype (fix_bug | wrong_approach | explain_60s | spot_edge_case)
2. Seed 15–20 short drills (5–15 min each) tagged to categories matching failure patterns
3. API: GET /api/v1/drills/today — returns 1–3 drills based on weak patterns and regression risks
4. Drill completion counts toward plan task progress (PATCH task complete or drill_attempts table)
5. Frontend: "Quick drill" entry on dashboard; mobile-friendly layout (simple, no heavy IDE required for explain_60s)

Part 4 — Workflow integration
1. Dashboard shows upcoming interview countdown when briefing exists
2. After debrief logged, show suggested next action (practice / adjust / re-diagnose)
3. Update useWorkflowStage if new stages needed (briefing, debrief)

Part 5 — Tests and validation
1. Unit tests for drill selection logic
2. Debrief -> plan propose integration test
3. Manual checklist: full journey practice -> briefing -> mock -> real debrief -> adjust

Return:
1. Migrations and API contracts
2. New AI prompts (briefing at minimum)
3. Frontend pages/components
4. Seed data for drills
5. Manual test checklist
6. Update docs/product/interview-preparation-app-innovation-plan.md success metrics if needed

Rules:
- Drills must be completable in 5–15 minutes
- Debrief is user-reported — do not require employer data
- Briefing is generated quickly; degrade gracefully if AI fails (show pattern-based static briefing)
- No job portal, CV/JD analysis, or auth
- Do not start Phase 16 (company lens, multi-user) in this phase
```

## Phase 16: Strategic Expansion (Company Lens + Platform)

**Prerequisite:** Phase 12 validation succeeded AND Phases 13–15 in use personally for at least 2 weeks.

**Reference:** Innovation Plan Initiative I and platform expansion in `docs/product/interview-preparation-app-innovation-plan.md`.

```text
You are implementing Phase 16 of Interview Preparation Coach: Strategic Expansion.

WARNING: Only proceed if 30-day personal validation (Phase 12) proved the product helps. Otherwise return a short report explaining why to defer Phase 16.

Context:
- Personal-first MVP has failure patterns, velocity, pressure mode, mock interview, briefing, debrief, drills
- target_companies stored on profile as text
- Question bank should be expanded from Phase 13

Task:
Split Phase 16 into two sub-deliverables. Implement only what validation justifies.

Part A — Company Readiness Lens (Initiative I)
1. Parse profile.target_companies into company or company_type tags (startup | big_tech | product_mid)
2. Map question categories and evaluation history to company interview patterns (config JSON, not hardcoded DB of companies)
3. API: GET /api/v1/progress/readiness-by-company
   - Returns readiness score per target company/type based on relevant attempt history
4. Frontend: progress or dashboard tab "Company readiness"
   - Show per-company readiness % and gap areas
5. Do not build job listings or application tracking

Part B — Multi-User / Auth (only if expanding beyond personal use)
1. Document decision criteria: when to add auth vs stay single-user
2. If implementing:
   - users table, profile.user_id FK, session or API key auth
   - migrate existing single profile to user 1
   - minimal login (email magic link or single password — keep simple)
3. If deferring:
   - write docs/product/interview-preparation-app-platform-deferral.md with rationale

Part C — Hardening for others
1. Rate limiting on AI endpoints
2. OpenAI cost logging per request
3. Error budgets and degraded mode audit across all AI features
4. Expanded README for self-hosters

Return:
1. Go/no-go recommendation for Part B
2. Implemented features with API and UI summary
3. company_readiness config format
4. Migration plan if auth added
5. What remains before public beta

Rules:
- Do not build billing, job portal, or community
- Company lens is interview-pattern fit only — not hiring probability claims
- Auth is optional in this phase — default to defer unless user explicitly requests multi-user
- Maintain backward compatibility for existing single-profile data
```

## Usage Reminder

Use one phase prompt at a time.

After each phase:

- review output
- fix gaps
- run tests if code was written
- commit or save progress
- then continue to the next phase
