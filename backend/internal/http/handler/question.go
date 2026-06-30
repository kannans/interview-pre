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

type QuestionHandler struct {
	questionService service.QuestionService
}

func NewQuestionHandler(questionService service.QuestionService) *QuestionHandler {
	return &QuestionHandler{questionService: questionService}
}

// List handles GET /api/v1/questions with optional type, category, and level filters.
func (h *QuestionHandler) List(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	filter := service.QuestionListFilter{
		Type:     r.URL.Query().Get("type"),
		Category: r.URL.Query().Get("category"),
		Level:    r.URL.Query().Get("level"),
	}

	result, err := h.questionService.List(r.Context(), filter)
	if err != nil {
		if isValidationError(err) {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		slog.Error("QuestionService.List failed", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}

	if err := json.NewEncoder(w).Encode(result); err != nil {
		slog.Error("Failed to encode questions list response", "error", err)
	}
}

// Get handles GET /api/v1/questions/{id}.
func (h *QuestionHandler) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := chi.URLParam(r, "id")
	var id pgtype.UUID
	if err := id.Scan(idStr); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid question id"})
		return
	}

	question, err := h.questionService.GetByID(r.Context(), id)
	if err != nil {
		if isValidationError(err) {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		if errors.Is(err, service.ErrQuestionNotFound) {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		slog.Error("QuestionService.GetByID failed", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}

	if err := json.NewEncoder(w).Encode(question); err != nil {
		slog.Error("Failed to encode question response", "error", err)
	}
}
