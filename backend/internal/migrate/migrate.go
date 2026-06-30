package migrate

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Up applies all pending SQL migrations from db/migrations before the API serves traffic.
func Up(databaseURL string) error {
	dir, err := ResolveMigrationsDir()
	if err != nil {
		return err
	}

	sourceURL, err := fileSourceURL(dir)
	if err != nil {
		return err
	}

	m, err := migrate.New(sourceURL, databaseURL)
	if err != nil {
		return fmt.Errorf("create migrator: %w", err)
	}
	defer func() {
		srcErr, dbErr := m.Close()
		if srcErr != nil {
			slog.Warn("migrate source close", "error", srcErr)
		}
		if dbErr != nil {
			slog.Warn("migrate database close", "error", dbErr)
		}
	}()

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			slog.Info("Database migrations up to date", "dir", dir)
			return nil
		}
		var dirty migrate.ErrDirty
		if errors.As(err, &dirty) {
			return fmt.Errorf(
				"database migration %d is dirty — fix the schema, then run: migrate -path %s -database \"$DATABASE_URL\" force %d",
				dirty.Version, dir, dirty.Version,
			)
		}
		return fmt.Errorf("apply migrations: %w", err)
	}

	version, dirty, verr := m.Version()
	if verr != nil && !errors.Is(verr, migrate.ErrNilVersion) {
		slog.Warn("could not read migration version after up", "error", verr)
	} else {
		slog.Info("Database migrations applied", "dir", dir, "version", version, "dirty", dirty)
	}

	return nil
}

// ResolveMigrationsDir locates db/migrations for local dev and container layouts.
func ResolveMigrationsDir() (string, error) {
	if custom := os.Getenv("MIGRATIONS_PATH"); custom != "" {
		abs, err := filepath.Abs(custom)
		if err != nil {
			return "", fmt.Errorf("resolve MIGRATIONS_PATH: %w", err)
		}
		if err := assertMigrationsDir(abs); err != nil {
			return "", err
		}
		return abs, nil
	}

	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("get working directory: %w", err)
	}

	dir := wd
	for i := 0; i < 8; i++ {
		candidate := filepath.Join(dir, "db", "migrations")
		if err := assertMigrationsDir(candidate); err == nil {
			return filepath.Abs(candidate)
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", fmt.Errorf("db/migrations not found (set MIGRATIONS_PATH or run from repo root)")
}

func assertMigrationsDir(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", path)
	}
	return nil
}

func fileSourceURL(dir string) (string, error) {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return "", fmt.Errorf("abs migrations path: %w", err)
	}
	// file source requires a forward-slash path; three slashes for absolute paths.
	return "file://" + filepath.ToSlash(abs), nil
}
