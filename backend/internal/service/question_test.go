package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kannans/interview-prep/backend/internal/db"
	"github.com/kannans/interview-prep/backend/internal/repository"
	"github.com/kannans/interview-prep/backend/internal/service"
	"github.com/kannans/interview-prep/backend/internal/testutil"
)

func TestQuestionListFilterValidation(t *testing.T) {
	_, err := service.NewQuestionService(&testutil.MockQuestionRepository{
		ListFn: func(context.Context, repository.QuestionListParams) ([]db.Question, error) {
			return nil, nil
		},
	}).List(context.Background(), service.QuestionListFilter{Type: "invalid"})
	if err == nil {
		t.Fatal("expected validation error for invalid type")
	}
}

func TestQuestionGetByIDNotFound(t *testing.T) {
	qid := testutil.MustUUID(t, testutil.TestUUID)
	repo := &testutil.MockQuestionRepository{
		GetByIDFn: func(context.Context, pgtype.UUID) (db.Question, error) {
			return db.Question{}, pgx.ErrNoRows
		},
	}
	_, err := service.NewQuestionService(repo).GetByID(context.Background(), qid)
	if !errors.Is(err, service.ErrQuestionNotFound) {
		t.Fatalf("expected ErrQuestionNotFound, got %v", err)
	}
}

func TestQuestionListSuccess(t *testing.T) {
	repo := &testutil.MockQuestionRepository{
		ListFn: func(context.Context, repository.QuestionListParams) ([]db.Question, error) {
			return []db.Question{{Title: "Q1"}}, nil
		},
	}
	result, err := service.NewQuestionService(repo).List(context.Background(), service.QuestionListFilter{Type: "dsa"})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(result.Questions) != 1 {
		t.Fatalf("expected 1 question, got %d", len(result.Questions))
	}
}

func TestQuestionGetByIDInvalidUUID(t *testing.T) {
	_, err := service.NewQuestionService(&testutil.MockQuestionRepository{}).GetByID(context.Background(), pgtype.UUID{})
	if err == nil {
		t.Fatal("expected validation error")
	}
}
