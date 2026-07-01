package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/kannans/interview-prep/backend/internal/service"
)

type RunHandler struct {
	runService service.RunService
}

func NewRunHandler(runService service.RunService) *RunHandler {
	return &RunHandler{runService: runService}
}

// Run handles POST /api/v1/run.
func (h *RunHandler) Run(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var input service.RunInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		slog.Error("Failed to decode run request body", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	result, err := h.runService.Run(r.Context(), input)
	if err != nil {
		if isValidationError(err) {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		slog.Error("RunService.Run failed", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}

	if err := json.NewEncoder(w).Encode(result); err != nil {
		slog.Error("Failed to encode run response", "error", err)
	}
}
