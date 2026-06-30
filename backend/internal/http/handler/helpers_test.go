package handler

import (
	"errors"
	"fmt"
	"testing"

	"github.com/kannans/interview-prep/backend/internal/service"
)

func TestIsValidationError(t *testing.T) {
	if !isValidationError(fmt.Errorf("validation error: bad field")) {
		t.Fatal("expected true for validation error prefix")
	}
	if isValidationError(errors.New("other error")) {
		t.Fatal("expected false for non-validation error")
	}
}

func TestIsAIError(t *testing.T) {
	if !isAIError(fmt.Errorf("wrapped: %w", service.ErrAIFailed)) {
		t.Fatal("expected true for wrapped ErrAIFailed")
	}
	if isAIError(errors.New("other")) {
		t.Fatal("expected false")
	}
}
