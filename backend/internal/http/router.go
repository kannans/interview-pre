package http

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kannans/interview-prep/backend/internal/ai"
	"github.com/kannans/interview-prep/backend/internal/db"
	"github.com/kannans/interview-prep/backend/internal/http/handler"
	"github.com/kannans/interview-prep/backend/internal/repository"
	"github.com/kannans/interview-prep/backend/internal/runner"
	"github.com/kannans/interview-prep/backend/internal/service"
)

func NewRouter(dbPool *pgxpool.Pool, aiService ai.Service, corsAllowedOrigins []string) *chi.Mux {
	r := chi.NewRouter()

	// Global Middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(120 * time.Second))

	// CORS Setup
	if len(corsAllowedOrigins) == 0 {
		corsAllowedOrigins = []string{"http://localhost:3000", "http://127.0.0.1:3000"}
	}
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   corsAllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	healthHandler := handler.NewHealthHandler(dbPool)

	queries := db.New(dbPool)
	profileRepo := repository.NewProfileRepository(queries)
	profileSvc := service.NewProfileService(profileRepo)
	profileHandler := handler.NewProfileHandler(profileSvc)

	diagnosisRepo := repository.NewDiagnosisRepository(queries)
	diagnosisSvc := service.NewDiagnosisService(profileRepo, diagnosisRepo, aiService)
	diagnosisHandler := handler.NewDiagnosisHandler(diagnosisSvc)

	planRepo := repository.NewPlanRepository(dbPool, queries)
	planSvc := service.NewPlanService(profileRepo, diagnosisRepo, planRepo, aiService)
	planHandler := handler.NewPlanHandler(planSvc)

	questionRepo := repository.NewQuestionRepository(queries)
	questionSvc := service.NewQuestionService(questionRepo)
	questionHandler := handler.NewQuestionHandler(questionSvc)

	attemptRepo := repository.NewAttemptRepository(queries)
	attemptSvc := service.NewAttemptService(profileRepo, questionRepo, attemptRepo)
	attemptHandler := handler.NewAttemptHandler(attemptSvc)

	evaluationRepo := repository.NewEvaluationRepository(queries)
	evaluationSvc := service.NewEvaluationService(profileRepo, attemptRepo, questionRepo, evaluationRepo, aiService)
	evaluationHandler := handler.NewEvaluationHandler(evaluationSvc)

	progressRepo := repository.NewProgressRepository(queries)
	progressSvc := service.NewProgressService(profileRepo, planRepo, progressRepo, aiService)
	progressHandler := handler.NewProgressHandler(progressSvc)

	runSvc := service.NewRunService(runner.NewLocalExecutor())
	runHandler := handler.NewRunHandler(runSvc)

	r.Get("/health", healthHandler.Check)

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

		r.Post("/run", runHandler.Run)

		r.Post("/attempts", attemptHandler.Submit)
		r.Get("/attempts", attemptHandler.List)
		r.Post("/attempts/{id}/evaluate", evaluationHandler.Evaluate)
		r.Get("/attempts/{id}/evaluation", evaluationHandler.GetByAttempt)
		r.Get("/attempts/{id}", attemptHandler.Get)

		r.Post("/progress/snapshot", progressHandler.CreateSnapshot)
		r.Get("/progress/latest", progressHandler.GetLatest)
		r.Get("/progress/snapshots", progressHandler.ListSnapshots)
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"status": "error", "message": "resource not found"}`))
	})

	return r
}
