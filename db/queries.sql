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
