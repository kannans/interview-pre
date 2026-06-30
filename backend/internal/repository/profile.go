package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kannans/interview-prep/backend/internal/db"
)

type ProfileRepository interface {
	Get(ctx context.Context) (db.Profile, error)
	Create(ctx context.Context, params db.CreateProfileParams) (db.Profile, error)
	Update(ctx context.Context, id pgtype.UUID, params db.CreateProfileParams) (db.Profile, error)
	Has(ctx context.Context) (bool, error)
}

type profileRepository struct {
	queries db.Querier
}

func NewProfileRepository(queries db.Querier) ProfileRepository {
	return &profileRepository{queries: queries}
}

func (r *profileRepository) Get(ctx context.Context) (db.Profile, error) {
	return r.queries.GetProfile(ctx)
}

func (r *profileRepository) Create(ctx context.Context, params db.CreateProfileParams) (db.Profile, error) {
	return r.queries.CreateProfile(ctx, params)
}

func (r *profileRepository) Update(ctx context.Context, id pgtype.UUID, params db.CreateProfileParams) (db.Profile, error) {
	return r.queries.UpdateProfile(ctx, db.UpdateProfileParams{
		ID:              id,
		ExperienceYears: params.ExperienceYears,
		TargetRole:      params.TargetRole,
		TargetCompanies: params.TargetCompanies,
		WeakAreasSelf:   params.WeakAreasSelf,
		Timeline:        params.Timeline,
	})
}

func (r *profileRepository) Has(ctx context.Context) (bool, error) {
	return r.queries.HasProfile(ctx)
}
