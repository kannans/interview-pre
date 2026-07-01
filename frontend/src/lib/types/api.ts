/** Mirrors Go backend JSON contracts (db + service layers). */

export interface Profile {
  id: string;
  experience_years: number;
  target_role: string;
  target_companies: string;
  weak_areas_self: string;
  timeline: string;
  created_at: string;
  updated_at: string;
}

export interface ProfileInput {
  experience_years: number;
  target_role: string;
  target_companies: string;
  weak_areas_self: string;
  timeline: string;
}

export interface DimensionScores {
  programming_basics: number;
  dsa: number;
  system_design_concepts: number;
  communication: number;
  behavioral: number;
}

export interface DiagnosticReport {
  overall_readiness: number;
  dimension_scores: DimensionScores;
  weak_areas: string[];
  strengths: string[];
  recommended_start_level: string;
  summary: string;
  priority_focus: string[];
  prompt_version?: string;
  model?: string;
}

export interface BaselineDiagnosis {
  id: string;
  profile_id: string;
  coding_answer: string;
  conceptual_answer: string;
  behavioral_answer: string;
  /** Go encodes JSONB []byte as base64; may also arrive as object. */
  diagnostic_report: string | DiagnosticReport;
  created_at: string;
}

export interface DiagnosisInput {
  coding_answer: string;
  conceptual_answer: string;
  behavioral_answer: string;
}

export interface RecoveryPlan {
  id: string;
  profile_id: string;
  status: string;
  generated_at: string;
}

export interface PlanTask {
  id: string;
  plan_id: string;
  category: string;
  level: string;
  title: string;
  sequence_order: number;
  is_completed: boolean;
  completed_at: string | null;
  created_at: string;
}

export interface PlanWithTasks {
  plan: RecoveryPlan;
  tasks: PlanTask[];
}

export interface TodayTaskResponse {
  plan: RecoveryPlan;
  task: PlanTask;
}

export interface Question {
  id: string;
  type: string;
  category: string;
  level: string;
  title: string;
  content: string;
  constraints: string | null;
  starter_code: string | null;
  solution_hint: string | null;
  created_at: string;
}

export interface ListQuestionsResponse {
  questions: Question[];
}

export interface QuestionListFilter {
  type?: string;
  category?: string;
  level?: string;
}

export interface Attempt {
  id: string;
  profile_id: string;
  question_id: string;
  submitted_content: string;
  time_taken_seconds: number;
  language: string;
  created_at: string;
}

export interface AttemptInput {
  question_id: string;
  submitted_content: string;
  time_taken_seconds: number;
  language: string;
}

export interface RunCodeInput {
  language: string;
  code: string;
  stdin?: string;
}

export interface RunCodeResult {
  stdout: string;
  stderr: string;
  exit_code: number;
  duration_ms: number;
  timed_out: boolean;
}

export interface ListAttemptsResponse {
  attempts: Attempt[];
}

export interface Evaluation {
  id: string;
  attempt_id: string;
  score_correctness: number;
  score_approach: number;
  score_edge_cases: number;
  score_communication: number;
  feedback_correctness: string;
  feedback_approach: string;
  feedback_edge_cases: string;
  feedback_communication: string;
  created_at: string;
}

export interface SnapshotAverageScores {
  correctness: number;
  approach: number;
  edge_cases: number;
  communication: number;
}

export interface SnapshotStats {
  attempts_count: number;
  evaluations_count: number;
  average_scores: SnapshotAverageScores;
  plan_tasks_completed: number;
  plan_tasks_total: number;
  readiness_index: number;
  weak_patterns: string[];
  status: string;
  weekly_summary?: string;
  wins?: string[];
  focus_areas?: string[];
  adjustment_recommendation?: string;
}

export interface ProgressSnapshot {
  id: string;
  profile_id: string;
  snapshot_date: string;
  stats: SnapshotStats;
  created_at: string;
}

export interface ListProgressSnapshotsResponse {
  snapshots: ProgressSnapshot[];
}

export interface ApiErrorBody {
  error: string;
}
