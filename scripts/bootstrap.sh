#!/usr/bin/env bash
#
# bootstrap.sh — DevOS local development environment setup.
#
# Idempotent. Run from the repository root:
#   ./scripts/bootstrap.sh
#
# Installs pinned toolchains via mise, creates a local .env from the example,
# downloads Go module dependencies, and installs git hooks.
# See Sprint 0, Component 1 (Monorepo Foundation) and
# governance/08-git-workflow.md.

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$REPO_ROOT"

log() { printf '\033[1;34m==>\033[0m %s\n' "$*"; }
err() { printf '\033[1;31mERROR:\033[0m %s\n' "$*" >&2; }

log "DevOS bootstrap starting (root: $REPO_ROOT)"

# 1. mise (toolchain manager)
if ! command -v mise >/dev/null 2>&1; then
  err "mise is not installed. Install it from https://mise.jdx.dev and re-run."
  err "Quick install: curl https://mise.jdx.dev/install.sh | sh"
  exit 1
fi
log "Installing pinned toolchains via mise..."
mise install

# 2. Local environment file
if [ ! -f .env ]; then
  if [ -f .env.example ]; then
    log "Creating .env from .env.example..."
    cp .env.example .env
  else
    err ".env.example not found; cannot create .env."
    exit 1
  fi
else
  log ".env already exists; leaving it untouched."
fi

# 3. Go module dependencies
if command -v go >/dev/null 2>&1; then
  log "Downloading Go module dependencies..."
  go work sync || true
  for mod in core packages/go apps/gateway; do
    if [ -f "$mod/go.mod" ]; then
      (cd "$mod" && go mod download)
    fi
  done
else
  log "go not found on PATH (mise should have provided it). Skipping go mod download."
fi

# 4. Git hooks (pre-commit + commit-msg)
if command -v pre-commit >/dev/null 2>&1; then
  log "Installing git hooks..."
  pre-commit install --install-hooks
  pre-commit install --hook-type commit-msg --install-hooks
else
  log "pre-commit not found; skipping git hook installation."
fi

# 5. Docker availability (infrastructure is a later component)
if command -v docker >/dev/null 2>&1; then
  log "Docker detected. Start local infrastructure later with: task infra-up"
else
  log "Docker not detected. Install Docker Desktop to run local infrastructure."
fi

log "Bootstrap complete. Try: task --list"
