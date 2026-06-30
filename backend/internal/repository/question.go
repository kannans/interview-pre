package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kannans/interview-prep/backend/internal/db"
)

type QuestionListParams struct {
	Type     pgtype.Text
	Category pgtype.Text
	Level    pgtype.Text
}

type QuestionRepository interface {
	List(ctx context.Context, params QuestionListParams) ([]db.Question, error)
	GetByID(ctx context.Context, id pgtype.UUID) (db.Question, error)
}

type questionRepository struct {
	queries db.Querier
}

func NewQuestionRepository(queries db.Querier) QuestionRepository {
	return &questionRepository{queries: queries}
}

func (r *questionRepository) List(ctx context.Context, params QuestionListParams) ([]db.Question, error) {
	questions, err := r.queries.ListQuestions(ctx, db.ListQuestionsParams{
		FilterType:     params.Type,
		FilterCategory: params.Category,
		FilterLevel:    params.Level,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list questions: %w", err)
	}
	return questions, nil
}

func (r *questionRepository) GetByID(ctx context.Context, id pgtype.UUID) (db.Question, error) {
	question, err := r.queries.GetQuestion(ctx, id)
	if err != nil {
		return db.Question{}, fmt.Errorf("failed to get question: %w", err)
	}
	return question, nil
}

var _ QuestionRepository = (*questionRepository)(nil)
