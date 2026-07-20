# bootstrap.ps1 — DevOS local development environment setup (Windows).
#
# Idempotent. Run from the repository root:
#   .\scripts\bootstrap.ps1
#
# See Sprint 0, Component 1 (Monorepo Foundation) and
# governance/08-git-workflow.md.

$ErrorActionPreference = "Stop"

$RepoRoot = (Resolve-Path (Join-Path $PSScriptRoot "..")).Path
Set-Location $RepoRoot

function Log($msg) { Write-Host "==> $msg" -ForegroundColor Cyan }
function Err($msg) { Write-Host "ERROR: $msg" -ForegroundColor Red }

Log "DevOS bootstrap starting (root: $RepoRoot)"

# 1. mise (toolchain manager)
if (-not (Get-Command mise -ErrorAction SilentlyContinue)) {
  Err "mise is not installed. Install from https://mise.jdx.dev and re-run."
  exit 1
}
Log "Installing pinned toolchains via mise..."
mise install

# 2. Local environment file
if (-not (Test-Path .env)) {
  if (Test-Path .env.example) {
    Log "Creating .env from .env.example..."
    Copy-Item .env.example .env
  } else {
    Err ".env.example not found; cannot create .env."
    exit 1
  }
} else {
  Log ".env already exists; leaving it untouched."
}

# 3. Go module dependencies
if (Get-Command go -ErrorAction SilentlyContinue) {
  Log "Downloading Go module dependencies..."
  go work sync | Out-Null
  foreach ($mod in @("core", "packages/go", "apps/gateway")) {
    if (Test-Path (Join-Path $mod "go.mod")) {
      Push-Location $mod
      try { go mod download } finally { Pop-Location }
    }
  }
} else {
  Log "go not found on PATH (mise should have provided it). Skipping go mod download."
}

# 4. Git hooks (pre-commit + commit-msg)
if (Get-Command pre-commit -ErrorAction SilentlyContinue) {
  Log "Installing git hooks..."
  pre-commit install --install-hooks
  pre-commit install --hook-type commit-msg --install-hooks
} else {
  Log "pre-commit not found; skipping git hook installation."
}

# 5. Docker availability
if (Get-Command docker -ErrorAction SilentlyContinue) {
  Log "Docker detected. Start local infrastructure later with: task infra-up"
} else {
  Log "Docker not detected. Install Docker Desktop to run local infrastructure."
}

Log "Bootstrap complete. Try: task --list"
