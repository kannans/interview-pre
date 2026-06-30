package service_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kannans/interview-prep/backend/internal/ai"
	"github.com/kannans/interview-prep/backend/internal/db"
	"github.com/kannans/interview-prep/backend/internal/repository"
	"github.com/kannans/interview-prep/backend/internal/service"
	"github.com/kannans/interview-prep/backend/internal/testutil"
)

func TestProgressCreateSnapshotSuccess(t *testing.T) {
	ctx := context.Background()
	profileID := testutil.MustUUID(t, testutil.TestUUID)
	profile := testutil.SampleProfile(profileID)

	profileRepo := &testutil.MockProfileRepository{
		GetFn: func(context.Context) (db.Profile, error) { return profile, nil },
	}
	planRepo := &testutil.MockPlanRepository{
		GetActiveFn: func(context.Context, pgtype.UUID) (db.RecoveryPlan, error) {
			return db.RecoveryPlan{}, pgx.ErrNoRows
		},
	}
	progressRepo := &testutil.MockProgressRepository{
		GetByProfileAndDateFn: func(context.Context, pgtype.UUID, pgtype.Date) (db.ProgressSnapshot, error) {
			return db.ProgressSnapshot{}, pgx.ErrNoRows
		},
		CountAttemptsFn: func(context.Context, pgtype.UUID) (int32, error) { return 5, nil },
		GetEvaluationStatsFn: func(context.Context, pgtype.UUID) (repository.EvaluationStats, error) {
			return repository.EvaluationStats{EvaluationsCount: 3, AvgCorrectness: 7}, nil
		},
		CreateFn: func(_ context.Context, _ pgtype.UUID, _ pgtype.Date, stats []byte) (db.ProgressSnapshot, error) {
			var parsed service.SnapshotStats
			if err := json.Unmarshal(stats, &parsed); err != nil {
				t.Fatalf("unmarshal stats: %v", err)
			}
			if parsed.Status != "complete" {
				t.Fatalf("expected status complete, got %q", parsed.Status)
			}
			return db.ProgressSnapshot{ID: profileID, Stats: stats}, nil
		},
	}
	aiSvc := &testutil.MockAIService{
		SummarizeWeeklyProgressFn: func(context.Context, ai.WeeklyProgressInput) (ai.AIResult[ai.WeeklyProgressOutput], error) {
			return ai.AIResult[ai.WeeklyProgressOutput]{
				Data: ai.WeeklyProgressOutput{
					WeeklySummary:            "Good week",
					FocusAreas:               []string{"edge cases"},
					AdjustmentRecommendation: "Slow down",
				},
			}, nil
		},
	}

	resp, created, err := service.NewProgressService(profileRepo, planRepo, progressRepo, aiSvc).CreateSnapshot(ctx)
	if err != nil {
		t.Fatalf("create snapshot: %v", err)
	}
	if !created {
		t.Fatal("expected new snapshot, got refresh")
	}
	if resp.Stats.Status != "complete" {
		t.Fatalf("status = %q", resp.Stats.Status)
	}
	if resp.Stats.WeeklySummary != "Good week" {
		t.Fatalf("summary = %q", resp.Stats.WeeklySummary)
	}
}

func TestProgressCreateSnapshotDegradedOnAIFailure(t *testing.T) {
	profileID := testutil.MustUUID(t, testutil.TestUUID)
	profileRepo := &testutil.MockProfileRepository{
		GetFn: func(context.Context) (db.Profile, error) { return testutil.SampleProfile(profileID), nil },
	}
	planRepo := &testutil.MockPlanRepository{
		GetActiveFn: func(context.Context, pgtype.UUID) (db.RecoveryPlan, error) {
			return db.RecoveryPlan{}, pgx.ErrNoRows
		},
	}
	progressRepo := &testutil.MockProgressRepository{
		GetByProfileAndDateFn: func(context.Context, pgtype.UUID, pgtype.Date) (db.ProgressSnapshot, error) {
			return db.ProgressSnapshot{}, pgx.ErrNoRows
		},
		CountAttemptsFn:      func(context.Context, pgtype.UUID) (int32, error) { return 1, nil },
		GetEvaluationStatsFn: func(context.Context, pgtype.UUID) (repository.EvaluationStats, error) { return repository.EvaluationStats{}, nil },
		CreateFn: func(_ context.Context, _ pgtype.UUID, _ pgtype.Date, stats []byte) (db.ProgressSnapshot, error) {
			var parsed service.SnapshotStats
			_ = json.Unmarshal(stats, &parsed)
			if parsed.Status != "degraded_ai_summary" {
				t.Fatalf("expected degraded status, got %q", parsed.Status)
			}
			return db.ProgressSnapshot{Stats: stats}, nil
		},
	}
	aiSvc := &testutil.MockAIService{
		SummarizeWeeklyProgressFn: func(context.Context, ai.WeeklyProgressInput) (ai.AIResult[ai.WeeklyProgressOutput], error) {
			return ai.AIResult[ai.WeeklyProgressOutput]{}, ai.ErrNotConfigured
		},
	}

	resp, _, err := service.NewProgressService(profileRepo, planRepo, progressRepo, aiSvc).CreateSnapshot(context.Background())
	if err != nil {
		t.Fatalf("expected degraded save without error, got %v", err)
	}
	if resp.Stats.Status != "degraded_ai_summary" {
		t.Fatalf("status = %q", resp.Stats.Status)
	}
}

func TestProgressCreateSnapshotRefreshesExisting(t *testing.T) {
	profileID := testutil.MustUUID(t, testutil.TestUUID)
	profileRepo := &testutil.MockProfileRepository{
		GetFn: func(context.Context) (db.Profile, error) { return testutil.SampleProfile(profileID), nil },
	}
	planRepo := &testutil.MockPlanRepository{
		GetActiveFn: func(context.Context, pgtype.UUID) (db.RecoveryPlan, error) {
			return db.RecoveryPlan{}, pgx.ErrNoRows
		},
	}
	progressRepo := &testutil.MockProgressRepository{
		GetByProfileAndDateFn: func(context.Context, pgtype.UUID, pgtype.Date) (db.ProgressSnapshot, error) {
			return db.ProgressSnapshot{ID: profileID}, nil
		},
		CountAttemptsFn:      func(context.Context, pgtype.UUID) (int32, error) { return 2, nil },
		GetEvaluationStatsFn: func(context.Context, pgtype.UUID) (repository.EvaluationStats, error) { return repository.EvaluationStats{EvaluationsCount: 1}, nil },
		UpdateFn: func(_ context.Context, _ pgtype.UUID, _ pgtype.Date, stats []byte) (db.ProgressSnapshot, error) {
			return db.ProgressSnapshot{ID: profileID, Stats: stats}, nil
		},
	}
	aiSvc := &testutil.MockAIService{
		SummarizeWeeklyProgressFn: func(context.Context, ai.WeeklyProgressInput) (ai.AIResult[ai.WeeklyProgressOutput], error) {
			return ai.AIResult[ai.WeeklyProgressOutput]{
				Data: ai.WeeklyProgressOutput{WeeklySummary: "Updated"},
			}, nil
		},
	}

	resp, created, err := service.NewProgressService(profileRepo, planRepo, progressRepo, aiSvc).CreateSnapshot(context.Background())
	if err != nil {
		t.Fatalf("expected refresh without error, got %v", err)
	}
	if created {
		t.Fatal("expected refresh, got create")
	}
	if resp.Stats.AttemptsCount != 2 {
		t.Fatalf("attempts_count = %d", resp.Stats.AttemptsCount)
	}
}

func TestProgressListSnapshotsInvalidLimit(t *testing.T) {
	profileID := testutil.MustUUID(t, testutil.TestUUID)
	profileRepo := &testutil.MockProfileRepository{
		GetFn: func(context.Context) (db.Profile, error) { return testutil.SampleProfile(profileID), nil },
	}

	_, err := service.NewProgressService(profileRepo, &testutil.MockPlanRepository{}, &testutil.MockProgressRepository{}, &testutil.MockAIService{}).
		ListSnapshots(context.Background(), "abc", "", "")
	if err == nil {
		t.Fatal("expected validation error for invalid limit")
	}
}
