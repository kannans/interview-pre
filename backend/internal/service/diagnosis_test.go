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
	"github.com/kannans/interview-prep/backend/internal/service"
	"github.com/kannans/interview-prep/backend/internal/testutil"
)

func TestDiagnosisInputValidate(t *testing.T) {
	input := service.DiagnosisInput{CodingAnswer: "a", ConceptualAnswer: "b"}
	if err := input.Validate(); err == nil {
		t.Fatal("expected missing behavioral_answer error")
	}
}

func TestDiagnosisSubmitSuccess(t *testing.T) {
	ctx := context.Background()
	profileID := testutil.MustUUID(t, testutil.TestUUID)
	profile := testutil.SampleProfile(profileID)

	profileRepo := &testutil.MockProfileRepository{
		GetFn: func(context.Context) (db.Profile, error) { return profile, nil },
	}
	diagnosisRepo := &testutil.MockDiagnosisRepository{
		CreateFn: func(_ context.Context, pid pgtype.UUID, _, _, _ string, report json.RawMessage) (db.BaselineDiagnosis, error) {
			return db.BaselineDiagnosis{ProfileID: pid, DiagnosticReport: report}, nil
		},
	}
	aiSvc := &testutil.MockAIService{
		EvaluateBaselineDiagnosisFn: func(_ context.Context, _ ai.BaselineDiagnosisInput) (ai.AIResult[ai.DiagnosticReport], error) {
			return ai.AIResult[ai.DiagnosticReport]{
				Data: ai.DiagnosticReport{OverallReadiness: 60, Summary: "ok"},
				Raw:  json.RawMessage(`{"overall_readiness":60,"summary":"ok"}`),
			}, nil
		},
	}

	svc := service.NewDiagnosisService(profileRepo, diagnosisRepo, aiSvc)
	_, err := svc.Submit(ctx, service.DiagnosisInput{
		CodingAnswer: "code", ConceptualAnswer: "concept", BehavioralAnswer: "behavior",
	})
	if err != nil {
		t.Fatalf("submit: %v", err)
	}
}

func TestDiagnosisSubmitAIFailure(t *testing.T) {
	profileID := testutil.MustUUID(t, testutil.TestUUID)
	profileRepo := &testutil.MockProfileRepository{
		GetFn: func(context.Context) (db.Profile, error) { return testutil.SampleProfile(profileID), nil },
	}
	aiSvc := &testutil.MockAIService{
		EvaluateBaselineDiagnosisFn: func(context.Context, ai.BaselineDiagnosisInput) (ai.AIResult[ai.DiagnosticReport], error) {
			return ai.AIResult[ai.DiagnosticReport]{}, ai.ErrNotConfigured
		},
	}

	svc := service.NewDiagnosisService(profileRepo, &testutil.MockDiagnosisRepository{}, aiSvc)
	_, err := svc.Submit(context.Background(), service.DiagnosisInput{
		CodingAnswer: "a", ConceptualAnswer: "b", BehavioralAnswer: "c",
	})
	if !errors.Is(err, service.ErrAIFailed) {
		t.Fatalf("expected ErrAIFailed, got %v", err)
	}
}

func TestDiagnosisSubmitProfileMissing(t *testing.T) {
	profileRepo := &testutil.MockProfileRepository{
		GetFn: func(context.Context) (db.Profile, error) { return db.Profile{}, pgx.ErrNoRows },
	}
	svc := service.NewDiagnosisService(profileRepo, &testutil.MockDiagnosisRepository{}, &testutil.MockAIService{})
	_, err := svc.Submit(context.Background(), service.DiagnosisInput{
		CodingAnswer: "a", ConceptualAnswer: "b", BehavioralAnswer: "c",
	})
	if err == nil {
		t.Fatal("expected error when profile missing")
	}
}
