package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kannans/interview-prep/backend/internal/service"
)

type EvaluationHandler struct {
	evaluationService service.EvaluationService
}

func NewEvaluationHandler(evaluationService service.EvaluationService) *EvaluationHandler {
	return &EvaluationHandler{evaluationService: evaluationService}
}

// Evaluate handles POST /api/v1/attempts/{id}/evaluate.
func (h *EvaluationHandler) Evaluate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	attemptID, ok := parseAttemptID(w, r)
	if !ok {
		return
	}

	evaluation, err := h.evaluationService.Evaluate(r.Context(), attemptID)
	if err != nil {
		if writeEvaluationError(w, err) {
			return
		}
		slog.Error("EvaluationService.Evaluate failed", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(evaluation); err != nil {
		slog.Error("Failed to encode evaluation response", "error", err)
	}
}

// GetByAttempt handles GET /api/v1/attempts/{id}/evaluation.
func (h *EvaluationHandler) GetByAttempt(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	attemptID, ok := parseAttemptID(w, r)
	if !ok {
		return
	}

	evaluation, err := h.evaluationService.GetByAttemptID(r.Context(), attemptID)
	if err != nil {
		if writeEvaluationError(w, err) {
			return
		}
		slog.Error("EvaluationService.GetByAttemptID failed", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}

	if err := json.NewEncoder(w).Encode(evaluation); err != nil {
		slog.Error("Failed to encode evaluation response", "error", err)
	}
}

func parseAttemptID(w http.ResponseWriter, r *http.Request) (pgtype.UUID, bool) {
	idStr := chi.URLParam(r, "id")
	var id pgtype.UUID
	if err := id.Scan(idStr); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid attempt id"})
		return pgtype.UUID{}, false
	}
	return id, true
}

func writeEvaluationError(w http.ResponseWriter, err error) bool {
	switch {
	case isValidationError(err):
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return true
	case errors.Is(err, service.ErrProfileMissing):
		w.WriteHeader(http.StatusPreconditionFailed)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return true
	case errors.Is(err, service.ErrAttemptNotFound),
		errors.Is(err, service.ErrAttemptNotOwned),
		errors.Is(err, service.ErrQuestionNotFound),
		errors.Is(err, service.ErrEvaluationNotFound):
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return true
	case errors.Is(err, service.ErrEvaluationAlreadyExists):
		w.WriteHeader(http.StatusConflict)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return true
	default:
		return false
	}
}
