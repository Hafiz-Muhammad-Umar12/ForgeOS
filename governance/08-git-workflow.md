# DevOS — Git Workflow

> **Status:** RELEASE CANDIDATE — For Ratification (no production code until ratified)
> **Version:** 1.0-rc1
> **Owner:** Head of Engineering
> **Companion:** Repository Standards (`06-repository-standards.md`), Coding (`05-coding-standards.md`), Release (`09-release-strategy.md`)

---

## 1. Model: Trunk-Based Development

- `main` is the **single long-lived branch** and is always releasable.
- Feature work happens on **short-lived branches** (≤ ~3 days) branched from `main`.
- Long-lived feature branches are discouraged; use **feature flags** (Eng §9) for in-progress work.
- **No** GitFlow `develop`/`release` long branches; releases are tags from `main` (see Release §9).

## 2. Branch Naming

```
feat/<slug>      feat/intent-ingress-discord
fix/<slug>       fix/workspace-coldstart
chore/<slug>     chore/bump-deps
docs/<slug>      docs/adr-009-observability
adr/<nnn-slug>   adr/009-observability-pipeline
```
Slug is lowercase kebab-case, descriptive.

## 3. Commits: Conventional Commits

```
<type>(<scope>): <subject>
```
Types: `feat`, `fix`, `docs`, `refactor`, `test`, `chore`, `perf`, `build`, `ci`, `revert`.
Scope = module (e.g., `agent-runtime`, `web`, `specs`). Subject ≤ 72 chars, imperative ("add", not "added").
Body explains **why**; references issue/PRD/ADR (`Refs: PRD §7`, `ADR-003`).

## 4. Pull Requests

1. Branch → `main`; PR title follows Conventional Commits.
2. PR description: **what/why**, linked PRD/ADR/spec section, test plan, screenshots for UI.
3. **Required:** CI green (lint, format, unit/integration, build, scripted agent-behavior, security scan).
4. **Required:** ≥ 1 approving review from CODEOWNERS of touched dirs.
5. **DCO / sign-off:** each commit signed-off (`git commit -s`); bot enforces.
6. **Squash merge** to `main` (one logical commit per PR); `main` history stays linear and readable.
7. PRs merged by the author only after approval; no self-merge of unreviewed PRs.

## 5. Hotfixes

- Branch `hotfix/<slug>` from the released tag; fix; PR to `main`; after merge, **cherry-pick** to the release tag and re-tag (Release §9).
- Hotfixes bypass feature-flag defaults only with on-call + eng-lead sign-off.

## 6. Protected Branches

- `main` protected: required reviews, required status checks, no force-push, no direct push.
- `governance/`, `specs/`, `core/` need additional designated-owner approval.

## 7. Issue Tracking

- Every PR references an issue (or PRD/ADR section). Orphan PRs discouraged.
- Labels: `bug`, `feature`, `spec`, `adr`, `docs`, `infra`, `security`, `blocked`.

## 8. Approval

| Role | Name | Accept? | Date |
|------|------|---------|------|
| Head of Engineering | __________ | ☐ Yes ☐ No | ______ |

*End of Git Workflow v1.0-rc1.*
