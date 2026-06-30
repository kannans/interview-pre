package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kannans/interview-prep/backend/internal/ai"
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
	aiService     ai.Service
}

func NewPlanService(
	profileRepo repository.ProfileRepository,
	diagnosisRepo repository.DiagnosisRepository,
	planRepo repository.PlanRepository,
	aiService ai.Service,
) PlanService {
	return &planService{
		profileRepo:   profileRepo,
		diagnosisRepo: diagnosisRepo,
		planRepo:      planRepo,
		aiService:     aiService,
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

	diagnosis, err := s.diagnosisRepo.Get(ctx, profile.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return PlanWithTasks{}, ErrDiagnosisMissing
		}
		return PlanWithTasks{}, fmt.Errorf("failed to get diagnosis: %w", err)
	}

	result, err := s.aiService.GenerateRecoveryPlan(ctx, ai.RecoveryPlanInput{
		Profile:          profileToAIContext(profile),
		DiagnosticReport: json.RawMessage(diagnosis.DiagnosticReport),
	})
	if err != nil {
		return PlanWithTasks{}, fmt.Errorf("%w: %v", ErrAIFailed, err)
	}

	tasks := make([]repository.PlanTaskInput, len(result.Data.Tasks))
	for i, t := range result.Data.Tasks {
		tasks[i] = repository.PlanTaskInput{
			Category:      t.Category,
			Level:         t.Level,
			Title:         t.Title,
			SequenceOrder: int32(i + 1),
		}
	}

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
