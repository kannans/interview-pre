package ai

import (
	"encoding/json"
	"fmt"
	"strings"
)

func clampScore(v int) error {
	if v < 0 || v > 10 {
		return fmt.Errorf("score %d out of range 0-10", v)
	}
	return nil
}

func requireNonEmpty(field, value string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%s is required", field)
	}
	return nil
}

func requireNonEmptySlice(field string, values []string) error {
	if len(values) == 0 {
		return fmt.Errorf("%s must have at least one item", field)
	}
	for i, v := range values {
		if strings.TrimSpace(v) == "" {
			return fmt.Errorf("%s[%d] is empty", field, i)
		}
	}
	return nil
}

func validateLevel(level string) error {
	switch level {
	case "level_1", "level_2", "level_3":
		return nil
	default:
		return fmt.Errorf("invalid level %q", level)
	}
}

func validatePlanCategory(category string) error {
	switch category {
	case "programming_basics", "dsa", "concept":
		return nil
	default:
		return fmt.Errorf("invalid category %q", category)
	}
}

// ParseDiagnosticReport unmarshals and validates baseline diagnosis output.
func ParseDiagnosticReport(raw json.RawMessage) (DiagnosticReport, error) {
	var report DiagnosticReport
	if err := json.Unmarshal(raw, &report); err != nil {
		return DiagnosticReport{}, fmt.Errorf("%w: %v", ErrInvalidJSON, err)
	}
	if report.OverallReadiness < 0 || report.OverallReadiness > 100 {
		return DiagnosticReport{}, fmt.Errorf("%w: overall_readiness out of range", ErrValidation)
	}
	scores := report.DimensionScores
	for name, score := range map[string]int{
		"programming_basics":     scores.ProgrammingBasics,
		"dsa":                      scores.DSA,
		"system_design_concepts": scores.SystemDesignConcept,
		"communication":            scores.Communication,
		"behavioral":               scores.Behavioral,
	} {
		if err := clampScore(score); err != nil {
			return DiagnosticReport{}, fmt.Errorf("%w: dimension_scores.%s %v", ErrValidation, name, err)
		}
	}
	if err := requireNonEmpty("summary", report.Summary); err != nil {
		return DiagnosticReport{}, fmt.Errorf("%w: %v", ErrValidation, err)
	}
	if err := requireNonEmptySlice("weak_areas", report.WeakAreas); err != nil {
		return DiagnosticReport{}, fmt.Errorf("%w: %v", ErrValidation, err)
	}
	if err := validateLevel(report.RecommendedStartLevel); err != nil {
		return DiagnosticReport{}, fmt.Errorf("%w: %v", ErrValidation, err)
	}
	return report, nil
}

// ParseRecoveryPlanOutput unmarshals and validates recovery plan output.
func ParseRecoveryPlanOutput(raw json.RawMessage) (RecoveryPlanOutput, error) {
	var out RecoveryPlanOutput
	if err := json.Unmarshal(raw, &out); err != nil {
		return RecoveryPlanOutput{}, fmt.Errorf("%w: %v", ErrInvalidJSON, err)
	}
	if len(out.Tasks) < 5 || len(out.Tasks) > 10 {
		return RecoveryPlanOutput{}, fmt.Errorf("%w: tasks count must be 5-10, got %d", ErrValidation, len(out.Tasks))
	}
	if err := requireNonEmpty("plan_summary", out.PlanSummary); err != nil {
		return RecoveryPlanOutput{}, fmt.Errorf("%w: %v", ErrValidation, err)
	}
	for i, task := range out.Tasks {
		if err := validatePlanCategory(task.Category); err != nil {
			return RecoveryPlanOutput{}, fmt.Errorf("%w: tasks[%d] %v", ErrValidation, i, err)
		}
		if err := validateLevel(task.Level); err != nil {
			return RecoveryPlanOutput{}, fmt.Errorf("%w: tasks[%d] %v", ErrValidation, i, err)
		}
		if err := requireNonEmpty("title", task.Title); err != nil {
			return RecoveryPlanOutput{}, fmt.Errorf("%w: tasks[%d] %v", ErrValidation, i, err)
		}
	}
	return out, nil
}

// ParseEvaluationOutput unmarshals and validates evaluation output.
func ParseEvaluationOutput(raw json.RawMessage) (EvaluationOutput, error) {
	var out EvaluationOutput
	if err := json.Unmarshal(raw, &out); err != nil {
		return EvaluationOutput{}, fmt.Errorf("%w: %v", ErrInvalidJSON, err)
	}
	for name, score := range map[string]int{
		"score_correctness":   out.ScoreCorrectness,
		"score_approach":      out.ScoreApproach,
		"score_edge_cases":    out.ScoreEdgeCases,
		"score_communication": out.ScoreCommunication,
	} {
		if err := clampScore(score); err != nil {
			return EvaluationOutput{}, fmt.Errorf("%w: %s %v", ErrValidation, name, err)
		}
	}
	for name, text := range map[string]string{
		"feedback_correctness":   out.FeedbackCorrectness,
		"feedback_approach":      out.FeedbackApproach,
		"feedback_edge_cases":    out.FeedbackEdgeCases,
		"feedback_communication": out.FeedbackCommunication,
		"next_action":            out.NextAction,
	} {
		if err := requireNonEmpty(name, text); err != nil {
			return EvaluationOutput{}, fmt.Errorf("%w: %v", ErrValidation, err)
		}
	}
	return out, nil
}

// ParseWeeklyProgressOutput unmarshals and validates weekly progress summary output.
func ParseWeeklyProgressOutput(raw json.RawMessage) (WeeklyProgressOutput, error) {
	var out WeeklyProgressOutput
	if err := json.Unmarshal(raw, &out); err != nil {
		return WeeklyProgressOutput{}, fmt.Errorf("%w: %v", ErrInvalidJSON, err)
	}
	if err := requireNonEmpty("weekly_summary", out.WeeklySummary); err != nil {
		return WeeklyProgressOutput{}, fmt.Errorf("%w: %v", ErrValidation, err)
	}
	if err := requireNonEmptySlice("focus_areas", out.FocusAreas); err != nil {
		return WeeklyProgressOutput{}, fmt.Errorf("%w: %v", ErrValidation, err)
	}
	if err := requireNonEmpty("adjustment_recommendation", out.AdjustmentRecommendation); err != nil {
		return WeeklyProgressOutput{}, fmt.Errorf("%w: %v", ErrValidation, err)
	}
	return out, nil
}
