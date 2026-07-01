package service_test

import (
	"context"
	"testing"

	"github.com/kannans/interview-prep/backend/internal/runner"
	"github.com/kannans/interview-prep/backend/internal/service"
)

type stubExecutor struct {
	result runner.RunResult
	err    error
}

func (s stubExecutor) Run(_ context.Context, _ runner.RunInput) (runner.RunResult, error) {
	return s.result, s.err
}

func TestRunServiceValidate(t *testing.T) {
	svc := service.NewRunService(stubExecutor{})

	_, err := svc.Run(context.Background(), service.RunInput{})
	if err == nil {
		t.Fatal("expected validation error for empty code")
	}

	_, err = svc.Run(context.Background(), service.RunInput{
		Code:     "print('hi')",
		Language: "rust",
	})
	if err == nil {
		t.Fatal("expected validation error for unsupported language")
	}
}

func TestRunServiceDelegates(t *testing.T) {
	svc := service.NewRunService(stubExecutor{
		result: runner.RunResult{Stdout: "ok", ExitCode: 0},
	})

	result, err := svc.Run(context.Background(), service.RunInput{
		Code:     "console.log('ok')",
		Language: "javascript",
	})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if result.Stdout != "ok" {
		t.Fatalf("stdout=%q", result.Stdout)
	}
}
