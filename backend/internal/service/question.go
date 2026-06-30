package service

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kannans/interview-prep/backend/internal/db"
	"github.com/kannans/interview-prep/backend/internal/repository"
)

var ErrQuestionNotFound = errors.New("question not found")

var (
	validQuestionTypes = []string{"programming_basics", "dsa", "concept"}
	validQuestionLevels = []string{"level_1", "level_2", "level_3"}
	validQuestionCategories = []string{
		"loops", "functions", "debugging",
		"arrays", "strings", "hash_maps", "two_pointer", "recursion", "trees",
		"complexity", "http", "sql",
	}
)

// QuestionListFilter holds optional query filters for listing questions.
type QuestionListFilter struct {
	Type     string
	Category string
	Level    string
}

// ListQuestionsResponse is the API response shape for a question list.
type ListQuestionsResponse struct {
	Questions []db.Question `json:"questions"`
}

// QuestionService provides read access to the practice question catalog.
type QuestionService interface {
	List(ctx context.Context, filter QuestionListFilter) (ListQuestionsResponse, error)
	GetByID(ctx context.Context, id pgtype.UUID) (db.Question, error)
}

type questionService struct {
	repo repository.QuestionRepository
}

func NewQuestionService(repo repository.QuestionRepository) QuestionService {
	return &questionService{repo: repo}
}

func (s *questionService) List(ctx context.Context, filter QuestionListFilter) (ListQuestionsResponse, error) {
	params, err := filter.toRepositoryParams()
	if err != nil {
		return ListQuestionsResponse{}, err
	}

	questions, err := s.repo.List(ctx, params)
	if err != nil {
		return ListQuestionsResponse{}, fmt.Errorf("failed to list questions: %w", err)
	}

	if questions == nil {
		questions = []db.Question{}
	}

	return ListQuestionsResponse{Questions: questions}, nil
}

func (s *questionService) GetByID(ctx context.Context, id pgtype.UUID) (db.Question, error) {
	if !id.Valid {
		return db.Question{}, fmt.Errorf("validation error: invalid question id")
	}

	question, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.Question{}, ErrQuestionNotFound
		}
		return db.Question{}, fmt.Errorf("failed to get question: %w", err)
	}

	return question, nil
}

func (f QuestionListFilter) toRepositoryParams() (repository.QuestionListParams, error) {
	params := repository.QuestionListParams{}

	if f.Type != "" {
		if !slices.Contains(validQuestionTypes, f.Type) {
			return repository.QuestionListParams{}, fmt.Errorf("validation error: type must be one of: programming_basics, dsa, concept")
		}
		params.Type = pgtype.Text{String: f.Type, Valid: true}
	}

	if f.Category != "" {
		if !slices.Contains(validQuestionCategories, f.Category) {
			return repository.QuestionListParams{}, fmt.Errorf("validation error: category must be one of: loops, functions, debugging, arrays, strings, hash_maps, two_pointer, recursion, trees, complexity, http, sql")
		}
		params.Category = pgtype.Text{String: f.Category, Valid: true}
	}

	if f.Level != "" {
		if !slices.Contains(validQuestionLevels, f.Level) {
			return repository.QuestionListParams{}, fmt.Errorf("validation error: level must be one of: level_1, level_2, level_3")
		}
		params.Level = pgtype.Text{String: f.Level, Valid: true}
	}

	return params, nil
}
