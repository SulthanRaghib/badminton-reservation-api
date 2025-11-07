#!/bin/sh
# Entry point for container: run migrations and seed, then start the main app

set -e

echo "[entrypoint] Running database migrations (if any)..."
if [ -x ./migrate ]; then
  ./migrate || echo "[entrypoint] migrate exited with non-zero status"
else
  echo "[entrypoint] migrate binary not found, skipping"
fi

echo "[entrypoint] Running database seed (if present)..."
if [ -x ./seed ]; then
  ./seed || echo "[entrypoint] seed exited with non-zero status"
else
  echo "[entrypoint] seed binary not found, skipping"
fi

echo "[entrypoint] Starting main application"
exec ./main
