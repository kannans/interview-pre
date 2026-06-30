# Interview Preparation Coach Master Implementation Prompts

## Purpose

This file contains detailed copy-paste prompts for building the Interview Preparation Coach app from planning to deployment.

Use these prompts one by one. Each prompt is intentionally detailed so the model has enough context to produce useful implementation work.

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

## Phase 11: EC2 Deployment

```text
Prepare Interview Preparation Coach for EC2 deployment.

Stack:
- Next.js frontend
- Go backend
- PostgreSQL
- OpenAI API key stored in environment
- Docker
- Nginx reverse proxy if needed

Return:
1. Dockerfile for backend
2. Dockerfile for frontend
3. docker-compose file for local or single-server deployment
4. environment variable list
5. EC2 setup steps
6. security group ports
7. Nginx reverse proxy config
8. systemd or Docker restart strategy
9. deployment commands
10. smoke test checklist after deployment

Rules:
- Keep it practical for a solo developer.
- Do not expose secrets in code.
- Explain what should be different later if using RDS instead of local PostgreSQL.
```

## Phase 12: Post-Launch Personal Validation

```text
Create a 30-day validation plan for using the app personally.

Goal:
Check whether the app actually helps me prepare better and improve interview readiness.

Return:
1. daily usage routine
2. weekly review checklist
3. success metrics
4. what to log after each session
5. how to decide whether the product is working
6. what to improve after 30 days

Rules:
- Focus on personal interview recovery.
- Do not focus on monetization yet.
```

## Usage Reminder

Use one phase prompt at a time.

After each phase:

- review output
- fix gaps
- run tests if code was written
- commit or save progress
- then continue to the next phase
