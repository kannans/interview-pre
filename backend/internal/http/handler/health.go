package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type HealthHandler struct {
	dbPool *pgxpool.Pool
}

func NewHealthHandler(dbPool *pgxpool.Pool) *HealthHandler {
	return &HealthHandler{dbPool: dbPool}
}

type HealthResponse struct {
	Status   string `json:"status"`
	Database string `json:"database"`
}

func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	dbStatus := "connected"
	statusCode := http.StatusOK

	if h.dbPool == nil {
		dbStatus = "not initialized"
		statusCode = http.StatusServiceUnavailable
	} else {
		// Ping database with a 2-second timeout context
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		if err := h.dbPool.Ping(ctx); err != nil {
			slog.Error("Database health check ping failed", "error", err)
			dbStatus = "disconnected"
			statusCode = http.StatusServiceUnavailable
		}
	}

	resp := HealthResponse{
		Status:   "ok",
		Database: dbStatus,
	}
	if statusCode != http.StatusOK {
		resp.Status = "error"
	}

	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error("Failed to encode health response", "error", err)
	}
}
