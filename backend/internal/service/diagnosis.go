package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/kannans/interview-prep/backend/internal/db"
	"github.com/kannans/interview-prep/backend/internal/repository"
)

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
// In Phase 7, this service will also trigger AI evaluation.
type DiagnosisService interface {
	Submit(ctx context.Context, input DiagnosisInput) (db.BaselineDiagnosis, error)
	Get(ctx context.Context) (db.BaselineDiagnosis, error)
}

type diagnosisService struct {
	profileRepo   repository.ProfileRepository
	diagnosisRepo repository.DiagnosisRepository
}

func NewDiagnosisService(profileRepo repository.ProfileRepository, diagnosisRepo repository.DiagnosisRepository) DiagnosisService {
	return &diagnosisService{
		profileRepo:   profileRepo,
		diagnosisRepo: diagnosisRepo,
	}
}

func (s *diagnosisService) Submit(ctx context.Context, input DiagnosisInput) (db.BaselineDiagnosis, error) {
	if err := input.Validate(); err != nil {
		return db.BaselineDiagnosis{}, fmt.Errorf("validation error: %w", err)
	}

	// A profile must exist before a diagnosis can be submitted
	profile, err := s.profileRepo.Get(ctx)
	if err != nil {
		return db.BaselineDiagnosis{}, fmt.Errorf("profile not found — create a profile first: %w", err)
	}

	// Placeholder diagnostic report: in Phase 7 this will be replaced with an OpenAI call
	placeholderReport := map[string]interface{}{
		"status":  "pending_ai_evaluation",
		"message": "AI evaluation will be triggered in Phase 7",
	}
	reportJSON, err := json.Marshal(placeholderReport)
	if err != nil {
		return db.BaselineDiagnosis{}, fmt.Errorf("failed to marshal report: %w", err)
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
