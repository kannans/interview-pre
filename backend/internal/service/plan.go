package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kannans/interview-prep/backend/internal/db"
	"github.com/kannans/interview-prep/backend/internal/repository"
)

var (
	ErrProfileMissing   = errors.New("profile not found — create a profile first")
	ErrDiagnosisMissing = errors.New("baseline diagnosis not found — submit a diagnosis first")
	ErrPlanNotFound     = errors.New("no active recovery plan found")
	ErrTaskNotInPlan    = errors.New("task does not belong to your active plan")
	ErrNoTasksRemaining = errors.New("all plan tasks are completed")
	ErrTaskAlreadyDone  = errors.New("task is already completed or not found")
)

// PlanWithTasks is the API response shape for a recovery plan and its ordered tasks.
type PlanWithTasks struct {
	Plan  db.RecoveryPlan `json:"plan"`
	Tasks []db.PlanTask   `json:"tasks"`
}

// TodayTaskResponse is the API response shape for today's recommended task.
type TodayTaskResponse struct {
	Plan db.RecoveryPlan `json:"plan"`
	Task db.PlanTask     `json:"task"`
}

// PlanService orchestrates recovery plan generation and task management.
type PlanService interface {
	Generate(ctx context.Context) (PlanWithTasks, error)
	Get(ctx context.Context) (PlanWithTasks, error)
	GetToday(ctx context.Context) (TodayTaskResponse, error)
	CompleteTask(ctx context.Context, taskID pgtype.UUID) (db.PlanTask, error)
}

type planService struct {
	profileRepo   repository.ProfileRepository
	diagnosisRepo repository.DiagnosisRepository
	planRepo      repository.PlanRepository
}

func NewPlanService(
	profileRepo repository.ProfileRepository,
	diagnosisRepo repository.DiagnosisRepository,
	planRepo repository.PlanRepository,
) PlanService {
	return &planService{
		profileRepo:   profileRepo,
		diagnosisRepo: diagnosisRepo,
		planRepo:      planRepo,
	}
}

func (s *planService) Generate(ctx context.Context) (PlanWithTasks, error) {
	profile, err := s.profileRepo.Get(ctx)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return PlanWithTasks{}, ErrProfileMissing
		}
		return PlanWithTasks{}, fmt.Errorf("failed to get profile: %w", err)
	}

	_, err = s.diagnosisRepo.Get(ctx, profile.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return PlanWithTasks{}, ErrDiagnosisMissing
		}
		return PlanWithTasks{}, fmt.Errorf("failed to get diagnosis: %w", err)
	}

	tasks := derivePlanTasks(profile.WeakAreasSelf)
	plan, createdTasks, err := s.planRepo.ReplacePlanWithTasks(ctx, profile.ID, tasks)
	if err != nil {
		return PlanWithTasks{}, fmt.Errorf("failed to generate recovery plan: %w", err)
	}

	return PlanWithTasks{Plan: plan, Tasks: createdTasks}, nil
}

func (s *planService) Get(ctx context.Context) (PlanWithTasks, error) {
	profile, err := s.profileRepo.Get(ctx)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return PlanWithTasks{}, ErrProfileMissing
		}
		return PlanWithTasks{}, fmt.Errorf("failed to get profile: %w", err)
	}

	plan, err := s.planRepo.GetActive(ctx, profile.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return PlanWithTasks{}, ErrPlanNotFound
		}
		return PlanWithTasks{}, fmt.Errorf("failed to get active plan: %w", err)
	}

	tasks, err := s.planRepo.ListTasks(ctx, plan.ID)
	if err != nil {
		return PlanWithTasks{}, fmt.Errorf("failed to list plan tasks: %w", err)
	}

	return PlanWithTasks{Plan: plan, Tasks: tasks}, nil
}

func (s *planService) GetToday(ctx context.Context) (TodayTaskResponse, error) {
	planWithTasks, err := s.Get(ctx)
	if err != nil {
		return TodayTaskResponse{}, err
	}

	task, err := s.planRepo.GetTodayTask(ctx, planWithTasks.Plan.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return TodayTaskResponse{}, ErrNoTasksRemaining
		}
		return TodayTaskResponse{}, fmt.Errorf("failed to get today's task: %w", err)
	}

	return TodayTaskResponse{Plan: planWithTasks.Plan, Task: task}, nil
}

func (s *planService) CompleteTask(ctx context.Context, taskID pgtype.UUID) (db.PlanTask, error) {
	profile, err := s.profileRepo.Get(ctx)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.PlanTask{}, ErrProfileMissing
		}
		return db.PlanTask{}, fmt.Errorf("failed to get profile: %w", err)
	}

	plan, err := s.planRepo.GetActive(ctx, profile.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.PlanTask{}, ErrPlanNotFound
		}
		return db.PlanTask{}, fmt.Errorf("failed to get active plan: %w", err)
	}

	task, err := s.planRepo.GetTask(ctx, taskID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.PlanTask{}, ErrTaskNotInPlan
		}
		return db.PlanTask{}, fmt.Errorf("failed to get task: %w", err)
	}

	if !bytes.Equal(task.PlanID.Bytes[:], plan.ID.Bytes[:]) {
		return db.PlanTask{}, ErrTaskNotInPlan
	}

	completed, err := s.planRepo.MarkTaskComplete(ctx, taskID, plan.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.PlanTask{}, ErrTaskAlreadyDone
		}
		return db.PlanTask{}, fmt.Errorf("failed to complete task: %w", err)
	}

	return completed, nil
}

type taskTemplate struct {
	Category string
	Level    string
	Title    string
}

// derivePlanTasks builds a static 5–10 task plan from self-identified weak areas.
// Phase 7 will replace this with AI-driven plan generation from diagnostic_report.
func derivePlanTasks(weakAreasSelf string) []repository.PlanTaskInput {
	weak := strings.ToLower(weakAreasSelf)

	emphasizeDSA := containsAny(weak, "dsa", "algorithm", "leetcode", "array", "tree", "graph", "dynamic programming", "dp")
	emphasizeBasics := containsAny(weak, "syntax", "fundamental", "basic", "javascript", "python", "go", "typescript", "language")
	emphasizeConcept := containsAny(weak, "system design", "concept", "api", "http", "database", "sql", "oop", "async")

	templates := []taskTemplate{
		{Category: "programming_basics", Level: "level_1", Title: "Review core data types, variables, and control flow"},
		{Category: "programming_basics", Level: "level_2", Title: "Practice writing clean functions with clear error handling"},
		{Category: "dsa", Level: "level_1", Title: "Arrays and string manipulation fundamentals"},
		{Category: "dsa", Level: "level_2", Title: "Hash maps and two-pointer techniques"},
		{Category: "concept", Level: "level_1", Title: "Explain time and space complexity for common operations"},
		{Category: "concept", Level: "level_2", Title: "REST APIs, HTTP methods, and status codes"},
		{Category: "dsa", Level: "level_3", Title: "Binary search, sorting, and recursion review"},
	}

	if emphasizeDSA {
		templates = append(templates,
			taskTemplate{Category: "dsa", Level: "level_2", Title: "Practice sliding window and prefix sum patterns"},
			taskTemplate{Category: "dsa", Level: "level_3", Title: "Tree traversal: BFS and DFS fundamentals"},
		)
	}
	if emphasizeBasics {
		templates = append(templates,
			taskTemplate{Category: "programming_basics", Level: "level_2", Title: "Refactor a small function for readability and edge cases"},
			taskTemplate{Category: "programming_basics", Level: "level_3", Title: "Debug a function with off-by-one and null/empty input bugs"},
		)
	}
	if emphasizeConcept {
		templates = append(templates,
			taskTemplate{Category: "concept", Level: "level_2", Title: "Compare SQL vs NoSQL trade-offs for a simple product"},
			taskTemplate{Category: "concept", Level: "level_3", Title: "Sketch a high-level design for a URL shortener"},
		)
	}
	if !emphasizeDSA && !emphasizeBasics && !emphasizeConcept {
		templates = append(templates,
			taskTemplate{Category: "concept", Level: "level_1", Title: "Walk through how you would debug a production issue step by step"},
		)
	}

	// Cap at 10 tasks for v1 static plans.
	if len(templates) > 10 {
		templates = templates[:10]
	}

	result := make([]repository.PlanTaskInput, len(templates))
	for i, t := range templates {
		result[i] = repository.PlanTaskInput{
			Category:      t.Category,
			Level:         t.Level,
			Title:         t.Title,
			SequenceOrder: int32(i + 1),
		}
	}
	return result
}

func containsAny(s string, keywords ...string) bool {
	for _, kw := range keywords {
		if strings.Contains(s, kw) {
			return true
		}
	}
	return false
}
