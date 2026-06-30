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

type AttemptHandler struct {
	attemptService service.AttemptService
}

func NewAttemptHandler(attemptService service.AttemptService) *AttemptHandler {
	return &AttemptHandler{attemptService: attemptService}
}

// Submit handles POST /api/v1/attempts.
func (h *AttemptHandler) Submit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var input service.AttemptInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		slog.Error("Failed to decode attempt request body", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	attempt, err := h.attemptService.Submit(r.Context(), input)
	if err != nil {
		if writeAttemptError(w, err) {
			return
		}
		slog.Error("AttemptService.Submit failed", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(attempt); err != nil {
		slog.Error("Failed to encode attempt response", "error", err)
	}
}

// List handles GET /api/v1/attempts with optional question_id filter.
func (h *AttemptHandler) List(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	result, err := h.attemptService.List(r.Context(), r.URL.Query().Get("question_id"))
	if err != nil {
		if writeAttemptError(w, err) {
			return
		}
		slog.Error("AttemptService.List failed", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}

	if err := json.NewEncoder(w).Encode(result); err != nil {
		slog.Error("Failed to encode attempts list response", "error", err)
	}
}

// Get handles GET /api/v1/attempts/{id}.
func (h *AttemptHandler) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := chi.URLParam(r, "id")
	var id pgtype.UUID
	if err := id.Scan(idStr); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid attempt id"})
		return
	}

	attempt, err := h.attemptService.GetByID(r.Context(), id)
	if err != nil {
		if writeAttemptError(w, err) {
			return
		}
		slog.Error("AttemptService.GetByID failed", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}

	if err := json.NewEncoder(w).Encode(attempt); err != nil {
		slog.Error("Failed to encode attempt response", "error", err)
	}
}

func writeAttemptError(w http.ResponseWriter, err error) bool {
	switch {
	case isValidationError(err):
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return true
	case errors.Is(err, service.ErrProfileMissing):
		w.WriteHeader(http.StatusPreconditionFailed)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return true
	case errors.Is(err, service.ErrQuestionNotFound),
		errors.Is(err, service.ErrAttemptNotFound),
		errors.Is(err, service.ErrAttemptNotOwned):
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return true
	default:
		return false
	}
}
