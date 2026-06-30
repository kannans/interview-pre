package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kannans/interview-prep/backend/internal/ai"
	"github.com/kannans/interview-prep/backend/internal/db"
	"github.com/kannans/interview-prep/backend/internal/service"
	"github.com/kannans/interview-prep/backend/internal/testutil"
)

func TestAttemptInputValidate(t *testing.T) {
	cases := []service.AttemptInput{
		{QuestionID: "", SubmittedContent: "x", TimeTakenSeconds: 1},
		{QuestionID: testutil.TestUUID, SubmittedContent: "", TimeTakenSeconds: 1},
		{QuestionID: testutil.TestUUID, SubmittedContent: "x", TimeTakenSeconds: 0},
	}
	for i, input := range cases {
		if err := input.Validate(); err == nil {
			t.Fatalf("case %d: expected validation error", i)
		}
	}
}

func TestAttemptSubmitSuccess(t *testing.T) {
	ctx := context.Background()
	profileID := testutil.MustUUID(t, testutil.TestUUID)
	questionID := testutil.MustUUID(t, testutil.TestUUID)

	profileRepo := &testutil.MockProfileRepository{
		GetFn: func(context.Context) (db.Profile, error) { return testutil.SampleProfile(profileID), nil },
	}
	questionRepo := &testutil.MockQuestionRepository{
		GetByIDFn: func(context.Context, pgtype.UUID) (db.Question, error) {
			return db.Question{ID: questionID}, nil
		},
	}
	attemptRepo := &testutil.MockAttemptRepository{
		CreateFn: func(_ context.Context, pid, qid pgtype.UUID, content string, secs int32) (db.Attempt, error) {
			return db.Attempt{ProfileID: pid, QuestionID: qid, SubmittedContent: content}, nil
		},
	}

	attempt, err := service.NewAttemptService(profileRepo, questionRepo, attemptRepo).Submit(ctx, service.AttemptInput{
		QuestionID: testutil.TestUUID, SubmittedContent: "answer", TimeTakenSeconds: 60,
	})
	if err != nil {
		t.Fatalf("submit: %v", err)
	}
	if attempt.SubmittedContent != "answer" {
		t.Fatalf("unexpected content: %q", attempt.SubmittedContent)
	}
}

func TestAttemptSubmitProfileMissing(t *testing.T) {
	profileRepo := &testutil.MockProfileRepository{
		GetFn: func(context.Context) (db.Profile, error) { return db.Profile{}, pgx.ErrNoRows },
	}
	_, err := service.NewAttemptService(profileRepo, &testutil.MockQuestionRepository{}, &testutil.MockAttemptRepository{}).Submit(
		context.Background(),
		service.AttemptInput{QuestionID: testutil.TestUUID, SubmittedContent: "a", TimeTakenSeconds: 1},
	)
	if !errors.Is(err, service.ErrProfileMissing) {
		t.Fatalf("expected ErrProfileMissing, got %v", err)
	}
}

func TestEvaluationEvaluateSuccess(t *testing.T) {
	ctx := context.Background()
	profileID := testutil.MustUUID(t, testutil.TestUUID)
	attemptID := testutil.MustUUID(t, "b2000002-0000-4000-8000-000000000002")
	questionID := testutil.MustUUID(t, testutil.TestUUID)

	profileRepo := &testutil.MockProfileRepository{
		GetFn: func(context.Context) (db.Profile, error) { return testutil.SampleProfile(profileID), nil },
	}
	attemptRepo := &testutil.MockAttemptRepository{
		GetByIDFn: func(context.Context, pgtype.UUID) (db.Attempt, error) {
			return db.Attempt{ID: attemptID, ProfileID: profileID, QuestionID: questionID}, nil
		},
	}
	questionRepo := &testutil.MockQuestionRepository{
		GetByIDFn: func(context.Context, pgtype.UUID) (db.Question, error) {
			return db.Question{ID: questionID, Type: "dsa", Title: "Max"}, nil
		},
	}
	evalRepo := &testutil.MockEvaluationRepository{
		GetByAttemptIDFn: func(context.Context, pgtype.UUID) (db.Evaluation, error) {
			return db.Evaluation{}, pgx.ErrNoRows
		},
		CreateFn: func(_ context.Context, params db.CreateEvaluationParams) (db.Evaluation, error) {
			return db.Evaluation{AttemptID: params.AttemptID, ScoreCorrectness: params.ScoreCorrectness}, nil
		},
	}
	aiSvc := &testutil.MockAIService{
		EvaluateAttemptFn: func(_ context.Context, qType string, _ ai.EvaluationInput) (ai.AIResult[ai.EvaluationOutput], error) {
			if qType != "dsa" {
				t.Fatalf("expected dsa, got %s", qType)
			}
			return ai.AIResult[ai.EvaluationOutput]{
				Data: ai.EvaluationOutput{
					ScoreCorrectness: 8, ScoreApproach: 7, ScoreEdgeCases: 6, ScoreCommunication: 7,
					FeedbackCorrectness: "ok", FeedbackApproach: "ok", FeedbackEdgeCases: "ok",
					FeedbackCommunication: "ok", NextAction: "retry",
				},
			}, nil
		},
	}

	eval, err := service.NewEvaluationService(profileRepo, attemptRepo, questionRepo, evalRepo, aiSvc).Evaluate(ctx, attemptID)
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	if eval.ScoreCorrectness != 8 {
		t.Fatalf("score = %d", eval.ScoreCorrectness)
	}
}

func TestEvaluationAlreadyExists(t *testing.T) {
	profileID := testutil.MustUUID(t, testutil.TestUUID)
	attemptID := testutil.MustUUID(t, "b2000002-0000-4000-8000-000000000002")

	profileRepo := &testutil.MockProfileRepository{
		GetFn: func(context.Context) (db.Profile, error) { return testutil.SampleProfile(profileID), nil },
	}
	attemptRepo := &testutil.MockAttemptRepository{
		GetByIDFn: func(context.Context, pgtype.UUID) (db.Attempt, error) {
			return db.Attempt{ID: attemptID, ProfileID: profileID}, nil
		},
	}
	evalRepo := &testutil.MockEvaluationRepository{
		GetByAttemptIDFn: func(context.Context, pgtype.UUID) (db.Evaluation, error) {
			return db.Evaluation{}, nil
		},
	}

	_, err := service.NewEvaluationService(profileRepo, attemptRepo, &testutil.MockQuestionRepository{}, evalRepo, &testutil.MockAIService{}).Evaluate(context.Background(), attemptID)
	if !errors.Is(err, service.ErrEvaluationAlreadyExists) {
		t.Fatalf("expected ErrEvaluationAlreadyExists, got %v", err)
	}
}
