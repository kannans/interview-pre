package migrate_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kannans/interview-prep/backend/internal/migrate"
)

func TestResolveMigrationsDirFromRepoRoot(t *testing.T) {
	root, err := findRepoRoot(t)
	if err != nil {
		t.Fatalf("find repo root: %v", err)
	}

	chdir(t, root)
	dir, err := migrate.ResolveMigrationsDir()
	if err != nil {
		t.Fatalf("ResolveMigrationsDir: %v", err)
	}

	want := filepath.Join(root, "db", "migrations")
	if dir != want {
		t.Fatalf("got %q want %q", dir, want)
	}
}

func TestResolveMigrationsDirFromBackendDir(t *testing.T) {
	root, err := findRepoRoot(t)
	if err != nil {
		t.Fatalf("find repo root: %v", err)
	}

	chdir(t, filepath.Join(root, "backend"))
	dir, err := migrate.ResolveMigrationsDir()
	if err != nil {
		t.Fatalf("ResolveMigrationsDir: %v", err)
	}

	want := filepath.Join(root, "db", "migrations")
	if dir != want {
		t.Fatalf("got %q want %q", dir, want)
	}
}

func chdir(t *testing.T, dir string) {
	t.Helper()
	prev, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir %s: %v", dir, err)
	}
	t.Cleanup(func() { _ = os.Chdir(prev) })
}

func findRepoRoot(t *testing.T) (string, error) {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	dir := wd
	for i := 0; i < 8; i++ {
		if _, err := os.Stat(filepath.Join(dir, "db", "migrations")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", os.ErrNotExist
}
