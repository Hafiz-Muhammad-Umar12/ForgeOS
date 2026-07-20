#!/usr/bin/env bash
#
# check-conventional-commit.sh — validate a commit message against
# Conventional Commits (governance/08-git-workflow.md §3).
# Invoked by pre-commit as a commit-msg hook.
# Usage: check-conventional-commit.sh <commit-msg-file>

set -euo pipefail

MSG_FILE="${1:-}"
if [ -z "$MSG_FILE" ] || [ ! -f "$MSG_FILE" ]; then
  echo "usage: check-conventional-commit.sh <commit-msg-file>" >&2
  exit 1
fi

# Strip comments and blank lines, then take the first subject line.
SUBJECT="$(grep -v '^#' "$MSG_FILE" | sed '/^$/d' | head -n 1)"

PATTERN='^(feat|fix|docs|refactor|test|chore|perf|build|ci|revert)(\([a-z0-9._-]+\))?(!)?: .+'

if ! printf '%s' "$SUBJECT" | grep -Eq "$PATTERN"; then
  echo "ERROR: commit subject does not follow Conventional Commits:" >&2
  echo "  $SUBJECT" >&2
  echo "Expected format: <type>(<scope>): <subject>" >&2
  echo "Types: feat|fix|docs|refactor|test|chore|perf|build|ci|revert" >&2
  echo "Example: feat(core): add NATS JetStream bus client" >&2
  exit 1
fi

exit 0
