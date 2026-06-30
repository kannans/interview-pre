package ai

import (
	"context"
	"encoding/json"
	"fmt"
)

// Service is the domain-facing AI interface used by backend services.
type Service interface {
	EvaluateBaselineDiagnosis(ctx context.Context, input BaselineDiagnosisInput) (AIResult[DiagnosticReport], error)
	GenerateRecoveryPlan(ctx context.Context, input RecoveryPlanInput) (AIResult[RecoveryPlanOutput], error)
	EvaluateAttempt(ctx context.Context, questionType string, input EvaluationInput) (AIResult[EvaluationOutput], error)
	SummarizeWeeklyProgress(ctx context.Context, input WeeklyProgressInput) (AIResult[WeeklyProgressOutput], error)

	// Stubs for Phase 8 mock interview module.
	GenerateMockInterviewQuestions(ctx context.Context, profile ProfileContext, focusAreas []string) (AIResult[MockInterviewQuestionsOutput], error)
	GradeMockInterview(ctx context.Context, transcript string) (AIResult[MockInterviewGradingOutput], error)
}

type service struct {
	client   *Client
	registry *Registry
}

func NewService(client *Client, registry *Registry) Service {
	return &service{client: client, registry: registry}
}

func (s *service) EvaluateBaselineDiagnosis(ctx context.Context, input BaselineDiagnosisInput) (AIResult[DiagnosticReport], error) {
	def, err := s.registry.Get(PromptBaselineDiagnosis)
	if err != nil {
		return AIResult[DiagnosticReport]{}, err
	}

	userJSON, err := json.Marshal(input)
	if err != nil {
		return AIResult[DiagnosticReport]{}, fmt.Errorf("marshal diagnosis input: %w", err)
	}

	raw, meta, err := s.client.CompleteStructured(ctx, StructuredRequest{
		Model:           def.Model,
		Instructions:    def.SystemInstruction,
		UserInput:       string(userJSON),
		SchemaName:      def.SchemaName,
		OutputSchema:    def.OutputSchema,
		Temperature:     def.Temperature,
		MaxOutputTokens: def.MaxOutputTokens,
		PromptID:        def.ID,
		PromptVersion:   def.Version,
	})
	if err != nil {
		return AIResult[DiagnosticReport]{}, err
	}

	validated, err := ParseDiagnosticReport(raw)
	if err != nil {
		return AIResult[DiagnosticReport]{}, err
	}
	validated.PromptVersion = def.Version
	validated.Model = meta.Model
	return AIResult[DiagnosticReport]{Data: validated, Raw: raw, Meta: meta}, nil
}

func (s *service) GenerateRecoveryPlan(ctx context.Context, input RecoveryPlanInput) (AIResult[RecoveryPlanOutput], error) {
	def, err := s.registry.Get(PromptRecoveryPlan)
	if err != nil {
		return AIResult[RecoveryPlanOutput]{}, err
	}

	userJSON, err := json.Marshal(input)
	if err != nil {
		return AIResult[RecoveryPlanOutput]{}, fmt.Errorf("marshal plan input: %w", err)
	}

	raw, meta, err := s.client.CompleteStructured(ctx, StructuredRequest{
		Model:           def.Model,
		Instructions:    def.SystemInstruction,
		UserInput:       string(userJSON),
		SchemaName:      def.SchemaName,
		OutputSchema:    def.OutputSchema,
		Temperature:     def.Temperature,
		MaxOutputTokens: def.MaxOutputTokens,
		PromptID:        def.ID,
		PromptVersion:   def.Version,
	})
	if err != nil {
		return AIResult[RecoveryPlanOutput]{}, err
	}

	validated, err := ParseRecoveryPlanOutput(raw)
	if err != nil {
		return AIResult[RecoveryPlanOutput]{}, err
	}
	validated.PromptVersion = def.Version
	validated.Model = meta.Model
	return AIResult[RecoveryPlanOutput]{Data: validated, Raw: raw, Meta: meta}, nil
}

func (s *service) EvaluateAttempt(ctx context.Context, questionType string, input EvaluationInput) (AIResult[EvaluationOutput], error) {
	promptID, err := evaluationPromptForType(questionType)
	if err != nil {
		return AIResult[EvaluationOutput]{}, err
	}

	def, err := s.registry.Get(promptID)
	if err != nil {
		return AIResult[EvaluationOutput]{}, err
	}

	userJSON, err := json.Marshal(input)
	if err != nil {
		return AIResult[EvaluationOutput]{}, fmt.Errorf("marshal evaluation input: %w", err)
	}

	raw, meta, err := s.client.CompleteStructured(ctx, StructuredRequest{
		Model:           def.Model,
		Instructions:    def.SystemInstruction,
		UserInput:       string(userJSON),
		SchemaName:      def.SchemaName,
		OutputSchema:    def.OutputSchema,
		Temperature:     def.Temperature,
		MaxOutputTokens: def.MaxOutputTokens,
		PromptID:        def.ID,
		PromptVersion:   def.Version,
	})
	if err != nil {
		return AIResult[EvaluationOutput]{}, err
	}

	validated, err := ParseEvaluationOutput(raw)
	if err != nil {
		return AIResult[EvaluationOutput]{}, err
	}
	validated.PromptVersion = def.Version
	validated.Model = meta.Model
	return AIResult[EvaluationOutput]{Data: validated, Raw: raw, Meta: meta}, nil
}

func (s *service) SummarizeWeeklyProgress(ctx context.Context, input WeeklyProgressInput) (AIResult[WeeklyProgressOutput], error) {
	def, err := s.registry.Get(PromptWeeklyProgressSummary)
	if err != nil {
		return AIResult[WeeklyProgressOutput]{}, err
	}

	userJSON, err := json.Marshal(input)
	if err != nil {
		return AIResult[WeeklyProgressOutput]{}, fmt.Errorf("marshal progress input: %w", err)
	}

	raw, meta, err := s.client.CompleteStructured(ctx, StructuredRequest{
		Model:           def.Model,
		Instructions:    def.SystemInstruction,
		UserInput:       string(userJSON),
		SchemaName:      def.SchemaName,
		OutputSchema:    def.OutputSchema,
		Temperature:     def.Temperature,
		MaxOutputTokens: def.MaxOutputTokens,
		PromptID:        def.ID,
		PromptVersion:   def.Version,
	})
	if err != nil {
		return AIResult[WeeklyProgressOutput]{}, err
	}

	validated, err := ParseWeeklyProgressOutput(raw)
	if err != nil {
		return AIResult[WeeklyProgressOutput]{}, err
	}
	validated.PromptVersion = def.Version
	validated.Model = meta.Model
	return AIResult[WeeklyProgressOutput]{Data: validated, Raw: raw, Meta: meta}, nil
}

func (s *service) GenerateMockInterviewQuestions(_ context.Context, _ ProfileContext, _ []string) (AIResult[MockInterviewQuestionsOutput], error) {
	return AIResult[MockInterviewQuestionsOutput]{}, ErrNotImplemented
}

func (s *service) GradeMockInterview(_ context.Context, _ string) (AIResult[MockInterviewGradingOutput], error) {
	return AIResult[MockInterviewGradingOutput]{}, ErrNotImplemented
}

func evaluationPromptForType(questionType string) (string, error) {
	switch questionType {
	case "programming_basics":
		return PromptProgrammingEval, nil
	case "dsa":
		return PromptDSAEval, nil
	case "concept":
		return PromptTechnicalAnswerEval, nil
	default:
		return "", fmt.Errorf("unsupported question type %q", questionType)
	}
}
