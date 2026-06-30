package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kannans/interview-prep/backend/internal/db"
	"github.com/kannans/interview-prep/backend/internal/repository"
)

var (
	ErrEvaluationNotFound      = errors.New("evaluation not found")
	ErrEvaluationAlreadyExists = errors.New("evaluation already exists for this attempt")
)

// EvaluationService orchestrates attempt scoring and feedback retrieval.
// In Phase 7, Evaluate will trigger a real OpenAI call.
type EvaluationService interface {
	Evaluate(ctx context.Context, attemptID pgtype.UUID) (db.Evaluation, error)
	GetByAttemptID(ctx context.Context, attemptID pgtype.UUID) (db.Evaluation, error)
}

type evaluationService struct {
	profileRepo    repository.ProfileRepository
	attemptRepo    repository.AttemptRepository
	questionRepo   repository.QuestionRepository
	evaluationRepo repository.EvaluationRepository
}

func NewEvaluationService(
	profileRepo repository.ProfileRepository,
	attemptRepo repository.AttemptRepository,
	questionRepo repository.QuestionRepository,
	evaluationRepo repository.EvaluationRepository,
) EvaluationService {
	return &evaluationService{
		profileRepo:    profileRepo,
		attemptRepo:    attemptRepo,
		questionRepo:   questionRepo,
		evaluationRepo: evaluationRepo,
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

	scores, feedback, rawAIResponse := derivePlaceholderEvaluation(attempt, question)

	params := db.CreateEvaluationParams{
		AttemptID:             attemptID,
		ScoreCorrectness:      scores.correctness,
		ScoreApproach:         scores.approach,
		ScoreEdgeCases:        scores.edgeCases,
		ScoreCommunication:    scores.communication,
		FeedbackCorrectness:   feedback.correctness,
		FeedbackApproach:      feedback.approach,
		FeedbackEdgeCases:     feedback.edgeCases,
		FeedbackCommunication: feedback.communication,
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

type placeholderScores struct {
	correctness   int32
	approach      int32
	edgeCases     int32
	communication int32
}

type placeholderFeedback struct {
	correctness   string
	approach      string
	edgeCases     string
	communication string
}

func derivePlaceholderEvaluation(attempt db.Attempt, question db.Question) (placeholderScores, placeholderFeedback, json.RawMessage) {
	contentLen := len(strings.TrimSpace(attempt.SubmittedContent))
	timeTaken := attempt.TimeTakenSeconds

	correctness := clampScore(int32(4 + contentLen/200))
	approach := scoreFromTimeTaken(timeTaken)
	edgeCases := clampScore(int32(3 + contentLen/400))
	communication := clampScore(int32(4 + strings.Count(attempt.SubmittedContent, "\n")))

	feedback := placeholderFeedback{
		correctness: fmt.Sprintf(
			"Placeholder: review your solution against the %s question \"%s\". Submission length suggests %s detail.",
			question.Type,
			question.Title,
			detailLevel(contentLen),
		),
		approach: fmt.Sprintf(
			"Placeholder: you spent %d seconds. %s",
			timeTaken,
			timeTakenFeedback(timeTaken),
		),
		edgeCases: fmt.Sprintf(
			"Placeholder: explicitly test empty input, boundary values, and invalid cases for \"%s\".",
			question.Title,
		),
		communication: fmt.Sprintf(
			"Placeholder: structure your explanation with problem restatement, approach, complexity, and trade-offs for %s questions.",
			question.Category,
		),
	}

	rawPayload := map[string]interface{}{
		"status":              "pending_ai_evaluation",
		"message":             "Real AI evaluation will be triggered in Phase 7",
		"placeholder":         true,
		"attempt_id":          uuidString(attempt.ID),
		"question_id":         uuidString(question.ID),
		"question_type":       question.Type,
		"time_taken_seconds":  timeTaken,
		"submission_length":   contentLen,
		"scores_placeholder":  map[string]int32{"correctness": correctness, "approach": approach, "edge_cases": edgeCases, "communication": communication},
	}
	rawAIResponse, _ := json.Marshal(rawPayload)

	return placeholderScores{
		correctness:   correctness,
		approach:      approach,
		edgeCases:     edgeCases,
		communication: communication,
	}, feedback, rawAIResponse
}

func clampScore(score int32) int32 {
	if score < 0 {
		return 0
	}
	if score > 10 {
		return 10
	}
	return score
}

func scoreFromTimeTaken(seconds int32) int32 {
	switch {
	case seconds < 60:
		return 4
	case seconds <= 900:
		return 7
	case seconds <= 1800:
		return 8
	default:
		return 6
	}
}

func detailLevel(contentLen int) string {
	switch {
	case contentLen < 100:
		return "minimal"
	case contentLen < 500:
		return "moderate"
	default:
		return "strong"
	}
}

func timeTakenFeedback(seconds int32) string {
	switch {
	case seconds < 60:
		return "Consider spending more time outlining your approach before coding."
	case seconds <= 900:
		return "Time usage looks reasonable for a focused practice session."
	case seconds <= 1800:
		return "Good depth; watch for diminishing returns on very long sessions."
	default:
		return "Session ran long — practice articulating a concise plan upfront."
	}
}

func uuidString(id pgtype.UUID) string {
	if !id.Valid {
		return ""
	}
	return fmt.Sprintf(
		"%x-%x-%x-%x-%x",
		id.Bytes[0:4],
		id.Bytes[4:6],
		id.Bytes[6:8],
		id.Bytes[8:10],
		id.Bytes[10:16],
	)
}
