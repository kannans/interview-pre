package service_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kannans/interview-prep/backend/internal/ai"
	"github.com/kannans/interview-prep/backend/internal/db"
	"github.com/kannans/interview-prep/backend/internal/repository"
	"github.com/kannans/interview-prep/backend/internal/service"
	"github.com/kannans/interview-prep/backend/internal/testutil"
)

func TestPlanGenerateSuccess(t *testing.T) {
	ctx := context.Background()
	profileID := testutil.MustUUID(t, testutil.TestUUID)
	planID := testutil.MustUUID(t, "d4000004-0000-4000-8000-000000000004")
	profile := testutil.SampleProfile(profileID)

	profileRepo := &testutil.MockProfileRepository{
		GetFn: func(context.Context) (db.Profile, error) { return profile, nil },
	}
	diagnosisRepo := &testutil.MockDiagnosisRepository{
		GetFn: func(context.Context, pgtype.UUID) (db.BaselineDiagnosis, error) {
			return db.BaselineDiagnosis{DiagnosticReport: json.RawMessage(`{"weak_areas":["dsa"]}`)}, nil
		},
	}
	planRepo := &testutil.MockPlanRepository{
		ReplacePlanWithTasksFn: func(_ context.Context, _ pgtype.UUID, tasks []repository.PlanTaskInput) (db.RecoveryPlan, []db.PlanTask, error) {
			if len(tasks) != 2 {
				t.Fatalf("expected 2 tasks, got %d", len(tasks))
			}
			return db.RecoveryPlan{ID: planID}, []db.PlanTask{{Title: tasks[0].Title}}, nil
		},
	}
	aiSvc := &testutil.MockAIService{
		GenerateRecoveryPlanFn: func(context.Context, ai.RecoveryPlanInput) (ai.AIResult[ai.RecoveryPlanOutput], error) {
			return ai.AIResult[ai.RecoveryPlanOutput]{
				Data: ai.RecoveryPlanOutput{
					Tasks: []ai.PlanTaskSuggestion{
						{Category: "dsa", Level: "level_1", Title: "Arrays"},
						{Category: "concept", Level: "level_1", Title: "HTTP"},
					},
				},
			}, nil
		},
	}

	result, err := service.NewPlanService(profileRepo, diagnosisRepo, planRepo, aiSvc).Generate(ctx)
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	if len(result.Tasks) != 1 {
		t.Fatalf("expected 1 created task, got %d", len(result.Tasks))
	}
}

func TestPlanGenerateDiagnosisMissing(t *testing.T) {
	profileID := testutil.MustUUID(t, testutil.TestUUID)
	profileRepo := &testutil.MockProfileRepository{
		GetFn: func(context.Context) (db.Profile, error) { return testutil.SampleProfile(profileID), nil },
	}
	diagnosisRepo := &testutil.MockDiagnosisRepository{
		GetFn: func(context.Context, pgtype.UUID) (db.BaselineDiagnosis, error) {
			return db.BaselineDiagnosis{}, pgx.ErrNoRows
		},
	}

	_, err := service.NewPlanService(profileRepo, diagnosisRepo, &testutil.MockPlanRepository{}, &testutil.MockAIService{}).Generate(context.Background())
	if !errors.Is(err, service.ErrDiagnosisMissing) {
		t.Fatalf("expected ErrDiagnosisMissing, got %v", err)
	}
}

func TestPlanGetNotFound(t *testing.T) {
	profileID := testutil.MustUUID(t, testutil.TestUUID)
	profileRepo := &testutil.MockProfileRepository{
		GetFn: func(context.Context) (db.Profile, error) { return testutil.SampleProfile(profileID), nil },
	}
	planRepo := &testutil.MockPlanRepository{
		GetActiveFn: func(context.Context, pgtype.UUID) (db.RecoveryPlan, error) {
			return db.RecoveryPlan{}, pgx.ErrNoRows
		},
	}

	_, err := service.NewPlanService(profileRepo, &testutil.MockDiagnosisRepository{}, planRepo, &testutil.MockAIService{}).Get(context.Background())
	if !errors.Is(err, service.ErrPlanNotFound) {
		t.Fatalf("expected ErrPlanNotFound, got %v", err)
	}
}

func TestPlanGetTodayExhausted(t *testing.T) {
	profileID := testutil.MustUUID(t, testutil.TestUUID)
	planID := testutil.MustUUID(t, "d4000004-0000-4000-8000-000000000004")
	profileRepo := &testutil.MockProfileRepository{
		GetFn: func(context.Context) (db.Profile, error) { return testutil.SampleProfile(profileID), nil },
	}
	planRepo := &testutil.MockPlanRepository{
		GetActiveFn: func(context.Context, pgtype.UUID) (db.RecoveryPlan, error) {
			return db.RecoveryPlan{ID: planID}, nil
		},
		ListTasksFn: func(context.Context, pgtype.UUID) ([]db.PlanTask, error) {
			return []db.PlanTask{{IsCompleted: true}}, nil
		},
		GetTodayTaskFn: func(context.Context, pgtype.UUID) (db.PlanTask, error) {
			return db.PlanTask{}, pgx.ErrNoRows
		},
	}

	_, err := service.NewPlanService(profileRepo, &testutil.MockDiagnosisRepository{}, planRepo, &testutil.MockAIService{}).GetToday(context.Background())
	if !errors.Is(err, service.ErrNoTasksRemaining) {
		t.Fatalf("expected ErrNoTasksRemaining, got %v", err)
	}
}
