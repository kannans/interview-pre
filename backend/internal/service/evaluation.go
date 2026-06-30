package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kannans/interview-prep/backend/internal/ai"
	"github.com/kannans/interview-prep/backend/internal/db"
	"github.com/kannans/interview-prep/backend/internal/repository"
)

var (
	ErrEvaluationNotFound      = errors.New("evaluation not found")
	ErrEvaluationAlreadyExists = errors.New("evaluation already exists for this attempt")
)

// EvaluationService orchestrates attempt scoring and feedback retrieval.
type EvaluationService interface {
	Evaluate(ctx context.Context, attemptID pgtype.UUID) (db.Evaluation, error)
	GetByAttemptID(ctx context.Context, attemptID pgtype.UUID) (db.Evaluation, error)
}

type evaluationService struct {
	profileRepo    repository.ProfileRepository
	attemptRepo    repository.AttemptRepository
	questionRepo   repository.QuestionRepository
	evaluationRepo repository.EvaluationRepository
	aiService      ai.Service
}

func NewEvaluationService(
	profileRepo repository.ProfileRepository,
	attemptRepo repository.AttemptRepository,
	questionRepo repository.QuestionRepository,
	evaluationRepo repository.EvaluationRepository,
	aiService ai.Service,
) EvaluationService {
	return &evaluationService{
		profileRepo:    profileRepo,
		attemptRepo:    attemptRepo,
		questionRepo:   questionRepo,
		evaluationRepo: evaluationRepo,
		aiService:      aiService,
	}
}

func (s *evaluationService) Evaluate(ctx context.Context, attemptID pgtype.UUID) (db.Evaluation, error) {
	if !attemptID.Valid {
		return db.Evaluation{}, fmt.Errorf("validation error: invalid attempt id")
	}

	attempt, err := s.loadOwnedAttempt(ctx, attemptID)
	if err != nil {
		return db.Evaluation{}, err
	}

	_, err = s.evaluationRepo.GetByAttemptID(ctx, attemptID)
	if err == nil {
		return db.Evaluation{}, ErrEvaluationAlreadyExists
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return db.Evaluation{}, fmt.Errorf("failed to check existing evaluation: %w", err)
	}

	question, err := s.questionRepo.GetByID(ctx, attempt.QuestionID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.Evaluation{}, ErrQuestionNotFound
		}
		return db.Evaluation{}, fmt.Errorf("failed to get question: %w", err)
	}

	profile, err := s.profileRepo.Get(ctx)
	if err != nil {
		return db.Evaluation{}, fmt.Errorf("failed to get profile: %w", err)
	}

	evalInput := ai.EvaluationInput{
		Profile:          profileToAIContext(profile),
		QuestionType:     question.Type,
		QuestionCategory: question.Category,
		QuestionLevel:    question.Level,
		QuestionTitle:    question.Title,
		QuestionContent:  question.Content,
		SubmittedContent: attempt.SubmittedContent,
		TimeTakenSeconds: attempt.TimeTakenSeconds,
	}
	if question.Constraints.Valid {
		evalInput.Constraints = question.Constraints.String
	}
	if question.StarterCode.Valid {
		evalInput.StarterCode = question.StarterCode.String
	}
	if question.SolutionHint.Valid {
		evalInput.SolutionHint = question.SolutionHint.String
	}

	result, err := s.aiService.EvaluateAttempt(ctx, question.Type, evalInput)
	if err != nil {
		return db.Evaluation{}, fmt.Errorf("%w: %v", ErrAIFailed, err)
	}

	out := result.Data
	rawAIResponse, err := json.Marshal(map[string]any{
		"evaluation":     out,
		"prompt_version": out.PromptVersion,
		"model":          out.Model,
		"usage": map[string]int{
			"input_tokens":  result.Meta.InputTokens,
			"output_tokens": result.Meta.OutputTokens,
			"total_tokens":  result.Meta.TotalTokens,
		},
		"response_id": result.Meta.ResponseID,
	})
	if err != nil {
		return db.Evaluation{}, fmt.Errorf("failed to marshal raw ai response: %w", err)
	}

	params := db.CreateEvaluationParams{
		AttemptID:             attemptID,
		ScoreCorrectness:      int32(out.ScoreCorrectness),
		ScoreApproach:         int32(out.ScoreApproach),
		ScoreEdgeCases:        int32(out.ScoreEdgeCases),
		ScoreCommunication:    int32(out.ScoreCommunication),
		FeedbackCorrectness:   out.FeedbackCorrectness,
		FeedbackApproach:      out.FeedbackApproach,
		FeedbackEdgeCases:     out.FeedbackEdgeCases,
		FeedbackCommunication: out.FeedbackCommunication,
		RawAiResponse:         rawAIResponse,
	}

	evaluation, err := s.evaluationRepo.Create(ctx, params)
	if err != nil {
		return db.Evaluation{}, fmt.Errorf("failed to create evaluation: %w", err)
	}

	return evaluation, nil
}

func (s *evaluationService) GetByAttemptID(ctx context.Context, attemptID pgtype.UUID) (db.Evaluation, error) {
	if !attemptID.Valid {
		return db.Evaluation{}, fmt.Errorf("validation error: invalid attempt id")
	}

	if _, err := s.loadOwnedAttempt(ctx, attemptID); err != nil {
		return db.Evaluation{}, err
	}

	evaluation, err := s.evaluationRepo.GetByAttemptID(ctx, attemptID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.Evaluation{}, ErrEvaluationNotFound
		}
		return db.Evaluation{}, fmt.Errorf("failed to get evaluation: %w", err)
	}

	return evaluation, nil
}

func (s *evaluationService) loadOwnedAttempt(ctx context.Context, attemptID pgtype.UUID) (db.Attempt, error) {
	profile, err := s.profileRepo.Get(ctx)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.Attempt{}, ErrProfileMissing
		}
		return db.Attempt{}, fmt.Errorf("failed to get profile: %w", err)
	}

	attempt, err := s.attemptRepo.GetByID(ctx, attemptID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.Attempt{}, ErrAttemptNotFound
		}
		return db.Attempt{}, fmt.Errorf("failed to get attempt: %w", err)
	}

	if !bytes.Equal(attempt.ProfileID.Bytes[:], profile.ID.Bytes[:]) {
		return db.Attempt{}, ErrAttemptNotOwned
	}

	return attempt, nil
}
