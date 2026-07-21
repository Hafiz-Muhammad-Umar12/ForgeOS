# DevOS — Implementation Planning Package

> **Status:** DRAFT — For Approval (no production code until ratified)
> **Version:** 1.0-draft
> **Owner:** CTO / Head of Engineering
> **Companion:** Governance (`/governance/`), PRD (`/product/PRD.md`), Specs (`/specs/`), SDD (`/sdd/`)
> **Role:** The master execution plan. Converts approved governance + PRD + architecture + SDD into an executable engineering roadmap.

---

## Traceability Convention (mandatory)

Every implementation task in this package references **three anchors**:

```
[PRD:<ref>][ADR-<nnn>][SDD:§<nn>]
```

- `PRD:<ref>` → a PRD requirement (Goal `G1`–`G6`, Functional `FR-1`–`FR-18`, MVP scope `§6.1`, Decision `D1`–`D10`, KPI `§11`).
- `ADR-<nnn>` → an Architecture Decision Record (`001`–`008`).
- `SDD:§<nn>` → an SDD service section (`01` Intent Ingress … `11` Kernel).

Example: `Build Intent Ingress webhook receiver [PRD:FR-1][ADR-006][SDD:§01]`.

This satisfies the binding rule: **no task without PRD + ADR + SDD references.**

## Package Contents (20 + index)

| # | Artifact | File | Batch |
|---|----------|------|-------|
| — | Master Index (this) | `README.md` | — |
| 1 | Product Milestones | [01-product-milestones.md](01-product-milestones.md) | **1** |
| 2 | Engineering Epics | [02-engineering-epics.md](02-engineering-epics.md) | **1** |
| 3 | Build Order | [03-build-order.md](03-build-order.md) | **1** |
| 4 | Dependency Graph | [04-dependency-graph.md](04-dependency-graph.md) | **1** |
| 5 | Sprint Planning | [05-sprint-planning.md](05-sprint-planning.md) | 2 |
| 6 | Team Responsibilities | [06-team-responsibilities.md](06-team-responsibilities.md) | 2 |
| 7 | Risk Register | [07-risk-register.md](07-risk-register.md) | 2 |
| 8 | Acceptance Criteria | [08-acceptance-criteria.md](08-acceptance-criteria.md) | 2 |
| 9 | Sprint Deliverables | [09-sprint-deliverables.md](09-sprint-deliverables.md) | 3 |
| 10 | Milestone Exit Criteria | [10-milestone-exit-criteria.md](10-milestone-exit-criteria.md) | 3 |
| 11 | Technical Backlog | [11-technical-backlog.md](11-technical-backlog.md) | 3 |
| 12 | Product Backlog | [12-product-backlog.md](12-product-backlog.md) | 3 |
| 13 | Infrastructure Backlog | [13-infrastructure-backlog.md](13-infrastructure-backlog.md) | 3 |
| 14 | Documentation Backlog | [14-documentation-backlog.md](14-documentation-backlog.md) | 4 |
| 15 | Testing Backlog | [15-testing-backlog.md](15-testing-backlog.md) | 4 |
| 16 | Security Backlog | [16-security-backlog.md](16-security-backlog.md) | 4 |
| 17 | CI/CD Roadmap | [17-cicd-roadmap.md](17-cicd-roadmap.md) | 4 |
| 18 | Deployment Roadmap | [18-deployment-roadmap.md](18-deployment-roadmap.md) | 4 |
| 19 | Repository Initialization Plan | [19-repository-init-plan.md](19-repository-init-plan.md) | 4 |
| 20 | Development Checklist | [20-development-checklist.md](20-development-checklist.md) | 4 |

## Milestone Summary

| Milestone | Sprints | Outcome |
|-----------|---------|---------|
| M1 Foundation | S0–S1 | Repo, CI, bus, kernel skeleton, auth, Claude adapter, one workspace |
| M2 MVP | S2–S10 | Full intent→plan→agents→deploy→notify (Web+REST+Discord, 8 agents, 2 providers, TS stack, Vercel+Fly) |
| M3 Beta | S11–S12 | Private beta: hardening, eval, chaos, security review |
| M4 Public Launch | S13–S14 | GA: canary, scale, monitoring |
| M5 Enterprise | S15+ | SAML, audit, more channels/providers/agents, A2A, voice |

## Batch Status
- **Batch 1 (this review):** README (this update), 01 Milestones, 02 Epics, 03 Build Order, 04 Dependency Graph.
- **Batch 2+ (pending your approval):** 05 Sprint Planning, 06 Team, 07 Risk, 08 Acceptance; then 09–20 backlogs/roadmaps/checklist.

## Approval
| Role | Name | Approve? | Date |
|------|------|---------|------|
| CTO | __________ | ☐ Yes ☐ No | ______ |
| Head of Engineering | __________ | ☐ Yes ☐ No | ______ |

*Once approved, execution begins per Sprint 0. No production code before then.*
