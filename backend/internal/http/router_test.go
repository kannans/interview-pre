package http_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kannans/interview-prep/backend/internal/db"
	"github.com/kannans/interview-prep/backend/internal/http/handler"
	"github.com/kannans/interview-prep/backend/internal/service"
	"github.com/kannans/interview-prep/backend/internal/testutil"
)

type testServices struct {
	profile    service.ProfileService
	diagnosis  service.DiagnosisService
	plan       service.PlanService
	question   service.QuestionService
	attempt    service.AttemptService
	evaluation service.EvaluationService
	progress   service.ProgressService
}

func newTestRouter(svc testServices) http.Handler {
	r := chi.NewRouter()

	r.Get("/health", handler.NewHealthHandler(nil).Check)

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/profile", handler.NewProfileHandler(svc.profile).Upsert)
		r.Get("/profile", handler.NewProfileHandler(svc.profile).Get)

		r.Post("/diagnosis", handler.NewDiagnosisHandler(svc.diagnosis).Submit)
		r.Get("/diagnosis", handler.NewDiagnosisHandler(svc.diagnosis).Get)

		r.Post("/plan", handler.NewPlanHandler(svc.plan).Generate)
		r.Get("/plan", handler.NewPlanHandler(svc.plan).Get)
		r.Get("/plan/today", handler.NewPlanHandler(svc.plan).GetToday)
		r.Patch("/plan/tasks/{id}/complete", handler.NewPlanHandler(svc.plan).CompleteTask)

		r.Get("/questions", handler.NewQuestionHandler(svc.question).List)
		r.Get("/questions/{id}", handler.NewQuestionHandler(svc.question).Get)

		r.Post("/attempts", handler.NewAttemptHandler(svc.attempt).Submit)
		r.Get("/attempts", handler.NewAttemptHandler(svc.attempt).List)
		r.Post("/attempts/{id}/evaluate", handler.NewEvaluationHandler(svc.evaluation).Evaluate)
		r.Get("/attempts/{id}/evaluation", handler.NewEvaluationHandler(svc.evaluation).GetByAttempt)
		r.Get("/attempts/{id}", handler.NewAttemptHandler(svc.attempt).Get)

		r.Post("/progress/snapshot", handler.NewProgressHandler(svc.progress).CreateSnapshot)
		r.Get("/progress/latest", handler.NewProgressHandler(svc.progress).GetLatest)
		r.Get("/progress/snapshots", handler.NewProgressHandler(svc.progress).ListSnapshots)
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"status": "error", "message": "resource not found"}`))
	})

	return r
}

func happyPathServices(t *testing.T) testServices {
	t.Helper()

	profileID := testutil.MustUUID(t, testutil.TestUUID)
	attemptID := testutil.MustUUID(t, "b2000002-0000-4000-8000-000000000002")
	taskID := testutil.MustUUID(t, "c3000003-0000-4000-8000-000000000003")
	planID := testutil.MustUUID(t, "d4000004-0000-4000-8000-000000000004")
	profile := testutil.SampleProfile(profileID)

	return testServices{
		profile: &testutil.MockProfileService{
			UpsertFn: func(_ context.Context, _ service.ProfileInput) (service.UpsertResult, error) {
				return service.UpsertResult{Profile: profile, Created: true}, nil
			},
			GetFn: func(_ context.Context) (db.Profile, error) {
				return profile, nil
			},
		},
		diagnosis: &testutil.MockDiagnosisService{
			SubmitFn: func(_ context.Context, _ service.DiagnosisInput) (db.BaselineDiagnosis, error) {
				return db.BaselineDiagnosis{
					ID:               profileID,
					ProfileID:        profileID,
					DiagnosticReport: json.RawMessage(`{"overall_readiness":60}`),
				}, nil
			},
			GetFn: func(_ context.Context) (db.BaselineDiagnosis, error) {
				return db.BaselineDiagnosis{ID: profileID, ProfileID: profileID}, nil
			},
		},
		plan: &testutil.MockPlanService{
			GenerateFn: func(_ context.Context) (service.PlanWithTasks, error) {
				return service.PlanWithTasks{
					Plan:  db.RecoveryPlan{ID: planID, ProfileID: profileID, Status: "active"},
					Tasks: []db.PlanTask{{ID: taskID, PlanID: planID, Title: "Task 1"}},
				}, nil
			},
			GetFn: func(_ context.Context) (service.PlanWithTasks, error) {
				return service.PlanWithTasks{
					Plan:  db.RecoveryPlan{ID: planID, ProfileID: profileID, Status: "active"},
					Tasks: []db.PlanTask{{ID: taskID, PlanID: planID, Title: "Task 1"}},
				}, nil
			},
			GetTodayFn: func(_ context.Context) (service.TodayTaskResponse, error) {
				return service.TodayTaskResponse{
					Plan: db.RecoveryPlan{ID: planID},
					Task: db.PlanTask{ID: taskID, Title: "Today task"},
				}, nil
			},
			CompleteTaskFn: func(_ context.Context, _ pgtype.UUID) (db.PlanTask, error) {
				return db.PlanTask{ID: taskID, IsCompleted: true}, nil
			},
		},
		question: &testutil.MockQuestionService{
			ListFn: func(_ context.Context, _ service.QuestionListFilter) (service.ListQuestionsResponse, error) {
				return service.ListQuestionsResponse{
					Questions: []db.Question{{ID: profileID, Title: "Q1", Type: "dsa"}},
				}, nil
			},
			GetByIDFn: func(_ context.Context, _ pgtype.UUID) (db.Question, error) {
				return db.Question{ID: profileID, Title: "Q1", Type: "dsa"}, nil
			},
		},
		attempt: &testutil.MockAttemptService{
			SubmitFn: func(_ context.Context, _ service.AttemptInput) (db.Attempt, error) {
				return db.Attempt{ID: attemptID, ProfileID: profileID}, nil
			},
			ListFn: func(_ context.Context, _ string) (service.ListAttemptsResponse, error) {
				return service.ListAttemptsResponse{Attempts: []db.Attempt{{ID: attemptID}}}, nil
			},
			GetByIDFn: func(_ context.Context, _ pgtype.UUID) (db.Attempt, error) {
				return db.Attempt{ID: attemptID, ProfileID: profileID}, nil
			},
		},
		evaluation: &testutil.MockEvaluationService{
			EvaluateFn: func(_ context.Context, _ pgtype.UUID) (db.Evaluation, error) {
				return db.Evaluation{ID: profileID, AttemptID: attemptID, ScoreCorrectness: 8}, nil
			},
			GetByAttemptIDFn: func(_ context.Context, _ pgtype.UUID) (db.Evaluation, error) {
				return db.Evaluation{ID: profileID, AttemptID: attemptID, ScoreCorrectness: 8}, nil
			},
		},
		progress: &testutil.MockProgressService{
			CreateSnapshotFn: func(_ context.Context) (service.ProgressSnapshotResponse, error) {
				return service.ProgressSnapshotResponse{
					ID:        profileID,
					ProfileID: profileID,
					Stats:     service.SnapshotStats{Status: "complete", ReadinessIndex: 55},
				}, nil
			},
			GetLatestFn: func(_ context.Context) (service.ProgressSnapshotResponse, error) {
				return service.ProgressSnapshotResponse{
					Stats: service.SnapshotStats{Status: "complete"},
				}, nil
			},
			ListSnapshotsFn: func(_ context.Context, _, _, _ string) (service.ListProgressSnapshotsResponse, error) {
				return service.ListProgressSnapshotsResponse{
					Snapshots: []service.ProgressSnapshotResponse{{Stats: service.SnapshotStats{Status: "complete"}}},
				}, nil
			},
		},
	}
}

func TestAllEndpointsCoverage(t *testing.T) {
	svc := happyPathServices(t)
	router := newTestRouter(svc)

	attemptID := "b2000002-0000-4000-8000-000000000002"
	taskID := "c3000003-0000-4000-8000-000000000003"

	tests := []struct {
		name       string
		method     string
		path       string
		body       any
		wantStatus int
	}{
		{"GET /health", http.MethodGet, "/health", nil, http.StatusServiceUnavailable},
		{"POST /api/v1/profile", http.MethodPost, "/api/v1/profile", map[string]any{
			"experience_years": 8, "target_role": "SWE", "target_companies": "Google",
			"weak_areas_self": "DSA", "timeline": "60 days",
		}, http.StatusCreated},
		{"GET /api/v1/profile", http.MethodGet, "/api/v1/profile", nil, http.StatusOK},
		{"POST /api/v1/diagnosis", http.MethodPost, "/api/v1/diagnosis", map[string]string{
			"coding_answer": "code", "conceptual_answer": "concept", "behavioral_answer": "behavior",
		}, http.StatusCreated},
		{"GET /api/v1/diagnosis", http.MethodGet, "/api/v1/diagnosis", nil, http.StatusOK},
		{"POST /api/v1/plan", http.MethodPost, "/api/v1/plan", nil, http.StatusCreated},
		{"GET /api/v1/plan", http.MethodGet, "/api/v1/plan", nil, http.StatusOK},
		{"GET /api/v1/plan/today", http.MethodGet, "/api/v1/plan/today", nil, http.StatusOK},
		{"PATCH /api/v1/plan/tasks/complete", http.MethodPatch, "/api/v1/plan/tasks/" + taskID + "/complete", nil, http.StatusOK},
		{"GET /api/v1/questions", http.MethodGet, "/api/v1/questions", nil, http.StatusOK},
		{"GET /api/v1/questions/id", http.MethodGet, "/api/v1/questions/" + testutil.TestUUID, nil, http.StatusOK},
		{"POST /api/v1/attempts", http.MethodPost, "/api/v1/attempts", map[string]any{
			"question_id": testutil.TestUUID, "submitted_content": "answer", "time_taken_seconds": 120,
		}, http.StatusCreated},
		{"GET /api/v1/attempts", http.MethodGet, "/api/v1/attempts", nil, http.StatusOK},
		{"GET /api/v1/attempts/id", http.MethodGet, "/api/v1/attempts/" + attemptID, nil, http.StatusOK},
		{"POST /api/v1/attempts/evaluate", http.MethodPost, "/api/v1/attempts/" + attemptID + "/evaluate", nil, http.StatusCreated},
		{"GET /api/v1/attempts/evaluation", http.MethodGet, "/api/v1/attempts/" + attemptID + "/evaluation", nil, http.StatusOK},
		{"POST /api/v1/progress/snapshot", http.MethodPost, "/api/v1/progress/snapshot", nil, http.StatusCreated},
		{"GET /api/v1/progress/latest", http.MethodGet, "/api/v1/progress/latest", nil, http.StatusOK},
		{"GET /api/v1/progress/snapshots", http.MethodGet, "/api/v1/progress/snapshots?limit=5", nil, http.StatusOK},
		{"unknown route", http.MethodGet, "/api/v1/unknown", nil, http.StatusNotFound},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rec := testutil.DoRequest(t, router, testutil.JSONRequest(tc.method, tc.path, tc.body))
			testutil.AssertStatus(t, rec, tc.wantStatus)
		})
	}
}

func TestEndpointErrorResponses(t *testing.T) {
	profileID := testutil.MustUUID(t, testutil.TestUUID)
	attemptID := testutil.MustUUID(t, "b2000002-0000-4000-8000-000000000002")

	base := happyPathServices(t)

	t.Run("POST /profile invalid body", func(t *testing.T) {
		router := newTestRouter(base)
		rec := testutil.DoRequest(t, router, testutil.JSONRequest(http.MethodPost, "/api/v1/profile", "not-json"))
		testutil.AssertStatus(t, rec, http.StatusBadRequest)
	})

	t.Run("GET /profile not found", func(t *testing.T) {
		svc := base
		svc.profile = &testutil.MockProfileService{
			GetFn: func(_ context.Context) (db.Profile, error) {
				return db.Profile{}, pgx.ErrNoRows
			},
		}
		rec := testutil.DoRequest(t, newTestRouter(svc), testutil.JSONRequest(http.MethodGet, "/api/v1/profile", nil))
		testutil.AssertStatus(t, rec, http.StatusNotFound)
	})

	t.Run("POST /diagnosis AI failure 502", func(t *testing.T) {
		svc := base
		svc.diagnosis = &testutil.MockDiagnosisService{
			SubmitFn: func(_ context.Context, _ service.DiagnosisInput) (db.BaselineDiagnosis, error) {
				return db.BaselineDiagnosis{}, fmt.Errorf("%w: timeout", service.ErrAIFailed)
			},
		}
		body := map[string]string{"coding_answer": "a", "conceptual_answer": "b", "behavioral_answer": "c"}
		rec := testutil.DoRequest(t, newTestRouter(svc), testutil.JSONRequest(http.MethodPost, "/api/v1/diagnosis", body))
		testutil.AssertStatus(t, rec, http.StatusBadGateway)
	})

	t.Run("POST /plan missing diagnosis 412", func(t *testing.T) {
		svc := base
		svc.plan = &testutil.MockPlanService{
			GenerateFn: func(_ context.Context) (service.PlanWithTasks, error) {
				return service.PlanWithTasks{}, service.ErrDiagnosisMissing
			},
		}
		rec := testutil.DoRequest(t, newTestRouter(svc), testutil.JSONRequest(http.MethodPost, "/api/v1/plan", nil))
		testutil.AssertStatus(t, rec, http.StatusPreconditionFailed)
	})

	t.Run("GET /questions invalid filter 400", func(t *testing.T) {
		svc := base
		svc.question = &testutil.MockQuestionService{
			ListFn: func(_ context.Context, _ service.QuestionListFilter) (service.ListQuestionsResponse, error) {
				return service.ListQuestionsResponse{}, fmt.Errorf("validation error: bad type")
			},
		}
		rec := testutil.DoRequest(t, newTestRouter(svc), testutil.JSONRequest(http.MethodGet, "/api/v1/questions?type=invalid", nil))
		testutil.AssertStatus(t, rec, http.StatusBadRequest)
	})

	t.Run("GET /attempts/{id} invalid uuid 400", func(t *testing.T) {
		rec := testutil.DoRequest(t, newTestRouter(base), testutil.JSONRequest(http.MethodGet, "/api/v1/attempts/not-a-uuid", nil))
		testutil.AssertStatus(t, rec, http.StatusBadRequest)
	})

	t.Run("POST /attempts/{id}/evaluate conflict 409", func(t *testing.T) {
		svc := base
		svc.evaluation = &testutil.MockEvaluationService{
			EvaluateFn: func(_ context.Context, _ pgtype.UUID) (db.Evaluation, error) {
				return db.Evaluation{}, service.ErrEvaluationAlreadyExists
			},
		}
		path := "/api/v1/attempts/" + uuidStr(attemptID) + "/evaluate"
		rec := testutil.DoRequest(t, newTestRouter(svc), testutil.JSONRequest(http.MethodPost, path, nil))
		testutil.AssertStatus(t, rec, http.StatusConflict)
	})

	t.Run("POST /progress/snapshot duplicate 409", func(t *testing.T) {
		svc := base
		svc.progress = &testutil.MockProgressService{
			CreateSnapshotFn: func(_ context.Context) (service.ProgressSnapshotResponse, error) {
				return service.ProgressSnapshotResponse{}, service.ErrProgressSnapshotAlreadyExists
			},
		}
		rec := testutil.DoRequest(t, newTestRouter(svc), testutil.JSONRequest(http.MethodPost, "/api/v1/progress/snapshot", nil))
		testutil.AssertStatus(t, rec, http.StatusConflict)
	})

	t.Run("GET /progress/latest not found 404", func(t *testing.T) {
		svc := base
		svc.progress = &testutil.MockProgressService{
			GetLatestFn: func(_ context.Context) (service.ProgressSnapshotResponse, error) {
				return service.ProgressSnapshotResponse{}, service.ErrProgressSnapshotNotFound
			},
		}
		rec := testutil.DoRequest(t, newTestRouter(svc), testutil.JSONRequest(http.MethodGet, "/api/v1/progress/latest", nil))
		testutil.AssertStatus(t, rec, http.StatusNotFound)
	})

	t.Run("POST /diagnosis validation 400", func(t *testing.T) {
		svc := base
		svc.diagnosis = &testutil.MockDiagnosisService{
			SubmitFn: func(_ context.Context, _ service.DiagnosisInput) (db.BaselineDiagnosis, error) {
				return db.BaselineDiagnosis{}, fmt.Errorf("validation error: coding_answer is required")
			},
		}
		rec := testutil.DoRequest(t, newTestRouter(svc), testutil.JSONRequest(http.MethodPost, "/api/v1/diagnosis", map[string]string{
			"coding_answer": "only one field",
		}))
		testutil.AssertStatus(t, rec, http.StatusBadRequest)
	})

	t.Run("PATCH /plan/tasks invalid uuid 400", func(t *testing.T) {
		rec := testutil.DoRequest(t, newTestRouter(base), testutil.JSONRequest(http.MethodPatch, "/api/v1/plan/tasks/bad-id/complete", nil))
		testutil.AssertStatus(t, rec, http.StatusBadRequest)
	})

	t.Run("GET /progress/snapshots invalid limit 400", func(t *testing.T) {
		svc := base
		svc.progress = &testutil.MockProgressService{
			ListSnapshotsFn: func(_ context.Context, _, _, _ string) (service.ListProgressSnapshotsResponse, error) {
				return service.ListProgressSnapshotsResponse{}, fmt.Errorf("validation error: limit must be a positive integer")
			},
		}
		rec := testutil.DoRequest(t, newTestRouter(svc), testutil.JSONRequest(http.MethodGet, "/api/v1/progress/snapshots?limit=abc", nil))
		testutil.AssertStatus(t, rec, http.StatusBadRequest)
	})

	t.Run("GET /diagnosis not found 404", func(t *testing.T) {
		svc := base
		svc.diagnosis = &testutil.MockDiagnosisService{
			GetFn: func(_ context.Context) (db.BaselineDiagnosis, error) {
				return db.BaselineDiagnosis{}, pgx.ErrNoRows
			},
		}
		rec := testutil.DoRequest(t, newTestRouter(svc), testutil.JSONRequest(http.MethodGet, "/api/v1/diagnosis", nil))
		testutil.AssertStatus(t, rec, http.StatusNotFound)
	})

	t.Run("POST /profile validation 400", func(t *testing.T) {
		svc := base
		svc.profile = &testutil.MockProfileService{
			UpsertFn: func(_ context.Context, _ service.ProfileInput) (service.UpsertResult, error) {
				return service.UpsertResult{}, fmt.Errorf("validation error: target_role is required")
			},
		}
		rec := testutil.DoRequest(t, newTestRouter(svc), testutil.JSONRequest(http.MethodPost, "/api/v1/profile", map[string]any{}))
		testutil.AssertStatus(t, rec, http.StatusBadRequest)
	})

	_ = profileID
}

func uuidStr(id pgtype.UUID) string {
	return fmt.Sprintf(
		"%x-%x-%x-%x-%x",
		id.Bytes[0:4], id.Bytes[4:6], id.Bytes[6:8], id.Bytes[8:10], id.Bytes[10:16],
	)
}
