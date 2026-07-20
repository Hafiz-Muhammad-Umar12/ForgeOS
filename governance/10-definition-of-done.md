# DevOS — Definition of Done (DoD)

> **Status:** RELEASE CANDIDATE — For Ratification (no production code until ratified)
> **Version:** 1.0-rc1
> **Owner:** Head of Engineering / Head of Product
> **Companion:** Engineering (`04-engineering-standards.md`), Coding (`05-coding-standards.md`), Git (`08-git-workflow.md`), Release (`09-release-strategy.md`)

---

## 1. Purpose

An item (feature, service, agent, fix) is **Done** only when **all** of the following are satisfied. Partial completion is "In Progress," never "Done."

## 2. The Checklist

### A. Requirements & Design
- [ ] Conforms to the relevant PRD section / engineering spec / ADR.
- [ ] Acceptance criteria (from PRD user stories or ticket) are explicit and met.
- [ ] New architectural decisions recorded as ADR (if applicable) and accepted.

### B. Implementation
- [ ] Code complete per scope; no stubs/TODOs left in shipping paths.
- [ ] Follows Coding Standards (`05`) and Repository Standards (`06`).
- [ ] Dependency direction respected (no cycles, ports/adapters).
- [ ] No hardcoded config or secrets; config injected; secrets via proxy (T4/T10).

### C. Quality
- [ ] Lint + format clean (gofmt/ruff/black/prettier/eslint/sqlfluff).
- [ ] Unit + integration tests pass; coverage ≥ 80% on new non-LLM logic.
- [ ] Agent behavior tests pass (scripted) where agents are involved.
- [ ] No flaky tests merged green.

### D. Review & Process
- [ ] PR approved by required CODEOWNERS (≥ 1).
- [ ] DCO sign-off present; Conventional Commits followed.
- [ ] CI green (build, scan/SBOM/vuln, all required checks).
- [ ] Squash-merged to `main`.

### E. Observability & Reliability (T6)
- [ ] OTel traces + metrics + structured logs emitted; `traceId` correlated.
- [ ] Agent spans present for any agent work; replayable.
- [ ] Dashboard panel + SLO alert added for the feature.
- [ ] Performance/reliability budgets met (ACK window, cold-start, stream render).

### F. Security (T10)
- [ ] No raw secrets logged or exposed; secret-proxy used where needed.
- [ ] AuthZ/scope enforced; audit log on mutations.
- [ ] Dependency scan clean; new deps allowlisted.
- [ ] Security review done for auth/secrets/network-egress changes.

### G. Documentation
- [ ] Relevant docs updated in same PR (README, spec, runbook, API docs).
- [ ] Diagrams (Mermaid) updated if architecture changed.

### H. Delivery Readiness
- [ ] Feature flag in place (default off) if risky/large.
- [ ] Rollback path verified (image tag, migration reversible).
- [ ] For UI: accessibility (WCAG 2.2 AA) spot-checked.
- [ ] For releases: changelog/release notes drafted.

### I. Product Sign-off
- [ ] Acceptance criteria demoed/validated (product or designated reviewer).
- [ ] For MVP-critical flows: end-to-end journey verified (e.g., NL → deploy → notify).

## 3. Exceptions

- Exceptions to any item require **explicit, recorded waiver** by the item's owner + Head of Engineering (e.g., temporarily lowering coverage on a spike, with a follow-up ticket).

## 4. Approval

| Role | Name | Accept? | Date |
|------|------|---------|------|
| Head of Engineering | __________ | ☐ Yes ☐ No | ______ |
| Head of Product | __________ | ☐ Yes ☐ No | ______ |

*End of Definition of Done v1.0-rc1.*
