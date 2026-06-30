package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/kannans/interview-prep/backend/internal/ai"
	"github.com/kannans/interview-prep/backend/internal/db"
	"github.com/kannans/interview-prep/backend/internal/repository"
)

var ErrAIFailed = errors.New("ai service failed")

// DiagnosisInput holds the three answer types from the baseline assessment.
type DiagnosisInput struct {
	CodingAnswer     string `json:"coding_answer"`
	ConceptualAnswer string `json:"conceptual_answer"`
	BehavioralAnswer string `json:"behavioral_answer"`
}

func (i DiagnosisInput) Validate() error {
	if i.CodingAnswer == "" {
		return errors.New("coding_answer is required")
	}
	if i.ConceptualAnswer == "" {
		return errors.New("conceptual_answer is required")
	}
	if i.BehavioralAnswer == "" {
		return errors.New("behavioral_answer is required")
	}
	return nil
}

// DiagnosisService orchestrates baseline diagnosis creation and retrieval.
type DiagnosisService interface {
	Submit(ctx context.Context, input DiagnosisInput) (db.BaselineDiagnosis, error)
	Get(ctx context.Context) (db.BaselineDiagnosis, error)
}

type diagnosisService struct {
	profileRepo   repository.ProfileRepository
	diagnosisRepo repository.DiagnosisRepository
	aiService     ai.Service
}

func NewDiagnosisService(profileRepo repository.ProfileRepository, diagnosisRepo repository.DiagnosisRepository, aiService ai.Service) DiagnosisService {
	return &diagnosisService{
		profileRepo:   profileRepo,
		diagnosisRepo: diagnosisRepo,
		aiService:     aiService,
	}
}

func (s *diagnosisService) Submit(ctx context.Context, input DiagnosisInput) (db.BaselineDiagnosis, error) {
	if err := input.Validate(); err != nil {
		return db.BaselineDiagnosis{}, fmt.Errorf("validation error: %w", err)
	}

	profile, err := s.profileRepo.Get(ctx)
	if err != nil {
		return db.BaselineDiagnosis{}, fmt.Errorf("profile not found — create a profile first: %w", err)
	}

	result, err := s.aiService.EvaluateBaselineDiagnosis(ctx, ai.BaselineDiagnosisInput{
		Profile:          profileToAIContext(profile),
		CodingAnswer:     input.CodingAnswer,
		ConceptualAnswer: input.ConceptualAnswer,
		BehavioralAnswer: input.BehavioralAnswer,
	})
	if err != nil {
		return db.BaselineDiagnosis{}, fmt.Errorf("%w: %v", ErrAIFailed, err)
	}

	reportJSON, err := json.Marshal(result.Data)
	if err != nil {
		return db.BaselineDiagnosis{}, fmt.Errorf("failed to marshal diagnostic report: %w", err)
	}

	diagnosis, err := s.diagnosisRepo.Create(ctx, profile.ID, input.CodingAnswer, input.ConceptualAnswer, input.BehavioralAnswer, reportJSON)
	if err != nil {
		return db.BaselineDiagnosis{}, fmt.Errorf("failed to save diagnosis: %w", err)
	}

	return diagnosis, nil
}

func (s *diagnosisService) Get(ctx context.Context) (db.BaselineDiagnosis, error) {
	profile, err := s.profileRepo.Get(ctx)
	if err != nil {
		return db.BaselineDiagnosis{}, fmt.Errorf("profile not found: %w", err)
	}

	return s.diagnosisRepo.Get(ctx, profile.ID)
}

func profileToAIContext(profile db.Profile) ai.ProfileContext {
	return ai.ProfileContext{
		ExperienceYears: profile.ExperienceYears,
		TargetRole:      profile.TargetRole,
		TargetCompanies: profile.TargetCompanies,
		WeakAreasSelf:   profile.WeakAreasSelf,
		Timeline:        profile.Timeline,
	}
}
