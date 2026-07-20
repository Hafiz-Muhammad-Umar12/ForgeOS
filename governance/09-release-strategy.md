# DevOS — Release Strategy

> **Status:** DRAFT — For Approval (no production code until ratified)
> **Version:** 1.0-draft
> **Owner:** Release Engineering / CTO
> **Companion:** Git Workflow (`08-git-workflow.md`), Engineering (`04-engineering-standards.md`), DoD (`10-definition-of-done.md`)

---

## 1. Versioning

| Artifact | Scheme | Example |
|----------|--------|---------|
| Services / platform | **SemVer** `MAJOR.MINOR.PATCH` | `1.4.2` |
| Container images | Git **SHA** tag (+ SemVer label) | `devos/agent-runtime:ab12cde` |
| Public REST API | **URL version** `/v1` (contract-stable) | — |
| Infrastructure | Tagged with platform version it ships | — |

- **MAJOR:** breaking API/contract change. **MINOR:** backward-compatible feature. **PATCH:** fix.
- API contract changes are backward-compatible ≥ 12 months (Spec `/specs/02-specification/01-api-contracts.md` §11).

## 2. Environments

```
dev → staging (pre-prod) → prod
```
- `dev`: continuous deploy from `main` (post-merge).
- `staging`: mirrors prod; full e2e + chaos nightly.
- `prod`: promoted only via GitOps (ArgoCD) from a **release tag**.

## 3. GitOps & Progressive Delivery

- Kubernetes manifests in `/infra/k8s`, reconciled by ArgoCD from git (GitOps).
- Releases use **canary** rollout (10% → 50% → 100%) with automated analysis on error-rate/SLA (Deploy Spec `/specs/09-deployment/`).
- **Auto-rollback** on SLO breach during rollout.

## 4. Feature Flags

- Risky/large features ship **default-off** (Eng §9); enabled per project/tenant after validation.
- Flags have an owner + removal plan; stale flags removed within 2 releases.

## 5. Release Cadence (indicative)

- **MVP (v1):** one coordinated GA release (see PRD roadmap).
- **Post-MVP:** minor releases every 2–4 weeks; patches as needed; majors rare (annual-ish).
- **Hotfix:** same-day, cherry-pick + retag (Git Workflow §5).

## 6. Release Process

1. Cut release branch/tag from `main` at a green `staging` state.
2. Changelog + Release Notes generated (Docs §2).
3. Promote to prod via GitOps canary; monitor SLO dashboards + alerts.
4. Announce (internal + users for minors/majors); update PRD/specs if behavior changed.

## 7. Rollback & Recovery

- Any release rollbackable to previous image tag (immutable images).
- DB migrations: expand/contract only; never destructive in a single release (DB Spec `/specs/04-database/`).
- RTO < 15 min, RPO < 5 min (Deploy Spec §9).

## 8. Communication

- Release Notes summarize user-impacting changes, breaking changes, and upgrade notes.
- Internal changelog for engineering; customer-facing notes for minors/majors.

## 9. Approval

| Role | Name | Accept? | Date |
|------|------|---------|------|
| Release Engineering | __________ | ☐ Yes ☐ No | ______ |
| CTO | __________ | ☐ Yes ☐ No | ______ |

*End of Release Strategy v1.0-draft.*
