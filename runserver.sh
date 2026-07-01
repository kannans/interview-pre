#!/usr/bin/env bash
# Start backend + frontend for local development.
# - Frontend: Next.js dev server (hot reload built in)
# - Backend: air (recommended) or go run with a file watcher
#
# Usage: ./runserver.sh
# Stop:  Ctrl+C

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
FRONTEND_PORT="${FRONTEND_PORT:-3000}"
BACKEND_PORT="${BACKEND_PORT:-8080}"
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

pids_on_port() {
  local port="$1"
  lsof -tiTCP:"$port" -sTCP:LISTEN 2>/dev/null || lsof -ti :"$port" 2>/dev/null || true
}

free_port() {
  local port="$1"
  local pids pid

  pids="$(pids_on_port "$port")"
  [[ -z "$pids" ]] && return 0

  log "Freeing port $port (PID: $(echo "$pids" | tr '\n' ' ' | sed 's/ $//'))..."
  while IFS= read -r pid; do
    [[ -z "$pid" ]] && continue
    kill -TERM "$pid" 2>/dev/null || true
  done <<< "$pids"

  sleep 0.5

  pids="$(pids_on_port "$port")"
  while IFS= read -r pid; do
    [[ -z "$pid" ]] && continue
    kill -KILL "$pid" 2>/dev/null || true
  done <<< "$pids"
}

cleanup() {
  log "Shutting down..."
  [[ -n "$BACKEND_PID" ]] && kill "$BACKEND_PID" 2>/dev/null || true
  [[ -n "$FRONTEND_PID" ]] && kill "$FRONTEND_PID" 2>/dev/null || true
  free_port "$BACKEND_PORT"
  free_port "$FRONTEND_PORT"
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
  free_port "$BACKEND_PORT"
}

start_backend_gorun() {
  (cd "$ROOT/backend" && go run ./cmd/api) &
  BACKEND_PID=$!
}

start_backend_air() {
  (cd "$ROOT/backend" && air) &
  BACKEND_PID=$!
}

restart_backend() {
  local reason="$1"
  warn "$reason — restarting backend in 2s (frontend keeps running)..."
  stop_backend
  sleep 2
}

watch_backend_gorun() {
  log "Starting backend with go run + file watcher (http://localhost:$BACKEND_PORT)..."
  log "Tip: install air for faster reloads: go install github.com/air-verse/air@latest"
  start_backend_gorun

  local hash
  hash="$(backend_source_hash)"
  while true; do
    if ! kill -0 "$FRONTEND_PID" 2>/dev/null; then
      die "Frontend exited unexpectedly"
    fi
    if [[ -n "$BACKEND_PID" ]] && ! kill -0 "$BACKEND_PID" 2>/dev/null; then
      restart_backend "Backend exited unexpectedly"
      start_backend_gorun
      hash="$(backend_source_hash)"
      continue
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
  log "Starting frontend (http://localhost:$FRONTEND_PORT)..."
  (cd "$ROOT/frontend" && npm run dev) &
  FRONTEND_PID=$!
}

run_backend_air() {
  log "Starting backend with air (http://localhost:$BACKEND_PORT)..."
  start_backend_air
  while true; do
    if ! kill -0 "$FRONTEND_PID" 2>/dev/null; then
      die "Frontend exited unexpectedly"
    fi

    wait "$BACKEND_PID" || true
    BACKEND_PID=""
    restart_backend "Backend (air) exited"
    start_backend_air
  done
}

run_backend() {
  if command -v air >/dev/null; then
    run_backend_air
    return
  fi
  watch_backend_gorun
}

main() {
  cd "$ROOT"
  ensure_prereqs
  free_port "$FRONTEND_PORT"
  free_port "$BACKEND_PORT"
  start_frontend
  log "Ready — frontend http://localhost:$FRONTEND_PORT · API http://localhost:$BACKEND_PORT/api/v1"
  log "Press Ctrl+C to stop both servers."
  run_backend
}

main "$@"
