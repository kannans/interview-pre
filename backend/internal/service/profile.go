package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/kannans/interview-prep/backend/internal/db"
	"github.com/kannans/interview-prep/backend/internal/repository"
)

// ProfileInput holds the candidate-provided data for creating or updating a profile.
type ProfileInput struct {
	ExperienceYears int32  `json:"experience_years"`
	TargetRole      string `json:"target_role"`
	TargetCompanies string `json:"target_companies"`
	WeakAreasSelf   string `json:"weak_areas_self"`
	Timeline        string `json:"timeline"`
}

func (i ProfileInput) Validate() error {
	if i.ExperienceYears < 0 {
		return errors.New("experience_years must be zero or greater")
	}
	if i.TargetRole == "" {
		return errors.New("target_role is required")
	}
	if i.TargetCompanies == "" {
		return errors.New("target_companies is required")
	}
	if i.WeakAreasSelf == "" {
		return errors.New("weak_areas_self is required")
	}
	if i.Timeline == "" {
		return errors.New("timeline is required")
	}
	return nil
}

// UpsertResult indicates whether a profile was created or updated.
type UpsertResult struct {
	Profile db.Profile
	Created bool // true if newly created, false if updated
}

type ProfileService interface {
	Upsert(ctx context.Context, input ProfileInput) (UpsertResult, error)
	Get(ctx context.Context) (db.Profile, error)
}

type profileService struct {
	repo repository.ProfileRepository
}

func NewProfileService(repo repository.ProfileRepository) ProfileService {
	return &profileService{repo: repo}
}

func (s *profileService) Upsert(ctx context.Context, input ProfileInput) (UpsertResult, error) {
	if err := input.Validate(); err != nil {
		return UpsertResult{}, fmt.Errorf("validation error: %w", err)
	}

	params := db.CreateProfileParams{
		ExperienceYears: input.ExperienceYears,
		TargetRole:      input.TargetRole,
		TargetCompanies: input.TargetCompanies,
		WeakAreasSelf:   input.WeakAreasSelf,
		Timeline:        input.Timeline,
	}

	has, err := s.repo.Has(ctx)
	if err != nil {
		return UpsertResult{}, fmt.Errorf("failed to check profile existence: %w", err)
	}

	if has {
		// Profile already exists — update it
		existing, err := s.repo.Get(ctx)
		if err != nil {
			return UpsertResult{}, fmt.Errorf("failed to get existing profile: %w", err)
		}
		updated, err := s.repo.Update(ctx, existing.ID, params)
		if err != nil {
			return UpsertResult{}, fmt.Errorf("failed to update profile: %w", err)
		}
		return UpsertResult{Profile: updated, Created: false}, nil
	}

	// No profile yet — create one
	created, err := s.repo.Create(ctx, params)
	if err != nil {
		return UpsertResult{}, fmt.Errorf("failed to create profile: %w", err)
	}
	return UpsertResult{Profile: created, Created: true}, nil
}

func (s *profileService) Get(ctx context.Context) (db.Profile, error) {
	profile, err := s.repo.Get(ctx)
	if err != nil {
		return db.Profile{}, fmt.Errorf("failed to get profile: %w", err)
	}
	return profile, nil
}
