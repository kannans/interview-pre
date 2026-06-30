package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/kannans/interview-prep/backend/internal/service"
)

type DiagnosisHandler struct {
	diagnosisService service.DiagnosisService
}

func NewDiagnosisHandler(diagnosisService service.DiagnosisService) *DiagnosisHandler {
	return &DiagnosisHandler{diagnosisService: diagnosisService}
}

// Submit handles POST /api/v1/diagnosis
func (h *DiagnosisHandler) Submit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var input service.DiagnosisInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		slog.Error("Failed to decode diagnosis request body", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	diagnosis, err := h.diagnosisService.Submit(r.Context(), input)
	if err != nil {
		slog.Error("DiagnosisService.Submit failed", "error", err)
		if isValidationError(err) {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		if isAIError(err) {
			w.WriteHeader(http.StatusBadGateway)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "ai service unavailable — try again later"})
			return
		}
		if errors.Is(err, pgx.ErrNoRows) {
			w.WriteHeader(http.StatusPreconditionFailed)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "profile not found — create a profile first"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(diagnosis); err != nil {
		slog.Error("Failed to encode diagnosis response", "error", err)
	}
}

// Get handles GET /api/v1/diagnosis
func (h *DiagnosisHandler) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	diagnosis, err := h.diagnosisService.Get(r.Context())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "no diagnosis found"})
			return
		}
		slog.Error("DiagnosisService.Get failed", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}

	if err := json.NewEncoder(w).Encode(diagnosis); err != nil {
		slog.Error("Failed to encode diagnosis response", "error", err)
	}
}


