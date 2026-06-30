package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/kannans/interview-prep/backend/internal/service"
)

type ProfileHandler struct {
	profileService service.ProfileService
}

func NewProfileHandler(profileService service.ProfileService) *ProfileHandler {
	return &ProfileHandler{profileService: profileService}
}

// Upsert handles POST /api/v1/profile
// Creates the profile if it doesn't exist, or updates the existing one.
func (h *ProfileHandler) Upsert(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var input service.ProfileInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		slog.Error("Failed to decode profile request body", "error", err)
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	result, err := h.profileService.Upsert(r.Context(), input)
	if err != nil {
		slog.Error("ProfileService.Upsert failed", "error", err)
		// Return validation errors as 400 Bad Request
		if isValidationError(err) {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}

	if result.Created {
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	if err := json.NewEncoder(w).Encode(result.Profile); err != nil {
		slog.Error("Failed to encode profile response", "error", err)
	}
}

// Get handles GET /api/v1/profile
// Returns the single candidate profile or 404 if not yet created.
func (h *ProfileHandler) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	profile, err := h.profileService.Get(r.Context())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "profile not found"})
			return
		}
		slog.Error("ProfileService.Get failed", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}

	if err := json.NewEncoder(w).Encode(profile); err != nil {
		slog.Error("Failed to encode profile response", "error", err)
	}
}


