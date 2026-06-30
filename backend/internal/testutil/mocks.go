package testutil

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kannans/interview-prep/backend/internal/ai"
	"github.com/kannans/interview-prep/backend/internal/db"
	"github.com/kannans/interview-prep/backend/internal/repository"
	"github.com/kannans/interview-prep/backend/internal/service"
)

// --- Service mocks ---

type MockProfileService struct {
	UpsertFn func(ctx context.Context, input service.ProfileInput) (service.UpsertResult, error)
	GetFn    func(ctx context.Context) (db.Profile, error)
}

func (m *MockProfileService) Upsert(ctx context.Context, input service.ProfileInput) (service.UpsertResult, error) {
	return m.UpsertFn(ctx, input)
}

func (m *MockProfileService) Get(ctx context.Context) (db.Profile, error) {
	return m.GetFn(ctx)
}

type MockDiagnosisService struct {
	SubmitFn func(ctx context.Context, input service.DiagnosisInput) (db.BaselineDiagnosis, error)
	GetFn    func(ctx context.Context) (db.BaselineDiagnosis, error)
}

func (m *MockDiagnosisService) Submit(ctx context.Context, input service.DiagnosisInput) (db.BaselineDiagnosis, error) {
	return m.SubmitFn(ctx, input)
}

func (m *MockDiagnosisService) Get(ctx context.Context) (db.BaselineDiagnosis, error) {
	return m.GetFn(ctx)
}

type MockPlanService struct {
	GenerateFn     func(ctx context.Context) (service.PlanWithTasks, error)
	GetFn          func(ctx context.Context) (service.PlanWithTasks, error)
	GetTodayFn     func(ctx context.Context) (service.TodayTaskResponse, error)
	CompleteTaskFn func(ctx context.Context, taskID pgtype.UUID) (db.PlanTask, error)
}

func (m *MockPlanService) Generate(ctx context.Context) (service.PlanWithTasks, error) {
	return m.GenerateFn(ctx)
}

func (m *MockPlanService) Get(ctx context.Context) (service.PlanWithTasks, error) {
	return m.GetFn(ctx)
}

func (m *MockPlanService) GetToday(ctx context.Context) (service.TodayTaskResponse, error) {
	return m.GetTodayFn(ctx)
}

func (m *MockPlanService) CompleteTask(ctx context.Context, taskID pgtype.UUID) (db.PlanTask, error) {
	return m.CompleteTaskFn(ctx, taskID)
}

type MockQuestionService struct {
	ListFn    func(ctx context.Context, filter service.QuestionListFilter) (service.ListQuestionsResponse, error)
	GetByIDFn func(ctx context.Context, id pgtype.UUID) (db.Question, error)
}

func (m *MockQuestionService) List(ctx context.Context, filter service.QuestionListFilter) (service.ListQuestionsResponse, error) {
	return m.ListFn(ctx, filter)
}

func (m *MockQuestionService) GetByID(ctx context.Context, id pgtype.UUID) (db.Question, error) {
	return m.GetByIDFn(ctx, id)
}

type MockAttemptService struct {
	SubmitFn  func(ctx context.Context, input service.AttemptInput) (db.Attempt, error)
	ListFn    func(ctx context.Context, questionIDStr string) (service.ListAttemptsResponse, error)
	GetByIDFn func(ctx context.Context, id pgtype.UUID) (db.Attempt, error)
}

func (m *MockAttemptService) Submit(ctx context.Context, input service.AttemptInput) (db.Attempt, error) {
	return m.SubmitFn(ctx, input)
}

func (m *MockAttemptService) List(ctx context.Context, questionIDStr string) (service.ListAttemptsResponse, error) {
	return m.ListFn(ctx, questionIDStr)
}

func (m *MockAttemptService) GetByID(ctx context.Context, id pgtype.UUID) (db.Attempt, error) {
	return m.GetByIDFn(ctx, id)
}

type MockEvaluationService struct {
	EvaluateFn       func(ctx context.Context, attemptID pgtype.UUID) (db.Evaluation, error)
	GetByAttemptIDFn func(ctx context.Context, attemptID pgtype.UUID) (db.Evaluation, error)
}

func (m *MockEvaluationService) Evaluate(ctx context.Context, attemptID pgtype.UUID) (db.Evaluation, error) {
	return m.EvaluateFn(ctx, attemptID)
}

func (m *MockEvaluationService) GetByAttemptID(ctx context.Context, attemptID pgtype.UUID) (db.Evaluation, error) {
	return m.GetByAttemptIDFn(ctx, attemptID)
}

type MockProgressService struct {
	CreateSnapshotFn func(ctx context.Context) (service.ProgressSnapshotResponse, bool, error)
	GetLatestFn      func(ctx context.Context) (service.ProgressSnapshotResponse, error)
	ListSnapshotsFn  func(ctx context.Context, limitStr, startDateStr, endDateStr string) (service.ListProgressSnapshotsResponse, error)
}

func (m *MockProgressService) CreateSnapshot(ctx context.Context) (service.ProgressSnapshotResponse, bool, error) {
	return m.CreateSnapshotFn(ctx)
}

func (m *MockProgressService) GetLatest(ctx context.Context) (service.ProgressSnapshotResponse, error) {
	return m.GetLatestFn(ctx)
}

func (m *MockProgressService) ListSnapshots(ctx context.Context, limitStr, startDateStr, endDateStr string) (service.ListProgressSnapshotsResponse, error) {
	return m.ListSnapshotsFn(ctx, limitStr, startDateStr, endDateStr)
}

// --- Repository mocks ---

type MockProfileRepository struct {
	HasFn     func(ctx context.Context) (bool, error)
	GetFn     func(ctx context.Context) (db.Profile, error)
	CreateFn  func(ctx context.Context, params db.CreateProfileParams) (db.Profile, error)
	UpdateFn  func(ctx context.Context, id pgtype.UUID, params db.CreateProfileParams) (db.Profile, error)
}

func (m *MockProfileRepository) Has(ctx context.Context) (bool, error)       { return m.HasFn(ctx) }
func (m *MockProfileRepository) Get(ctx context.Context) (db.Profile, error)   { return m.GetFn(ctx) }
func (m *MockProfileRepository) Create(ctx context.Context, params db.CreateProfileParams) (db.Profile, error) {
	return m.CreateFn(ctx, params)
}
func (m *MockProfileRepository) Update(ctx context.Context, id pgtype.UUID, params db.CreateProfileParams) (db.Profile, error) {
	return m.UpdateFn(ctx, id, params)
}

type MockDiagnosisRepository struct {
	CreateFn func(ctx context.Context, profileID pgtype.UUID, coding, conceptual, behavioral string, report json.RawMessage) (db.BaselineDiagnosis, error)
	GetFn    func(ctx context.Context, profileID pgtype.UUID) (db.BaselineDiagnosis, error)
	HasFn    func(ctx context.Context, profileID pgtype.UUID) (bool, error)
}

func (m *MockDiagnosisRepository) Create(ctx context.Context, profileID pgtype.UUID, coding, conceptual, behavioral string, report json.RawMessage) (db.BaselineDiagnosis, error) {
	return m.CreateFn(ctx, profileID, coding, conceptual, behavioral, report)
}
func (m *MockDiagnosisRepository) Get(ctx context.Context, profileID pgtype.UUID) (db.BaselineDiagnosis, error) {
	return m.GetFn(ctx, profileID)
}
func (m *MockDiagnosisRepository) Has(ctx context.Context, profileID pgtype.UUID) (bool, error) {
	return m.HasFn(ctx, profileID)
}

type MockQuestionRepository struct {
	ListFn    func(ctx context.Context, params repository.QuestionListParams) ([]db.Question, error)
	GetByIDFn func(ctx context.Context, id pgtype.UUID) (db.Question, error)
}

func (m *MockQuestionRepository) List(ctx context.Context, params repository.QuestionListParams) ([]db.Question, error) {
	return m.ListFn(ctx, params)
}
func (m *MockQuestionRepository) GetByID(ctx context.Context, id pgtype.UUID) (db.Question, error) {
	return m.GetByIDFn(ctx, id)
}

type MockAttemptRepository struct {
	CreateFn  func(ctx context.Context, profileID, questionID pgtype.UUID, content string, timeTaken int32) (db.Attempt, error)
	ListFn    func(ctx context.Context, params repository.AttemptListParams) ([]db.Attempt, error)
	GetByIDFn func(ctx context.Context, id pgtype.UUID) (db.Attempt, error)
}

func (m *MockAttemptRepository) Create(ctx context.Context, profileID, questionID pgtype.UUID, content string, timeTaken int32) (db.Attempt, error) {
	return m.CreateFn(ctx, profileID, questionID, content, timeTaken)
}
func (m *MockAttemptRepository) List(ctx context.Context, params repository.AttemptListParams) ([]db.Attempt, error) {
	return m.ListFn(ctx, params)
}
func (m *MockAttemptRepository) GetByID(ctx context.Context, id pgtype.UUID) (db.Attempt, error) {
	return m.GetByIDFn(ctx, id)
}

type MockEvaluationRepository struct {
	CreateFn         func(ctx context.Context, params db.CreateEvaluationParams) (db.Evaluation, error)
	GetByAttemptIDFn func(ctx context.Context, attemptID pgtype.UUID) (db.Evaluation, error)
	GetByIDFn        func(ctx context.Context, id pgtype.UUID) (db.Evaluation, error)
}

func (m *MockEvaluationRepository) Create(ctx context.Context, params db.CreateEvaluationParams) (db.Evaluation, error) {
	return m.CreateFn(ctx, params)
}
func (m *MockEvaluationRepository) GetByAttemptID(ctx context.Context, attemptID pgtype.UUID) (db.Evaluation, error) {
	return m.GetByAttemptIDFn(ctx, attemptID)
}
func (m *MockEvaluationRepository) GetByID(ctx context.Context, id pgtype.UUID) (db.Evaluation, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(ctx, id)
	}
	return db.Evaluation{}, nil
}

type MockPlanRepository struct {
	HasActiveFn            func(ctx context.Context, profileID pgtype.UUID) (bool, error)
	GetActiveFn            func(ctx context.Context, profileID pgtype.UUID) (db.RecoveryPlan, error)
	ListTasksFn            func(ctx context.Context, planID pgtype.UUID) ([]db.PlanTask, error)
	GetTodayTaskFn         func(ctx context.Context, planID pgtype.UUID) (db.PlanTask, error)
	GetTaskFn              func(ctx context.Context, taskID pgtype.UUID) (db.PlanTask, error)
	MarkTaskCompleteFn     func(ctx context.Context, taskID, planID pgtype.UUID) (db.PlanTask, error)
	ReplacePlanWithTasksFn func(ctx context.Context, profileID pgtype.UUID, tasks []repository.PlanTaskInput) (db.RecoveryPlan, []db.PlanTask, error)
}

func (m *MockPlanRepository) HasActive(ctx context.Context, profileID pgtype.UUID) (bool, error) {
	return m.HasActiveFn(ctx, profileID)
}
func (m *MockPlanRepository) GetActive(ctx context.Context, profileID pgtype.UUID) (db.RecoveryPlan, error) {
	return m.GetActiveFn(ctx, profileID)
}
func (m *MockPlanRepository) ListTasks(ctx context.Context, planID pgtype.UUID) ([]db.PlanTask, error) {
	return m.ListTasksFn(ctx, planID)
}
func (m *MockPlanRepository) GetTodayTask(ctx context.Context, planID pgtype.UUID) (db.PlanTask, error) {
	return m.GetTodayTaskFn(ctx, planID)
}
func (m *MockPlanRepository) GetTask(ctx context.Context, taskID pgtype.UUID) (db.PlanTask, error) {
	return m.GetTaskFn(ctx, taskID)
}
func (m *MockPlanRepository) MarkTaskComplete(ctx context.Context, taskID, planID pgtype.UUID) (db.PlanTask, error) {
	return m.MarkTaskCompleteFn(ctx, taskID, planID)
}
func (m *MockPlanRepository) ReplacePlanWithTasks(ctx context.Context, profileID pgtype.UUID, tasks []repository.PlanTaskInput) (db.RecoveryPlan, []db.PlanTask, error) {
	return m.ReplacePlanWithTasksFn(ctx, profileID, tasks)
}

type MockProgressRepository struct {
	CreateFn              func(ctx context.Context, profileID pgtype.UUID, snapshotDate pgtype.Date, stats []byte) (db.ProgressSnapshot, error)
	UpdateFn              func(ctx context.Context, profileID pgtype.UUID, snapshotDate pgtype.Date, stats []byte) (db.ProgressSnapshot, error)
	GetLatestFn           func(ctx context.Context, profileID pgtype.UUID) (db.ProgressSnapshot, error)
	GetByIDFn             func(ctx context.Context, id pgtype.UUID) (db.ProgressSnapshot, error)
	GetByProfileAndDateFn func(ctx context.Context, profileID pgtype.UUID, snapshotDate pgtype.Date) (db.ProgressSnapshot, error)
	ListFn                func(ctx context.Context, params repository.ProgressSnapshotListParams) ([]db.ProgressSnapshot, error)
	CountAttemptsFn       func(ctx context.Context, profileID pgtype.UUID) (int32, error)
	GetEvaluationStatsFn  func(ctx context.Context, profileID pgtype.UUID) (repository.EvaluationStats, error)
}

func (m *MockProgressRepository) Create(ctx context.Context, profileID pgtype.UUID, snapshotDate pgtype.Date, stats []byte) (db.ProgressSnapshot, error) {
	return m.CreateFn(ctx, profileID, snapshotDate, stats)
}
func (m *MockProgressRepository) Update(ctx context.Context, profileID pgtype.UUID, snapshotDate pgtype.Date, stats []byte) (db.ProgressSnapshot, error) {
	if m.UpdateFn != nil {
		return m.UpdateFn(ctx, profileID, snapshotDate, stats)
	}
	return db.ProgressSnapshot{Stats: stats}, nil
}
func (m *MockProgressRepository) GetLatest(ctx context.Context, profileID pgtype.UUID) (db.ProgressSnapshot, error) {
	return m.GetLatestFn(ctx, profileID)
}
func (m *MockProgressRepository) GetByID(ctx context.Context, id pgtype.UUID) (db.ProgressSnapshot, error) {
	return m.GetByIDFn(ctx, id)
}
func (m *MockProgressRepository) GetByProfileAndDate(ctx context.Context, profileID pgtype.UUID, snapshotDate pgtype.Date) (db.ProgressSnapshot, error) {
	return m.GetByProfileAndDateFn(ctx, profileID, snapshotDate)
}
func (m *MockProgressRepository) List(ctx context.Context, params repository.ProgressSnapshotListParams) ([]db.ProgressSnapshot, error) {
	return m.ListFn(ctx, params)
}
func (m *MockProgressRepository) CountAttempts(ctx context.Context, profileID pgtype.UUID) (int32, error) {
	return m.CountAttemptsFn(ctx, profileID)
}
func (m *MockProgressRepository) GetEvaluationStats(ctx context.Context, profileID pgtype.UUID) (repository.EvaluationStats, error) {
	return m.GetEvaluationStatsFn(ctx, profileID)
}

// --- AI mock ---

type MockAIService struct {
	EvaluateBaselineDiagnosisFn      func(ctx context.Context, input ai.BaselineDiagnosisInput) (ai.AIResult[ai.DiagnosticReport], error)
	GenerateRecoveryPlanFn           func(ctx context.Context, input ai.RecoveryPlanInput) (ai.AIResult[ai.RecoveryPlanOutput], error)
	EvaluateAttemptFn                func(ctx context.Context, questionType string, input ai.EvaluationInput) (ai.AIResult[ai.EvaluationOutput], error)
	SummarizeWeeklyProgressFn        func(ctx context.Context, input ai.WeeklyProgressInput) (ai.AIResult[ai.WeeklyProgressOutput], error)
	GenerateMockInterviewQuestionsFn func(ctx context.Context, profile ai.ProfileContext, focusAreas []string) (ai.AIResult[ai.MockInterviewQuestionsOutput], error)
	GradeMockInterviewFn             func(ctx context.Context, transcript string) (ai.AIResult[ai.MockInterviewGradingOutput], error)
}

func (m *MockAIService) EvaluateBaselineDiagnosis(ctx context.Context, input ai.BaselineDiagnosisInput) (ai.AIResult[ai.DiagnosticReport], error) {
	return m.EvaluateBaselineDiagnosisFn(ctx, input)
}
func (m *MockAIService) GenerateRecoveryPlan(ctx context.Context, input ai.RecoveryPlanInput) (ai.AIResult[ai.RecoveryPlanOutput], error) {
	return m.GenerateRecoveryPlanFn(ctx, input)
}
func (m *MockAIService) EvaluateAttempt(ctx context.Context, questionType string, input ai.EvaluationInput) (ai.AIResult[ai.EvaluationOutput], error) {
	return m.EvaluateAttemptFn(ctx, questionType, input)
}
func (m *MockAIService) SummarizeWeeklyProgress(ctx context.Context, input ai.WeeklyProgressInput) (ai.AIResult[ai.WeeklyProgressOutput], error) {
	return m.SummarizeWeeklyProgressFn(ctx, input)
}
func (m *MockAIService) GenerateMockInterviewQuestions(ctx context.Context, profile ai.ProfileContext, focusAreas []string) (ai.AIResult[ai.MockInterviewQuestionsOutput], error) {
	return m.GenerateMockInterviewQuestionsFn(ctx, profile, focusAreas)
}
func (m *MockAIService) GradeMockInterview(ctx context.Context, transcript string) (ai.AIResult[ai.MockInterviewGradingOutput], error) {
	return m.GradeMockInterviewFn(ctx, transcript)
}

// SampleProfile returns a minimal profile fixture.
func SampleProfile(id pgtype.UUID) db.Profile {
	return db.Profile{
		ID:              id,
		ExperienceYears: 8,
		TargetRole:      "Senior Backend Engineer",
		TargetCompanies: "Google",
		WeakAreasSelf:   "DSA",
		Timeline:        "60 days",
	}
}
