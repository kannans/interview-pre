package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kannans/interview-prep/backend/internal/db"
)

type DiagnosisRepository interface {
	Create(ctx context.Context, profileID pgtype.UUID, coding, conceptual, behavioral string, report json.RawMessage) (db.BaselineDiagnosis, error)
	Get(ctx context.Context, profileID pgtype.UUID) (db.BaselineDiagnosis, error)
	Has(ctx context.Context, profileID pgtype.UUID) (bool, error)
}

type diagnosisRepository struct {
	queries db.Querier
}

func NewDiagnosisRepository(queries db.Querier) DiagnosisRepository {
	return &diagnosisRepository{queries: queries}
}

func (r *diagnosisRepository) Create(ctx context.Context, profileID pgtype.UUID, coding, conceptual, behavioral string, report json.RawMessage) (db.BaselineDiagnosis, error) {
	d, err := r.queries.CreateBaselineDiagnosis(ctx, db.CreateBaselineDiagnosisParams{
		ProfileID:        profileID,
		CodingAnswer:     coding,
		ConceptualAnswer: conceptual,
		BehavioralAnswer: behavioral,
		DiagnosticReport: report,
	})
	if err != nil {
		return db.BaselineDiagnosis{}, fmt.Errorf("failed to create diagnosis: %w", err)
	}
	return d, nil
}

func (r *diagnosisRepository) Get(ctx context.Context, profileID pgtype.UUID) (db.BaselineDiagnosis, error) {
	d, err := r.queries.GetBaselineDiagnosis(ctx, profileID)
	if err != nil {
		return db.BaselineDiagnosis{}, fmt.Errorf("failed to get diagnosis: %w", err)
	}
	return d, nil
}

func (r *diagnosisRepository) Has(ctx context.Context, profileID pgtype.UUID) (bool, error) {
	return r.queries.HasBaselineDiagnosis(ctx, profileID)
}
