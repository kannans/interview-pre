# Interview Preparation Coach Master Build Plan

## Purpose

This document is the single reference plan for building the product from start to finish.

It connects the existing product, strategy, and finance docs into one practical execution path:

- personal-first MVP
- backend implementation
- frontend implementation
- AI prompt design
- testing and iteration
- deployment to EC2

This plan is meant to be followed step by step. Each phase should be completed before moving to the next.

## Reference Docs

Use these docs as the source of truth while building:

- [Personal-First MVP](/Users/kannans/Documents/Codex/2026-06-26/codex-security-plugin-codex-security-openai/outputs/docs/product/interview-preparation-app-personal-first-mvp.md)
- [Product Requirements Document](/Users/kannans/Documents/Codex/2026-06-26/codex-security-plugin-codex-security-openai/outputs/docs/product/interview-preparation-app-prd.md)
- [Technical Design](/Users/kannans/Documents/Codex/2026-06-26/codex-security-plugin-codex-security-openai/outputs/docs/product/interview-preparation-app-technical-design.md)
- [Roadmap](/Users/kannans/Documents/Codex/2026-06-26/codex-security-plugin-codex-security-openai/outputs/docs/product/interview-preparation-app-roadmap.md)
- [Master Implementation Prompts](/Users/kannans/Documents/Codex/2026-06-26/codex-security-plugin-codex-security-openai/outputs/docs/prompts/interview-preparation-app-master-implementation-prompts.md)
- [Master Strategy](/Users/kannans/Documents/Codex/2026-06-26/codex-security-plugin-codex-security-openai/outputs/docs/strategy/interview-preparation-app-master-strategy.md)

## Product Direction

The product is a personal-first interview recovery app for experienced engineers.

The goal is not revenue first. The goal is:

- prepare better
- track progress
- improve coding confidence
- answer interview questions better
- perform better in real interviews

Revenue and startup scaling can come later.

## Recommended Stack

### Frontend

- `Next.js`
- `TypeScript`
- `Tailwind CSS`

### Backend

- `Go`
- `Chi`
- `pgx`
- `sqlc`
- `golang-migrate`

### Database

- `PostgreSQL`

### AI

- `OpenAI Responses API` through the Go backend

### Deployment

- `Frontend` on EC2 or another deploy target later
- `Backend` on EC2
- `PostgreSQL` on RDS or EC2 depending on stage
- `Nginx` as reverse proxy if needed
- `Docker` for local and server deployment

## Build Principles

- Build one flow end to end before adding the next.
- Keep the first version small and useful.
- Store every attempt, score, and mistake tag.
- Keep AI calls behind the Go backend.
- Do not add job portal features in v1.
- Do not add CV/JD analysis in the first recovery loop unless explicitly needed later.

## Implementation Phases

### Phase 1: Scope Freeze

Goal:

- lock the exact MVP and avoid scope creep.

Outputs:

- final MVP scope
- user flow
- feature boundary list
- success metrics

Prompt to use:

```text
Turn the personal-first interview recovery idea into a frozen MVP scope.
Include what to build, what not to build, success metrics, and the first user flow.
Do not write code yet.
```

### Phase 2: Technical Architecture

Goal:

- define the system structure before coding.

Outputs:

- frontend pages
- Go backend modules
- PostgreSQL tables
- AI integration design
- API contracts
- workflow diagrams

Prompt to use:

```text
Design the technical architecture for the personal-first MVP.
Use Next.js frontend, Go backend, PostgreSQL, and OpenAI through the backend.
Return architecture, data flow, and API structure.
```

### Phase 3: Folder Structure

Goal:

- create a clean repo layout for frontend, backend, docs, prompts, and deployment.

Outputs:

- file tree
- backend service layout
- frontend page layout
- docs folder organization

Prompt to use:

```text
Create the full folder structure for the app.
Split frontend, backend, database, AI prompts, deployment, and docs cleanly.
Return only the file tree and short notes.
```

### Phase 4: Database Design

Goal:

- define the schema before writing application logic.

Outputs:

- PostgreSQL tables
- relationships
- indexes
- migration plan

Prompt to use:

```text
Design the PostgreSQL schema for the personal-first MVP.
Include user profile, baseline diagnosis, recovery plan, practice attempts, evaluations, and progress snapshots.
Return SQL-ready table definitions.
```

### Phase 5: Go Backend Foundation

Goal:

- set up the backend skeleton and first API endpoints.

Outputs:

- Go app entrypoint
- config loading
- router
- DB connection
- health endpoint

Prompt to use:

```text
Now implement the Go backend foundation.
Create the app entrypoint, router, config, PostgreSQL connection, and a health endpoint.
Keep it simple and production-friendly.
```

### Phase 6: Core Backend Modules

Goal:

- build the main recovery loop on the backend.

Modules:

- profile
- baseline diagnosis
- recovery plan
- practice attempts
- evaluations
- progress tracking
- AI prompt execution

Prompt to use:

```text
Implement the backend module by module in this order:
profile, baseline diagnosis, recovery plan, practice attempts, evaluations, and progress tracking.
After each module, show what is complete and what is left.
```

### Phase 7: AI Prompt Layer

Goal:

- create stable prompts for each app function.

Prompt areas:

- baseline diagnosis
- recovery plan generation
- basic programming evaluation
- DSA evaluation
- technical answer evaluation
- mock question generation
- mock grading
- progress summary

Prompt to use:

```text
Create versioned AI prompt templates for each section of the app.
Each prompt should return structured JSON and be suitable for use through the backend.
```

### Phase 8: Next.js Frontend Foundation

Goal:

- build the app UI around the recovery workflow.

Pages:

- onboarding
- baseline diagnosis
- recovery plan
- daily practice
- mock interview
- progress dashboard

Prompt to use:

```text
Implement the Next.js frontend foundation for the personal-first MVP.
Build the main pages and connect them to the backend APIs.
Keep the UI clean, practical, and easy to scan.
```

### Phase 9: End-to-End Workflow

Goal:

- make the full loop work from profile to progress tracking.

Workflow:

- user creates profile
- user completes diagnosis
- app generates plan
- user practices
- AI evaluates
- progress updates
- next task is recommended

Prompt to use:

```text
Connect the frontend and backend into one working loop.
Make sure the full recovery workflow works end to end.
```

### Phase 10: Testing

Goal:

- validate the app behaves correctly.

Test areas:

- API tests
- DB migration tests
- prompt output checks
- workflow validation
- frontend state handling

Prompt to use:

```text
Review and test the current implementation.
List bugs, edge cases, missing tests, and workflow gaps.
Focus on user-facing issues and data integrity.
```

### Phase 11: Deployment to EC2

Goal:

- deploy frontend and backend to EC2 in a clean and repeatable way.

Recommended approach:

- Dockerize backend and frontend
- use EC2 for app hosting
- use RDS or managed PostgreSQL if possible
- configure Nginx if needed
- set environment variables securely
- set up logs and restart policy

Prompt to use:

```text
Prepare the app for deployment to EC2.
Show the Docker setup, environment variables, server setup, and deployment steps for frontend and backend.
```

## Build Order

Use this order strictly:

1. freeze MVP scope
2. design architecture
3. create folder structure
4. design database
5. build Go backend foundation
6. build backend modules
7. create AI prompt layer
8. build Next.js frontend
9. connect full workflow
10. test and fix
11. deploy to EC2

## What To Avoid In v1

- job portal features
- vacancy search
- large company database
- public community features
- heavy gamification
- voice interview system
- broad role support
- unnecessary CV/JD analysis in the first loop

## Success Criteria

The first version is successful if it helps you:

- practice consistently
- solve basic coding questions better
- answer more clearly
- track progress
- see improvement week by week
- feel more prepared for real interviews

## Final Build Mindset

Build this like a recovery system first.

Only after it helps you personally should it grow into a product for others and eventually into monetization.

The product should start with this loop:

`profile -> diagnosis -> plan -> practice -> evaluate -> progress -> adjust`
