package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kannans/interview-prep/backend/internal/db"
)

type PlanTaskInput struct {
	Category      string
	Level         string
	Title         string
	SequenceOrder int32
}

type PlanRepository interface {
	HasActive(ctx context.Context, profileID pgtype.UUID) (bool, error)
	GetActive(ctx context.Context, profileID pgtype.UUID) (db.RecoveryPlan, error)
	ListTasks(ctx context.Context, planID pgtype.UUID) ([]db.PlanTask, error)
	GetTodayTask(ctx context.Context, planID pgtype.UUID) (db.PlanTask, error)
	GetTask(ctx context.Context, taskID pgtype.UUID) (db.PlanTask, error)
	MarkTaskComplete(ctx context.Context, taskID, planID pgtype.UUID) (db.PlanTask, error)
	ReplacePlanWithTasks(ctx context.Context, profileID pgtype.UUID, tasks []PlanTaskInput) (db.RecoveryPlan, []db.PlanTask, error)
}

type planRepository struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

func NewPlanRepository(pool *pgxpool.Pool, queries *db.Queries) PlanRepository {
	return &planRepository{pool: pool, queries: queries}
}

func (r *planRepository) HasActive(ctx context.Context, profileID pgtype.UUID) (bool, error) {
	return r.queries.HasActiveRecoveryPlan(ctx, profileID)
}

func (r *planRepository) GetActive(ctx context.Context, profileID pgtype.UUID) (db.RecoveryPlan, error) {
	plan, err := r.queries.GetActiveRecoveryPlan(ctx, profileID)
	if err != nil {
		return db.RecoveryPlan{}, fmt.Errorf("failed to get active plan: %w", err)
	}
	return plan, nil
}

func (r *planRepository) ListTasks(ctx context.Context, planID pgtype.UUID) ([]db.PlanTask, error) {
	tasks, err := r.queries.ListPlanTasks(ctx, planID)
	if err != nil {
		return nil, fmt.Errorf("failed to list plan tasks: %w", err)
	}
	return tasks, nil
}

func (r *planRepository) GetTodayTask(ctx context.Context, planID pgtype.UUID) (db.PlanTask, error) {
	task, err := r.queries.GetTodayPlanTask(ctx, planID)
	if err != nil {
		return db.PlanTask{}, fmt.Errorf("failed to get today's task: %w", err)
	}
	return task, nil
}

func (r *planRepository) GetTask(ctx context.Context, taskID pgtype.UUID) (db.PlanTask, error) {
	task, err := r.queries.GetPlanTask(ctx, taskID)
	if err != nil {
		return db.PlanTask{}, fmt.Errorf("failed to get plan task: %w", err)
	}
	return task, nil
}

func (r *planRepository) MarkTaskComplete(ctx context.Context, taskID, planID pgtype.UUID) (db.PlanTask, error) {
	task, err := r.queries.MarkPlanTaskComplete(ctx, db.MarkPlanTaskCompleteParams{
		ID:     taskID,
		PlanID: planID,
	})
	if err != nil {
		return db.PlanTask{}, fmt.Errorf("failed to mark task complete: %w", err)
	}
	return task, nil
}

// ReplacePlanWithTasks supersedes any existing active plan and creates a new one with tasks atomically.
func (r *planRepository) ReplacePlanWithTasks(ctx context.Context, profileID pgtype.UUID, tasks []PlanTaskInput) (db.RecoveryPlan, []db.PlanTask, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return db.RecoveryPlan{}, nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	q := r.queries.WithTx(tx)

	if err := q.SupersedeActiveRecoveryPlan(ctx, profileID); err != nil {
		return db.RecoveryPlan{}, nil, fmt.Errorf("failed to supersede active plan: %w", err)
	}

	plan, err := q.CreateRecoveryPlan(ctx, profileID)
	if err != nil {
		return db.RecoveryPlan{}, nil, fmt.Errorf("failed to create recovery plan: %w", err)
	}

	createdTasks := make([]db.PlanTask, 0, len(tasks))
	for _, t := range tasks {
		task, err := q.CreatePlanTask(ctx, db.CreatePlanTaskParams{
			PlanID:        plan.ID,
			Category:      t.Category,
			Level:         t.Level,
			Title:         t.Title,
			SequenceOrder: t.SequenceOrder,
		})
		if err != nil {
			return db.RecoveryPlan{}, nil, fmt.Errorf("failed to create plan task: %w", err)
		}
		createdTasks = append(createdTasks, task)
	}

	if err := tx.Commit(ctx); err != nil {
		return db.RecoveryPlan{}, nil, fmt.Errorf("failed to commit plan transaction: %w", err)
	}

	return plan, createdTasks, nil
}

var _ PlanRepository = (*planRepository)(nil)
