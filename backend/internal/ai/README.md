# AI Prompt Layer

Interview Preparation Coach calls OpenAI only through `backend/internal/ai`. The frontend never talks to OpenAI directly.

## Layout

- `client.go` — OpenAI Responses API HTTP wrapper (`POST /v1/responses`) with strict JSON schema output
- `service.go` — domain-facing `Service` interface used by diagnosis, plan, evaluation, and progress services
- `prompts/` — versioned prompt documents (`{id}.v{N}.json`)
- `schemas.go` — strict JSON parsing and validation after each AI call

## Prompt versioning

Prompts are versioned in **both** the filename and an embedded `version` field:

- Filename: `baseline_diagnosis.v1.json` → canonical id `baseline_diagnosis`, version `v1`
- To ship breaking output changes, add `baseline_diagnosis.v2.json` and update the registry lookup (or bump the constant) in code
- Old versions remain in git for audit and A/B comparison

Each prompt JSON includes: `system_instruction`, `user_input_shape`, `output_schema`, `example_input`, `example_output`.

## Logging

Every OpenAI call logs (via `slog`):

- `prompt_id`, `prompt_version`
- `model`
- `latency_ms`
- `input_tokens`, `output_tokens`, `total_tokens` (when returned by API)
- `response_id`
- Errors log HTTP status and response body (no API key)

Raw AI responses are stored on domain records (`diagnostic_report`, evaluation `raw_ai_response`, progress snapshot stats metadata).

## Fallback behavior

| Operation | On AI failure |
|-----------|----------------|
| Diagnosis Submit | **Fail the request** — no diagnosis row with fake data |
| Plan Generate | **Fail the request** — caller must retry |
| Evaluation Evaluate | **Fail the request** — caller must retry |
| Progress CreateSnapshot | **Degraded save** — persist computed numeric stats with `status: "degraded_ai_summary"` and omit AI narrative fields |

Rationale: diagnosis, plan, and evaluation are AI-core; progress tracking should not block on summary generation.

## Environment

- `OPENAI_API_KEY` — required for AI endpoints
- `OPENAI_MODEL` — optional override (default `gpt-4o-mini`)

## Stubs (Phase 8)

- `mock_interview_questions.v1.json` — prompt defined; `GenerateMockInterviewQuestions` returns `ErrNotImplemented`
- `mock_interview_grading.v1.json` — prompt defined; `GradeMockInterview` returns `ErrNotImplemented`
