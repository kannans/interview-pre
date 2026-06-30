package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kannans/interview-prep/backend/internal/db"
	"github.com/kannans/interview-prep/backend/internal/repository"
)

var (
	ErrAttemptNotFound = errors.New("attempt not found")
	ErrAttemptNotOwned = errors.New("attempt does not belong to your profile")
)

// AttemptInput holds the payload for submitting a practice attempt.
type AttemptInput struct {
	QuestionID       string `json:"question_id"`
	SubmittedContent string `json:"submitted_content"`
	TimeTakenSeconds int32  `json:"time_taken_seconds"`
}

func (i AttemptInput) Validate() error {
	if strings.TrimSpace(i.QuestionID) == "" {
		return errors.New("question_id is required")
	}
	if strings.TrimSpace(i.SubmittedContent) == "" {
		return errors.New("submitted_content is required")
	}
	if i.TimeTakenSeconds <= 0 {
		return errors.New("time_taken_seconds must be greater than zero")
	}
	return nil
}

// ListAttemptsResponse is the API response shape for listing attempts.
type ListAttemptsResponse struct {
	Attempts []db.Attempt `json:"attempts"`
}

// AttemptService orchestrates practice attempt submission and retrieval.
type AttemptService interface {
	Submit(ctx context.Context, input AttemptInput) (db.Attempt, error)
	List(ctx context.Context, questionIDStr string) (ListAttemptsResponse, error)
	GetByID(ctx context.Context, id pgtype.UUID) (db.Attempt, error)
}

type attemptService struct {
	profileRepo  repository.ProfileRepository
	questionRepo repository.QuestionRepository
	attemptRepo  repository.AttemptRepository
}

func NewAttemptService(
	profileRepo repository.ProfileRepository,
	questionRepo repository.QuestionRepository,
	attemptRepo repository.AttemptRepository,
) AttemptService {
	return &attemptService{
		profileRepo:  profileRepo,
		questionRepo: questionRepo,
		attemptRepo:  attemptRepo,
	}
}

func (s *attemptService) Submit(ctx context.Context, input AttemptInput) (db.Attempt, error) {
	if err := input.Validate(); err != nil {
		return db.Attempt{}, fmt.Errorf("validation error: %w", err)
	}

	var questionID pgtype.UUID
	if err := questionID.Scan(input.QuestionID); err != nil {
		return db.Attempt{}, fmt.Errorf("validation error: invalid question_id")
	}

	profile, err := s.profileRepo.Get(ctx)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.Attempt{}, ErrProfileMissing
		}
		return db.Attempt{}, fmt.Errorf("failed to get profile: %w", err)
	}

	_, err = s.questionRepo.GetByID(ctx, questionID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.Attempt{}, ErrQuestionNotFound
		}
		return db.Attempt{}, fmt.Errorf("failed to get question: %w", err)
	}

	attempt, err := s.attemptRepo.Create(ctx, profile.ID, questionID, input.SubmittedContent, input.TimeTakenSeconds)
	if err != nil {
		return db.Attempt{}, fmt.Errorf("failed to submit attempt: %w", err)
	}

	return attempt, nil
}

func (s *attemptService) List(ctx context.Context, questionIDStr string) (ListAttemptsResponse, error) {
	profile, err := s.profileRepo.Get(ctx)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ListAttemptsResponse{}, ErrProfileMissing
		}
		return ListAttemptsResponse{}, fmt.Errorf("failed to get profile: %w", err)
	}

	params := repository.AttemptListParams{ProfileID: profile.ID}

	if questionIDStr != "" {
		var questionID pgtype.UUID
		if err := questionID.Scan(questionIDStr); err != nil {
			return ListAttemptsResponse{}, fmt.Errorf("validation error: invalid question_id")
		}
		params.FilterQuestionID = questionID
	}

	attempts, err := s.attemptRepo.List(ctx, params)
	if err != nil {
		return ListAttemptsResponse{}, fmt.Errorf("failed to list attempts: %w", err)
	}

	if attempts == nil {
		attempts = []db.Attempt{}
	}

	return ListAttemptsResponse{Attempts: attempts}, nil
}

func (s *attemptService) GetByID(ctx context.Context, id pgtype.UUID) (db.Attempt, error) {
	if !id.Valid {
		return db.Attempt{}, fmt.Errorf("validation error: invalid attempt id")
	}

	profile, err := s.profileRepo.Get(ctx)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.Attempt{}, ErrProfileMissing
		}
		return db.Attempt{}, fmt.Errorf("failed to get profile: %w", err)
	}

	attempt, err := s.attemptRepo.GetByID(ctx, id)
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
