package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	Env          string
	DatabaseURL  string
	OpenAIAPIKey string
	OpenAIModel  string
}

func Load() (*Config, error) {
	// Attempt to load .env file; ignore error if it doesn't exist (e.g. in containerized envs)
	if err := godotenv.Load(); err != nil {
		slog.Info("No .env file found, relying on system environment variables")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is required")
	}

	openaiAPIKey := os.Getenv("OPENAI_API_KEY")
	if openaiAPIKey == "" {
		slog.Warn("OPENAI_API_KEY is not set. AI features will fail.")
	}

	openaiModel := os.Getenv("OPENAI_MODEL")
	if openaiModel == "" {
		openaiModel = "gpt-4o-mini"
	}

	return &Config{
		Port:         port,
		Env:          env,
		DatabaseURL:  databaseURL,
		OpenAIAPIKey: openaiAPIKey,
		OpenAIModel:  openaiModel,
	}, nil
}
