#!/usr/bin/env bash
# Start backend + frontend for local development.
# - Frontend: Next.js dev server (hot reload built in)
# - Backend: air (recommended) or go run with a file watcher
#
# Usage: ./runserver.sh
# Stop:  Ctrl+C

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
FRONTEND_PID=""
BACKEND_PID=""

log() {
  printf '\033[1;36m[runserver]\033[0m %s\n' "$*"
}

warn() {
  printf '\033[1;33m[runserver]\033[0m %s\n' "$*" >&2
}

die() {
  printf '\033[1;31m[runserver]\033[0m %s\n' "$*" >&2
  exit 1
}

cleanup() {
  log "Shutting down..."
  [[ -n "$BACKEND_PID" ]] && kill "$BACKEND_PID" 2>/dev/null || true
  [[ -n "$FRONTEND_PID" ]] && kill "$FRONTEND_PID" 2>/dev/null || true
  wait 2>/dev/null || true
}

trap cleanup EXIT INT TERM

ensure_prereqs() {
  command -v go >/dev/null || die "Go is not installed"
  command -v npm >/dev/null || die "npm is not installed"

  if [[ ! -f "$ROOT/backend/.env" ]]; then
    if [[ -f "$ROOT/backend/.env.example" ]]; then
      cp "$ROOT/backend/.env.example" "$ROOT/backend/.env"
      warn "Created backend/.env from .env.example — add OPENAI_API_KEY before using AI endpoints."
    else
      die "Missing backend/.env"
    fi
  fi

  if [[ ! -f "$ROOT/frontend/.env.local" ]]; then
    if [[ -f "$ROOT/frontend/.env.local.example" ]]; then
      cp "$ROOT/frontend/.env.local.example" "$ROOT/frontend/.env.local"
      log "Created frontend/.env.local from .env.local.example"
    fi
  fi

  if [[ ! -d "$ROOT/frontend/node_modules" ]]; then
    log "Installing frontend dependencies..."
    (cd "$ROOT/frontend" && npm install)
  fi
}

backend_source_hash() {
  {
    find "$ROOT/backend" -name '*.go' -type f 2>/dev/null || true
    find "$ROOT/db/migrations" -name '*.sql' -type f 2>/dev/null || true
  } | LC_ALL=C sort | while IFS= read -r file; do
    stat -f '%m %N' "$file" 2>/dev/null || stat -c '%Y %n' "$file" 2>/dev/null
  done | shasum | awk '{print $1}'
}

stop_backend() {
  if [[ -n "$BACKEND_PID" ]]; then
    kill "$BACKEND_PID" 2>/dev/null || true
    wait "$BACKEND_PID" 2>/dev/null || true
    BACKEND_PID=""
  fi
}

start_backend_gorun() {
  (cd "$ROOT/backend" && go run ./cmd/api) &
  BACKEND_PID=$!
}

watch_backend_gorun() {
  log "Starting backend with go run + file watcher (http://localhost:8080)..."
  log "Tip: install air for faster reloads: go install github.com/air-verse/air@latest"
  start_backend_gorun

  local hash
  hash="$(backend_source_hash)"
  while true; do
    if ! kill -0 "$FRONTEND_PID" 2>/dev/null; then
      die "Frontend exited unexpectedly"
    fi
    if [[ -n "$BACKEND_PID" ]] && ! kill -0 "$BACKEND_PID" 2>/dev/null; then
      die "Backend exited unexpectedly"
    fi

    sleep 1
    local next_hash
    next_hash="$(backend_source_hash)"
    if [[ "$next_hash" != "$hash" ]]; then
      hash="$next_hash"
      log "Backend code changed — restarting..."
      stop_backend
      start_backend_gorun
    fi
  done
}

start_frontend() {
  log "Starting frontend (http://localhost:3000)..."
  (cd "$ROOT/frontend" && npm run dev) &
  FRONTEND_PID=$!
}

run_backend() {
  if command -v air >/dev/null; then
    log "Starting backend with air (http://localhost:8080)..."
    (cd "$ROOT/backend" && air) &
    BACKEND_PID=$!
    wait "$BACKEND_PID"
    return
  fi
  watch_backend_gorun
}

main() {
  cd "$ROOT"
  ensure_prereqs
  start_frontend
  log "Ready — frontend http://localhost:3000 · API http://localhost:8080/api/v1"
  log "Press Ctrl+C to stop both servers."
  run_backend
}

main "$@"
