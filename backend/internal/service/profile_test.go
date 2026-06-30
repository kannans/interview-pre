package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kannans/interview-prep/backend/internal/db"
	"github.com/kannans/interview-prep/backend/internal/service"
	"github.com/kannans/interview-prep/backend/internal/testutil"
)

func TestProfileInputValidate(t *testing.T) {
	valid := service.ProfileInput{
		ExperienceYears: 5,
		TargetRole:      "Engineer",
		TargetCompanies: "Co",
		WeakAreasSelf:   "DSA",
		Timeline:        "30 days",
	}
	if err := valid.Validate(); err != nil {
		t.Fatalf("expected valid input: %v", err)
	}

	cases := []struct {
		name  string
		input service.ProfileInput
	}{
		{"missing role", service.ProfileInput{ExperienceYears: 1, TargetCompanies: "x", WeakAreasSelf: "y", Timeline: "z"}},
		{"negative years", service.ProfileInput{ExperienceYears: -1, TargetRole: "r", TargetCompanies: "x", WeakAreasSelf: "y", Timeline: "z"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.input.Validate(); err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}

func TestProfileServiceUpsertCreateAndUpdate(t *testing.T) {
	ctx := context.Background()
	id := testutil.MustUUID(t, testutil.TestUUID)
	created := false

	repo := &testutil.MockProfileRepository{
		HasFn: func(context.Context) (bool, error) { return created, nil },
		GetFn: func(context.Context) (db.Profile, error) {
			return testutil.SampleProfile(id), nil
		},
		CreateFn: func(_ context.Context, params db.CreateProfileParams) (db.Profile, error) {
			created = true
			return db.Profile{ID: id, TargetRole: params.TargetRole}, nil
		},
		UpdateFn: func(_ context.Context, _ pgtype.UUID, params db.CreateProfileParams) (db.Profile, error) {
			return db.Profile{ID: id, TargetRole: params.TargetRole}, nil
		},
	}

	svc := service.NewProfileService(repo)
	input := service.ProfileInput{
		ExperienceYears: 8, TargetRole: "SWE", TargetCompanies: "G",
		WeakAreasSelf: "DSA", Timeline: "60d",
	}

	result, err := svc.Upsert(ctx, input)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if !result.Created {
		t.Fatal("expected created=true")
	}

	result, err = svc.Upsert(ctx, input)
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if result.Created {
		t.Fatal("expected created=false on update")
	}
}

func TestProfileServiceGetNotFound(t *testing.T) {
	repo := &testutil.MockProfileRepository{
		GetFn: func(context.Context) (db.Profile, error) { return db.Profile{}, pgx.ErrNoRows },
	}
	_, err := service.NewProfileService(repo).Get(context.Background())
	if !errors.Is(err, pgx.ErrNoRows) {
		t.Fatalf("expected ErrNoRows, got %v", err)
	}
}
