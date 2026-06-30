package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kannans/interview-prep/backend/internal/db"
)

type ProgressSnapshotListParams struct {
	ProfileID       pgtype.UUID
	FilterStartDate pgtype.Date
	FilterEndDate   pgtype.Date
	RowLimit        pgtype.Int4
}

type EvaluationStats struct {
	EvaluationsCount  int32
	AvgCorrectness    float64
	AvgApproach       float64
	AvgEdgeCases      float64
	AvgCommunication  float64
}

type ProgressRepository interface {
	Create(ctx context.Context, profileID pgtype.UUID, snapshotDate pgtype.Date, stats []byte) (db.ProgressSnapshot, error)
	Update(ctx context.Context, profileID pgtype.UUID, snapshotDate pgtype.Date, stats []byte) (db.ProgressSnapshot, error)
	GetLatest(ctx context.Context, profileID pgtype.UUID) (db.ProgressSnapshot, error)
	GetByID(ctx context.Context, id pgtype.UUID) (db.ProgressSnapshot, error)
	GetByProfileAndDate(ctx context.Context, profileID pgtype.UUID, snapshotDate pgtype.Date) (db.ProgressSnapshot, error)
	List(ctx context.Context, params ProgressSnapshotListParams) ([]db.ProgressSnapshot, error)
	CountAttempts(ctx context.Context, profileID pgtype.UUID) (int32, error)
	GetEvaluationStats(ctx context.Context, profileID pgtype.UUID) (EvaluationStats, error)
}

type progressRepository struct {
	queries db.Querier
}

func NewProgressRepository(queries db.Querier) ProgressRepository {
	return &progressRepository{queries: queries}
}

func (r *progressRepository) Create(ctx context.Context, profileID pgtype.UUID, snapshotDate pgtype.Date, stats []byte) (db.ProgressSnapshot, error) {
	snapshot, err := r.queries.CreateProgressSnapshot(ctx, db.CreateProgressSnapshotParams{
		ProfileID:    profileID,
		SnapshotDate: snapshotDate,
		Stats:        stats,
	})
	if err != nil {
		return db.ProgressSnapshot{}, fmt.Errorf("failed to create progress snapshot: %w", err)
	}
	return snapshot, nil
}

func (r *progressRepository) Update(ctx context.Context, profileID pgtype.UUID, snapshotDate pgtype.Date, stats []byte) (db.ProgressSnapshot, error) {
	snapshot, err := r.queries.UpdateProgressSnapshot(ctx, db.UpdateProgressSnapshotParams{
		ProfileID:    profileID,
		SnapshotDate: snapshotDate,
		Stats:        stats,
	})
	if err != nil {
		return db.ProgressSnapshot{}, fmt.Errorf("failed to update progress snapshot: %w", err)
	}
	return snapshot, nil
}

func (r *progressRepository) GetLatest(ctx context.Context, profileID pgtype.UUID) (db.ProgressSnapshot, error) {
	snapshot, err := r.queries.GetLatestProgressSnapshot(ctx, profileID)
	if err != nil {
		return db.ProgressSnapshot{}, fmt.Errorf("failed to get latest progress snapshot: %w", err)
	}
	return snapshot, nil
}

func (r *progressRepository) GetByID(ctx context.Context, id pgtype.UUID) (db.ProgressSnapshot, error) {
	snapshot, err := r.queries.GetProgressSnapshotByID(ctx, id)
	if err != nil {
		return db.ProgressSnapshot{}, fmt.Errorf("failed to get progress snapshot: %w", err)
	}
	return snapshot, nil
}

func (r *progressRepository) GetByProfileAndDate(ctx context.Context, profileID pgtype.UUID, snapshotDate pgtype.Date) (db.ProgressSnapshot, error) {
	snapshot, err := r.queries.GetProgressSnapshotByProfileAndDate(ctx, db.GetProgressSnapshotByProfileAndDateParams{
		ProfileID:    profileID,
		SnapshotDate: snapshotDate,
	})
	if err != nil {
		return db.ProgressSnapshot{}, fmt.Errorf("failed to get progress snapshot by date: %w", err)
	}
	return snapshot, nil
}

func (r *progressRepository) List(ctx context.Context, params ProgressSnapshotListParams) ([]db.ProgressSnapshot, error) {
	snapshots, err := r.queries.ListProgressSnapshots(ctx, db.ListProgressSnapshotsParams{
		ProfileID:       params.ProfileID,
		FilterStartDate: params.FilterStartDate,
		FilterEndDate:   params.FilterEndDate,
		RowLimit:        params.RowLimit,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list progress snapshots: %w", err)
	}
	return snapshots, nil
}

func (r *progressRepository) CountAttempts(ctx context.Context, profileID pgtype.UUID) (int32, error) {
	count, err := r.queries.CountAttemptsByProfile(ctx, profileID)
	if err != nil {
		return 0, fmt.Errorf("failed to count attempts: %w", err)
	}
	return count, nil
}

func (r *progressRepository) GetEvaluationStats(ctx context.Context, profileID pgtype.UUID) (EvaluationStats, error) {
	stats, err := r.queries.GetEvaluationStatsForProfile(ctx, profileID)
	if err != nil {
		return EvaluationStats{}, fmt.Errorf("failed to get evaluation stats: %w", err)
	}
	return EvaluationStats{
		EvaluationsCount: stats.EvaluationsCount,
		AvgCorrectness:   stats.AvgCorrectness,
		AvgApproach:      stats.AvgApproach,
		AvgEdgeCases:     stats.AvgEdgeCases,
		AvgCommunication: stats.AvgCommunication,
	}, nil
}

var _ ProgressRepository = (*progressRepository)(nil)
