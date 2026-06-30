package ai

import (
	"encoding/json"
	"testing"
)

func TestParseDiagnosticReport(t *testing.T) {
	raw := json.RawMessage(`{
		"overall_readiness": 58,
		"dimension_scores": {
			"programming_basics": 7,
			"dsa": 5,
			"system_design_concepts": 6,
			"communication": 6,
			"behavioral": 7
		},
		"weak_areas": ["hash maps"],
		"strengths": ["clarity"],
		"recommended_start_level": "level_2",
		"summary": "Needs DSA work.",
		"priority_focus": ["two-pointer"]
	}`)

	report, err := ParseDiagnosticReport(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if report.OverallReadiness != 58 {
		t.Fatalf("expected readiness 58, got %d", report.OverallReadiness)
	}
}

func TestRegistryLoadsAllPrompts(t *testing.T) {
	reg, err := NewRegistry()
	if err != nil {
		t.Fatalf("registry load failed: %v", err)
	}
	ids := []string{
		PromptBaselineDiagnosis,
		PromptRecoveryPlan,
		PromptProgrammingEval,
		PromptDSAEval,
		PromptTechnicalAnswerEval,
		PromptMockInterviewQuestions,
		PromptMockInterviewGrading,
		PromptWeeklyProgressSummary,
	}
	for _, id := range ids {
		if _, err := reg.Get(id); err != nil {
			t.Fatalf("missing prompt %s: %v", id, err)
		}
	}
}
