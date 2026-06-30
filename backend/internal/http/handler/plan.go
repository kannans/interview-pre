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

type PlanHandler struct {
	planService service.PlanService
}

func NewPlanHandler(planService service.PlanService) *PlanHandler {
	return &PlanHandler{planService: planService}
}

// Generate handles POST /api/v1/plan — creates or replaces the active recovery plan.
func (h *PlanHandler) Generate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	result, err := h.planService.Generate(r.Context())
	if err != nil {
		slog.Error("PlanService.Generate failed", "error", err)
		if writePlanError(w, err) {
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		slog.Error("Failed to encode plan response", "error", err)
	}
}

// Get handles GET /api/v1/plan — returns the active plan with all tasks.
func (h *PlanHandler) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	result, err := h.planService.Get(r.Context())
	if err != nil {
		if writePlanError(w, err) {
			return
		}
		slog.Error("PlanService.Get failed", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}

	if err := json.NewEncoder(w).Encode(result); err != nil {
		slog.Error("Failed to encode plan response", "error", err)
	}
}

// GetToday handles GET /api/v1/plan/today — returns the first incomplete task.
func (h *PlanHandler) GetToday(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	result, err := h.planService.GetToday(r.Context())
	if err != nil {
		if writePlanError(w, err) {
			return
		}
		slog.Error("PlanService.GetToday failed", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}

	if err := json.NewEncoder(w).Encode(result); err != nil {
		slog.Error("Failed to encode today task response", "error", err)
	}
}

// CompleteTask handles PATCH /api/v1/plan/tasks/{id}/complete.
func (h *PlanHandler) CompleteTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	taskIDStr := chi.URLParam(r, "id")
	var taskID pgtype.UUID
	if err := taskID.Scan(taskIDStr); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid task id"})
		return
	}

	task, err := h.planService.CompleteTask(r.Context(), taskID)
	if err != nil {
		slog.Error("PlanService.CompleteTask failed", "error", err)
		if writePlanError(w, err) {
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}

	if err := json.NewEncoder(w).Encode(task); err != nil {
		slog.Error("Failed to encode completed task response", "error", err)
	}
}

func writePlanError(w http.ResponseWriter, err error) bool {
	switch {
	case errors.Is(err, service.ErrProfileMissing),
		errors.Is(err, service.ErrDiagnosisMissing):
		w.WriteHeader(http.StatusPreconditionFailed)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return true
	case errors.Is(err, service.ErrPlanNotFound),
		errors.Is(err, service.ErrNoTasksRemaining):
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return true
	case errors.Is(err, service.ErrTaskNotInPlan):
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return true
	case errors.Is(err, service.ErrTaskAlreadyDone):
		w.WriteHeader(http.StatusConflict)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return true
	default:
		return false
	}
}
