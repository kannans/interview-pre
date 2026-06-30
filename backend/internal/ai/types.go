package ai

import (
	"encoding/json"
)

const (
	PromptBaselineDiagnosis      = "baseline_diagnosis"
	PromptRecoveryPlan           = "recovery_plan"
	PromptProgrammingEval        = "programming_basics_eval"
	PromptDSAEval                = "dsa_eval"
	PromptTechnicalAnswerEval    = "technical_answer_eval"
	PromptMockInterviewQuestions = "mock_interview_questions"
	PromptMockInterviewGrading   = "mock_interview_grading"
	PromptWeeklyProgressSummary  = "weekly_progress_summary"
)

// ProfileContext is shared profile fields passed to AI prompts.
type ProfileContext struct {
	ExperienceYears int32  `json:"experience_years"`
	TargetRole      string `json:"target_role"`
	TargetCompanies string `json:"target_companies"`
	WeakAreasSelf   string `json:"weak_areas_self"`
	Timeline        string `json:"timeline"`
}

// BaselineDiagnosisInput is the user payload for baseline diagnosis.
type BaselineDiagnosisInput struct {
	Profile          ProfileContext `json:"profile"`
	CodingAnswer     string         `json:"coding_answer"`
	ConceptualAnswer string         `json:"conceptual_answer"`
	BehavioralAnswer string         `json:"behavioral_answer"`
}

// DimensionScores holds per-dimension 0–10 scores.
type DimensionScores struct {
	ProgrammingBasics   int `json:"programming_basics"`
	DSA                 int `json:"dsa"`
	SystemDesignConcept int `json:"system_design_concepts"`
	Communication       int `json:"communication"`
	Behavioral          int `json:"behavioral"`
}

// DiagnosticReport is the validated AI output for baseline diagnosis.
type DiagnosticReport struct {
	OverallReadiness      int             `json:"overall_readiness"`
	DimensionScores       DimensionScores `json:"dimension_scores"`
	WeakAreas             []string        `json:"weak_areas"`
	Strengths             []string        `json:"strengths"`
	RecommendedStartLevel string          `json:"recommended_start_level"`
	Summary               string          `json:"summary"`
	PriorityFocus         []string        `json:"priority_focus"`
	PromptVersion         string          `json:"prompt_version,omitempty"`
	Model                 string          `json:"model,omitempty"`
}

// RecoveryPlanInput is the user payload for plan generation.
type RecoveryPlanInput struct {
	Profile          ProfileContext  `json:"profile"`
	DiagnosticReport json.RawMessage `json:"diagnostic_report"`
}

// PlanTaskSuggestion is one AI-generated recovery task.
type PlanTaskSuggestion struct {
	Category  string `json:"category"`
	Level     string `json:"level"`
	Title     string `json:"title"`
	Rationale string `json:"rationale"`
}

// RecoveryPlanOutput is the validated AI output for plan generation.
type RecoveryPlanOutput struct {
	Tasks        []PlanTaskSuggestion `json:"tasks"`
	PlanSummary  string               `json:"plan_summary"`
	PromptVersion string              `json:"prompt_version,omitempty"`
	Model         string              `json:"model,omitempty"`
}

// EvaluationInput is the user payload for attempt evaluation.
type EvaluationInput struct {
	Profile          ProfileContext `json:"profile"`
	QuestionType     string         `json:"question_type"`
	QuestionCategory string         `json:"question_category"`
	QuestionLevel    string         `json:"question_level"`
	QuestionTitle    string         `json:"question_title"`
	QuestionContent  string         `json:"question_content"`
	Constraints      string         `json:"constraints,omitempty"`
	StarterCode      string         `json:"starter_code,omitempty"`
	SolutionHint     string         `json:"solution_hint,omitempty"`
	SubmittedContent string         `json:"submitted_content"`
	TimeTakenSeconds int32          `json:"time_taken_seconds"`
}

// EvaluationOutput is the validated AI output for practice evaluation.
type EvaluationOutput struct {
	ScoreCorrectness      int    `json:"score_correctness"`
	ScoreApproach         int    `json:"score_approach"`
	ScoreEdgeCases        int    `json:"score_edge_cases"`
	ScoreCommunication    int    `json:"score_communication"`
	FeedbackCorrectness   string `json:"feedback_correctness"`
	FeedbackApproach      string `json:"feedback_approach"`
	FeedbackEdgeCases     string `json:"feedback_edge_cases"`
	FeedbackCommunication string `json:"feedback_communication"`
	NextAction            string `json:"next_action"`
	PromptVersion         string `json:"prompt_version,omitempty"`
	Model                 string `json:"model,omitempty"`
}

// WeeklyProgressInput is the user payload for progress summary.
type WeeklyProgressInput struct {
	Profile             ProfileContext        `json:"profile"`
	AttemptsCount       int32                 `json:"attempts_count"`
	EvaluationsCount    int32                 `json:"evaluations_count"`
	AverageScores       SnapshotAverageScores `json:"average_scores"`
	PlanTasksCompleted  int32                 `json:"plan_tasks_completed"`
	PlanTasksTotal      int32                 `json:"plan_tasks_total"`
	ReadinessIndex      int32                 `json:"readiness_index"`
	WeakPatterns        []string              `json:"weak_patterns"`
	SnapshotDate        string                `json:"snapshot_date"`
}

// SnapshotAverageScores mirrors service-layer snapshot averages.
type SnapshotAverageScores struct {
	Correctness   float64 `json:"correctness"`
	Approach      float64 `json:"approach"`
	EdgeCases     float64 `json:"edge_cases"`
	Communication float64 `json:"communication"`
}

// WeeklyProgressOutput is the validated AI output for weekly summary.
type WeeklyProgressOutput struct {
	WeeklySummary            string   `json:"weekly_summary"`
	Wins                     []string `json:"wins"`
	FocusAreas               []string `json:"focus_areas"`
	AdjustmentRecommendation string   `json:"adjustment_recommendation"`
	PromptVersion            string   `json:"prompt_version,omitempty"`
	Model                    string   `json:"model,omitempty"`
}

// MockInterviewQuestionsOutput is a stub schema for future mock interview generation.
type MockInterviewQuestionsOutput struct {
	Questions []MockInterviewQuestion `json:"questions"`
}

type MockInterviewQuestion struct {
	Type     string `json:"type"`
	Category string `json:"category"`
	Level    string `json:"level"`
	Title    string `json:"title"`
	Content  string `json:"content"`
}

// MockInterviewGradingOutput is a stub schema for future mock interview grading.
type MockInterviewGradingOutput struct {
	OverallScore int      `json:"overall_score"`
	Strengths    []string `json:"strengths"`
	Gaps         []string `json:"gaps"`
	Feedback     string   `json:"feedback"`
	NextSteps    []string `json:"next_steps"`
}

// AIResult wraps a typed result with raw JSON and call metadata.
type AIResult[T any] struct {
	Data T
	Raw  json.RawMessage
	Meta CallMeta
}
