package service_test

import (
	"testing"

	"github.com/kannans/interview-prep/backend/internal/service"
)

func TestNormalizePracticeLanguage(t *testing.T) {
	lang, err := service.NormalizePracticeLanguage("")
	if err != nil || lang != "javascript" {
		t.Fatalf("empty: got %q, %v", lang, err)
	}

	lang, err = service.NormalizePracticeLanguage("Python")
	if err != nil || lang != "python" {
		t.Fatalf("python: got %q, %v", lang, err)
	}

	_, err = service.NormalizePracticeLanguage("rust")
	if err == nil {
		t.Fatal("expected error for unsupported language")
	}
}
