-- Enable UUID extension if not active
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 1. Profiles Table
CREATE TABLE IF NOT EXISTS profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    experience_years INT NOT NULL,
    target_role VARCHAR(255) NOT NULL,
    target_companies TEXT NOT NULL,          -- Comma separated or simple list
    weak_areas_self TEXT NOT NULL,           -- Self-identified weak areas
    timeline VARCHAR(100) NOT NULL,          -- e.g. "30 days", "3 months"
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- 2. Baseline Diagnoses Table
CREATE TABLE IF NOT EXISTS baseline_diagnoses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    profile_id UUID NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    coding_answer TEXT NOT NULL,
    conceptual_answer TEXT NOT NULL,
    behavioral_answer TEXT NOT NULL,
    diagnostic_report JSONB NOT NULL,        -- AI parsed diagnosis scorecard
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- 3. Recovery Plans Table
CREATE TABLE IF NOT EXISTS recovery_plans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    profile_id UUID NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    status VARCHAR(50) NOT NULL DEFAULT 'active', -- 'active', 'completed', 'paused'
    generated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- 4. Plan Tasks Table
CREATE TABLE IF NOT EXISTS plan_tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plan_id UUID NOT NULL REFERENCES recovery_plans(id) ON DELETE CASCADE,
    category VARCHAR(100) NOT NULL,          -- 'programming_basics', 'dsa', 'system_design'
    level VARCHAR(50) NOT NULL,              -- 'level_1', 'level_2', etc.
    title VARCHAR(255) NOT NULL,
    sequence_order INT NOT NULL,
    is_completed BOOLEAN NOT NULL DEFAULT FALSE,
    completed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- 5. Questions Catalog Table
CREATE TABLE IF NOT EXISTS questions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type VARCHAR(50) NOT NULL,               -- 'programming_basics', 'dsa', 'concept'
    category VARCHAR(100) NOT NULL,          -- 'arrays', 'loops', 'closures', 'react_hooks'
    level VARCHAR(50) NOT NULL,              -- 'level_1', 'level_2', 'level_3'
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,                   -- Markdown problem description
    constraints TEXT,                        -- Markdown string of constraints
    starter_code TEXT,                       -- Sandbox starter files/code boilerplate
    solution_hint TEXT,                      -- Helpful hints to expose in UI
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- 6. Attempts Log Table
CREATE TABLE IF NOT EXISTS attempts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    profile_id UUID NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    question_id UUID NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    submitted_content TEXT NOT NULL,
    time_taken_seconds INT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- 7. AI Evaluations Table
CREATE TABLE IF NOT EXISTS evaluations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    attempt_id UUID NOT NULL UNIQUE REFERENCES attempts(id) ON DELETE CASCADE,
    score_correctness INT NOT NULL CHECK (score_correctness BETWEEN 0 AND 10),
    score_approach INT NOT NULL CHECK (score_approach BETWEEN 0 AND 10),
    score_edge_cases INT NOT NULL CHECK (score_edge_cases BETWEEN 0 AND 10),
    score_communication INT NOT NULL CHECK (score_communication BETWEEN 0 AND 10),
    feedback_correctness TEXT NOT NULL,
    feedback_approach TEXT NOT NULL,
    feedback_edge_cases TEXT NOT NULL,
    feedback_communication TEXT NOT NULL,
    raw_ai_response JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- 8. Progress Snapshots Table
CREATE TABLE IF NOT EXISTS progress_snapshots (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    profile_id UUID NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    snapshot_date DATE NOT NULL DEFAULT CURRENT_DATE,
    stats JSONB NOT NULL,                    -- Completed counts, weak patterns, readiness index
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE(profile_id, snapshot_date)
);

-- Useful Indexes
CREATE INDEX IF NOT EXISTS idx_baseline_diagnoses_profile ON baseline_diagnoses(profile_id);
CREATE INDEX IF NOT EXISTS idx_plan_tasks_plan_seq ON plan_tasks(plan_id, sequence_order);
CREATE INDEX IF NOT EXISTS idx_questions_lookup ON questions(type, category, level);
CREATE INDEX IF NOT EXISTS idx_attempts_profile_question ON attempts(profile_id, question_id);
CREATE INDEX IF NOT EXISTS idx_progress_snapshots_date ON progress_snapshots(profile_id, snapshot_date);
