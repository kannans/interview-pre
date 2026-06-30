-- name: GetProfile :one
SELECT * FROM profiles LIMIT 1;

-- name: CreateProfile :one
INSERT INTO profiles (experience_years, target_role, target_companies, weak_areas_self, timeline)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateProfile :one
UPDATE profiles
SET experience_years = $2,
    target_role = $3,
    target_companies = $4,
    weak_areas_self = $5,
    timeline = $6,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: HasProfile :one
SELECT EXISTS(SELECT 1 FROM profiles);

-- name: CreateBaselineDiagnosis :one
INSERT INTO baseline_diagnoses (profile_id, coding_answer, conceptual_answer, behavioral_answer, diagnostic_report)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetBaselineDiagnosis :one
SELECT * FROM baseline_diagnoses WHERE profile_id = $1 ORDER BY created_at DESC LIMIT 1;

-- name: HasBaselineDiagnosis :one
SELECT EXISTS(SELECT 1 FROM baseline_diagnoses WHERE profile_id = $1);

-- name: CreateRecoveryPlan :one
INSERT INTO recovery_plans (profile_id, status)
VALUES ($1, 'active')
RETURNING *;

-- name: GetActiveRecoveryPlan :one
SELECT * FROM recovery_plans
WHERE profile_id = $1 AND status = 'active'
ORDER BY generated_at DESC
LIMIT 1;

-- name: HasActiveRecoveryPlan :one
SELECT EXISTS(
    SELECT 1 FROM recovery_plans WHERE profile_id = $1 AND status = 'active'
);

-- name: SupersedeActiveRecoveryPlan :exec
UPDATE recovery_plans
SET status = 'replaced'
WHERE profile_id = $1 AND status = 'active';

-- name: CreatePlanTask :one
INSERT INTO plan_tasks (plan_id, category, level, title, sequence_order)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: ListPlanTasks :many
SELECT * FROM plan_tasks
WHERE plan_id = $1
ORDER BY sequence_order ASC;

-- name: GetPlanTask :one
SELECT * FROM plan_tasks WHERE id = $1;

-- name: GetTodayPlanTask :one
SELECT * FROM plan_tasks
WHERE plan_id = $1 AND is_completed = FALSE
ORDER BY sequence_order ASC
LIMIT 1;

-- name: MarkPlanTaskComplete :one
UPDATE plan_tasks
SET is_completed = TRUE, completed_at = NOW()
WHERE id = $1 AND plan_id = $2 AND is_completed = FALSE
RETURNING *;

-- name: ListQuestions :many
SELECT *
FROM questions
WHERE (sqlc.narg('filter_type')::text IS NULL OR type = sqlc.narg('filter_type'))
  AND (sqlc.narg('filter_category')::text IS NULL OR category = sqlc.narg('filter_category'))
  AND (sqlc.narg('filter_level')::text IS NULL OR level = sqlc.narg('filter_level'))
ORDER BY type ASC, category ASC, level ASC, created_at ASC;

-- name: GetQuestion :one
SELECT * FROM questions WHERE id = $1;

-- name: CreateAttempt :one
INSERT INTO attempts (profile_id, question_id, submitted_content, time_taken_seconds, language)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetAttempt :one
SELECT * FROM attempts WHERE id = $1;

-- name: ListAttempts :many
SELECT *
FROM attempts
WHERE profile_id = $1
  AND (sqlc.narg('filter_question_id')::uuid IS NULL OR question_id = sqlc.narg('filter_question_id'))
ORDER BY created_at DESC;

-- name: CreateEvaluation :one
INSERT INTO evaluations (
    attempt_id,
    score_correctness,
    score_approach,
    score_edge_cases,
    score_communication,
    feedback_correctness,
    feedback_approach,
    feedback_edge_cases,
    feedback_communication,
    raw_ai_response
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;

-- name: GetEvaluationByAttemptID :one
SELECT * FROM evaluations WHERE attempt_id = $1;

-- name: GetEvaluationByID :one
SELECT * FROM evaluations WHERE id = $1;

-- name: CreateProgressSnapshot :one
INSERT INTO progress_snapshots (profile_id, snapshot_date, stats)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateProgressSnapshot :one
UPDATE progress_snapshots
SET stats = $3
WHERE profile_id = $1 AND snapshot_date = $2
RETURNING *;

-- name: GetLatestProgressSnapshot :one
SELECT * FROM progress_snapshots
WHERE profile_id = $1
ORDER BY snapshot_date DESC
LIMIT 1;

-- name: GetProgressSnapshotByID :one
SELECT * FROM progress_snapshots WHERE id = $1;

-- name: GetProgressSnapshotByProfileAndDate :one
SELECT * FROM progress_snapshots
WHERE profile_id = $1 AND snapshot_date = $2;

-- name: ListProgressSnapshots :many
SELECT *
FROM progress_snapshots
WHERE profile_id = $1
  AND (sqlc.narg('filter_start_date')::date IS NULL OR snapshot_date >= sqlc.narg('filter_start_date')::date)
  AND (sqlc.narg('filter_end_date')::date IS NULL OR snapshot_date <= sqlc.narg('filter_end_date')::date)
ORDER BY snapshot_date DESC
LIMIT COALESCE(sqlc.narg('row_limit')::int, 100);

-- name: CountAttemptsByProfile :one
SELECT COUNT(*)::int AS count FROM attempts WHERE profile_id = $1;

-- name: GetEvaluationStatsForProfile :one
SELECT
    COUNT(e.id)::int AS evaluations_count,
    COALESCE(AVG(e.score_correctness), 0)::float8 AS avg_correctness,
    COALESCE(AVG(e.score_approach), 0)::float8 AS avg_approach,
    COALESCE(AVG(e.score_edge_cases), 0)::float8 AS avg_edge_cases,
    COALESCE(AVG(e.score_communication), 0)::float8 AS avg_communication
FROM attempts a
INNER JOIN evaluations e ON e.attempt_id = a.id
WHERE a.profile_id = $1;
