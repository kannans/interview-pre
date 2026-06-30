package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kannans/interview-prep/backend/internal/db"
)

type EvaluationRepository interface {
	Create(ctx context.Context, params db.CreateEvaluationParams) (db.Evaluation, error)
	GetByAttemptID(ctx context.Context, attemptID pgtype.UUID) (db.Evaluation, error)
	GetByID(ctx context.Context, id pgtype.UUID) (db.Evaluation, error)
}

type evaluationRepository struct {
	queries db.Querier
}

func NewEvaluationRepository(queries db.Querier) EvaluationRepository {
	return &evaluationRepository{queries: queries}
}

func (r *evaluationRepository) Create(ctx context.Context, params db.CreateEvaluationParams) (db.Evaluation, error) {
	evaluation, err := r.queries.CreateEvaluation(ctx, params)
	if err != nil {
		return db.Evaluation{}, fmt.Errorf("failed to create evaluation: %w", err)
	}
	return evaluation, nil
}

func (r *evaluationRepository) GetByAttemptID(ctx context.Context, attemptID pgtype.UUID) (db.Evaluation, error) {
	evaluation, err := r.queries.GetEvaluationByAttemptID(ctx, attemptID)
	if err != nil {
		return db.Evaluation{}, fmt.Errorf("failed to get evaluation by attempt id: %w", err)
	}
	return evaluation, nil
}

func (r *evaluationRepository) GetByID(ctx context.Context, id pgtype.UUID) (db.Evaluation, error) {
	evaluation, err := r.queries.GetEvaluationByID(ctx, id)
	if err != nil {
		return db.Evaluation{}, fmt.Errorf("failed to get evaluation: %w", err)
	}
	return evaluation, nil
}

var _ EvaluationRepository = (*evaluationRepository)(nil)
