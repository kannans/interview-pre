package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kannans/interview-prep/backend/internal/config"
	"github.com/kannans/interview-prep/backend/internal/ai"
	"github.com/kannans/interview-prep/backend/internal/db"
	internalHttp "github.com/kannans/interview-prep/backend/internal/http"
)

func main() {
	// Initialize structured logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	slog.Info("Starting Interview Preparation Coach Backend...")

	// 1. Load Config
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Configuration load failed", "error", err)
		os.Exit(1)
	}

	// 2. Establish Database Connection Pool
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbPool, err := db.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("Database connection failed", "error", err)
		os.Exit(1)
	}
	defer dbPool.Close()

	promptRegistry, err := ai.NewRegistry()
	if err != nil {
		slog.Error("Failed to load AI prompts", "error", err)
		os.Exit(1)
	}
	openaiClient := ai.NewClient(ai.ClientConfig{
		APIKey: cfg.OpenAIAPIKey,
		Model:  cfg.OpenAIModel,
	})
	aiService := ai.NewService(openaiClient, promptRegistry)

	// 3. Set up HTTP router & server
	router := internalHttp.NewRouter(dbPool, aiService)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 4. Run HTTP server in background goroutine
	go func() {
		slog.Info(fmt.Sprintf("HTTP Server running on port %s", cfg.Port), "env", cfg.Env)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("HTTP Server failed to listen", "error", err)
			os.Exit(1)
		}
	}()

	// 5. Intercept Interrupt signals for graceful shutdown
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	sig := <-shutdownChan
	slog.Info("Shutdown signal received", "signal", sig.String())

	// Shutdown context with 5s timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("Forced HTTP Server shutdown", "error", err)
	} else {
		slog.Info("HTTP Server cleanly stopped")
	}

	slog.Info("Graceful shutdown sequence complete")
}
