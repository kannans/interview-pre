package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"sort"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kannans/interview-prep/backend/internal/ai"
	"github.com/kannans/interview-prep/backend/internal/db"
	"github.com/kannans/interview-prep/backend/internal/repository"
)

var (
	ErrProgressSnapshotNotFound      = errors.New("no progress snapshot found")
	ErrProgressSnapshotAlreadyExists = errors.New("progress snapshot already exists for today")
	ErrProgressSnapshotNotOwned      = errors.New("progress snapshot does not belong to your profile")
)

// SnapshotAverageScores holds mean evaluation dimension scores (0–10 scale).
type SnapshotAverageScores struct {
	Correctness   float64 `json:"correctness"`
	Approach      float64 `json:"approach"`
	EdgeCases     float64 `json:"edge_cases"`
	Communication float64 `json:"communication"`
}

// SnapshotStats is the JSONB payload stored on each progress snapshot.
type SnapshotStats struct {
	AttemptsCount            int32                 `json:"attempts_count"`
	EvaluationsCount         int32                 `json:"evaluations_count"`
	AverageScores            SnapshotAverageScores `json:"average_scores"`
	PlanTasksCompleted       int32                 `json:"plan_tasks_completed"`
	PlanTasksTotal           int32                 `json:"plan_tasks_total"`
	ReadinessIndex           int32                 `json:"readiness_index"`
	WeakPatterns             []string              `json:"weak_patterns"`
	Status                   string                `json:"status"`
	WeeklySummary            string                `json:"weekly_summary,omitempty"`
	Wins                     []string              `json:"wins,omitempty"`
	FocusAreas               []string              `json:"focus_areas,omitempty"`
	AdjustmentRecommendation string                `json:"adjustment_recommendation,omitempty"`
	PromptVersion            string                `json:"prompt_version,omitempty"`
	Model                    string                `json:"model,omitempty"`
}

// ProgressSnapshotResponse is the API shape for a single snapshot.
type ProgressSnapshotResponse struct {
	ID           pgtype.UUID        `json:"id"`
	ProfileID    pgtype.UUID        `json:"profile_id"`
	SnapshotDate pgtype.Date        `json:"snapshot_date"`
	Stats        SnapshotStats      `json:"stats"`
	CreatedAt    pgtype.Timestamptz `json:"created_at"`
}

// ListProgressSnapshotsResponse is the API shape for listing snapshots.
type ListProgressSnapshotsResponse struct {
	Snapshots []ProgressSnapshotResponse `json:"snapshots"`
}

// ProgressService orchestrates progress snapshot computation and retrieval.
type ProgressService interface {
	CreateSnapshot(ctx context.Context) (ProgressSnapshotResponse, error)
	GetLatest(ctx context.Context) (ProgressSnapshotResponse, error)
	ListSnapshots(ctx context.Context, limitStr, startDateStr, endDateStr string) (ListProgressSnapshotsResponse, error)
}

type progressService struct {
	profileRepo  repository.ProfileRepository
	planRepo     repository.PlanRepository
	progressRepo repository.ProgressRepository
	aiService    ai.Service
}

func NewProgressService(
	profileRepo repository.ProfileRepository,
	planRepo repository.PlanRepository,
	progressRepo repository.ProgressRepository,
	aiService ai.Service,
) ProgressService {
	return &progressService{
		profileRepo:  profileRepo,
		planRepo:     planRepo,
		progressRepo: progressRepo,
		aiService:    aiService,
	}
}

func (s *progressService) CreateSnapshot(ctx context.Context) (ProgressSnapshotResponse, error) {
	profile, err := s.profileRepo.Get(ctx)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ProgressSnapshotResponse{}, ErrProfileMissing
		}
		return ProgressSnapshotResponse{}, fmt.Errorf("failed to get profile: %w", err)
	}

	snapshotDate := todayDate()

	stats, err := s.computeStats(ctx, profile)
	if err != nil {
		return ProgressSnapshotResponse{}, err
	}

	statsJSON, err := json.Marshal(stats)
	if err != nil {
		return ProgressSnapshotResponse{}, fmt.Errorf("failed to marshal snapshot stats: %w", err)
	}

	_, err = s.progressRepo.GetByProfileAndDate(ctx, profile.ID, snapshotDate)
	if err == nil {
		snapshot, err := s.progressRepo.Update(ctx, profile.ID, snapshotDate, statsJSON)
		if err != nil {
			return ProgressSnapshotResponse{}, fmt.Errorf("failed to refresh progress snapshot: %w", err)
		}
		return toProgressSnapshotResponse(snapshot, stats), nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return ProgressSnapshotResponse{}, fmt.Errorf("failed to check existing snapshot: %w", err)
	}

	snapshot, err := s.progressRepo.Create(ctx, profile.ID, snapshotDate, statsJSON)
	if err != nil {
		return ProgressSnapshotResponse{}, fmt.Errorf("failed to create progress snapshot: %w", err)
	}

	return toProgressSnapshotResponse(snapshot, stats), nil
}

func (s *progressService) GetLatest(ctx context.Context) (ProgressSnapshotResponse, error) {
	profile, err := s.profileRepo.Get(ctx)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ProgressSnapshotResponse{}, ErrProfileMissing
		}
		return ProgressSnapshotResponse{}, fmt.Errorf("failed to get profile: %w", err)
	}

	snapshot, err := s.progressRepo.GetLatest(ctx, profile.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ProgressSnapshotResponse{}, ErrProgressSnapshotNotFound
		}
		return ProgressSnapshotResponse{}, fmt.Errorf("failed to get latest progress snapshot: %w", err)
	}

	return decodeProgressSnapshotResponse(snapshot)
}

func (s *progressService) ListSnapshots(ctx context.Context, limitStr, startDateStr, endDateStr string) (ListProgressSnapshotsResponse, error) {
	profile, err := s.profileRepo.Get(ctx)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ListProgressSnapshotsResponse{}, ErrProfileMissing
		}
		return ListProgressSnapshotsResponse{}, fmt.Errorf("failed to get profile: %w", err)
	}

	params := repository.ProgressSnapshotListParams{ProfileID: profile.ID}

	if limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			return ListProgressSnapshotsResponse{}, fmt.Errorf("validation error: limit must be a positive integer")
		}
		params.RowLimit = pgtype.Int4{Int32: int32(limit), Valid: true}
	}

	if startDateStr != "" {
		startDate, err := parseSnapshotDate(startDateStr)
		if err != nil {
			return ListProgressSnapshotsResponse{}, fmt.Errorf("validation error: invalid start_date — use YYYY-MM-DD")
		}
		params.FilterStartDate = startDate
	}

	if endDateStr != "" {
		endDate, err := parseSnapshotDate(endDateStr)
		if err != nil {
			return ListProgressSnapshotsResponse{}, fmt.Errorf("validation error: invalid end_date — use YYYY-MM-DD")
		}
		params.FilterEndDate = endDate
	}

	snapshots, err := s.progressRepo.List(ctx, params)
	if err != nil {
		return ListProgressSnapshotsResponse{}, fmt.Errorf("failed to list progress snapshots: %w", err)
	}

	if snapshots == nil {
		snapshots = []db.ProgressSnapshot{}
	}

	result := make([]ProgressSnapshotResponse, 0, len(snapshots))
	for _, snapshot := range snapshots {
		item, err := decodeProgressSnapshotResponse(snapshot)
		if err != nil {
			return ListProgressSnapshotsResponse{}, err
		}
		result = append(result, item)
	}

	return ListProgressSnapshotsResponse{Snapshots: result}, nil
}

func (s *progressService) computeStats(ctx context.Context, profile db.Profile) (SnapshotStats, error) {
	profileID := profile.ID

	attemptsCount, err := s.progressRepo.CountAttempts(ctx, profileID)
	if err != nil {
		return SnapshotStats{}, err
	}

	evalStats, err := s.progressRepo.GetEvaluationStats(ctx, profileID)
	if err != nil {
		return SnapshotStats{}, err
	}

	planTasksCompleted, planTasksTotal := s.planTaskCounts(ctx, profileID)

	avgScores := SnapshotAverageScores{
		Correctness:   roundToOneDecimal(evalStats.AvgCorrectness),
		Approach:      roundToOneDecimal(evalStats.AvgApproach),
		EdgeCases:     roundToOneDecimal(evalStats.AvgEdgeCases),
		Communication: roundToOneDecimal(evalStats.AvgCommunication),
	}

	readinessIndex := computeReadinessIndex(
		avgScores,
		planTasksCompleted,
		planTasksTotal,
		attemptsCount,
		evalStats.EvaluationsCount,
	)

	weakPatterns := deriveWeakPatterns(avgScores, evalStats.EvaluationsCount)

	stats := SnapshotStats{
		AttemptsCount:      attemptsCount,
		EvaluationsCount:   evalStats.EvaluationsCount,
		AverageScores:      avgScores,
		PlanTasksCompleted: planTasksCompleted,
		PlanTasksTotal:     planTasksTotal,
		ReadinessIndex:     readinessIndex,
		WeakPatterns:       weakPatterns,
	}

	snapshotDate := todayDate()
	dateStr := formatSnapshotDate(snapshotDate)

	summaryInput := ai.WeeklyProgressInput{
		Profile:            profileToAIContext(profile),
		AttemptsCount:      attemptsCount,
		EvaluationsCount:   evalStats.EvaluationsCount,
		AverageScores: ai.SnapshotAverageScores{
			Correctness:   avgScores.Correctness,
			Approach:      avgScores.Approach,
			EdgeCases:     avgScores.EdgeCases,
			Communication: avgScores.Communication,
		},
		PlanTasksCompleted: planTasksCompleted,
		PlanTasksTotal:     planTasksTotal,
		ReadinessIndex:     readinessIndex,
		WeakPatterns:       weakPatterns,
		SnapshotDate:       dateStr,
	}

	result, err := s.aiService.SummarizeWeeklyProgress(ctx, summaryInput)
	if err != nil {
		slog.Warn("AI weekly summary failed; saving degraded snapshot",
			"error", err,
			"profile_id", uuidString(profileID),
		)
		stats.Status = "degraded_ai_summary"
		return stats, nil
	}

	out := result.Data
	stats.Status = "complete"
	stats.WeeklySummary = out.WeeklySummary
	stats.Wins = out.Wins
	stats.FocusAreas = out.FocusAreas
	stats.AdjustmentRecommendation = out.AdjustmentRecommendation
	stats.PromptVersion = out.PromptVersion
	stats.Model = out.Model

	return stats, nil
}

func formatSnapshotDate(d pgtype.Date) string {
	if !d.Valid {
		return time.Now().UTC().Format("2006-01-02")
	}
	t := time.Date(int(d.Time.Year()), time.Month(d.Time.Month()), int(d.Time.Day()), 0, 0, 0, 0, time.UTC)
	return t.Format("2006-01-02")
}

func uuidString(id pgtype.UUID) string {
	if !id.Valid {
		return ""
	}
	return fmt.Sprintf(
		"%x-%x-%x-%x-%x",
		id.Bytes[0:4],
		id.Bytes[4:6],
		id.Bytes[6:8],
		id.Bytes[8:10],
		id.Bytes[10:16],
	)
}

func (s *progressService) planTaskCounts(ctx context.Context, profileID pgtype.UUID) (completed int32, total int32) {
	plan, err := s.planRepo.GetActive(ctx, profileID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, 0
		}
		return 0, 0
	}

	tasks, err := s.planRepo.ListTasks(ctx, plan.ID)
	if err != nil {
		return 0, 0
	}

	total = int32(len(tasks))
	for _, task := range tasks {
		if task.IsCompleted {
			completed++
		}
	}
	return completed, total
}

func computeReadinessIndex(
	avgScores SnapshotAverageScores,
	planTasksCompleted, planTasksTotal, attemptsCount, evaluationsCount int32,
) int32 {
	scoreSum := avgScores.Correctness + avgScores.Approach + avgScores.EdgeCases + avgScores.Communication
	evaluationComponent := float64(0)
	if evaluationsCount > 0 {
		evaluationComponent = (scoreSum / 40.0) * 100.0
	}

	planComponent := float64(0)
	if planTasksTotal > 0 {
		planComponent = (float64(planTasksCompleted) / float64(planTasksTotal)) * 100.0
	}

	attemptComponent := math.Min(float64(attemptsCount)*12.5, 100.0)

	// Weighted placeholder: evaluations 45%, plan completion 35%, practice volume 20%.
	readiness := evaluationComponent*0.45 + planComponent*0.35 + attemptComponent*0.20
	return int32(math.Round(math.Min(math.Max(readiness, 0), 100)))
}

type scoreDimension struct {
	name  string
	value float64
}

func deriveWeakPatterns(avgScores SnapshotAverageScores, evaluationsCount int32) []string {
	if evaluationsCount == 0 {
		return []string{}
	}

	dimensions := []scoreDimension{
		{name: "correctness", value: avgScores.Correctness},
		{name: "approach", value: avgScores.Approach},
		{name: "edge_cases", value: avgScores.EdgeCases},
		{name: "communication", value: avgScores.Communication},
	}

	sort.Slice(dimensions, func(i, j int) bool {
		if dimensions[i].value == dimensions[j].value {
			return dimensions[i].name < dimensions[j].name
		}
		return dimensions[i].value < dimensions[j].value
	})

	weakCount := 2
	if len(dimensions) < weakCount {
		weakCount = len(dimensions)
	}

	patterns := make([]string, weakCount)
	for i := 0; i < weakCount; i++ {
		patterns[i] = dimensions[i].name
	}
	return patterns
}

func todayDate() pgtype.Date {
	now := time.Now().UTC()
	var d pgtype.Date
	_ = d.Scan(now.Format("2006-01-02"))
	return d
}

func parseSnapshotDate(value string) (pgtype.Date, error) {
	var d pgtype.Date
	if err := d.Scan(value); err != nil {
		return pgtype.Date{}, err
	}
	return d, nil
}

func roundToOneDecimal(value float64) float64 {
	return math.Round(value*10) / 10
}

func decodeProgressSnapshotResponse(snapshot db.ProgressSnapshot) (ProgressSnapshotResponse, error) {
	var stats SnapshotStats
	if err := json.Unmarshal(snapshot.Stats, &stats); err != nil {
		return ProgressSnapshotResponse{}, fmt.Errorf("failed to decode snapshot stats: %w", err)
	}
	return toProgressSnapshotResponse(snapshot, stats), nil
}

func toProgressSnapshotResponse(snapshot db.ProgressSnapshot, stats SnapshotStats) ProgressSnapshotResponse {
	return ProgressSnapshotResponse{
		ID:           snapshot.ID,
		ProfileID:    snapshot.ProfileID,
		SnapshotDate: snapshot.SnapshotDate,
		Stats:        stats,
		CreatedAt:    snapshot.CreatedAt,
	}
}