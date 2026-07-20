# DevOS — Coding Standards

> **Status:** DRAFT — For Approval (no production code until ratified)
> **Version:** 1.0-draft
> **Owner:** Head of Engineering
> **Companion:** Engineering Standards (`04-engineering-standards.md`), Repo (`06-repository-standards.md`)

---

## 1. General Principles (all languages)

1. **Readability over cleverness.** Code is read more than written.
2. **Single Responsibility.** Small, focused functions/modules.
3. **Explicit over implicit.** No hidden side effects; fail loud, not silent.
4. **No dead code.** No commented-out blocks, no unused exports.
5. **Types are contracts.** Prefer strong, explicit types over `any`/loose typing.
6. **No hardcoded config or secrets.** Everything via injected config / secret proxy (Constitution T4, T10).

## 2. Naming

| Concept | Convention |
|---------|-----------|
| Functions/methods | `camelCase` (Go: `PascalCase` exported, `camelCase` private) |
| Types/interfaces | `PascalCase` |
| Constants | `UPPER_SNAKE_CASE` (Go: `PascalCase` exported) |
| Files | `kebab-case` or language idiomatic (`snake_case` py, `camelCase` ts) |
| DB columns/tables | `snake_case` |
| Event types | `domain.entity.action` (e.g., `intent.created`) |

## 3. Functions & Modules

- Keep functions ≤ ~40 lines where reasonable; extract, don't nest deeply.
- Return errors explicitly; do not swallow them.
- Pure logic separated from I/O for testability.

## 4. Error Handling

1. Catch at boundaries; let errors propagate with context (wrap, not replace).
2. Distinguish **retryable** (429/5xx/timeout) from **deterministic** (compile error, invalid input).
3. Never return raw provider errors to users containing secrets/stack internals.
4. Every error logged with `traceId` + `orgId` + `intentId` for correlation.

## 5. Logging & Secrets (T10, Eng §6)

1. Structured JSON logs only.
2. **Never log:** raw secrets, full agent prompts containing secrets, PII beyond what's needed.
3. Log levels: `debug` (dev), `info` (flow), `warn` (retryable), `error` (failure with action).
4. Include `traceId` on every log line.

## 6. Language-Specific

### 6.1 Go (backend services, ingress, gateway, workspace-mgr)
- Format: `gofmt` + `goimports`. Lint: `golangci-lint`.
- Errors: wrap with `fmt.Errorf("...: %w", err)`.
- Context: propagate `context.Context` as first param.
- No `panic` in request paths.
- Interfaces defined at consumer side.

### 6.2 Python (agent runtime, provider adapters)
- Format: `black` + `isort`. Lint: `ruff`. Type-check: `mypy` (strict-ish).
- Pin deps; virtualenv/uv per service.
- Async: `asyncio`/`anyio`; never block the event loop on I/O.
- Type hints on all public functions.

### 6.3 TypeScript (frontend, SDK, UI-kit)
- Format: `prettier`. Lint: `eslint` (strict). Type: `tsc --strict`.
- No `any`; use `unknown` + narrowing.
- React: functional components + hooks; no class components.
- State via stores; no prop-drilling beyond 2 levels.

### 6.4 SQL (PostgreSQL)
- Format: `sqlfluff` (dialect postgres). Snake_case.
- Explicit column lists (no `SELECT *`) in app queries.
- Migrations reviewed for zero-downtime (Spec `/specs/04-database/`).

## 7. Comments & Documentation

1. **Why, not what.** Comment non-obvious intent, invariants, and tradeoffs.
2. Public APIs/interfaces: doc comment stating purpose + contract.
3. No redundant comments restating code.

## 8. Testing in Code

- Unit tests co-located (`*_test.go`, `*.test.ts`, `test_*.py`).
- Use fakes for ports (LLMProvider, Workspace) — never call real providers in unit tests.
- Deterministic seeds; no `Date.now()`/`Math.random()` in test logic without injection.

## 9. Approval

| Role | Name | Accept? | Date |
|------|------|---------|------|
| Head of Engineering | __________ | ☐ Yes ☐ No | ______ |

*End of Coding Standards v1.0-draft.*
