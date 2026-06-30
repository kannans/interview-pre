package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kannans/interview-prep/backend/internal/db"
)

type AttemptListParams struct {
	ProfileID        pgtype.UUID
	FilterQuestionID pgtype.UUID
}

type AttemptRepository interface {
	Create(ctx context.Context, profileID, questionID pgtype.UUID, submittedContent string, timeTakenSeconds int32, language string) (db.Attempt, error)
	GetByID(ctx context.Context, id pgtype.UUID) (db.Attempt, error)
	List(ctx context.Context, params AttemptListParams) ([]db.Attempt, error)
}

type attemptRepository struct {
	queries db.Querier
}

func NewAttemptRepository(queries db.Querier) AttemptRepository {
	return &attemptRepository{queries: queries}
}

func (r *attemptRepository) Create(ctx context.Context, profileID, questionID pgtype.UUID, submittedContent string, timeTakenSeconds int32, language string) (db.Attempt, error) {
	attempt, err := r.queries.CreateAttempt(ctx, db.CreateAttemptParams{
		ProfileID:        profileID,
		QuestionID:       questionID,
		SubmittedContent: submittedContent,
		TimeTakenSeconds: timeTakenSeconds,
		Language:         language,
	})
	if err != nil {
		return db.Attempt{}, fmt.Errorf("failed to create attempt: %w", err)
	}
	return attempt, nil
}

func (r *attemptRepository) GetByID(ctx context.Context, id pgtype.UUID) (db.Attempt, error) {
	attempt, err := r.queries.GetAttempt(ctx, id)
	if err != nil {
		return db.Attempt{}, fmt.Errorf("failed to get attempt: %w", err)
	}
	return attempt, nil
}

func (r *attemptRepository) List(ctx context.Context, params AttemptListParams) ([]db.Attempt, error) {
	attempts, err := r.queries.ListAttempts(ctx, db.ListAttemptsParams{
		ProfileID:        params.ProfileID,
		FilterQuestionID: params.FilterQuestionID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list attempts: %w", err)
	}
	return attempts, nil
}

var _ AttemptRepository = (*attemptRepository)(nil)
