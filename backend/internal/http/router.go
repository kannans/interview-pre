package http

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kannans/interview-prep/backend/internal/db"
	"github.com/kannans/interview-prep/backend/internal/http/handler"
	"github.com/kannans/interview-prep/backend/internal/repository"
	"github.com/kannans/interview-prep/backend/internal/service"
)

func NewRouter(dbPool *pgxpool.Pool) *chi.Mux {
	r := chi.NewRouter()

	// Global Middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger) // Simple logger; can be replaced with structured logs later
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// CORS Setup
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://127.0.0.1:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not exceeded
	}))

	// Initializing handlers
	healthHandler := handler.NewHealthHandler(dbPool)

	// Wire up Profile module
	queries := db.New(dbPool)
	profileRepo := repository.NewProfileRepository(queries)
	profileSvc := service.NewProfileService(profileRepo)
	profileHandler := handler.NewProfileHandler(profileSvc)

	// Wire up Diagnosis module
	diagnosisRepo := repository.NewDiagnosisRepository(queries)
	diagnosisSvc := service.NewDiagnosisService(profileRepo, diagnosisRepo)
	diagnosisHandler := handler.NewDiagnosisHandler(diagnosisSvc)

	// Wire up Recovery Plan module
	planRepo := repository.NewPlanRepository(dbPool, queries)
	planSvc := service.NewPlanService(profileRepo, diagnosisRepo, planRepo)
	planHandler := handler.NewPlanHandler(planSvc)

	// Wire up Practice Questions module
	questionRepo := repository.NewQuestionRepository(queries)
	questionSvc := service.NewQuestionService(questionRepo)
	questionHandler := handler.NewQuestionHandler(questionSvc)

	// Wire up Practice Attempts module
	attemptRepo := repository.NewAttemptRepository(queries)
	attemptSvc := service.NewAttemptService(profileRepo, questionRepo, attemptRepo)
	attemptHandler := handler.NewAttemptHandler(attemptSvc)

	// Health Route
	r.Get("/health", healthHandler.Check)

	// API v1 Routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/profile", profileHandler.Upsert)
		r.Get("/profile", profileHandler.Get)

		r.Post("/diagnosis", diagnosisHandler.Submit)
		r.Get("/diagnosis", diagnosisHandler.Get)

		r.Post("/plan", planHandler.Generate)
		r.Get("/plan", planHandler.Get)
		r.Get("/plan/today", planHandler.GetToday)
		r.Patch("/plan/tasks/{id}/complete", planHandler.CompleteTask)

		r.Get("/questions", questionHandler.List)
		r.Get("/questions/{id}", questionHandler.Get)

		r.Post("/attempts", attemptHandler.Submit)
		r.Get("/attempts", attemptHandler.List)
		r.Get("/attempts/{id}", attemptHandler.Get)
	})

	// Fallback Route for undefined paths
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"status": "error", "message": "resource not found"}`))
	})

	return r
}
