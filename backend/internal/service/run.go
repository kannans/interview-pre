package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/kannans/interview-prep/backend/internal/runner"
)

// RunInput is the API payload for executing practice code locally.
type RunInput struct {
	Language string `json:"language"`
	Code     string `json:"code"`
	Stdin    string `json:"stdin,omitempty"`
}

func (i RunInput) Validate() error {
	if strings.TrimSpace(i.Code) == "" {
		return fmt.Errorf("validation error: code is required")
	}
	if _, err := NormalizePracticeLanguage(i.Language); err != nil {
		return err
	}
	return nil
}

func (i RunInput) ResolvedLanguage() string {
	lang, _ := NormalizePracticeLanguage(i.Language)
	return lang
}

// RunService executes user code for local debugging before AI evaluation.
type RunService interface {
	Run(ctx context.Context, input RunInput) (runner.RunResult, error)
}

type runService struct {
	executor runner.Executor
}

func NewRunService(executor runner.Executor) RunService {
	return &runService{executor: executor}
}

func (s *runService) Run(ctx context.Context, input RunInput) (runner.RunResult, error) {
	if err := input.Validate(); err != nil {
		return runner.RunResult{}, err
	}

	result, err := s.executor.Run(ctx, runner.RunInput{
		Language: input.ResolvedLanguage(),
		Code:     input.Code,
		Stdin:    input.Stdin,
	})
	if err != nil {
		return runner.RunResult{}, fmt.Errorf("failed to run code: %w", err)
	}

	return result, nil
}
