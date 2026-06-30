package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/kannans/interview-prep/backend/internal/service"
)

type ProgressHandler struct {
	progressService service.ProgressService
}

func NewProgressHandler(progressService service.ProgressService) *ProgressHandler {
	return &ProgressHandler{progressService: progressService}
}

// CreateSnapshot handles POST /api/v1/progress/snapshot.
func (h *ProgressHandler) CreateSnapshot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	snapshot, created, err := h.progressService.CreateSnapshot(r.Context())
	if err != nil {
		if writeProgressError(w, err) {
			return
		}
		slog.Error("ProgressService.CreateSnapshot failed", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}

	if created {
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	if err := json.NewEncoder(w).Encode(snapshot); err != nil {
		slog.Error("Failed to encode progress snapshot response", "error", err)
	}
}

// GetLatest handles GET /api/v1/progress/latest.
func (h *ProgressHandler) GetLatest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	snapshot, err := h.progressService.GetLatest(r.Context())
	if err != nil {
		if writeProgressError(w, err) {
			return
		}
		slog.Error("ProgressService.GetLatest failed", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}

	if err := json.NewEncoder(w).Encode(snapshot); err != nil {
		slog.Error("Failed to encode latest progress snapshot response", "error", err)
	}
}

// ListSnapshots handles GET /api/v1/progress/snapshots.
func (h *ProgressHandler) ListSnapshots(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query := r.URL.Query()
	result, err := h.progressService.ListSnapshots(
		r.Context(),
		query.Get("limit"),
		query.Get("start_date"),
		query.Get("end_date"),
	)
	if err != nil {
		if writeProgressError(w, err) {
			return
		}
		slog.Error("ProgressService.ListSnapshots failed", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}

	if err := json.NewEncoder(w).Encode(result); err != nil {
		slog.Error("Failed to encode progress snapshots list response", "error", err)
	}
}

func writeProgressError(w http.ResponseWriter, err error) bool {
	switch {
	case isValidationError(err):
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return true
	case errors.Is(err, service.ErrProfileMissing):
		w.WriteHeader(http.StatusPreconditionFailed)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return true
	case errors.Is(err, service.ErrProgressSnapshotNotFound):
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return true
	default:
		return false
	}
}
